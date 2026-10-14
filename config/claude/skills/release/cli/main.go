// macos-release is a CLI tool that automates the full release pipeline for
// native macOS apps distributed via GitHub Releases with Sparkle auto-update.
//
// It reads configuration from a release.json file in the project root and
// handles: DMG creation, EdDSA signing, appcast.xml updates, git push, and
// GitHub release creation.
//
// Usage:
//
//	go run github.com/fayazara/macos-app-skills/release/cli@latest
//
// Or build locally:
//
//	go build -o macos-release ./release/cli && ./macos-release
package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Terminal colors
const (
	red    = "\033[0;31m"
	green  = "\033[0;32m"
	yellow = "\033[0;33m"
	cyan   = "\033[0;36m"
	bold   = "\033[1m"
	reset  = "\033[0m"
)

func step(msg string)    { fmt.Printf("\n%s%s==> %s%s\n", cyan, bold, msg, reset) }
func success(msg string) { fmt.Printf("%s  OK %s%s\n", green, msg, reset) }
func warn(msg string)    { fmt.Printf("%s  WARN %s%s\n", yellow, msg, reset) }

func fail(msg string) {
	fmt.Printf("%s  ERROR %s%s\n", red, msg, reset)
	os.Exit(1)
}

// Config holds the release configuration loaded from release.json.
type Config struct {
	// AppName is the display name (e.g., "MyApp"). Used for the DMG volume
	// name and release titles.
	AppName string `json:"app_name"`

	// BundleName is the .app bundle filename (e.g., "MyApp.app").
	// Defaults to AppName + ".app" if omitted.
	BundleName string `json:"bundle_name"`

	// DMGName is the output DMG filename (e.g., "MyApp.dmg").
	// Defaults to AppName + ".dmg" if omitted.
	DMGName string `json:"dmg_name"`

	// GitHubRepo is the owner/repo string (e.g., "fayazara/myapp").
	GitHubRepo string `json:"github_repo"`

	// GitBranch is the branch to push appcast changes to. Defaults to "main".
	GitBranch string `json:"git_branch"`

	// MinSystemVersion is the sparkle:minimumSystemVersion value.
	// Defaults to "14.0".
	MinSystemVersion string `json:"min_system_version"`

	// AppcastFile is the filename of the appcast XML. Defaults to "appcast.xml".
	AppcastFile string `json:"appcast_file"`

	// DerivedDataPrefixes is a list of DerivedData directory prefixes to
	// search for Sparkle's sign_update binary. Defaults to [AppName + "-"].
	DerivedDataPrefixes []string `json:"derived_data_prefixes"`
}

func (c *Config) applyDefaults() {
	if c.BundleName == "" {
		c.BundleName = c.AppName + ".app"
	}
	if c.DMGName == "" {
		c.DMGName = c.AppName + ".dmg"
	}
	if c.GitBranch == "" {
		c.GitBranch = "main"
	}
	if c.MinSystemVersion == "" {
		c.MinSystemVersion = "14.0"
	}
	if c.AppcastFile == "" {
		c.AppcastFile = "appcast.xml"
	}
	if len(c.DerivedDataPrefixes) == 0 {
		c.DerivedDataPrefixes = []string{c.AppName + "-"}
	}
}

func (c *Config) validate() {
	if c.AppName == "" {
		fail("release.json: app_name is required")
	}
	if c.GitHubRepo == "" {
		fail("release.json: github_repo is required")
	}
}

func (c *Config) appcastURL() string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s",
		c.GitHubRepo, c.GitBranch, c.AppcastFile)
}

// Sparkle appcast XML types
type Appcast struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title    string `xml:"title"`
	Link     string `xml:"link"`
	Language string `xml:"language"`
	Items    []Item `xml:"item"`
}

type Item struct {
	Title              string    `xml:"title"`
	Version            string    `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle version"`
	ShortVersionString string    `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle shortVersionString"`
	MinSystemVersion   string    `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle minimumSystemVersion"`
	PubDate            string    `xml:"pubDate"`
	Description        string    `xml:"description"`
	Enclosure          Enclosure `xml:"enclosure"`
}

type Enclosure struct {
	URL         string `xml:"url,attr"`
	Type        string `xml:"type,attr"`
	EdSignature string `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle edSignature,attr"`
	Length      string `xml:"length,attr"`
}

