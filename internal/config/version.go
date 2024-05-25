package config

import (
	"runtime/debug"
	"strings"
)

func VersionString() string {
	version := "v0.0.0-unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			switch s.Key {
			case "-tags":
				for _, t := range strings.Split(s.Value, ",") {
					if len(t) > 0 && t[0] == 'v' {
						return t
					}
				}
			case "vcs":
				version = s.Value
			case "vcs.revision":
				if version == "git" && len(s.Value) >= 10 {
					version += "-" + s.Value[:10]
				}
			case "vcs.modified":
				if s.Value == "true" {
					return (version + "-modified")
				}
			}
		}
	}
	return version
}
