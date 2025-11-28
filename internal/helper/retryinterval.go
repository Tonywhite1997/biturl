package helper

import "time"

func RetryInterval(stage int) time.Duration {
	switch stage {
	case 0:
		return 5 * time.Second
	case 1:
		return 30 * time.Second
	case 2:
		return 5 * time.Minute
	default:
		return 1 * time.Second
	}
}
