package worker

import "time"

// RetryDelay returns how long to wait before next attempt
// attempt is 1-indexed (1 = first failure)
func RetryDelay(attempt int) time.Duration {
	switch attempt {
	case 1:
		return 5 * time.Second
	case 2:
		return 30 * time.Second
	case 3:
		return 2 * time.Minute
	case 4:
		return 10 * time.Minute
	default:
		return 0 // dead letter
	}
}