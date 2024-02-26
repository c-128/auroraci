package cron

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var dayOfWeeks = map[string]int{
	"sun": 0, "sunday": 0,
	"mon": 1, "monday": 1,
	"tue": 2, "tuesday": 2,
	"wed": 3, "wednesday": 3,
	"thu": 4, "thursday": 4,
	"fri": 5, "friday": 5,
	"sat": 6, "saturday": 6,
}

func IsExpressionTime(expr string, time time.Time) (bool, error) {
	expr = strings.TrimSpace(expr)
	fields := strings.Split(expr, " ")
	if len(fields) != 3 {
		return false, errors.New("more than 3 fields")
	}

	isMinute := isField(fields[0], time.Minute(), nil)
	isHour := isField(fields[1], time.Hour(), nil)
	isDayOfWeek := isField(fields[2], int(time.Weekday()), dayOfWeeks)
	isNow := isMinute && isHour && isDayOfWeek

	return isNow, nil
}

func isField(field string, value int, keywords map[string]int) bool {
	if keywords != nil {
		keywordValue, found := keywords[field]
		if found {
			return keywordValue == value
		}
	}

	switch field {
	case "*":
		return true
	default:
		fieldInt, err := strconv.Atoi(field)
		if err != nil {
			return false
		}

		return fieldInt == value
	}
}
