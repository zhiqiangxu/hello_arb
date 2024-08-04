package bot

import (
	"github.com/zhiqiangxu/arbbot/pkg/arb/bot/eth"
	"github.com/zhiqiangxu/litenode"
)

type Bot struct {
	Eth *eth.Bot
}

func New(config Config, lite *litenode.Node) *Bot {
	bot := &Bot{}
	if config.Eth != nil {
		bot.Eth = eth.NewBot(config.Eth, lite)
	}
	return bot
}

func (b *Bot) Start() (err error) {
	if b.Eth != nil {
		err = b.Eth.Start()
	}
	return
}

func (b *Bot) Stop() {
	if b.Eth != nil {
		b.Eth.Stop()
	}
}
