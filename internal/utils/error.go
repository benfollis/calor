package utils

import "log"

// CheckPanic tests if the error is nil
// and if not panics
func CheckPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// CheckLog tests if the error is nil
// and if not logs a fatal error
func CheckLog(err error) {
	if err != nil {
		log.Println(err);
	}
}
