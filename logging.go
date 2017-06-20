package lazycache

import (
	kitlog "github.com/go-kit/kit/log"
	"os"
)

var Logger kitlog.Logger

func init() {
	Logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
}
