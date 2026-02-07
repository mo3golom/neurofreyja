package time

import "time"

func DeleteAtMinutes(minutes int) time.Time {
	return time.Now().UTC().Add(time.Duration(minutes) * time.Minute)
}
