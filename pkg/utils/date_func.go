package utils

import "time"

// expired accepts a `RFC3339`/`ISO 8601` formatted time and
// compares it to the current time to tell if it expired or not.
func Expired(now time.Time, v string) (bool, error) {
	// fmt.Println("NOW:", now.Format(time.RFC3339))
	// fmt.Println("YOU:", v)

	exp, err := time.ParseInLocation(time.RFC3339, v, now.Location())
	if err != nil {
		return false, err
	}

	if now.Sub(exp).Seconds() <= 0 {
		return false, nil
	}

	return true, nil
}
