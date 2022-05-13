package service

import (
	"fmt"
	"math"
	"time"

	"github.com/keiko233/V2Board-Bot/lib/rand"
	"github.com/keiko233/V2Board-Bot/model"
	"github.com/keiko233/V2Board-Bot/utils"
)

var fortuneTraffer map[model.FortuneType][]int64
var fortuneList []model.FortuneType
var pool []float64

func GetFortune(tgid int64) (model.FortuneType, error) {
	m, err := setFortune(tgid)
	if err != nil {
		return "", err
	}

	return m.TodayFortune, nil
}

// 设置运势
func setFortune(tgid int64) (*model.Fortune, error) {
	today := utils.TodayStart()
	key := fmt.Sprintf("%d:%d", tgid, today.Unix())
	var f model.Fortune
	found, err := model.Cache.Exists(key)
	if err != nil {
		return nil, err
	}
	if found {
		// var m map[string]interface{}
		if err := model.Cache.GetStruct(key, &f); err != nil {
			return nil, err
		}
		// f.TgID = int64(m["TgID"].(float64))
		// f.TodayFortune = model.FortuneType(m["TodayFortune"].(string))
		return &f, nil
	}
	f.TgID = tgid
	f.TodayFortune = getFortune(tgid)
	if err := model.Cache.Set(key, f, time.Hour*24); err != nil {
		return nil, err
	}

	return &f, nil
}

func getFortune(tgid int64) model.FortuneType {
	switch rand.RandIntWithSeed(tgid, 5, 0) {
	case 1:
		return model.FortuneVeryLuck
	case 2:
		return model.FortuneLuck
	case 3:
		return model.FortuneUnfavourable
	default:
		return model.FortuneVeryUnfavourable
	}
}

func t(f model.FortuneType, tgid int64) model.FortuneType {
	if fortuneList == nil {
		fortuneList = make([]model.FortuneType, 0)
		fortuneList = append(fortuneList, model.FortuneVeryLuck)
		fortuneList = append(fortuneList, model.FortuneLuck)
		fortuneList = append(fortuneList, model.FortuneUnfavourable)
		fortuneList = append(fortuneList, model.FortuneVeryUnfavourable)
	}

	nl := make([]model.FortuneType, 0)
	for _, i := range fortuneList {
		if i != f {
			nl = append(nl, i)
		}
	}

	switch rand.RandIntWithSeed(tgid, 10, 1) {
	case 1, 2, 3, 4, 5, 6:
		return f
	case 7:
		return nl[0]
	case 8:
		return nl[1]
	default:
		return nl[2]
	}

}

// 获取真实签到流量
func GetTraffer(f model.FortuneType, tgid int64) (int64, error) {
	if fortuneTraffer == nil {
		fortuneTraffer = make(map[model.FortuneType][]int64)
		all := model.Config.Bot.MaxByte - model.Config.Bot.MinByte
		vlmin := all/4*3 + model.Config.Bot.MinByte
		vlmax := all + model.Config.Bot.MinByte
		lmin := all/4*2 + model.Config.Bot.MinByte
		umin := all/4 + model.Config.Bot.MinByte
		vumin := model.Config.Bot.MinByte
		fortuneTraffer[model.FortuneVeryLuck] = []int64{vlmax, vlmin}
		fortuneTraffer[model.FortuneLuck] = []int64{vlmin, lmin}
		fortuneTraffer[model.FortuneUnfavourable] = []int64{lmin, umin}
		fortuneTraffer[model.FortuneVeryUnfavourable] = []int64{umin, vumin}
	}
	f = t(f, tgid)
	l, ok := fortuneTraffer[f]
	if !ok {
		return 0, fmt.Errorf("not found fortune")
	}
	return rand.RandIntWithSeed(tgid, l[0], l[1]), nil
}

func PassPool(n int64, f model.FortuneType, tgid int64) (int64, bool, error) {
	key := fmt.Sprintf("%s:%d", "pool", utils.TodayStart().Unix())
	if pool == nil {
		ok, err := model.Cache.Exists(key)
		if err != nil {
			return 0, false, err
		}
		if ok {
			// var p []interface{}
			if err := model.Cache.GetStruct(key, &pool); err != nil {
				return 0, false, err
			}
			// pool = make([]float64, 0)
			// pool = append(pool, p[0].(float64))
			// pool = append(pool, p[1].(float64))
			// pool = append(pool, p[2].(float64))
		} else {

			pool = []float64{0, float64(rand.RandIntWithSeed(tgid, 15, 10)), 0}
			if err := model.Cache.Set(key, pool, time.Hour*24); err != nil {
				return 0, false, err
			}
		}
	}
	if pool[2] == 1 {
		return 0, false, nil
	}

	if pool[1] <= 0 {
		n = int64(pool[0])
		pool[2] = 1
		if err := model.Cache.Set(key, pool, time.Hour*24); err != nil {
			return 0, false, err
		}
		return n, false, nil
	}
	yes := passTraffer(n, f)
	if yes {
		pool[0] = pool[0] + math.Abs(float64(n))
	} else {
		pool[0] = pool[0] - math.Abs(float64(n))
	}
	pool[1] = pool[1] - 1
	if err := model.Cache.Set(key, pool, time.Hour*24); err != nil {
		return 0, false, err
	}
	return 0, yes, nil
}

// 还愿
//
// 如果运势为凶，流量不在该区间，则还怨
//
// 如果运势为吉，流量在该区间，则还愿
//
func passTraffer(n int64, f model.FortuneType) bool {
	l := fortuneTraffer[f]
	if n <= l[0] && n > l[1] {
		return true
	}
	return false
}
