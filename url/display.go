package url

func Display(label, url string) string {
	return "\033]8;;" + url + "\a" + label + "\033]8;;\a"
}
