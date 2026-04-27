package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	. "github.com/treysu/velocity/internal/app/discover"
)

func main() {
	// Parse environment variables
	event := os.Getenv("GITHUB_EVENT_NAME")
	branch := os.Getenv("GITHUB_REF_NAME")
	dryRunStr := os.Getenv("DRY_RUN")
	dryRun, err := strconv.ParseBool(dryRunStr)
	if err != nil {
		dryRun = false
	}
	fmt.Println("Trigger event:", event)
	fmt.Println("Branch:", branch)
	fmt.Println("Dry run:", dryRun)

	client := resty.New()

	// Get existing tags and find max build number
	var tags []string
	dockerTags := GetExistingTags(DockerRepository)
	for _, dockerTag := range dockerTags {
		tags = append(tags, dockerTag.Name)
	}
	maxBuild := 0
	for _, tag := range tags {
		if strings.Contains(tag, "-") {
			parts := strings.Split(tag, "-")
			if len(parts) >= 2 {
				if build, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
					if build > maxBuild {
						maxBuild = build
					}
				}
			}
		}
	}

	// Get velocity versions
	var versionsResp VersionsResponse
	url := fmt.Sprintf("%s/projects/%s/versions", BaseURL, PROJECT)
	resp, err := client.R().Get(url)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(resp.Body(), &versionsResp)
	if err != nil {
		panic(err)
	}

	// Collect new builds from supported versions
	buildToVersion := make(map[int]string)
	var allBuilds []int
	for _, entry := range versionsResp.Versions {
		if entry.Version.Support.Status == "SUPPORTED" {
			for _, build := range entry.Builds {
				if build > maxBuild {
					allBuilds = append(allBuilds, build)
					buildToVersion[build] = entry.Version.ID
				}
			}
		}
	}
	sort.Ints(allBuilds)

	// Fetch build details for new builds
	var versionFamilyBuilds []VersionFamilyBuild
	for _, buildNum := range allBuilds {
		version := buildToVersion[buildNum]
		buildURL := fmt.Sprintf("%s/projects/%s/versions/%s/builds/%d", BaseURL, PROJECT, version, buildNum)
		resp, err := client.R().Get(buildURL)
		if err != nil {
			panic(err)
		}
		var buildResp BuildResponse
		err = json.Unmarshal(resp.Body(), &buildResp)
		if err != nil {
			panic(err)
		}

		// Create VersionFamilyBuild
		download := Download{
			Name:   buildResp.Downloads[DownloadsKey].Name,
			Sha256: buildResp.Downloads[DownloadsKey].Checksums.Sha256,
			URL:    buildResp.Downloads[DownloadsKey].URL,
		}
		downloads := map[string]Download{DownloadsKey: download}
		promoted := buildResp.Channel == "STABLE"
		var changes []Change
		for _, commit := range buildResp.Commits {
			changes = append(changes, Change{
				Commit:  commit.Sha,
				Summary: commit.Message,
				Message: commit.Message,
			})
		}
		timeParsed, err := time.Parse(time.RFC3339, buildResp.Time)
		if err != nil {
			timeParsed = time.Time{}
		}
		vb := VersionBuild{
			Build:     buildResp.ID,
			Time:      timeParsed,
			Channel:   buildResp.Channel,
			Promoted:  promoted,
			Changes:   changes,
			Downloads: downloads,
		}
		vfb := VersionFamilyBuild{
			VersionBuild: vb,
			Version:      version,
		}
		versionFamilyBuilds = append(versionFamilyBuilds, vfb)
	}

	// build promotions
	var eventForPromotions Event
	if event == "push" && branch == "main" {
		eventForPromotions = Rebuild
	} else {
		eventForPromotions = Cron
	}
	promotions := BuildPromotions(versionFamilyBuilds, tags, eventForPromotions)

	// Print tags to promotion
	fmt.Println("\nTags to build:")
	for _, promotion := range promotions {
		fmt.Println(strings.Join(strings.Split(promotion.DockerTags, "\\n"), ","))
	}

	// Build promote commands and write to scripts/promote.sh
	cmd := "#!/bin/sh\n\n"
	for _, promotion := range promotions {
		cmd += BuildCommand(DockerBuildWorkflow, promotion) + "\n"
	}

	// Create scripts folder
	err = os.MkdirAll("scripts", 0700)
	if err != nil {
		panic(err)
	}

	if dryRun {
		// Create empty scripts/dispatch.sh
		err := os.WriteFile("scripts/dispatch.sh", []byte("#!/bin/sh\n"), 0700)
		if err != nil {
			panic(err)
		}

		// print scripts content
		fmt.Println("\nThis is a dry run")
		fmt.Println("We generate the following script but not write to scripts/dispatch.sh")
		fmt.Println("\n" + cmd)
	} else {
		// Write to scripts/dispatch.sh
		err = os.WriteFile("scripts/dispatch.sh", []byte(cmd), 0700)
		if err != nil {
			panic(err)
		}

		fmt.Println("\nShell script has been generated to scripts/dispatch.sh")
	}
}
