#!/usr/bin/env bash
# Рабочее окружение goshop в tmux.
#
# Использование: scripts/dev-tmux.sh [start|local|stop]
#
#   start  (по умолчанию) — сервисы работают в контейнерах (task up-services),
#          в tmux три окна:
#            code     — nvim, lazygit и свободный shell;
#            services — логи iam, inventory, payment, order, assembly, notification
#                       (docker logs -f, JSON-строки сворачиваются в «время уровень сообщение»);
#            ops      — shell и шпаргалка по адресам UI (Kibana, Grafana, Jaeger, ...).
#   local  — старый режим: сервисы запускаются через `go run` прямо в панелях одного окна
#          (контейнеры сервисов при этом должны быть остановлены: task down-services).
#   stop   — убить tmux-сессию.
set -euo pipefail

SESSION="${TMUX_SESSION:-goshop}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SERVICES=(iam inventory payment order assembly notification)

attach() {
  # NO_ATTACH=1 нужен для проверки скрипта без интерактивного терминала
  if [ "${NO_ATTACH:-}" = "1" ]; then
    return 0
  fi
  exec tmux attach -t "$SESSION"
}

# Создаёт в окне $1 стопку из N панелей справа от $2 (pane id) равной высоты.
# Печатает id созданных панелей по одному в строке, сверху вниз.
make_column() {
  local first=$1 count=$2 current=$1 k m pct new
  echo "$first"
  for ((k = 1; k < count; k++)); do
    m=$((count - k + 1)) # сколько панелей осталось разместить, включая текущую
    pct=$((100 * (m - 1) / m))
    new=$(tmux split-window -v -p "$pct" -t "$current" -c "$ROOT" -P -F '#{pane_id}')
    echo "$new"
    # новая панель — нижняя часть; делить дальше нужно именно её
    current=$new
  done
}

new_session() {
  tmux new-session -d -s "$SESSION" -c "$ROOT" -n "$1"
  tmux set-option -t "$SESSION" -g mouse on
  tmux set-option -t "$SESSION" -g pane-border-status top
  # Имя панели хранится в собственной опции @name: заголовок (pane_title) перезаписывает shell
  # при запуске команды, и вместо «order» на рамке оказалась бы строка docker logs ...
  tmux set-option -t "$SESSION" -g pane-border-format ' #{?@name,#{@name},#{pane_title}} '
}

# name_pane <pane> <имя> — подписывает панель на рамке
name_pane() {
  tmux set-option -p -t "$1" @name "$2"
}

# Команда просмотра логов контейнера: JSON zap → «ЧЧ:ММ:СС УРОВЕНЬ сообщение»,
# не-JSON строки (goose, лог Telegram-библиотеки) выводятся как есть.
# `. as $l`: внутри catch точка — это текст ошибки, а не исходная строка.
logs_cmd() {
  local name=$1
  if command -v jq >/dev/null 2>&1; then
    # shellcheck disable=SC2016 # фильтр jq — в одинарных кавычках намеренно
    printf '%s' "docker logs -f --tail 30 $name 2>&1 | jq --unbuffered -Rr '. as \$l | try (fromjson | \"\\(.timestamp[11:19]) \\(.level) \\(.message)\") catch \$l'"
  else
    printf '%s' "docker logs -f --tail 30 $name 2>&1"
  fi
}

