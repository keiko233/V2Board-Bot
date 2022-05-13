package model

type FortuneType string

const (
	FortuneLuck             FortuneType = "吉"
	FortuneVeryLuck         FortuneType = "大吉"
	FortuneUnfavourable     FortuneType = "寄"
	FortuneVeryUnfavourable FortuneType = "大寄"
)

type Fortune struct {
	TgID         int64
	UserID       int64
	TodayFortune FortuneType
}
