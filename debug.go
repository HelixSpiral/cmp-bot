package main

import "log"

// debugPrint prints a message prepended with 'DEBUG: ' if DEBUG is set to true
func debugPrint(msg string, args ...interface{}) {
	if DEBUG {
		log.Printf(msg, args...)
	}
}
