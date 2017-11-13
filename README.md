# go-scheduler
To schedule a background task using a fixed interval or a cron expression

### Example
```golang
valueToTest := 1
// schedule a background task to increment a test value every 5 seconds
scheduler := scheduler.Every("5s").Schedule(func() {
  valueToTest++
})

// After sometime, if the task needs to be stopped
scheduler.Stop()
```

