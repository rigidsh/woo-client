package protocol

import (
	"bytes"
	"testing"
	"time"
)

func TestReadEvent(t *testing.T) {
	eventSuccessTest(
		"Stop recoding event test",
		[]byte{0x58, 0x00, 0x00, 0x01, 0x01, 0x80, 0xCC, 0xCE, 0x01, 0x00, 0x00, 0xF6, 0x01},
		RecordingEventType,
		&RecordingEvent{
			Start:          false,
			NumberOfEvents: 257,
			Counter:        30329984,
		},
		t)
	eventSuccessTest(
		"Start recoding event test",
		[]byte{0x58, 0x01, 0x00, 0x01, 0x01, 0x80, 0xCC, 0xCE, 0x01, 0x00, 0x00, 0xF6, 0x01},
		RecordingEventType,
		&RecordingEvent{
			Start:          true,
			NumberOfEvents: 257,
			Counter:        30329984,
		},
		t)
	eventSuccessTest(
		"Jump event test",
		[]byte{
			0x44, 0x02, 0x07, 0x00, 0x00, 0x1e, 0x19, 0x01, 0x00, 0x1b,
			0x20, 0x02, 0x00, 0x00, 0x3b, 0x01, 0x00, 0xb5, 0x4b, 0xfb,
			0xff, 0x24, 0x05, 0x00, 0x00, 0x11, 0x01, 0x28, 0x0f, 0x07,
			0x45, 0x08, 0x31, 0x05, 0x25, 0x8a, 0x11, 0x04, 0x00, 0x00,
			0x00, 0x58, 0x6c, 0x66, 0x05, 0xad, 0x17, 0x32, 0x08,
		},
		JumpEventType,
		&JumpEvent{
			JumpType:   0x02,
			JumpNumber: 7,
			JumpTime:   ToWooTime(time.Date(2025, 5, 31, 8, 45, 7, 0, time.UTC)),
		},
		t)
}

func eventSuccessTest(name string, srcData []byte, expectedEvenType EventType, expectedEvent interface{}, t *testing.T) {
	t.Run(name, func(t *testing.T) {
		eventType, event, err := ReadEvent(bytes.NewReader(srcData))
		if err != nil {
			t.Errorf("error on read: %s", err.Error())
			return
		}
		if eventType != expectedEvenType {
			t.Errorf("event type is not correct")
		}

		switch expectedEvenType {
		case JumpEventType:
			testJumpEvent(event, expectedEvent, t)
		case RecordingEventType:
			testRecordingEvent(event, expectedEvent, t)
		}
	})
}

func testRecordingEvent(testEvent interface{}, expectedEvent interface{}, t *testing.T) {
	event, ok := testEvent.(*RecordingEvent)
	if !ok {
		t.Errorf("event is not RecordingEvent")
		return
	}
	if event.Start != expectedEvent.(*RecordingEvent).Start {
		t.Errorf("recording is not correct")
	}
	if event.NumberOfEvents != expectedEvent.(*RecordingEvent).NumberOfEvents {
		t.Errorf("number of events is not correct")
	}
	if event.Counter != expectedEvent.(*RecordingEvent).Counter {
		t.Errorf("counter is not correct")
	}
}

func testJumpEvent(testEvent interface{}, expectedEvent interface{}, t *testing.T) {
	event, ok := testEvent.(*JumpEvent)
	if !ok {
		t.Errorf("event is not JumpEvent")
		return
	}
	if event.JumpType != expectedEvent.(*JumpEvent).JumpType {
		t.Errorf("jump type is not correct")
	}
	if event.JumpNumber != expectedEvent.(*JumpEvent).JumpNumber {
		t.Errorf("jump number is not correct")
	}
	if event.JumpTime != expectedEvent.(*JumpEvent).JumpTime {
		t.Errorf("jump time is not correct")
	}
}
