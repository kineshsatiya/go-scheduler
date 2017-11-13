package scheduler

// CronScheduler to run tasks based on a cron expression
type CronScheduler struct {
	cronExpression string
}

// TODO
func (scheduler *CronScheduler) Schedule(task func()) Scheduler {
	return scheduler
}

// TODO
func (scheduler *CronScheduler) ScheduleWithDelay(task func()) Scheduler {
	return scheduler
}

// TODO
func (scheduler *CronScheduler) Stop() {
}
