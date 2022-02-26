package config

type Config struct {
	JaegerUrl string
	Urls      *UrlSet
}

type UrlSet struct {
	External string
	Core     string
}
