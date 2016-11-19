package signup

import (
	"os"
	"os/signal"
)

const (
	resourcePath = "data/xdmsc.club/signup"
)

var cleanUp = make(chan os.Signal)

func init() {
	signal.Notify(cleanUp, os.Interrupt, os.Kill)
}

// 用于清理数据库相关的操作
func CleanUp() {
	close(cleanUp)
}
