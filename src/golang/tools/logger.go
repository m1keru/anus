package tools

import (
	"log"
	dbg "runtime/debug"
)

//CheckErr - проверка на ошибки с выводом stacktrace
func CheckErr(err error) {
	if err != nil {
		dbg.PrintStack()
		log.Printf("ERROR:: %s", err.Error())
	}
}
