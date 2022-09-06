package service

import (
	"context"
	"log"

	"telegram/pb"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ pb.TelegramServer = (*Service)(nil)

type Service struct {
	bot *tgbotapi.BotAPI
}

func New(token string) *Service {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	return &Service{bot: bot}
}

func (s *Service) SendDocument(ctx context.Context, req *pb.SendDocumentRequest) (*pb.SendDocumentReply, error) {
	var (
		chatID   = req.GetChatId()
		caption  = req.GetCaption()
		file     = req.GetFile()
		fileName = req.GetFileName()
	)

	doc := tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: file,
	})
	doc.Caption = caption

	_, err := s.bot.Send(doc)
	if err != nil {
		return &pb.SendDocumentReply{}, err
	}
	return &pb.SendDocumentReply{}, nil
}
