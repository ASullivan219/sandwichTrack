package notifier

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

const (
	MAIN_PATH = "../../cmd/main/.env"
)

func TestSendEmail(t *testing.T) {
	godotenv.Load(MAIN_PATH)
	from := os.Getenv("FROM_EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	notifier := NewEmailNotifier(from, password, host, port)
	err := notifier.NotifyOne("This is a test Message", "alexander.sullivan219@gmail.com")
	if err != nil {
		t.Fatal("Error sending email")
	}
}
