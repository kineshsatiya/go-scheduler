package scheduler

import (
	"time"
)

// SimpleScheduler to run tasks based on time units.
type SimpleScheduler struct {
	interval       time.Duration
	scheduled      bool
	backgroundTask *BackgroundTask
	taskQuit       chan int
}

// BackgroundTask is created upon starting a schedule
type BackgroundTask struct {
	task            func()
	interval        time.Duration
	delay           bool
	stopTaskChannel chan int
	stopCompleted   chan int
}

// Schedule methods schedules the task to be run in background based on scheduler created.
func (scheduler *SimpleScheduler) Schedule(task func()) Scheduler {
	scheduler.schedule(task, false)
	return scheduler
}

// ScheduleWithDelay schedules the task to execute after the duration provided.
func (scheduler *SimpleScheduler) ScheduleWithDelay(task func()) Scheduler {
	scheduler.schedule(task, true)
	return scheduler
}

func (scheduler *SimpleScheduler) schedule(task func(), delay bool) {
	scheduler.taskQuit = make(chan int)
	scheduler.backgroundTask = createBackgroundTask(task, scheduler.interval, true, scheduler.taskQuit)
	go scheduler.backgroundTask.start()
	scheduler.scheduled = true
}

func createBackgroundTask(task func(), interval time.Duration, delay bool, stopCompleted chan int) *BackgroundTask {
	backgroundTask := &BackgroundTask{}
	backgroundTask.task = task
	backgroundTask.interval = interval
	backgroundTask.delay = delay
	backgroundTask.stopTaskChannel = make(chan int)
	backgroundTask.stopCompleted = stopCompleted
	return backgroundTask
}

// Stop the backgrond task. This call is blocking if the task is currently executing, and only exits once the current task completes.
func (scheduler *SimpleScheduler) Stop() {
	if !scheduler.scheduled {
		panic("Cannot pause without starting the scheduling")
	}
	scheduler.backgroundTask.stopTaskChannel <- 1
	<-scheduler.taskQuit
}

func (bg *BackgroundTask) start() {
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
