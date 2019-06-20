package common

import (
	"bytes"
	"errors"
	"fmt"
	debug "server/common/debug"
	log "server/leaf/log"
)

func PrintBuf(buf *bytes.Buffer) {
	fmt.Printf("\tbuf.Len() == %d\n", buf.Len())
	fmt.Printf("\tbuf.Cap() == %d\n", buf.Cap())
	fmt.Printf("\tbuf.String() == '%s'\n", buf.String())
	fmt.Println("")
}

func PrintLog(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println("")
}

func LogError(format string, v ...interface{}) {
	//debug.Print(v...)
	st := debug.StackTrace(0)
	log.Error(st.String("    "))
	log.Error(format, v...)
}

func LogDebug(format string, v ...interface{}) {
	log.Debug(format, v...)
}

func LogInfo(format string, v ...interface{}) {
	log.Release(format, v...)
}

func LogPanic(format string, v ...interface{}) {
	st := debug.StackTrace(0)
	log.Error(st.String("    "))
	panic(errors.New(fmt.Sprintf(format, v...)))
}

func LogFatal(format string, v ...interface{}) {
	st := debug.StackTrace(0)
	log.Error(st.String("    "))
	log.Fatal(format, v...)
}
