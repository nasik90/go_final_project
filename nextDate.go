package main

import (
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	err := checkNextDateArgs(date, repeat)
	if err != nil {
		return "", err
	}

	dateTime, _ := time.Parse(DateTemplate, date)

	m_repeat := strings.Split(repeat, " ")

	firstSymbol := m_repeat[0]

	var newDate time.Time

	if firstSymbol == "y" {
		newDate = dateTime.AddDate(1, 0, 0)
		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
	} else if firstSymbol == "d" {
		newDate = nextDateD(now, dateTime, m_repeat)
	} else if firstSymbol == "w" {
		newDate = nextDateW(now, dateTime, m_repeat)
	} else if firstSymbol == "m" {
		newDate = nextDateM(now, dateTime, m_repeat)
	}

	return newDate.Format(DateTemplate), nil
}

func nextDateD(now, dateTime time.Time, m_repeat []string) time.Time {
	days := m_repeat[1]
	daysInt, _ := strconv.Atoi(days)
	newDate := dateTime.AddDate(0, 0, daysInt)
	for newDate.Before(now) {
		newDate = newDate.AddDate(0, 0, daysInt)
	}
	return newDate
}

func nextDateW(now, dateTime time.Time, m_repeat []string) time.Time {
	days := m_repeat[1]
	mWeekDays := strings.Split(days, ",")
	var mWeekDaysInt []int
	for _, v := range mWeekDays {
		vInt, _ := strconv.Atoi(v)
		if vInt == 7 {
			vInt = 0 // воскресенье 0 день
		}
		mWeekDaysInt = append(mWeekDaysInt, vInt)
	}
	newDate := dateTime
	weekday := int(newDate.Weekday())
	for {
		if newDate.After(now) && contains(mWeekDaysInt, weekday) {
			break
		}
		newDate = newDate.AddDate(0, 0, 1)
		weekday = int(newDate.Weekday())
	}
	return newDate
}

func nextDateM(now, dateTime time.Time, m_repeat []string) time.Time {
	var (
		mDaysInt   []int
		months     string
		mMonths    []string
		mMonthsInt []int
	)

	days := m_repeat[1]
	mDays := strings.Split(days, ",")
	for _, v := range mDays {
		vInt, _ := strconv.Atoi(v)
		mDaysInt = append(mDaysInt, vInt)
	}

	if len(m_repeat) > 2 {
		months = m_repeat[2]
		mMonths = strings.Split(months, ",")
		for _, v := range mMonths {
			vInt, _ := strconv.Atoi(v)
			mMonthsInt = append(mMonthsInt, vInt)
		}
	}

	newDate := dateTime
	for {
		if newDate.After(now) {
			nextMonth := time.Date(newDate.Year(), newDate.Month()+1, 1, 0, 0, 0, 0, time.UTC)
			dayContain := contains(mDaysInt, newDate.Day())
			day_1Contain := contains(mDaysInt, -1) && newDate == nextMonth.AddDate(0, 0, -1)
			day_2Contain := contains(mDaysInt, -2) && newDate == nextMonth.AddDate(0, 0, -2)
			monthContain := len(mMonthsInt) == 0 || contains(mMonthsInt, int(newDate.Month()))
			if monthContain && (dayContain || day_1Contain || day_2Contain) {
				break
			}
		}
		newDate = newDate.AddDate(0, 0, 1)
	}
	return newDate
}

func contains(s []int, str int) bool {
	for _, el := range s {
		if el == str {
			return true
		}
	}
	return false
}
