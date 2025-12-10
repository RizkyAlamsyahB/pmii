package handlers

import "strconv"

// getRoleName convert role code ke nama role
func getRoleName(role int) string {
	switch role {
	case 1:
		return "admin"
	case 2:
		return "author"
	default:
		return "author"
	}
}

// getStatusName convert isActive bool ke nama status
func getStatusName(isActive bool) string {
	if isActive {
		return "active"
	}
	return "inactive"
}

// getRoleString convert int role to string for JWT (keeps numeric values)
func getRoleString(role int) string {
	return strconv.Itoa(role)
}
