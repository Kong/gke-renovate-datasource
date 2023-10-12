package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/samber/lo"
)

var (
	channelsReleaseNotes = map[string]string{
		"stable":  "https://cloud.google.com/feeds/gke-stable-channel-release-notes.xml",
		"regular": "https://cloud.google.com/feeds/gke-regular-channel-release-notes.xml",
		"rapid":   "https://cloud.google.com/feeds/gke-rapid-channel-release-notes.xml",
	}

	httpClient = &http.Client{}
)

type Channel struct {
	Releases []Release `json:"releases"`
}

type Release struct {
	Version string `json:"version"`
}

// This program is used to generate Google Kubernetes Engine JSON versions for Renovate custom datasource.
// Custom datasource docs: https://docs.renovatebot.com/modules/datasource/custom/
func main() {
	requestedChannel := flag.String("channel", "stable", "Channel to scrape")
	flag.Parse()

	channel, err := scrapeChannel(*requestedChannel)
	if err != nil {
		fmt.Println("Error scraping channel:", err)
		os.Exit(1)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", " ")
	if err := encoder.Encode(channel); err != nil {
		fmt.Println("Error encoding JSON:", err)
		os.Exit(1)
	}
}

func scrapeChannel(channelRequested string) (Channel, error) {
	channelURL, ok := channelsReleaseNotes[channelRequested]
	if !ok {
		return Channel{}, fmt.Errorf("unknown channel %q", channelRequested)
	}

	releases, err := scrapeReleases(channelURL)
	if err != nil {
		return Channel{}, fmt.Errorf("scraping %q: %w", channelURL, err)
	}

	return Channel{Releases: releases}, nil
}

func scrapeReleases(channelURL string) ([]Release, error) {
	resp, err := httpClient.Get(channelURL)
	if err != nil {
		return nil, fmt.Errorf("GET %q: %w", channelURL, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %q: status %d", channelURL, resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	versions, err := extractVersionsFromReleaseNotes(content)
	if err != nil {
		return nil, fmt.Errorf("extracting versions from release notes: %w", err)
	}

	var releases []Release
	for _, v := range lo.Uniq(versions) {
		releases = append(releases, Release{Version: v})
	}

	return releases, nil
}

var releaseRegexp = regexp.MustCompile(`<a href="https://github\.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-\d.\d+\.md#v\d+" class="external">(\d\.\d+\.\d+)-gke\.\d+</a>`)

func extractVersionsFromReleaseNotes(content []byte) ([]string, error) {
	var versions []string
	for _, version := range releaseRegexp.FindAllSubmatch(content, -1) {
		if len(version) != 2 {
			continue
		}
		versions = append(versions, string(version[1]))
	}

	return versions, nil
}
