package config

type Socnets struct {
	Telegram Telegram
}

type Telegram struct {
	Application Application
}

type Application struct {
	Title           string
	ShortName       string
	AppID           int
	AppHash         string
	ServerHost      string
	PublicKeyFile   string
	InitWarnChannel bool
}
