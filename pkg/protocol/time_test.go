package protocol

import (
	"testing"
	"time"
)

func TestTime_ToTime(t *testing.T) {
	wooTime := Time{0x27, 0x25, 0x18, 0x20, 0x05, 0x25}
	goTime, err := wooTime.ToTime()
	if err != nil {
		t.Errorf("error on convert: %s", err.Error())
		return
	}
	if goTime != time.Date(2025, time.May, 20, 18, 25, 27, 0, time.UTC) {
		t.Errorf("time is not correct")
	}
}

func TestTime_ToWooTime(t *testing.T) {
	goTime := time.Date(2025, time.May, 20, 18, 25, 27, 0, time.UTC)
	wooTime := ToWooTime(goTime)

	if wooTime != (Time{0x27, 0x25, 0x18, 0x20, 0x05, 0x25}) {
		t.Errorf("time is not correct")
	}
}
