package scheduler

import "time"

// Scheduler interface
type Scheduler interface {
	// Schedule the given task, that will start the execution immediately
	Schedule(task func()) Scheduler
	// ScheduleWithDelay schedules the execution of the task by calculating next execution time
	ScheduleWithDelay(task func()) Scheduler
	// Stop the scheduled task. This is a blocking call, if the task is currently executing.
	Stop()
}

// Every returns a Scheduler that will execute the task on fixed interval expressed as a string.
// A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func Every(duration string) Scheduler {
	fixedInterval, error := time.ParseDuration(duration)
	if error != nil {
		panic(error.Error())
	}
	return &SimpleScheduler{interval: fixedInterval}
}

// Cron returns a Scheduler that will execute the task based on the cron expression
func Cron(cron string) Scheduler {
	return &CronScheduler{cronExpression: cron}
}
