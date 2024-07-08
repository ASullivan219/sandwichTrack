package main

import (
	"github.com/ASullivan219/freeSandwich/internal/cronmanager"
	"github.com/ASullivan219/freeSandwich/internal/lafc"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cronMan := cronmanager.New()
	for _, job := range lafc.BuildCronnables() {
		cronMan.AddJob(job.Sf.ExecutionTime, &job, job.Sf.Description)
	}
	cronMan.Start()
}
