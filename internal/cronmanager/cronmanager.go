package cronmanager

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

// Interface to allow facilitate making a cronnable even from an object
type ICronnable interface {
	setCronString()
	getCronString() string
	setFunction() func()
	getFunction()
}

// Cron manager to mak eusing the cron library easier
type CronManager struct {
	Cron *cron.Cron
}

// Create a new cron manager
func New() *CronManager {
	ce := CronManager{
		Cron: cron.New(),
	}
	return &ce
}

func (cm *CronManager) Start() {
	slog.Info("Starting cron Manager")
	cm.Cron.Start()
	select {}
}

func (cm *CronManager) NextFive() []cron.Entry {
	if len(cm.Cron.Entries()) < 5 {
		return cm.Cron.Entries()[:len(cm.Cron.Entries())]
	}
	return cm.Cron.Entries()[:5]
}

func (cm *CronManager) AddJob(time time.Time, job cron.Job, description string) {
	_, err := cm.Cron.AddJob(timeToCronStr(time), job)
	if err != nil {
		slog.Error("Error creating Cronnable for: " + description)
	}
	slog.Info("Created Cronnable for " + description)
}

func timeToCronStr(time time.Time) string {
	return fmt.Sprintf("%d %d %d %d *", time.Minute(), time.Hour(), time.Day(), time.Month())
}

func (ce *CronManager) addCronnable(cronnable ICronnable) cron.EntryID {
	entry, err := ce.Cron.AddFunc(cronnable.getCronString(), cronnable.getFunction)
	if err != nil {
		slog.Error("ERROR: " + err.Error())
	}
	return entry
}
