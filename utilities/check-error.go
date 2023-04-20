package utilities

import (
	"log"
)

// A function to check error
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
