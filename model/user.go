package model

type User struct {
	Id                uint
	TelegramId        uint
	Email             string
	Token             string
	U                 int64
	D                 int64
	PlanId            int64
	Balance           int64
	TransferEnable    int64
	CommissionBalance int64
	ExpiredAt         int64
	CreatedAt         int64
}

type Plan struct {
	Id   uint
	Name string
}