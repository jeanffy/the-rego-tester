package colors

const DIM = "\x1b[2m"
const RED = "\033[31m"
const RESET = "\x1b[0m"

func Dimmed(msg string) string {
	return DIM + msg + RESET
}

func Error(msg string) string {
	return RED + msg + RESET
}
