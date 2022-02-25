package config

type Config struct {
	Urls *UrlSet
}

type UrlSet struct {
	Core string
}