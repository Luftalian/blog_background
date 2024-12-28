package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var slackWebhookForLogger string

// InitSlack は環境変数から Webhook URL を取得します
func InitSlackForLogger() {
	log.Println("InitSlackForLogger")
	slackWebhookForLogger = os.Getenv("SLACK_WEBHOOK_URL_LOGGER")
}

// postToSlack は Webhook URL が設定されていればメッセージを POST 送信します
func postToSlack(msg string) {
	if slackWebhookForLogger == "" {
		log.Fatal("slackWebhookForLogger is not set")
		return
	}
	payload := map[string]string{"text": msg}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", slackWebhookForLogger, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func Print(v ...interface{}) {
	postToSlack(fmt.Sprint(v...))
	log.Print(v...)
}

func Printf(format string, v ...interface{}) {
	postToSlack(fmt.Sprintf(format, v...))
	log.Printf(format, v...)
}

func Println(v ...interface{}) {
	postToSlack(fmt.Sprintln(v...))
	log.Println(v...)
}

func Fatal(v ...interface{}) {
	postToSlack(fmt.Sprintln(v...))
	log.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	postToSlack(fmt.Sprintf(format, v...))
	log.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	postToSlack(fmt.Sprintln(v...))
	log.Fatalln(v...)
}

func Panic(v ...interface{}) {
	postToSlack(fmt.Sprintln(v...))
	log.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	postToSlack(fmt.Sprintf(format, v...))
	log.Panicf(format, v...)
}

func Panicln(v ...interface{}) {
	postToSlack(fmt.Sprintln(v...))
	log.Panicln(v...)
}
