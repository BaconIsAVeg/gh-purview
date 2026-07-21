package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BaconIsAVeg/gh-purview/internal/github"
	"gopkg.in/yaml.v3"
)

// Filter is a named saved GitHub search query.
type Filter struct {
	Query string `yaml:"query"`
}

// Config is the user-authored configuration loaded from config.yml.
// The application never writes this file, so user comments are preserved.
type Config struct {
	Default *Filter           `yaml:"default,omitempty"`
	Filters map[string]Filter `yaml:"filters,omitempty"`
	path    string            `yaml:"-"`
}

// Last is the app-managed resume state written to last.yml after every
// in-app filter edit. It lives in the cache directory, not the config
// directory, so it is disposable derived state rather than user config.
type Last struct {
	Query  string `yaml:"query"`
	Filter string `yaml:"filter,omitempty"`
}

// Dir returns the configuration directory (XDG_CONFIG_HOME/gh-purview).
func Dir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("locate config dir: %w", err)
	}
	return filepath.Join(base, "gh-purview"), nil
}

// Path returns the full path to config.yml.
func Path() (string, error) {
	d, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "config.yml"), nil
}

// Load reads config.yml. If the file does not exist, it is created with the
// built-in default filter so the user can simply edit it in place rather than
// building the path and contents by hand. A missing file is therefore never a
// non-error outcome: either a real Config is returned or an error is.
func Load() (*Config, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			if seedErr := seedDefaultConfig(p); seedErr != nil {
				return nil, fmt.Errorf("create default config: %w", seedErr)
			}
			data, err = os.ReadFile(p)
			if err != nil {
				return nil, fmt.Errorf("read seeded config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", p, err)
	}
	c.path = p
	return &c, nil
}

// seedDefaultConfig writes a commented config.yml containing the built-in
// default filter to path, creating parent directories as needed.
func seedDefaultConfig(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	cfg := Config{
		Default: &Filter{Query: github.DefaultQuery},
	}
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("marshal default config: %w", err)
	}
	header := []byte("# gh-purview configuration. Edit this file to customize startup filters.\n" +
		"# Run gh-purview --filter <name> to use a named filter below.\n\n")
	return os.WriteFile(path, append(header, data...), 0o644)
}

// cacheDir returns the cache directory (XDG_CACHE_HOME/gh-purview).
func cacheDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("locate cache dir: %w", err)
	}
	return filepath.Join(base, "gh-purview"), nil
}

// LastPath returns the full path to last.yml.
func LastPath() (string, error) {
	d, err := cacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "last.yml"), nil
}

// LoadLast reads last.yml. A missing file returns (nil, nil).
func LoadLast() (*Last, error) {
	p, err := LastPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read last: %w", err)
	}
	var l Last
	if err := yaml.Unmarshal(data, &l); err != nil {
		return nil, fmt.Errorf("parse last %s: %w", p, err)
	}
	return &l, nil
}

// SaveLast writes last.yml, creating the cache directory as needed.
func SaveLast(l *Last) error {
	p, err := LastPath()
	if err != nil {
		return err
	}
	d := filepath.Dir(p)
	if err := os.MkdirAll(d, 0o755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}
	data, err := yaml.Marshal(l)
	if err != nil {
		return fmt.Errorf("marshal last: %w", err)
	}
	if err := os.WriteFile(p, data, 0o644); err != nil {
		return fmt.Errorf("write last: %w", err)
	}
	return nil
}

// Resolve picks the startup query by precedence:
//  1. filterFlag ("--filter <name>") -> named filter (hard error if missing)
//  2. last.Query (resume prior in-app edit)
//  3. Config.Default.Query
//  4. "" (caller falls back to the built-in default)
//
// activeFilter is the name of the filter the query came from, or "" when the
// source is the default section or the built-in fallback. It is recorded into
// last.Filter on the next in-app edit so the resume state remembers its basis.
func Resolve(cfg *Config, last *Last, filterFlag string) (query, activeFilter string, err error) {
	if filterFlag != "" {
		if filterFlag == "default" && cfg != nil && cfg.Default != nil && cfg.Default.Query != "" {
			return cfg.Default.Query, "default", nil
		}
		if cfg == nil {
			return "", "", fmt.Errorf("filter %q not found: no config file at %s", filterFlag, configPathOrUnknown(cfg))
		}
		if cfg.Filters == nil {
			return "", "", fmt.Errorf("filter %q not found: no filters defined in %s", filterFlag, configPathOrUnknown(cfg))
		}
		f, ok := cfg.Filters[filterFlag]
		if !ok {
			return "", "", fmt.Errorf("filter %q not found in %s", filterFlag, configPathOrUnknown(cfg))
		}
		return f.Query, filterFlag, nil
	}
	if last != nil && last.Query != "" {
		return last.Query, last.Filter, nil
	}
	if cfg != nil && cfg.Default != nil && cfg.Default.Query != "" {
		return cfg.Default.Query, "", nil
	}
	return "", "", nil
}

func configPathOrUnknown(cfg *Config) string {
	if cfg != nil && cfg.path != "" {
		return cfg.path
	}
	p, err := Path()
	if err != nil {
		return "config.yml"
	}
	return p
}
