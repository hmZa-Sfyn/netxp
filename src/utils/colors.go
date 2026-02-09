package colors

// Color code constants
const (
	CReset   = "\033[0m"
	CBold    = "\033[1m"
	CDim     = "\033[2m"
	CRed     = "\033[31m"
	CGreen   = "\033[32m"
	CYellow  = "\033[33m"
	CBlue    = "\033[34m"
	CMagenta = "\033[35m"
	CCyan    = "\033[36m"
	CWhite   = "\033[37m"

	CBgRed   = "\033[41m"
	CBgGreen = "\033[42m"
)

// Colorize wraps text with ANSI color codes
func Colorize(text string, color string) string {
	return color + text + CReset
}

// ColorizeError wraps error message in red
func ColorizeError(text string) string {
	return CBold + CRed + text + CReset
}

// ColorizeSuccess wraps success message in green
func ColorizeSuccess(text string) string {
	return CBold + CGreen + text + CReset
}

// ColorizeInfo wraps info message in cyan
func ColorizeInfo(text string) string {
	return CCyan + text + CReset
}

// ColorizeWarn wraps warning message in yellow
func ColorizeWarn(text string) string {
	return CBold + CYellow + text + CReset
}
