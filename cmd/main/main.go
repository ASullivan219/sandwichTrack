package main

import (
	"os"

	"github.com/ASullivan219/freeSandwich/internal/cronmanager"
	"github.com/ASullivan219/freeSandwich/internal/lafc"
	"github.com/ASullivan219/freeSandwich/internal/notifier"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cronMan := cronmanager.New()
	emailNotifier := notifier.NewEmailNotifier(
		os.Getenv("FROM_EMAIL"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_HOST"),
		os.Getenv("EMAIL_PORT"),
	)

	for _, job := range lafc.BuildCronnables() {
		cronMan.AddJob(job.Sf.ExecutionTime, &job, job.Sf.Description)
	}
	cronMan.Start()
}
