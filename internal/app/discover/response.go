package discover

import "time"

type ProjectBase struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
}

type ProjectResponse struct {
	ProjectBase
	VersionGroups []string `json:"version_groups"`
	Versions      []string `json:"versions"`
}

type VersionFamilyResponse struct {
	ProjectBase
	VersionGroup string   `json:"version_group"`
	Versions     []string `json:"versions"`
}

type BuildsResponse struct {
	ProjectBase
	Version string         `json:"version"`
	Builds  []VersionBuild `json:"builds"`
}

type VersionBuild struct {
	Build     int                 `json:"build"`
	Time      time.Time           `json:"time"`
	Channel   string              `json:"channel"`
	Promoted  bool                `json:"promoted"`
	Changes   []Change            `json:"change"`
	Downloads map[string]Download `json:"downloads"`
}

type Change struct {
	Commit  string `json:"commit"`
	Summary string `json:"summary"`
	Message string `json:"message"`
}

type Download struct {
	Name   string `json:"name"`
	Sha256 string `json:"sha256"`
	URL    string `json:"url,omitempty"`
}

type VersionFamilyBuildsResponse struct {
	VersionFamilyResponse
	Builds []VersionFamilyBuild `json:"builds"`
}

type VersionFamilyBuild struct {
	VersionBuild
	Version string `json:"version"`
}

type VersionsResponse struct {
	Versions []VersionEntry `json:"versions"`
}

type VersionEntry struct {
	Version VersionInfo `json:"version"`
	Builds  []int       `json:"builds"`
}

type VersionInfo struct {
	ID      string         `json:"id"`
	Support SupportInfo    `json:"support"`
	Java    JavaInfo       `json:"java"`
}

type SupportInfo struct {
	Status string `json:"status"`
}

type JavaInfo struct {
	Version VersionDetails `json:"version"`
	Flags   FlagInfo       `json:"flags"`
}

type VersionDetails struct {
	Minimum int `json:"minimum"`
}

type FlagInfo struct {
	Recommended []string `json:"recommended"`
}

type BuildResponse struct {
	ID        int                    `json:"id"`
	Time      string                 `json:"time"`
	Channel   string                 `json:"channel"`
	Commits   []Commit               `json:"commits"`
	Downloads map[string]DownloadV3  `json:"downloads"`
}

type Commit struct {
	Sha     string `json:"sha"`
	Time    string `json:"time"`
	Message string `json:"message"`
}

type DownloadV3 struct {
	Name      string            `json:"name"`
	Checksums ChecksumInfo      `json:"checksums"`
	Size      int               `json:"size"`
	URL       string            `json:"url"`
}

type ChecksumInfo struct {
	Sha256 string `json:"sha256"`
}
