package tgbot

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Engine struct {
	bot *tb.Bot
	ms  []Middleware
}

type Middleware func(HandleFunc) HandleFunc

type HandleFunc func(*Context) error

func NewDefaultEngine(token string) (*Engine, error) {
	bot, err := tb.NewBot(tb.Settings{
		URL:    "https://api.telegram.org",
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}
	return &Engine{
		bot: bot,
		ms:  make([]Middleware, 0),
	}, nil
}

func (e *Engine) Run() {
	e.bot.Start()
}

// handle *tb.Message
func (e *Engine) Handle(h HandleFunc, rs ...interface{}) {
	for _, r := range rs {
		for _, i := range e.ms {
			h = i(h)
		}
		e.bot.Handle(r, e.handle(h))
	}
}

// handle *tb.Callback
func (e *Engine) HandleCallback(h HandleFunc, rs ...string) {
	for _, r := range rs {
		for _, i := range e.ms {
			h = i(h)
		}
		e.bot.Handle(fmt.Sprintf("\f%s", r), e.handleCallback(h))
	}
}

// use middleware
func (e *Engine) Use(middleware ...Middleware) *Engine {
	e.ms = append(e.ms, middleware...)
	return e
}

func (e *Engine) handle(h HandleFunc) func(*tb.Message) {
	return func(m *tb.Message) {
		h(&Context{
			Message: m,
			bot:     e.bot,
		})
	}
}

func (e *Engine) handleCallback(h HandleFunc) func(*tb.Callback) {
	return func(c *tb.Callback) {
		h(&Context{
			Message:  c.Message,
			Callback: c,
			bot:      e.bot,
		})
	}
}
