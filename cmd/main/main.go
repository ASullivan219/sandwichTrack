package main

import (
	"log/slog"
	"os"

	"github.com/ASullivan219/freeSandwich/internal/angels"
	"github.com/ASullivan219/freeSandwich/internal/cronmanager"
	"github.com/ASullivan219/freeSandwich/internal/lafc"
	"github.com/ASullivan219/freeSandwich/internal/notifier"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cronMan := cronmanager.New()
	emailNotifier := notifier.NewEmailNotifier(
		os.Getenv("FROM_EMAIL"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_HOST"),
		os.Getenv("EMAIL_PORT"),
	)

	emailNotifier.NotifyAll("Starting the Sandwich tracker lookout")

	for _, job := range lafc.BuildCronnables(&emailNotifier) {
		cronMan.AddJob(job.Sf.ExecutionTime, &job, job.Sf.Description)
	}
	for _, job := range angels.BuildCronnables(&emailNotifier, *cronMan) {
		cronMan.AddJob(job.ExecutionTime, job, job.SimpleGame.Description)
	}
	cronMan.Start()
}
