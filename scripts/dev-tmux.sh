#!/usr/bin/env bash
# Запускает nvim, inventory, order, payment и lazygit в одном окне tmux.
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

# Слева (50% ширины): nvim на всю высоту.
# Справа (50%), сверху вниз: inventory, order, payment, lazygit.
NVIM=$(tmux display-message -p -t "$SESSION:services.0" '#{pane_id}')
INVENTORY=$(tmux split-window -h -p 50 -t "$NVIM" -c "$ROOT" -P -F '#{pane_id}')
LAZYGIT=$(tmux split-window -v -p 25 -t "$INVENTORY" -c "$ROOT" -P -F '#{pane_id}')
PAYMENT=$(tmux split-window -v -p 33 -t "$INVENTORY" -c "$ROOT" -P -F '#{pane_id}')
ORDER=$(tmux split-window -v -p 50 -t "$INVENTORY" -c "$ROOT" -P -F '#{pane_id}')

tmux select-pane -t "$NVIM" -T nvim
tmux select-pane -t "$INVENTORY" -T inventory
tmux select-pane -t "$ORDER" -T order
tmux select-pane -t "$PAYMENT" -T payment
tmux select-pane -t "$LAZYGIT" -T lazygit

tmux send-keys -t "$NVIM" "nvim ." C-m
tmux send-keys -t "$INVENTORY" "go run ./inventory/cmd" C-m
tmux send-keys -t "$PAYMENT" "go run ./payment/cmd" C-m
# order ходит в inventory и payment по gRPC при старте — даём им подняться
tmux send-keys -t "$ORDER" "sleep 3 && go run ./order/cmd" C-m
tmux send-keys -t "$LAZYGIT" "lazygit" C-m

tmux select-pane -t "$NVIM"
exec tmux attach -t "$SESSION"
