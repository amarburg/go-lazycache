package lazycache

import (
	kitlog  "github.com/go-kit/kit/log"
  "os"
)

var DefaultLogger kitlog.Logger

func init() {
	 DefaultLogger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
}
