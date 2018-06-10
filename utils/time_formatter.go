package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// StringTimeToTime ...
func StringTimeToTime(stringTime string) (time.Time, error) {
	dateParts := strings.Split(stringTime, "-")
	years, err := strconv.Atoi(dateParts[0])
	if err != nil {
		return time.Time{}, errors.New("Invalid value for year")
	}
	months, err := strconv.Atoi(dateParts[1])
	if err != nil {
		return time.Time{}, errors.New("Invalid value for month")
	}
	days, err := strconv.Atoi(dateParts[2])
	if err != nil {
		return time.Time{}, errors.New("Invalid value for day")
	}
	return time.Date(years, time.Month(months), days, 0, 0, 0, 0, time.UTC), nil
}
