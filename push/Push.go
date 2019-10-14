package push

var funcMap = make(map[string]func(msg string))

type PushToken struct {
	Label string `mapstructure:"label"`
	Value string `mapstructure:"value"`
}

func (token PushToken) Push(msg string) bool {
	sendDdMsg(token.Value, msg)
	return true
}
