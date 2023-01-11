package notifications

import "fmt"

func StatusColor(status string) string {
	const SUCCESS = "34A853"
	const ERROR = "de5246"

	var statusColor string

	switch status {
	case "success":
		fmt.Println(SUCCESS, "success")
		statusColor = SUCCESS
	case "error":
		fmt.Println(ERROR, "error")
		statusColor = ERROR
	default:
		fmt.Println("default")
		statusColor = ""
	}
	return statusColor
}
