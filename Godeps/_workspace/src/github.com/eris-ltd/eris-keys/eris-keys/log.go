package keys

import (
	. "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/log"
)

var logger *Logger

func initLog() {
	logger = AddLogger("commands")
}
