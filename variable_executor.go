package inventory_notifier

import (
	"time"
)

type VariableExecutor struct {
	ticker   time.Ticker
	done     chan bool
	callback func()
}

func (ve *VariableExecutor) Run() {
}
