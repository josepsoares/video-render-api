package utils

import (
	"fmt"
	"log"
)

func CheckError(msg string, err error) {
	if err != nil {
		log.Printf("❗ ERROR DETECTED ❗ %s => %s", msg, err)
	}
}

func FailOnError(msg string, err error) {
	if err != nil {
		log.Fatalf(fmt.Sprintf("❗❗ PANIC ❗❗ %s => %s", msg, err))
	}
}
