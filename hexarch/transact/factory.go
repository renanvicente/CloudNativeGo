package transact

import (
	"EllenCourse/CloudNativeGo/Chapter8/hexarch/core"
	"os"
)

func NewTransactionLogger(logger string) (core.TransactionLogger, error) {
	switch logger {
	case "file":
		return NewFileTransactionLogger(os.Getenv("TLOG_FILENAME"))
	}
}
