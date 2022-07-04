package bot

import (
	"fmt"
	"testing"

	"github.com/kirontoo/rxkiro/config"
)

// var Bot := &RxKiro {

// }

var MockConfig = &config.Config{
	BotName:    "rxkiro",
	Streamer:   "kironto",
	AuthToken:  "",
	DbUrl:      "",
	DbToken:    "",
	DbHost:     "",
	DbUser:     "",
	DbPassword: "",
	DbPort:     "",
	DbName:     "",
}

func TestIsPositiveInt(t *testing.T) {
	t.Run("Should be a positive int", func(t *testing.T) {
		got := isPositiveInt(5)
		expected := true

		if got != expected {
			fmt.Errorf("Expected %v, but got %v", expected, got)
		}
	})

	t.Run("Should not be a positive int", func(t *testing.T) {
		got := isPositiveInt(-5)
		expected := false

		if got != expected {
			fmt.Errorf("Expected %v, but got %v", expected, got)
		}
	})
}
