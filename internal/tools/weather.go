package tools

import (
	"log"

	"github.com/suyashmohan/go-aiagent/internal"
)

func WeatherTool() (string, internal.AgentTool) {
	return "get_weater", internal.AgentTool{
		Description: "Get weather at a given location",
		Parameters: map[string]string{
			"location": "string",
		},
		Required: []string{"location"},
		Fn: func(m map[string]interface{}) string {
			log.Println("Inside get_weather")
			// Return a dummy weather info
			return "Sunny, 25C"
		},
	}
}
