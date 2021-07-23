package hlog

import (
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"hanaBFT/utils"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var ID peer.ID

type logger struct {
	sync.Mutex
	debug *log.Logger
	warn  *log.Logger
	error *log.Logger
	fatal *log.Logger
	sign  *log.Logger
	level int
}

var HLOG logger

func init() {
	os.Mkdir("logs", os.ModePerm)
}

func Setup(p peer.ID) {
	ID = p
	format := log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	fname := fmt.Sprintf("%s.%s.log", filepath.Base(os.Args[0]), utils.ShortPeerID(p))
	f, err := os.Create(filepath.Join("logs", fname))
	if err != nil {
		panic(err)
	}
	HLOG.debug = log.New(f, "[DEBUG]", format)
	HLOG.warn = log.New(f, "[WARN]", format)
	HLOG.error = log.New(f, "[ERROR]", format)
	HLOG.fatal = log.New(f, "[FATAL]", format)
	HLOG.sign = log.New(f, "[SIGN]", format)
}

func Debug(v ...interface{}) {
	HLOG.debug.Output(2, fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	HLOG.Lock()
	defer HLOG.Unlock()
	HLOG.debug.Output(2, fmt.Sprintf(format, v...))
}

func Warn(v ...interface{}) {
	HLOG.warn.Output(2, fmt.Sprint(v...))
}

func Warnf(format string, v ...interface{}) {
	HLOG.warn.Output(2, fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	HLOG.error.Output(2, fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	HLOG.Lock()
	defer HLOG.Unlock()
	HLOG.error.Output(2, fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	HLOG.fatal.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	HLOG.Lock()
	defer HLOG.Unlock()
	HLOG.fatal.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Sign(v ...interface{}) {
	HLOG.sign.Output(2, fmt.Sprint(v...))
}

func Signf(format string, v ...interface{}) {
	HLOG.sign.Output(2, fmt.Sprintf(format, v...))
}
