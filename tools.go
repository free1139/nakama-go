package nakama

func checkStr(str *string) bool {
	return str != nil && *str != ""
}

func checkBool(b *bool) bool {
	return b != nil && *b
}
