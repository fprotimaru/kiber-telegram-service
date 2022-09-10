package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"telegram/internal/entity"
	"telegram/internal/repository"
	"telegram/pb"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ pb.TelegramServer = (*Service)(nil)

type Service struct {
	bot              *tgbotapi.BotAPI
	telegramUserRepo *repository.TelegramUserRepository
}

func New(token string, telegramUserRepo *repository.TelegramUserRepository) *Service {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	return &Service{
		bot:              bot,
		telegramUserRepo: telegramUserRepo,
	}
}

func (s *Service) SendDocument(ctx context.Context, req *pb.SendDocumentRequest) (*pb.SendDocumentReply, error) {
	var (
		phone    = req.GetPhone()
		caption  = req.GetCaption()
		file     = req.GetFile()
		fileName = req.GetFileName()
	)

	telegramUser, err := s.telegramUserRepo.GetByPhone(ctx, phone)
	if err != nil {
		log.Printf("telegramUser.GetByPhone error: %v\n", err)
		if errors.Is(err, sql.ErrNoRows) {
			return &pb.SendDocumentReply{}, nil
		}
		return &pb.SendDocumentReply{}, err
	}

	doc := tgbotapi.NewDocument(telegramUser.ChatID, tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: file,
	})
	doc.Caption = caption
	doc.ParseMode = "MarkdownV2"

	_, err = s.bot.Send(doc)
	if err != nil {
		log.Printf("botSend error: %v\n", err)
		return &pb.SendDocumentReply{}, err
	}
	return &pb.SendDocumentReply{}, nil
}

func (s *Service) Listen(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.Contact != nil {
				telegramUser := entity.TelegramUser{
					ChatID: update.Message.Chat.ID,
					Phone:  update.Message.Contact.PhoneNumber,
				}

				if !strings.HasPrefix(telegramUser.Phone, "+") {
					telegramUser.Phone = fmt.Sprintf("+%s", telegramUser.Phone)
				}

				err := s.telegramUserRepo.Create(ctx, &telegramUser)
				if err != nil {
					log.Printf("s.telegramUserRepo.Create error: %v\n", err)
				}

				msg := tgbotapi.NewMessage(telegramUser.ChatID, "Siz ro'yhatdan o'tdingiz")
				msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
					RemoveKeyboard: true,
				}
				_, err = s.bot.Send(msg)
				if err != nil {
					log.Printf("s.bot.Send error: %v\n", err)
				}
			}
			_, err := s.telegramUserRepo.GetByChatID(ctx, update.Message.Chat.ID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					// btn := tgbotapi.NewKeyboardButtonContact("Ro'yhatdan o'tish")
					// btn.RequestContact = true
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ro'yhatdan o'ting")
					msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
						Keyboard: [][]tgbotapi.KeyboardButton{
							{{RequestContact: true, Text: "Ro'yhatdan o'tish"}},
						},
						ResizeKeyboard: true,
					}
					_, err = s.bot.Send(msg)
					if err != nil {
						log.Printf("s.bot.Send error: %v\n", err)
					}
					continue
				}
				log.Printf("s.telegramUserRepo.GetByChatID error: %v\n", err)
			}
		}
	}
}
