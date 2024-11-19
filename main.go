package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

func init() {
	// โหลดไฟล์ .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	var err error
	// ใช้ข้อมูลจาก .env
	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelAccessToken := os.Getenv("CHANNEL_ACCESS_TOKEN")

	// สร้าง bot client
	bot, err = linebot.New(channelSecret, channelAccessToken)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", callbackHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// ลูปผ่าน events
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			// เรียกฟังก์ชันตอบกลับ QuickReply เมื่อได้รับข้อความ
			replyQuickReplyMessage(event)
		}
	}
}

func replyQuickReplyMessage(event *linebot.Event) {
	// สร้าง QuickReply Items
	quickReplyItems := linebot.NewQuickReplyItems(
		linebot.NewQuickReplyButton(
			"ตัวเลือก 1",
			linebot.NewMessageAction("Message 1", "Message 1"),
		),
		linebot.NewQuickReplyButton(
			"ตัวเลือก 2",
			linebot.NewMessageAction("Message 2", "Message 2"),
		),
	)

	// ส่งข้อความพร้อม QuickReply
	if _, err := bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTextMessage("เลือกคำตอบของคุณ").WithQuickReplies(quickReplyItems),
	).Do(); err != nil {
		log.Print(err)
	}
}
