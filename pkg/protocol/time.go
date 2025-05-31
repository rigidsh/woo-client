package protocol

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Time [6]byte

func (t Time) ToTime() (time.Time, error) {
	second, err := hexToNumber(t[0])
	if err != nil {
		return time.Time{}, errors.New("wrong seconds value")
	}
	minute, err := hexToNumber(t[1])
	if err != nil {
		return time.Time{}, errors.New("wrong minutes value")
	}
	hour, err := hexToNumber(t[2])
	if err != nil {
		return time.Time{}, errors.New("wrong hour value")
	}
	day, err := hexToNumber(t[3])
	if err != nil {
		return time.Time{}, errors.New("wrong day value")
	}
	month, err := hexToNumber(t[4])
	if err != nil {
		return time.Time{}, errors.New("wrong month value")
	}
	year, err := hexToNumber(t[5])
	if err != nil {
		return time.Time{}, errors.New("wrong year value")
	}

	return time.Date(2000+year, time.Month(month), day, hour, minute, second, 0, time.UTC), nil
}

func (t Time) String() string {
	tt, err := t.ToTime()
	if err != nil {
		panic("wrong time format: " + err.Error() + " " + fmt.Sprintf("% x", [6]byte(t)))
	}
	return tt.String()
}

func hexToNumber(b byte) (int, error) {
	hexNumber := fmt.Sprintf("%x", b)
	return strconv.Atoi(hexNumber)
}

func numberToHex(b int) (byte, error) {
	number := fmt.Sprintf("%d", b)
	result, err := strconv.ParseInt(number, 16, 8)
	if err != nil {
		return 0, err
	}
	return byte(result), nil
}

func ToWooTime(t time.Time) Time {
	second, _ := numberToHex(t.Second())
	minute, _ := numberToHex(t.Minute())
	hour, _ := numberToHex(t.Hour())
	day, _ := numberToHex(t.Day())
	month, _ := numberToHex(int(t.Month()))
	year, _ := numberToHex(t.Year() - 2000)

	return Time{second, minute, hour, day, month, year}
}
