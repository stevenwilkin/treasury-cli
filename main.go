package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gorilla/websocket"
	_ "github.com/joho/godotenv/autoload"
)

type authMessage struct {
	Auth string `json:"auth"`
}

type stateMessage struct {
	Error           string  `json:"error"`
	Exposure        float64 `json:"exposure"`
	LeverageDeribit float64 `json:"leverage_deribit"`
	LeverageBybit   float64 `json:"leverage_bybit"`
}

func display(sm stateMessage) {
	fmt.Println("\033[2J\033[H\033[?25l") // clear screen, move cursor to top of screen, hide cursor
	fmt.Printf("  Exposure: %.8f\n", sm.Exposure)
	fmt.Printf("  Leverage: %.2f %.2f\n", sm.LeverageDeribit, sm.LeverageBybit)
}

func main() {
	url := os.Getenv("TREASURY_WS_URL")
	auth_token := os.Getenv("TREASURY_AUTH_TOKEN")

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = c.WriteJSON(authMessage{Auth: auth_token}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		var sm stateMessage
		json.Unmarshal(message, &sm)

		if sm.Error != "" {
			fmt.Println(sm.Error)
			os.Exit(1)
		}

		display(sm)
	}
}
