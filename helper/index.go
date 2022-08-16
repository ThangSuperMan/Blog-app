package helper

import "log"

func HaltOn(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
