package worker

import (
	"fmt"
	"time"
)

func ExecuteDailyDeleteWorker() {
	ticker := time.Tick(24 * time.Hour)

	fmt.Println(ticker)
}
