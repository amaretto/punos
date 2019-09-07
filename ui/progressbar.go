package ui

// GetProgressbar is
func GetProgressbar(width, position, length int) string {
	current := float64(width) * float64(position) / float64(length)
	str := "["
	for i := 0; i < width; i++ {
		if i > int(current) {
			str = str + "-"
		} else if i == int(current) {
			str = str + ">"
		} else {
			str = str + "="
		}
	}
	str = str + "]"
	return str
}
