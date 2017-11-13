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
	return &simpleScheduler{interval: fixedInterval}
}

// Cron returns a Scheduler that will execute the task based on the cron expression
func Cron(cron string) Scheduler {
	return &CronScheduler{cronExpression: cron}
}

// simpleScheduler to run tasks based on time units.
type simpleScheduler struct {
	interval       time.Duration
	scheduled      bool
	backgroundTask *backgroundTask
	taskQuit       chan int
}

// backgroundTask is created upon starting a schedule
type backgroundTask struct {
	task            func()
	interval        time.Duration
	delay           bool
	stopTaskChannel chan int
	stopCompleted   chan int
}

// Schedule methods schedules the task to be run in background based on scheduler created.
func (scheduler *simpleScheduler) Schedule(task func()) Scheduler {
	scheduler.schedule(task, false)
	return scheduler
}

// ScheduleWithDelay schedules the task to execute after the duration provided.
func (scheduler *simpleScheduler) ScheduleWithDelay(task func()) Scheduler {
	scheduler.schedule(task, true)
	return scheduler
}

func (scheduler *simpleScheduler) schedule(task func(), delay bool) {
	scheduler.taskQuit = make(chan int)
	scheduler.backgroundTask = createbackgroundTask(task, scheduler.interval, true, scheduler.taskQuit)
	go scheduler.backgroundTask.start()
	scheduler.scheduled = true
}

func createbackgroundTask(task func(), interval time.Duration, delay bool, stopCompleted chan int) *backgroundTask {
	backgroundTask := &backgroundTask{}
	backgroundTask.task = task
	backgroundTask.interval = interval
	backgroundTask.delay = delay
	backgroundTask.stopTaskChannel = make(chan int)
	backgroundTask.stopCompleted = stopCompleted
	return backgroundTask
}

// Stop the backgrond task. This call is blocking if the task is currently executing, and only exits once the current task completes.
func (scheduler *simpleScheduler) Stop() {
	if !scheduler.scheduled {
		panic("Cannot pause without starting the scheduling")
	}
	scheduler.backgroundTask.stopTaskChannel <- 1
	<-scheduler.taskQuit
}

func (bg *backgroundTask) start() {
	if !bg.delay {
		bg.task()
	}
	for {
		select {
		case <-bg.stopTaskChannel:
			bg.stopCompleted <- 1
			break
		default:
			start := time.Now()
			next := start.Add(bg.interval)
			sleepTime := next.Sub(start)
			time.Sleep(sleepTime)
			bg.task()
		}
	}
}

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