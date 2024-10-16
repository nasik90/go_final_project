package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func checkNextDateArgs(date string, repeat string) error {

	if repeat == "" {
		return errors.New("repeat is empty")
	}

	_, err := time.Parse(DateTemplate, date)
	if err != nil {
		return err
	}

	repeatSplitted := strings.Split(repeat, " ")
	firstSymbol := repeatSplitted[0]
	secondSymbol := ""
	thirdSymbol := ""
	if len(repeatSplitted) > 1 {
		secondSymbol = repeatSplitted[1]
	}
	if len(repeatSplitted) > 2 {
		thirdSymbol = repeatSplitted[2]
	}

	if !strings.Contains("ydwm", firstSymbol) {
		return errors.New("repeat format error")
	}

	if firstSymbol == "d" {
		err = checkDateArgs(firstSymbol, secondSymbol)
	} else if firstSymbol == "w" {
		err = checkWeekArgs(firstSymbol, secondSymbol)
	} else if firstSymbol == "m" {
		err = checkMonthArgs(firstSymbol, secondSymbol, thirdSymbol)
	}

	return err
}

func checkDateArgs(firstSymbol, secondSymbol string) error {

	if secondSymbol == "" {
		return errors.New("no day in repeat")
	}
	secondSymbolInt, err := strconv.Atoi(secondSymbol)
	if err != nil {
		return err
	}
	if secondSymbolInt > 400 {
		return errors.New("day more than 400")
	}

	return nil
}

func checkWeekArgs(firstSymbol, secondSymbol string) error {

	if secondSymbol == "" {
		return errors.New("no month in repeat")
	}
	mWeekDays := strings.Split(secondSymbol, ",")
	for _, v := range mWeekDays {
		vInt, err := strconv.Atoi(v)
		if err != nil || vInt < 1 || vInt > 7 {
			return errors.New("week day is not valide")
		}
	}
	return nil
}

func checkMonthArgs(firstSymbol, secondSymbol, thirdSymbol string) error {

	if secondSymbol == "" {
		return errors.New("no days in repeat")
	}
	mMonthDays := strings.Split(secondSymbol, ",")
	for _, v := range mMonthDays {
		vInt, err := strconv.Atoi(v)
		if err != nil || vInt < -2 || vInt > 31 || vInt == 0 {
			return errors.New("day is not valide")
		}
	}
	if thirdSymbol != "" {
		mYears := strings.Split(thirdSymbol, ",")
		for _, v := range mYears {
			vInt, err := strconv.Atoi(v)
			if err != nil || vInt < 1 || vInt > 12 {
				return errors.New("month is not valide")
			}
		}
	}
	return nil
}
