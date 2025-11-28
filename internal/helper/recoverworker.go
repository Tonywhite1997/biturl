package helper

import (
	"fmt"
)

func RecoverWorker() {
	if r := recover(); r != nil {
		fmt.Println("worker panicked:", r)
	}
}
