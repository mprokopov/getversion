package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"reflect"
)

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func FindVersionStringInFile(re *regexp.Regexp, file string) string {
	dat, err := ioutil.ReadFile(file)
	CheckIfError(err)
	version := re.FindSubmatch(dat)

	if version == nil {
		log.Fatal("version is not found in " + file)
	}

	return string(version[1])
}

func GetGitBranch() string {
	res, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	CheckIfError(err)
	return strings.TrimSpace(string(res))
}

func GetPreReleaseLabel(branch string) string {
	// alpha for develop
	// beta for release or hotfix
	// feature name for feature
	// otherwise the same as input
	isFeature, _ := regexp.MatchString(`feature/`, branch)

	if isFeature {
		return strings.TrimLeft(branch, "feature/")
	}
	isDevelop, _ := regexp.MatchString(`^develop.*`, branch)

	if isDevelop {
		return "alpha"
	}

	isRelease, _ := regexp.MatchString(`^release.*|^hotfix.*`, branch)

	if isRelease {
		return "beta"
	}

	isMaster, _ := regexp.MatchString(`^master`, branch)

	if isMaster {
		return ""
	}

	return branch
}

func GetBaseVersion(source *string) *Version {
	var version *Version

	switch *source {
	case "gradle":
		re := regexp.MustCompile(`(?m)^version=(\d+.\d+.\d+).*`)
		str := FindVersionStringInFile(re, "gradle.properties")

		version = StrToVersion(str)
	case "node":
		re := regexp.MustCompile(`"version": "(\d+.\d+.\d+).*"`)
		str := FindVersionStringInFile(re, "package.json")

		version = StrToVersion(str)
		// case "git-tag":
	}
	return version
}

type Version struct {
	Major                           int
	Minor                           int
	Patch                           int
	CommitsSinceVersionSource       int
	CommitsSinceVersionSourcePadded string
	PreReleaseTag                   string
	PreReleaseTagWithDash           string
	PreReleaseLabel                 string
	PreReleaseNumber                int
	SemVer                          string
	BranchName                      string
	AssemblySemVer                  string
	BuildMetaData                   string
}

func StrToVersion(tag string) *Version {
	arr := strings.Split(tag, ".")
	major, _ := strconv.Atoi(arr[0])
	minor, _ := strconv.Atoi(arr[1])
	patch, _ := strconv.Atoi(arr[2])
	version := Version{Major: major, Minor: minor, Patch: patch}
	return &version
}

func VersionToA(version *Version) string {
	major := strconv.Itoa(version.Major)
	minor := strconv.Itoa(version.Minor)
	patch := strconv.Itoa(version.Patch)
	return major + "." + minor + "." + patch
}

func getField(v *Version, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

func main() {
	source := flag.String("source", "gradle", "version source")
	build_id := flag.String("build-id", "0", "build id")
	showvariable:= flag.String("showvariable", "", "variable to show")
	flag.Parse()

	var version *Version
	var isRelease = true

	switch *source {
	case "gradle":
		re := regexp.MustCompile(`(?m)^version=(\d+.\d+.\d+).*`)
		str := FindVersionStringInFile(re, "gradle.properties")

		version = StrToVersion(str)
	case "node":
		re := regexp.MustCompile(`"version": "(\d+.\d+.\d+).*"`)
		str := FindVersionStringInFile(re, "package.json")

		version = StrToVersion(str)
	case "git-tag":
		res, err := exec.Command("git", "rev-list", "--tags", "--no-walk", "--max-count=1").Output()
		if err != nil {
			log.Fatal("git rev-list does not have tags")
		}
		sha := strings.TrimSpace(string(res))
		res, err = exec.Command("git", "tag", "--points-at="+string(sha)).Output()
		if err != nil {
			log.Fatal("git tag failed")
		}
		tag := strings.TrimSpace(string(res))

		match, _ := regexp.Match(`\d+.\d+.\d+`, res)

		if !match {
			log.Fatal("version is not in proper format 0.0.0")
		}

		version = StrToVersion(tag)

		res, err = exec.Command("git", "rev-list", tag+"..", "--count").Output()
		if err != nil {
			log.Fatal("git rev-list doesn't return proper count")
		}

		commits_count := strings.TrimSpace(string(res))

		version.CommitsSinceVersionSource, err = strconv.Atoi(commits_count)
		version.CommitsSinceVersionSourcePadded = fmt.Sprintf("%04d", version.CommitsSinceVersionSource)

		CheckIfError(err)

		// Jenkins uses GIT_BRANCH for pipeline
		// and BRANCH_NAME for pultibranch pipeline
		if os.Getenv("BRANCH_NAME") != "" {
			version.BranchName = os.Getenv("BRANCH_NAME") // Jenkins sets this
		} else {
			version.BranchName = GetGitBranch()
		}

		isRelease, _ = regexp.MatchString(`^release.*|master|^hotfix.*`, version.BranchName)

	default:
		log.Fatal("wrong source for version, should be gradle or node")
	}

	version.PreReleaseLabel = GetPreReleaseLabel(version.BranchName)
	version.PreReleaseTag = version.PreReleaseLabel + "." + strconv.Itoa(version.CommitsSinceVersionSource)
	version.PreReleaseTagWithDash = "-" + version.PreReleaseTag

	isMaster, _ := regexp.MatchString(`^master`, version.BranchName)

	isHotfix, _ := regexp.MatchString(`^hotfix.*`, version.BranchName)

	if isHotfix {
		version.Patch = version.Patch + 1
	}

	if !isRelease {
		version.Minor = version.Minor + 1
	}

	if isMaster {
		version.SemVer = VersionToA(version)
	} else {
		version.SemVer = VersionToA(version) + version.PreReleaseTagWithDash
	}

	version.BuildMetaData = *build_id
	version.AssemblySemVer = version.SemVer + "." + *build_id
	jsonOutput, _ := json.Marshal(version)

	if *showvariable == "" {
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println(getField(version, *showvariable))
	}
}

// {
// 	"Major":0,
// 	"Minor":1,
// 	"Patch":0,
// 	"PreReleaseTag":"",
// 	"PreReleaseTagWithDash":"",
// 	"PreReleaseLabel":"",
// 	"PreReleaseNumber":"",
// 	"WeightedPreReleaseNumber":"",
// 	"BuildMetaData":"",
// 	"BuildMetaDataPadded":"",
// 	"FullBuildMetaData":"Branch.master.Sha.566280a7343576dc79a156f640473c5091f9244f",
// 	"MajorMinorPatch":"0.1.0",
// 	"SemVer":"0.1.0",
// 	"LegacySemVer":"0.1.0",
// 	"LegacySemVerPadded":"0.1.0",
// 	"AssemblySemVer":"0.1.0.0",
// 	"AssemblySemFileVer":"0.1.0.0",
// 	"FullSemVer":"0.1.0",
// 	"InformationalVersion":"0.1.0+Branch.master.Sha.566280a7343576dc79a156f640473c5091f9244f",
// 	"BranchName":"master",
// 	"Sha":"566280a7343576dc79a156f640473c5091f9244f",
// 	"ShortSha":"566280a",
// 	"NuGetVersionV2":"0.1.0",
// 	"NuGetVersion":"0.1.0",
// 	"NuGetPreReleaseTagV2":"",
// 	"NuGetPreReleaseTag":"",
// 	"VersionSourceSha":"566280a7343576dc79a156f640473c5091f9244f",
// 	"CommitsSinceVersionSource":1,
// 	"CommitsSinceVersionSourcePadded":"0001",
// 	"CommitDate":"2020-02-17"
// }
