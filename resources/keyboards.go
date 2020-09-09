package resources

import (
	"encoding/json"
	"fmt"
)

func AcceptDeclineKeyboard() string {
	keyboard_first := map[string]interface{}{
		"inline_keyboard": [][]map[string]string{
			{
				{
					"text":          "Accept",
					"callback_data": "Please choose truth or dare",
				},
			},

			{
				{
					"text":          "Decline",
					"callback_data": "You have lost ten points",
				},
			},
		},
	}

	jsonkeyboard, _ := json.Marshal(keyboard_first)
	return string(jsonkeyboard)
}


func TruthOrDareKeyboard() string {
	keyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]string{
			{
				{
					"text":          "Truth",
					"callback_data": "truth",
				},
			},

			{
				{
					"text":          "Dare",
					"callback_data": "dare",
				},
			},
		},
	}

	jsonkeyboard, _ := json.Marshal(keyboard)
	return string(jsonkeyboard)
}



func PlayerCountKeyboard(gameId int64) string {

	keyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]string{
			{
				{
					"text":          "Two",
					"callback_data": fmt.Sprintf("players2-%d", gameId),
				},
			},

			{
				{
					"text":          "Three",
					"callback_data": fmt.Sprintf("players3-%d", gameId),
				},
			},
			{
				{
					"text":          "Four",
					"callback_data": fmt.Sprintf("players3-%d", gameId),
				},
			},

			{
				{
					"text":          "Five",
					"callback_data": fmt.Sprintf("players5-%d",gameId),
				},
			},
		},
	}
	jsonkeyboard, _ := json.Marshal(keyboard)
	return string(jsonkeyboard)
}