func main() {
	homeDir, _ := os.UserHomeDir()

	repoDir := findRepoDir()
	cfg := loadConfig(repoDir)

	appPath := filepath.Join(homeDir, "Downloads", cfg.BundleName)
	dmgPath := filepath.Join(homeDir, "Downloads", cfg.DMGName)
	appcastPath := filepath.Join(repoDir, cfg.AppcastFile)

	fmt.Printf("\n%s=======================================%s\n", bold, reset)
	fmt.Printf("%s  %s Release Manager%s\n", bold, cfg.AppName, reset)
	fmt.Printf("%s=======================================%s\n", bold, reset)

	// -- Preflight --

	step("Checking prerequisites...")

	requireCommand("create-dmg", "Install with: brew install create-dmg")
	requireCommand("gh", "Install with: brew install gh")
	requireCommand("git", "")
	requireCommand("plutil", "")

	signUpdate := findSignUpdate(homeDir, cfg.DerivedDataPrefixes)
	if signUpdate == "" {
		fail("Sparkle sign_update not found in DerivedData. Build the project once first.")
	}
	success("All tools found")

	// -- Validate app --

	step("Validating " + appPath + "...")

	info, err := os.Stat(appPath)
	if err != nil || !info.IsDir() {
		fail(cfg.BundleName + " not found in ~/Downloads. Export it from Xcode first.")
	}

	plist := filepath.Join(appPath, "Contents", "Info.plist")
	version, err := plistValue(plist, "CFBundleShortVersionString")
	if err != nil {
		fail("Could not read version from Info.plist")
	}

	build, err := plistValue(plist, "CFBundleVersion")
	if err != nil {
		fail("Could not read build number from Info.plist")
	}

	if feedURL, err := plistValue(plist, "SUFeedURL"); err != nil {
		warn("SUFeedURL not found in Info.plist")
	} else if feedURL != cfg.appcastURL() {
		warn("SUFeedURL is " + feedURL + ", expected " + cfg.appcastURL())
	}

	if _, err := plistValue(plist, "SUPublicEDKey"); err != nil {
		warn("SUPublicEDKey not found in Info.plist")
	}

	fmt.Printf("  Version: %s%s%s  Build: %s%s%s\n", bold, version, reset, bold, build, reset)

	existingData, _ := os.ReadFile(appcastPath)
	if strings.Contains(string(existingData), "sparkle:version>"+build+"<") {
		warn(fmt.Sprintf("Build %s already exists in appcast.xml", build))
		if !confirm("Continue anyway?", false) {
			os.Exit(0)
		}
	}

	success("App validated")

	// -- Release notes --

	step("Release notes (one bullet point per line, empty line to finish):")
	fmt.Printf("  %sEnter your release notes below:%s\n", yellow, reset)

	var notes []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break
		}
		notes = append(notes, line)
	}
	if err := scanner.Err(); err != nil {
		fail("Could not read release notes: " + err.Error())
	}
	if len(notes) == 0 {
		fail("No release notes provided")
	}

	// -- Confirm --

	fmt.Printf("\n  %sRelease summary:%s\n", bold, reset)
	fmt.Printf("  App:     %s\n", cfg.AppName)
	fmt.Printf("  Version: %s (build %s)\n", version, build)
	fmt.Printf("  Tag:     v%s\n", version)
	fmt.Println("  Notes:")
	for _, n := range notes {
		fmt.Printf("    - %s\n", n)
	}
	fmt.Println()

	if !confirm("Proceed with release?", true) {
		os.Exit(0)
	}

	// -- Create DMG --

	step("Creating DMG...")

	_ = os.Remove(dmgPath)
	dmgArgs := []string{
		"--volname", cfg.AppName,
		"--window-pos", "200", "120",
		"--window-size", "600", "400",
		"--icon-size", "100",
		"--icon", cfg.BundleName, "150", "185",
		"--app-drop-link", "450", "185",
		dmgPath,
		appPath,
	}

	out, err := runCmd("create-dmg", dmgArgs...)
	if err != nil {
		if _, statErr := os.Stat(dmgPath); statErr != nil {
			fail(fmt.Sprintf("DMG creation failed: %s\n%s", err, out))
		}
	}

	if _, err := os.Stat(dmgPath); err != nil {
		fail("DMG creation failed: file not found")
	}
	success("DMG created at " + dmgPath)

	// -- Sign --

	step("Signing DMG with Sparkle...")

	signOut, err := runCmd(signUpdate, dmgPath)
	if err != nil {
		fail(fmt.Sprintf("sign_update failed: %s\n%s", err, signOut))
	}

	signature, length := parseSparkleSignature(signOut)
	success(fmt.Sprintf("Signed (length: %s bytes)", length))

	// -- Update appcast --

	step("Updating " + cfg.AppcastFile + "...")

	appcastData, err := os.ReadFile(appcastPath)
	if err != nil {
		fail("Could not read " + cfg.AppcastFile + ": " + err.Error())
	}

	var appcast Appcast
	if err := xml.Unmarshal(appcastData, &appcast); err != nil {
		fail("Could not parse " + cfg.AppcastFile + ": " + err.Error())
	}

	pubDate := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 +0000")
	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/v%s/%s",
		cfg.GitHubRepo, version, cfg.DMGName)

	newItem := Item{
		Title:              fmt.Sprintf("Version %s", version),
		Version:            build,
		ShortVersionString: version,
		MinSystemVersion:   cfg.MinSystemVersion,
		PubDate:            pubDate,
		Description:        buildDescription(version, notes),
		Enclosure: Enclosure{
			URL:         downloadURL,
			Type:        "application/octet-stream",
			EdSignature: signature,
			Length:      length,
		},
	}

	allItems := make([]Item, 0, len(appcast.Channel.Items)+1)
	allItems = append(allItems, newItem)
	allItems = append(allItems, appcast.Channel.Items...)

	if err := writeAppcast(appcastPath, cfg, allItems); err != nil {
		fail("Could not write " + cfg.AppcastFile + ": " + err.Error())
	}
	success(fmt.Sprintf("Appcast updated with v%s", version))

	// -- Push --

	step("Pushing appcast to GitHub...")

	if _, err := runCmd("git", "-C", repoDir, "add", cfg.AppcastFile); err != nil {
		fail("git add failed: " + err.Error())
	}

	commitMsg := fmt.Sprintf("Release v%s appcast", version)
	if out, err := runCmd("git", "-C", repoDir, "commit", "--only", cfg.AppcastFile, "-m", commitMsg); err != nil {
		fail(fmt.Sprintf("git commit failed: %s\n%s", err, out))
	}

	if out, err := runCmd("git", "-C", repoDir, "push", "origin", cfg.GitBranch); err != nil {
		fail(fmt.Sprintf("git push failed: %s\n%s", err, out))
	}
	success("Pushed to " + cfg.GitBranch)

	// -- GitHub release --

	step("Creating GitHub release...")

	var mdNotes strings.Builder
	mdNotes.WriteString("## What's New\n\n")
	for _, n := range notes {
		mdNotes.WriteString(fmt.Sprintf("- %s\n", n))
	}

	releaseURL, err := runCmd("gh", "release", "create",
		"v"+version,
		dmgPath,
		"--repo", cfg.GitHubRepo,
		"--title", "v"+version,
		"--notes", mdNotes.String(),
	)
	if err != nil {
		fail(fmt.Sprintf("gh release create failed: %s\n%s", err, releaseURL))
	}
	success("Release created")

	fmt.Printf("\n%s%s=======================================%s\n", green, bold, reset)
	fmt.Printf("%s%s  Released %s v%s%s\n", green, bold, cfg.AppName, version, reset)
	fmt.Printf("%s%s  %s%s\n", green, bold, releaseURL, reset)
	fmt.Printf("%s%s=======================================%s\n\n", green, bold, reset)
}

