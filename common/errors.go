package common

type ErrorMsg string

func (e ErrorMsg) Error() string {
	return string(e)
}

const (
	ErrHandle ErrorMsg = "❌ 操作失败，请稍后重试或联系管理员"

	ErrNotFoundPlan         ErrorMsg = "👀 当前暂无订阅计划,该功能需要订阅后使用～"
	ErrNotBindUser          ErrorMsg = "👀 当前未绑定账户\n请私聊发送 /bind <订阅地址> 绑定账户"
	ErrAlreadyCheckin       ErrorMsg = "✅ 今天已经签到过啦！明天再来哦～"
	ErrMakeImageError       ErrorMsg = "👀 生成图片失败"
	ErrNotFoundCheckinUsers ErrorMsg = "👀 今天还没有人签到哦~~"
	ErrMustPrivateChat      ErrorMsg = "👀 请私聊我命令哦~~"
	ErrBindFormatError      ErrorMsg = "👀 ️账户绑定格式: /bind <订阅地址>"
	ErrBindTokenInvalid     ErrorMsg = "❌ 订阅无效,请前往官网复制最新订阅地址!"
	ErrInvalidPlan          ErrorMsg = "👀 订阅套餐不存在，请稍后重试或联系管理员"
)
