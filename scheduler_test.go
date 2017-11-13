package scheduler_test

import (
	scheduler "go-scheduler"
	"testing"
	"time"
)

func Test_Schedule(t *testing.T) {
	valueToTest := 1
	scheduler := scheduler.Every("5ms").Schedule(func() {
		valueToTest++
	})

	duration, _ := time.ParseDuration("1s")
	time.Sleep(duration)
	scheduler.Stop()

	if valueToTest > 1000 || valueToTest == 1 {
		t.Errorf("Simple scheduler is not executing the scheduled task correctly, test value is %v", valueToTest)
	}
}

func Test_InvalidInterval(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.FailNow()
		}
	}()
	scheduler.Every("55gg")
}