// --- Config loading ---

func loadConfig(repoDir string) Config {
	configPath := filepath.Join(repoDir, "release.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		fail("release.json not found in " + repoDir + "\nCreate one with at minimum: {\"app_name\": \"MyApp\", \"github_repo\": \"owner/repo\"}")
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		fail("Could not parse release.json: " + err.Error())
	}

	cfg.validate()
	cfg.applyDefaults()
	return cfg
}

func findRepoDir() string {
	// Walk up from cwd looking for release.json
	if cwd, err := os.Getwd(); err == nil {
		for dir := cwd; ; dir = filepath.Dir(dir) {
			if fileExists(filepath.Join(dir, "release.json")) {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}
	}

	fail("Could not find release.json. Run this command from your project directory, or create a release.json file.")
	return ""
}

// --- Tool discovery ---

func findSignUpdate(homeDir string, prefixes []string) string {
	derivedData := filepath.Join(homeDir, "Library", "Developer", "Xcode", "DerivedData")
	entries, err := os.ReadDir(derivedData)
	if err != nil {
		return ""
	}

	for _, prefix := range prefixes {
		for _, entry := range entries {
			if !strings.HasPrefix(entry.Name(), prefix) {
				continue
			}

			candidate := filepath.Join(derivedData, entry.Name(),
				"SourcePackages", "artifacts", "sparkle", "Sparkle", "bin", "sign_update")
			if fileExists(candidate) {
				return candidate
			}
		}
	}
	return ""
}

func requireCommand(name, installHint string) {
	if commandExists(name) {
		return
	}
	if installHint == "" {
		fail(name + " not found")
	}
	fail(name + " not found. " + installHint)
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func plistValue(plistPath, key string) (string, error) {
	out, err := runCmd("plutil", "-extract", key, "raw", "-o", "-", plistPath)
	if err != nil {
		return "", fmt.Errorf("key %q not found", key)
	}
	return out, nil
}

func confirm(prompt string, defaultYes bool) bool {
	hint := "(Y/n)"
	if !defaultYes {
		hint = "(y/N)"
	}
	fmt.Printf("  %s %s ", prompt, hint)

	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))
	if line == "" {
		return defaultYes
	}
	return line == "y" || line == "yes"
}

