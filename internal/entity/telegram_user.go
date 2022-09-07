package entity

type TelegramUser struct {
	ID     int    `bun:"id,pk,autoincrement"`
	ChatID int64  `bun:"chat_id"`
	Phone  string `bun:"phone"`
}
