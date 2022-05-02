package tgbot

import tb "gopkg.in/tucnak/telebot.v2"

type Context struct {
	Message  *tb.Message
	Callback *tb.Callback
	bot      *tb.Bot
}

func (c *Context) Reply(msg interface{}, options ...interface{}) error {
	_, err := c.bot.Reply(c.Message, msg, options...)
	return err
}

func (c *Context) Send(to tb.Recipient, msg interface{}, options ...interface{}) error {
	_, err := c.bot.Send(to, msg, options...)
	return err
}

func (c *Context) Edit(msg interface{}, options ...interface{}) error {
	_, err := c.bot.Edit(c.Message, msg, options...)
	return err
}

func (c *Context) IsCallback() bool {
	return c.Callback != nil
}

func (c *Context) ChatMemberOf(userID int64) (*tb.ChatMember, error) {
	return c.bot.ChatMemberOf(c.Message.Chat, &tb.User{ID: userID})
}
