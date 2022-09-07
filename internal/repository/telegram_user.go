package repository

import (
	"context"
	"database/sql"

	"telegram/internal/entity"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type TelegramUserRepository struct {
	db *bun.DB
}

func NewTelegramUserRepository(url string) (*TelegramUserRepository, error) {
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(url)))
	db := bun.NewDB(sqlDB, pgdialect.New())
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &TelegramUserRepository{db: db}, nil
}

func (r *TelegramUserRepository) GetByPhone(ctx context.Context, phone string) (*entity.TelegramUser, error) {
	var telegramUser entity.TelegramUser
	err := r.db.NewSelect().Model(&telegramUser).Where("phone = ?", phone).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &telegramUser, nil
}

func (r *TelegramUserRepository) GetByChatID(ctx context.Context, chatID int64) (*entity.TelegramUser, error) {
	var telegramUser entity.TelegramUser
	err := r.db.NewSelect().Model(&telegramUser).Where("chat_id = ?", chatID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &telegramUser, nil
}

func (r *TelegramUserRepository) Create(ctx context.Context, telegramUser *entity.TelegramUser) error {
	_, err := r.db.NewInsert().Model(telegramUser).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
