package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolve(t *testing.T) {
	cfg := &Config{
		Default: &Filter{Query: "default-query"},
		Filters: map[string]Filter{
			"mine":   {Query: "author:@me"},
			"drafts": {Query: "draft:true"},
		},
	}

	tests := []struct {
		name       string
		cfg        *Config
		last       *Last
		flag       string
		wantQuery  string
		wantFilter string
		wantErr    bool
	}{
		{
			name:       "flag selects named filter",
			cfg:        cfg,
			flag:       "mine",
			wantQuery:  "author:@me",
			wantFilter: "mine",
		},
		{
			name:       "flag default resolves to default section",
			cfg:        cfg,
			flag:       "default",
			wantQuery:  "default-query",
			wantFilter: "default",
		},
		{
			name:    "flag missing filter hard errors",
			cfg:     cfg,
			flag:    "nope",
			wantErr: true,
		},
		{
			name:      "no flag, no last falls back to default section",
			cfg:       cfg,
			wantQuery: "default-query",
		},
		{
			name:       "no flag, last present resumes last query",
			cfg:        cfg,
			last:       &Last{Query: "edited-query", Filter: "mine"},
			wantQuery:  "edited-query",
			wantFilter: "mine",
		},
		{
			name:      "no config, no last, no flag yields empty (built-in fallback)",
			wantQuery: "",
		},
		{
			name:    "flag with nil config errors",
			flag:    "mine",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotFilter, err := Resolve(tt.cfg, tt.last, tt.flag)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotQuery != tt.wantQuery {
				t.Errorf("query = %q, want %q", gotQuery, tt.wantQuery)
			}
			if gotFilter != tt.wantFilter {
				t.Errorf("filter = %q, want %q", gotFilter, tt.wantFilter)
			}
		})
	}
}

func TestSaveLoadLast(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", dir)

	want := &Last{Query: "author:@me", Filter: "mine"}
	if err := SaveLast(want); err != nil {
		t.Fatalf("SaveLast: %v", err)
	}

	got, err := LoadLast()
	if err != nil {
		t.Fatalf("LoadLast: %v", err)
	}
	if got == nil {
		t.Fatal("LoadLast returned nil")
	}
	if got.Query != want.Query || got.Filter != want.Filter {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestLoadLastMissing(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", dir)

	got, err := LoadLast()
	if err != nil {
		t.Fatalf("LoadLast on missing: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestLoadMissingConfigSeedsDefault(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load on missing: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil Config")
	}
	if cfg.Default == nil || cfg.Default.Query == "" {
		t.Fatal("expected seeded default filter")
	}
	// file should now exist on disk
	if _, err := os.Stat(filepath.Join(dir, "gh-purview", "config.yml")); err != nil {
		t.Errorf("seeded config file not created: %v", err)
	}
	q, _, err := Resolve(cfg, nil, "")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if q == "" {
		t.Error("expected seeded default query, got empty")
	}
}