// --- Sparkle signing ---

func parseSparkleSignature(output string) (string, string) {
	sigRe := regexp.MustCompile(`sparkle:edSignature="([^"]+)"`)
	lenRe := regexp.MustCompile(`length="([^"]+)"`)

	sigMatch := sigRe.FindStringSubmatch(output)
	lenMatch := lenRe.FindStringSubmatch(output)

	if len(sigMatch) < 2 {
		fail("Could not parse signature from sign_update output:\n" + output)
	}
	if len(lenMatch) < 2 {
		fail("Could not parse length from sign_update output:\n" + output)
	}

	return sigMatch[1], lenMatch[1]
}

// --- Appcast writing ---

func buildDescription(version string, notes []string) string {
	var htmlItems strings.Builder
	for _, note := range notes {
		htmlItems.WriteString(fmt.Sprintf("          <li>%s</li>\n", xmlEscapeText(note)))
	}

	return fmt.Sprintf("<![CDATA[\n        <h2>What's New in %s</h2>\n        <ul>\n%s        </ul>\n      ]]>",
		xmlEscapeText(version), htmlItems.String())
}

func descriptionToCDATA(desc string) string {
	trimmed := strings.TrimSpace(desc)
	if strings.HasPrefix(trimmed, "<![CDATA[") {
		return desc
	}
	return "<![CDATA[\n        " + trimmed + "\n      ]]>"
}

func writeAppcast(path string, cfg Config, items []Item) error {
	var b strings.Builder

	b.WriteString(`<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0"
  xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle"
  xmlns:dc="http://purl.org/dc/elements/1.1/">
  <channel>
`)
	b.WriteString(fmt.Sprintf("    <title>%s Updates</title>\n", xmlEscapeText(cfg.AppName)))
	b.WriteString(fmt.Sprintf("    <link>%s</link>\n", xmlEscapeText(cfg.appcastURL())))
	b.WriteString("    <language>en</language>\n")

	for _, item := range items {
		desc := descriptionToCDATA(item.Description)

		b.WriteString("\n    <item>\n")
		b.WriteString(fmt.Sprintf("      <title>%s</title>\n", xmlEscapeText(item.Title)))
		b.WriteString(fmt.Sprintf("      <sparkle:version>%s</sparkle:version>\n", xmlEscapeText(item.Version)))
		b.WriteString(fmt.Sprintf("      <sparkle:shortVersionString>%s</sparkle:shortVersionString>\n", xmlEscapeText(item.ShortVersionString)))
		b.WriteString(fmt.Sprintf("      <sparkle:minimumSystemVersion>%s</sparkle:minimumSystemVersion>\n", xmlEscapeText(item.MinSystemVersion)))
		b.WriteString(fmt.Sprintf("      <pubDate>%s</pubDate>\n", xmlEscapeText(item.PubDate)))
		b.WriteString(fmt.Sprintf("      <description>%s</description>\n", desc))
		b.WriteString("      <enclosure\n")
		b.WriteString(fmt.Sprintf("        url=\"%s\"\n", xmlEscapeAttr(item.Enclosure.URL)))
		b.WriteString(fmt.Sprintf("        type=\"%s\"\n", xmlEscapeAttr(item.Enclosure.Type)))
		b.WriteString(fmt.Sprintf("        sparkle:edSignature=\"%s\"\n", xmlEscapeAttr(item.Enclosure.EdSignature)))
		b.WriteString(fmt.Sprintf("        length=\"%s\"\n", xmlEscapeAttr(item.Enclosure.Length)))
		b.WriteString("      />\n")
		b.WriteString("    </item>\n")
	}

	b.WriteString("\n  </channel>\n</rss>\n")
	return os.WriteFile(path, []byte(b.String()), 0644)
}

func xmlEscapeText(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func xmlEscapeAttr(s string) string {
	s = xmlEscapeText(s)
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
