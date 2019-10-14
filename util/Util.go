package util

func Add(l int, msg string) string {
	for len(msg) < l {
		msg = "0" + msg
	}
	return msg
}
