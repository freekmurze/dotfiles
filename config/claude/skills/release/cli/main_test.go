package main

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseSparkleSignature(t *testing.T) {
	output := `<enclosure sparkle:edSignature="abc123==" length="456789" />`

	signature, length := parseSparkleSignature(output)

	if signature != "abc123==" {
		t.Fatalf("signature = %q, want %q", signature, "abc123==")
	}
	if length != "456789" {
		t.Fatalf("length = %q, want %q", length, "456789")
	}
}

func TestConfigDefaults(t *testing.T) {
	cfg := Config{
		AppName:    "TestApp",
		GitHubRepo: "owner/repo",
	}
	cfg.applyDefaults()

	if cfg.BundleName != "TestApp.app" {
		t.Fatalf("BundleName = %q, want %q", cfg.BundleName, "TestApp.app")
	}
	if cfg.DMGName != "TestApp.dmg" {
		t.Fatalf("DMGName = %q, want %q", cfg.DMGName, "TestApp.dmg")
	}
	if cfg.GitBranch != "main" {
		t.Fatalf("GitBranch = %q, want %q", cfg.GitBranch, "main")
	}
	if cfg.MinSystemVersion != "14.0" {
		t.Fatalf("MinSystemVersion = %q, want %q", cfg.MinSystemVersion, "14.0")
	}
	if cfg.AppcastFile != "appcast.xml" {
		t.Fatalf("AppcastFile = %q, want %q", cfg.AppcastFile, "appcast.xml")
	}
	if len(cfg.DerivedDataPrefixes) != 1 || cfg.DerivedDataPrefixes[0] != "TestApp-" {
		t.Fatalf("DerivedDataPrefixes = %v, want [TestApp-]", cfg.DerivedDataPrefixes)
	}
}

func TestConfigAppcastURL(t *testing.T) {
	cfg := Config{
		AppName:    "TestApp",
		GitHubRepo: "owner/repo",
	}
	cfg.applyDefaults()

	want := "https://raw.githubusercontent.com/owner/repo/main/appcast.xml"
	if got := cfg.appcastURL(); got != want {
		t.Fatalf("appcastURL = %q, want %q", got, want)
	}
}

func TestLoadConfigFromJSON(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{
		AppName:          "MyApp",
		GitHubRepo:       "me/myapp",
		MinSystemVersion: "15.0",
	}
	data, _ := json.Marshal(cfg)
	os.WriteFile(filepath.Join(dir, "release.json"), data, 0644)

	raw, _ := os.ReadFile(filepath.Join(dir, "release.json"))
	var loaded Config
	if err := json.Unmarshal(raw, &loaded); err != nil {
		t.Fatal(err)
	}
	loaded.applyDefaults()

	if loaded.AppName != "MyApp" {
		t.Fatalf("AppName = %q, want MyApp", loaded.AppName)
	}
	if loaded.MinSystemVersion != "15.0" {
		t.Fatalf("MinSystemVersion = %q, want 15.0", loaded.MinSystemVersion)
	}
	if loaded.BundleName != "MyApp.app" {
		t.Fatalf("BundleName = %q, want MyApp.app", loaded.BundleName)
	}
}

func TestWriteAppcast(t *testing.T) {
	cfg := Config{
		AppName:    "TestApp",
		GitHubRepo: "owner/repo",
	}
	cfg.applyDefaults()

	path := filepath.Join(t.TempDir(), "appcast.xml")
	items := []Item{
		{
			Title:              "Version 1.0",
			Version:            "1",
			ShortVersionString: "1.0",
			MinSystemVersion:   cfg.MinSystemVersion,
			PubDate:            "Tue, 12 May 2026 09:30:00 +0000",
			Description:        buildDescription("1.0", []string{"Initial release"}),
			Enclosure: Enclosure{
				URL:         "https://github.com/owner/repo/releases/download/v1.0/TestApp.dmg",
				Type:        "application/octet-stream",
				EdSignature: "sig==",
				Length:      "123",
			},
		},
	}

	if err := writeAppcast(path, cfg, items); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	text := string(data)
	if !strings.Contains(text, "<title>TestApp Updates</title>") {
		t.Fatal("appcast title was not written correctly")
	}

	var appcast Appcast
	if err := xml.Unmarshal(data, &appcast); err != nil {
		t.Fatal(err)
	}
	if got := len(appcast.Channel.Items); got != 1 {
		t.Fatalf("item count = %d, want 1", got)
	}
	if appcast.Channel.Items[0].Version != "1" {
		t.Fatalf("version = %q, want 1", appcast.Channel.Items[0].Version)
	}
}
