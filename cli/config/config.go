package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Core   CoreConfig   `toml:"core"`
	Gist   GistConfig   `toml:"gist"`
	Flag   FlagConfig   `toml:"flag"`
	Screen ScreenConfig `toml:"screen"`
}

type CoreConfig struct {
	Editor    string `toml:"editor"`
	SelectCmd string `toml:"selectcmd"`
	TomlFile  string `toml:"tomlfile"`
	User      string `toml:"user"`
}

type GistConfig struct {
	Token       string        `toml:"token"`
	BaseURL     string        `toml:"base_url"`
	ApiURL      string        `toml:"api_url"`
	Dir         string        `toml:"dir"`
	FileExt     string        `toml:"file_ext"`
	UseCache    bool          `toml:"use_cache"`
	CacheTTL    time.Duration `toml:"cache_ttl"`
	RunnableExt []string      `toml:"runnable_ext"`
}

type FlagConfig struct {
	OpenURL      bool `toml:"open_url"`
	BlogMode     bool `toml:"blog_mode"`
	StarredItems bool `toml:"starred"`

	NewPrivate  bool `toml:"-"`
	OpenBaseURL bool `toml:"-"`
}

type ScreenConfig struct {
	Columns []string `toml:"columns"`
}

var Conf Config

func GetDefaultDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	default:
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	case "windows":
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data")
		}
	}
	dir = filepath.Join(dir, "gist")

	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return dir, fmt.Errorf("cannot create directory: %v", err)
	}

	return dir, nil
}

func (cfg *Config) LoadFile(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		_, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return err
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	cfg.Core.Editor = os.Getenv("EDITOR")
	if cfg.Core.Editor == "" {
		cfg.Core.Editor = "vim"
	}
	cfg.Core.SelectCmd = "fzf-tmux --multi:fzf --multi:peco"
	cfg.Core.TomlFile = file
	cfg.Core.User = os.Getenv("USER")

	cfg.Gist.Token = os.Getenv("GITHUB_TOKEN")
	cfg.Gist.BaseURL = "https://gist.github.com"
	cfg.Gist.ApiURL = "https://api.github.com/api/v3/"
	dir := filepath.Join(filepath.Dir(file), "files")
	os.MkdirAll(dir, 0700)
	cfg.Gist.Dir = dir
	cfg.Gist.FileExt = ".patch"
	cfg.Gist.UseCache = true
	cfg.Gist.CacheTTL = time.Hour * 24
	cfg.Gist.RunnableExt = []string{"sh", "rb", "py", "pl", "php"}

	cfg.Flag.OpenURL = true
	cfg.Flag.BlogMode = true
	cfg.Flag.StarredItems = false

	cfg.Screen.Columns = []string{
		"{{.ShortID}}",
		"{{.PrivateMark}} {{.Filename}}",
		"{{.Description}}",
	}

	return toml.NewEncoder(f).Encode(cfg)
}
