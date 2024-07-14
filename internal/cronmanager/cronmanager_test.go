package cronmanager

import (
	"testing"
	"time"
)

type testResult struct {
	Expected string
	Actual   string
}

func TestCronEngine(t *testing.T) {
	cm := New()
	if cm == nil {
		t.Fatalf("Error Creating a new cron engine, cron engine is nil")
	}
}

func TestTimeToCronString(t *testing.T) {
	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	tm, _ := time.Parse(layout, "Jan 5, 2024 at 2:30pm (GMT)")
	tr := testResult{Expected: "30 14 5 1 *", Actual: timeToCronStr(tm)}
	if tr.Expected != tr.Actual {
		t.Fatalf("Failed Cron string test expected: %s got %s", tr.Expected, tr.Actual)
	}
}

func TestGetFive(t *testing.T) {
	tj := &testJob{TestBool: false}
	now := time.Now()
	oneMinDur, _ := time.ParseDuration("1m")
	cm := New()
	cm.AddJob(now.Add(oneMinDur), tj, "TEST JOB")
	cm.AddJob(now.Add(oneMinDur), tj, "TEST JOB")
	cm.AddJob(now.Add(oneMinDur), tj, "TEST JOB")

	five := cm.NextFive()

	if len(five) != 3 {
		t.Fatal("Expected 3 Entries")
	}

	cm.AddJob(now.Add(oneMinDur), tj, "TEST JOB")
	cm.AddJob(now.Add(oneMinDur), tj, "TEST JOB")
	cm.Cron.Start()
	five = cm.NextFive()

	if len(five) != 5 {
		t.Fatal("Expected 5 entries")
	}

	t.Logf("Now: %s", now)
	for _, entry := range five {
		t.Logf("Entry.next: %s", entry.Next)
	}
}

// Test data structure
type testJob struct {
	TestBool bool
}

func (tj *testJob) Run() {
	tj.TestBool = true
}

// This test can take up to a minute to execute, as the smallest granularity
// possible with the cron package is one minute.
func TestCronFunctionality(t *testing.T) {
	tj := &testJob{TestBool: false}
	cm := New()
	now := time.Now()
	oneMinDur, _ := time.ParseDuration("1m")
	twoMinDur, _ := time.ParseDuration("2m")
	cm.AddJob(now.Add(oneMinDur), tj, "TEST JOB")
	cm.Cron.Start()
	for time.Now().Before(now.Add(twoMinDur)) {
		if tj.TestBool {
			cm.Cron.Stop()
			return
		}
	}
	t.Fatal("ERROR: Test timed out without changing the internal value of testJob")
}
