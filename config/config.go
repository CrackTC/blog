package config

type Config struct {
	Port              int
	StaticPath        string
	TemplatePath      string
	IgnoredPaths      []string
	BlogDirectoryName string
	RedisURL          string
	RedisExpiration   string
}

var config *Config

func init() {
	config = &Config{
		Port:              8880,
		StaticPath:        "web/static",
		TemplatePath:      "web/template",
		BlogDirectoryName: "blog",
		IgnoredPaths:      []string{".git", ".github", ".gitignore", ".obsidian", ".obsidian.vimrc", "img", "cedict_ts.u8"},
		RedisURL:          "redis://localhost:6379",
		RedisExpiration:   "240h"}
}

func Get() *Config {
	return config
}
