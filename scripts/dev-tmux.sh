#!/usr/bin/env bash
# Запускает inventory, payment, order и свободную оболочку в одном окне tmux.
# Использование: scripts/dev-tmux.sh [start|stop]
set -euo pipefail

SESSION="goshop"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

case "${1:-start}" in
stop)
  tmux kill-session -t "$SESSION" 2>/dev/null && echo "Сессия $SESSION остановлена" || echo "Сессия $SESSION не запущена"
  exit 0
  ;;
start) ;;
*)
  echo "Использование: $0 [start|stop]" >&2
  exit 1
  ;;
esac

if tmux has-session -t "$SESSION" 2>/dev/null; then
  exec tmux attach -t "$SESSION"
fi

tmux new-session -d -s "$SESSION" -c "$ROOT" -n services
tmux set-option -t "$SESSION" -g mouse on
tmux set-option -t "$SESSION" pane-border-status top

# 2x2: inventory | payment
#      order     | shell
tmux split-window -h -t "$SESSION:services.0" -c "$ROOT"
tmux split-window -v -t "$SESSION:services.0" -c "$ROOT"
tmux split-window -v -t "$SESSION:services.2" -c "$ROOT"

# Порядок панелей: 0 inventory, 1 order, 2 payment, 3 shell
tmux select-pane -t "$SESSION:services.0" -T inventory
tmux select-pane -t "$SESSION:services.1" -T order
tmux select-pane -t "$SESSION:services.2" -T payment
tmux select-pane -t "$SESSION:services.3" -T shell

tmux send-keys -t "$SESSION:services.0" "go run ./inventory/cmd" C-m
tmux send-keys -t "$SESSION:services.2" "go run ./payment/cmd" C-m
# order ходит в inventory и payment по gRPC при старте — даём им подняться
tmux send-keys -t "$SESSION:services.1" "sleep 3 && go run ./order/cmd" C-m
tmux send-keys -t "$SESSION:services.3" "# task test-api — проверка API; scripts/dev-tmux.sh stop — остановить всё" C-m

tmux select-pane -t "$SESSION:services.3"
exec tmux attach -t "$SESSION"
