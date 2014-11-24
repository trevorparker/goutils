package main

import (
	"testing"
	"time"
)

func TestSleepDurations(t *testing.T) {
	durations := [...]string{"0.5s", "0.0083m", "0.000138h"}
	dm, _ := time.ParseDuration("10ms")
	for _, duration := range durations {
		t0 := time.Now()
		sleep(duration)
		t1 := time.Now()
		de, _ := time.ParseDuration(duration)
		if d := t1.Sub(t0); d >= de+dm {
			t.Errorf("Sleep(%v) took %v, longer than %v", duration, d, de+dm)
		}
		if d := t1.Sub(t0); d <= de-dm {
			t.Errorf("Sleep(%v) took %v, shorter than %v", duration, d, de-dm)
		}
	}
}
