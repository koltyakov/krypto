package main

// settings structure
type settings struct {
	UpdateFrequency string `json:"updateFrequency"` // possible values: "10s", "30s", ...
}

// getSettings retrieves setting from disk or returns defaults
func getSettings() (settings, error) {
	var defaults = settings{
		UpdateFrequency: "30s",
	}

	return defaults, nil
}

// getAppVersion gets application version number
func getAppVersion() string {
	if len(version) == 0 {
		return "0.0.0-SNAPSHOT"
	}
	return version
}
