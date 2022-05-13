package controller

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/keiko233/V2Board-Bot/common"
	"github.com/keiko233/V2Board-Bot/dao"
	"github.com/keiko233/V2Board-Bot/lib/tgbot"
	"github.com/keiko233/V2Board-Bot/service"
)

func Bind(ctx *tgbot.Context) error {
	m := ctx.Message
	if m.Chat.ID < 0 {
		return common.ErrMustPrivateChat
	}
	user, notfound, err := dao.GetUserByTelegramID(nil, m.Sender.ID)
	if err != nil {
		log.Printf("QueryUser tgid = %d error, %s\n", m.Sender.ID, err)
		return err
	}
	if !notfound {
		return ctx.Send(m.Chat, fmt.Sprintf("✅ 当前绑定账户: %s\n若需要修改绑定,需要解绑当前账户。", user.Email))
	}

	format := strings.Index(m.Text, "token=")
	if format <= 0 {
		return common.ErrBindFormatError
	}

	user, err = service.BindUser(m.Text[format:][6:38], m.Sender.ID)
	if errors.Is(err, common.ErrBindAlready) {
		return err
	}

	if err != nil {
		log.Printf("Bind User token=%s and tgid=%d err %s\n", m.Text[6:38], m.Sender.ID, err)
		return common.ErrBindTokenInvalid
	}
	log.Printf("【Info】user %d bind to %s success\n", m.Sender.ID, user.Email)
	return ctx.Send(m.Chat, fmt.Sprintf("✅ 账户绑定成功: %s", user.Email))
}

func Unbind(ctx *tgbot.Context) error {
	m := ctx.Message
	if m.Chat.ID < 0 {
		return common.ErrMustPrivateChat
	}

	// 强制解绑token
	if m.Text != "" {
		format := strings.Index(m.Text, "token=")
		if format <= 0 {
			return common.ErrBindFormatError
		}

		notfound ,err := service.UnbindToken(m.Text[format:][6:38])
		if err != nil {
			return err
		}

		if notfound {
			return common.ErrNotBindUser
		}

		return ctx.Send(m.Chat, "✅ 账户解绑成功")
	}

	notfound, err := service.UnbindUser(m.Sender.ID)
	if err != nil {
		log.Printf("unbind user by tgid=%d error %s\n", m.Sender.ID, err)
		return err
	}

	if notfound {
		return common.ErrNotBindUser
	}

	return ctx.Send(m.Chat, "✅ 账户解绑成功")
}