start_docker() {
  new_session code

  # Окно code: слева nvim, справа lazygit (70%) и shell (30%)
  local nvim lazygit shell
  nvim=$(tmux display-message -p -t "$SESSION:code.0" '#{pane_id}')
  lazygit=$(tmux split-window -h -p 50 -t "$nvim" -c "$ROOT" -P -F '#{pane_id}')
  shell=$(tmux split-window -v -p 30 -t "$lazygit" -c "$ROOT" -P -F '#{pane_id}')
  name_pane "$nvim" nvim
  name_pane "$lazygit" lazygit
  name_pane "$shell" shell
  tmux send-keys -t "$nvim" "nvim ." C-m
  tmux send-keys -t "$lazygit" "lazygit" C-m

  # Окно services: по панели на каждый контейнер, сетка 2x3
  tmux new-window -t "$SESSION" -n services -c "$ROOT"
  local first
  first=$(tmux display-message -p -t "$SESSION:services.0" '#{pane_id}')
  local panes=("$first")
  local i
  for ((i = 1; i < ${#SERVICES[@]}; i++)); do
    panes+=("$(tmux split-window -t "$SESSION:services" -c "$ROOT" -P -F '#{pane_id}')")
    tmux select-layout -t "$SESSION:services" tiled >/dev/null
  done
  for i in "${!SERVICES[@]}"; do
    name_pane "${panes[$i]}" "${SERVICES[$i]}"
    tmux send-keys -t "${panes[$i]}" "$(logs_cmd "${SERVICES[$i]}")" C-m
  done

  # Окно ops: shell и шпаргалка
  tmux new-window -t "$SESSION" -n ops -c "$ROOT"
  name_pane "$SESSION:ops.0" ops
  tmux send-keys -t "$SESSION:ops.0" "clear; cat <<'EOF'
Стек:        task up-all | down-all (down-all удаляет тома!)
Сервисы:     task up-services | down-services   (пересборка: task up-services)
Тесты:       task test-api | test-alert | test-api-auth-rejection

Order API    http://localhost:8080          Kibana    http://localhost:5601
Grafana      http://localhost:3000 (admin)  Jaeger    http://localhost:16686
Prometheus   http://localhost:9090          Alertmgr  http://localhost:9093
EOF" C-m

  tmux select-window -t "$SESSION:code"
  tmux select-pane -t "$nvim"
  attach
}

# Старый режим: всё в одном окне, сервисы через go run.
# METRICS_ADDR разный у каждого процесса: иначе все займут :9100 на хосте.
# Трейсинг включён у order и payment (Jaeger должен быть запущен: task up-tracing).
start_local() {
  if docker ps --format '{{.Names}}' 2>/dev/null | grep -qx -E 'iam|inventory|payment|order|assembly|notification'; then
    echo "Контейнеры сервисов запущены и заняли порты. Остановите их: task down-services" >&2
    exit 1
  fi

  new_session services

  local nvim
  nvim=$(tmux display-message -p -t "$SESSION:services.0" '#{pane_id}')
  local first
  first=$(tmux split-window -h -p 50 -t "$nvim" -c "$ROOT" -P -F '#{pane_id}')

  # Справа сверху вниз: iam, inventory, payment, assembly, order, lazygit
  local names=(iam inventory payment assembly order lazygit)
  local panes=()
  while IFS= read -r p; do panes+=("$p"); done < <(make_column "$first" "${#names[@]}")

  name_pane "$nvim" nvim
  local i
  for i in "${!names[@]}"; do
    name_pane "${panes[$i]}" "${names[$i]}"
  done

  tmux send-keys -t "$nvim" "nvim ." C-m
  # iam: Postgres и Redis поднимаются отдельно (task up-iam); сначала он, остальные ходят к нему за сессиями
  tmux send-keys -t "${panes[0]}" "METRICS_ADDR=:9101 go run ./iam/cmd" C-m
  tmux send-keys -t "${panes[1]}" "METRICS_ADDR=:9102 go run ./inventory/cmd" C-m
  tmux send-keys -t "${panes[2]}" "METRICS_ADDR=:9103 OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317 go run ./payment/cmd" C-m
  # assembly нужен только Kafka (поднимается в core), стартует сразу
  tmux send-keys -t "${panes[3]}" "METRICS_ADDR=:9104 go run ./assembly/cmd" C-m
  # order ходит в iam, inventory и payment по gRPC — даём им подняться
  tmux send-keys -t "${panes[4]}" "sleep 5 && METRICS_ADDR=:9105 OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317 go run ./order/cmd" C-m
  tmux send-keys -t "${panes[5]}" "lazygit" C-m

  tmux select-pane -t "$nvim"
  attach
}

case "${1:-start}" in
stop)
  tmux kill-session -t "$SESSION" 2>/dev/null && echo "Сессия $SESSION остановлена" || echo "Сессия $SESSION не запущена"
  exit 0
  ;;
start | local) ;;
*)
  echo "Использование: $0 [start|local|stop]" >&2
  exit 1
  ;;
esac

if tmux has-session -t "$SESSION" 2>/dev/null; then
  attach
  exit 0
fi

case "${1:-start}" in
local) start_local ;;
*) start_docker ;;
esac
