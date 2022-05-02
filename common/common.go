package common

import (
	"log"

	"github.com/keiko233/V2Board-Bot/lib/tgbot"
)

func ErrorResult(c *tgbot.Context, err error) {
	e, ok := err.(ErrorMsg)
	if ok {
		c.Reply(e.Error())
	} else {
		log.Printf("【Error】%s\n", err)
		c.Reply(ErrHandle)
	}
}
