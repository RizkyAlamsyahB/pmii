package handlers

// getRoleName convert level code ke nama role
func getRoleName(level string) string {
	switch level {
	case "1":
		return "admin"
	case "2":
		return "user"
	default:
		return "user"
	}
}

// getStatusName convert status code ke nama status
func getStatusName(status string) string {
	switch status {
	case "1":
		return "active"
	case "0":
		return "inactive"
	default:
		return "unknown"
	}
}
