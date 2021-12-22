package utils

func GetUserID(token string) string {
	user, _, _, _ := TokenProcess(token)
	return user.ID.Hex()
}
