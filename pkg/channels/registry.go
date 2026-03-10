package channels

import (
	"fmt"
	"sync"

	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/config"
)

// ChannelFactory is a constructor function that creates a Channel from config and message bus.
// Each channel subpackage registers one or more factories via init().
type ChannelFactory func(cfg *config.Config, bus *bus.MessageBus) (Channel, error)

// TelegramBotFactory creates a Telegram channel from a specific TelegramConfig.
// Set by the telegram subpackage's init() function.
type TelegramBotFactory func(telegramCfg config.TelegramConfig, cfg *config.Config, bus *bus.MessageBus) (Channel, error)

var (
	factoriesMu        sync.RWMutex
	factories          = map[string]ChannelFactory{}
	telegramBotFactory TelegramBotFactory
)

// RegisterFactory registers a named channel factory. Called from subpackage init() functions.
func RegisterFactory(name string, f ChannelFactory) {
	factoriesMu.Lock()
	defer factoriesMu.Unlock()
	factories[name] = f
}

// getFactory looks up a channel factory by name.
func getFactory(name string) (ChannelFactory, bool) {
	factoriesMu.RLock()
	defer factoriesMu.RUnlock()
	f, ok := factories[name]
	return f, ok
}

// RegisterTelegramBotFactory registers the factory for creating per-config Telegram bots.
func RegisterTelegramBotFactory(f TelegramBotFactory) {
	factoriesMu.Lock()
	defer factoriesMu.Unlock()
	telegramBotFactory = f
}

// newTelegramFromConfig creates a Telegram channel from a specific TelegramConfig.
func newTelegramFromConfig(telegramCfg config.TelegramConfig, cfg *config.Config, msgBus *bus.MessageBus) (Channel, error) {
	factoriesMu.RLock()
	f := telegramBotFactory
	factoriesMu.RUnlock()
	if f == nil {
		return nil, fmt.Errorf("telegram bot factory not registered")
	}
	return f(telegramCfg, cfg, msgBus)
}
