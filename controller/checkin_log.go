package controller

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"strconv"
	"strings"

	"github.com/keiko233/V2Board-Bot/dao"
	"github.com/keiko233/V2Board-Bot/lib/image"
	"github.com/keiko233/V2Board-Bot/model"
	"github.com/keiko233/V2Board-Bot/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

func HistoryQuery(q *tb.Callback) {
	list := strings.Split(q.Data, ":")
	n, _ := strconv.Atoi(list[0])
	m, _ := strconv.Atoi(list[1])
	id, _ := strconv.Atoi(list[2])

	count, out, err := dao.GetCheckLogsByTelegramID(int64(id), n, 5)
	if err != nil {
		log.Println("test2 err", err)
		model.Bot.Reply(q.Message, "获取失败")
		return
	}
	s := fmt.Sprintf("当前位于第%d页, 总条数%d, 总页数%d", n, count, m)
	ss := make([][]string, 0)
	s1 := make([]string, 0)
	s2 := make([]string, 0)
	s1 = append(s1, "签到时间")
	s2 = append(s2, "获得流量")
	for _, i := range out {
		s1 = append(s1, i.CreatedAt.Format("2006-01-02 15:04:05"))
		s2 = append(s2, utils.ByteSize(i.CheckinTraffic))
	}
	ss = append(ss, s1, s2)
	img, err := image.NewDefaultTable(ss, "/usr/UUBot/微软雅黑.ttf")
	if err != nil {
		log.Println("test2 err", err)
		model.Bot.Reply(q.Message, "生成图片失败")
		return
	}
	var b []byte
	bf := bytes.NewBuffer(b)
	err = png.Encode(bf, img.GetImage())
	if err != nil {
		log.Println("test3 err", err)
		_, err = model.Bot.Reply(q.Message, "生成图片失败")
		return
	}
	model.Bot.Edit(q.Message, &tb.Photo{
		File:    tb.FromReader(bf),
		Caption: s,
	}, page(n-1, n+1, m, id))
}

func page(perv, next, max, id int) *tb.ReplyMarkup {
	r := make([][]tb.InlineButton, 0)
	r1 := make([]tb.InlineButton, 0)
	r2 := tb.InlineButton{
		Unique: "history_page",
		Data:   strconv.Itoa(perv) + ":" + strconv.Itoa(max) + ":" + strconv.Itoa(id),
		Text:   "上一页",
	}
	if perv > 0 {
		r1 = append(r1, r2)
	}
	r2.Data = strconv.Itoa(next) + ":" + strconv.Itoa(max) + ":" + strconv.Itoa(id)
	r2.Text = "下一页"
	if max != 0 && next <= max {
		r1 = append(r1, r2)
	}
	r = append(r, r1)
	return &tb.ReplyMarkup{
		InlineKeyboard: r,
	}
}

func GetCheckinHistory(m *tb.Message) {

	count, out, err := dao.GetCheckLogsByTelegramID(m.Sender.ID, 1, 5)
	if err != nil {
		_, err = model.Bot.Reply(m, "获取失败")
		if err != nil {
			log.Println("test err", err)
		}
		return
	}

	max := (count / 5) + 1
	if count == 5 {
		max = 1
	}
	s := fmt.Sprintf("当前位于第1页, 总条数%d, 总页数%d", count, max)
	ss := make([][]string, 0)
	s1 := make([]string, 0)
	s2 := make([]string, 0)
	s1 = append(s1, "签到时间")
	s2 = append(s2, "获得流量")
	for _, i := range out {
		s1 = append(s1, i.CreatedAt.Format("2006-01-02 15:04:05"))
		s2 = append(s2, utils.ByteSize(i.CheckinTraffic))
	}
	ss = append(ss, s1, s2)
	img, err := image.NewDefaultTable(ss, "/usr/UUBot/微软雅黑.ttf")
	if err != nil {
		log.Println("test2 err", err)
		model.Bot.Reply(m, "生成图片失败")
		return
	}

	var b []byte
	bf := bytes.NewBuffer(b)
	err = png.Encode(bf, img.GetImage())
	if err != nil {
		_, err = model.Bot.Reply(m, "生成图片失败")
		if err != nil {
			log.Println("test3 err", err)
		}
		return
	}

	_, err = model.Bot.Reply(m, &tb.Photo{
		File:    tb.FromReader(bf),
		Caption: s,
	}, page(0, 2, int(max), int(m.Sender.ID)))
	if err != nil {
		log.Println("test err", err)
	}
}
