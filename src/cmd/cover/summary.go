package main

import (
	"fmt"
	"path"
	"slices"

	"golang.org/x/tools/cover"
)

func packageOutput(profile, outfile string) error {
	profiles, err := cover.ParseProfiles(profile)
	if err != nil {
		return err
	}
	profilesByPkg := make(map[string][]*cover.Profile)
	for _, profile := range profiles {
		pkg := path.Dir(profile.FileName)
		profilesByPkg[pkg] = append(profilesByPkg[pkg], profile)
	}

	sortedPkgs := make([]string, 0, len(profilesByPkg))
	pkgCoverages := make(map[string]float64)

	for pkg, profiles := range profilesByPkg {
		sortedPkgs = append(sortedPkgs, pkg)
		pkgCoverages[pkg] = packagePercentCovered(profiles)
	}
	slices.Sort(sortedPkgs)

	for _, pkg := range sortedPkgs {
		fmt.Printf("%s,%.1f\n", path.Base(pkg), pkgCoverages[pkg])
	}

	return nil
}

func packagePercentCovered(profiles []*cover.Profile) float64 {
	var total, covered int64
	for _, p := range profiles {
		for _, b := range p.Blocks {
			total += int64(b.NumStmt)
			if b.Count > 0 {
				covered += int64(b.NumStmt)
			}
		}
	}
	if total == 0 {
		return 0
	}
	return float64(covered) / float64(total) * 100
}
