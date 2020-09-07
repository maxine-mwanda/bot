package resources

import "encoding/json"

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



func PlayerCountKeyboard() string {

	keyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]string{
			{
				{
					"text":          "Two",
					"callback_data": "players2",
				},
			},

			{
				{
					"text":          "Three",
					"callback_data": "players3",
				},
			},
			{
				{
					"text":          "Four",
					"callback_data": "players3",
				},
			},

			{
				{
					"text":          "Five",
					"callback_data": "players5",
				},
			},
		},
	}
	jsonkeyboard, _ := json.Marshal(keyboard)
	return string(jsonkeyboard)
}
