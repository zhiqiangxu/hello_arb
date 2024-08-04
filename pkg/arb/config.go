package arb

import (
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot"
	"github.com/zhiqiangxu/litenode"
)

type Config struct {
	Lite litenode.Config
	Bot  bot.Config
}
