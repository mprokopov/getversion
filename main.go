package main

import (
	"strconv"
	"flag"
	"regexp"
	"fmt"
	"log"
	"io/ioutil"
	"os/exec"
	"strings"
	"encoding/json"
	"os"
)
// func check(e error) {
// 	if e != nil {
// 		log.Fatal(e)
// 	}
// }

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func FindInFile(re *regexp.Regexp, file string) string {
	dat, err := ioutil.ReadFile(file)
	CheckIfError(err)
	version := re.FindSubmatch(dat)

	if version == nil {
		log.Fatal("version is not found in gradle.properties")
	}

	return string(version[1])
}

func GetGitBranch() string{
	res, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	CheckIfError(err)
	return strings.TrimSpace(string(res))
}

type Version struct  {
	Major int
	Minor int
	Patch int
	CommitsSinceVersionSource int
	PreReleaseTag string
	PreReleaseLabel string
	PreReleaseNumber int
	SemVer string
	BranchName string
	AssemblySemVer string
}

func SplitTag(tag string) *Version {
	arr := strings.Split(tag, ".")
	major, _ := strconv.Atoi(arr[0])
	minor, _ := strconv.Atoi(arr[1])
	patch, _ := strconv.Atoi(arr[2])
	version := Version{ Major: major, Minor: minor, Patch: patch }
	return &version
}

func main() {
	source := flag.String("source", "gradle", "version source")
	flag.Parse()

	switch *source {
	case "gradle":
		re := regexp.MustCompile(`(?m)^version=(\d+.\d+.\d+).*`)
		version := FindInFile(re, "gradle.properties")

		fmt.Print(version)
	case "node":
		re := regexp.MustCompile(`"version": "(\d+.\d+.\d+).*"`)
		version := FindInFile(re, "package.json")

		fmt.Print(version)
	case "git":
		res, err := exec.Command("git", "rev-list", "--tags", "--no-walk", "--max-count=1").Output()
		if err != nil {
			log.Fatal("git rev-list does not have tags")
		}
		sha := strings.TrimSpace(string(res))
		res, err = exec.Command("git", "tag","--points-at="+string(sha)).Output()
		if err != nil {
			log.Fatal("git tag failed")
		}
		tag := strings.TrimSpace(string(res))

		match,_ := regexp.Match(`\d+.\d+.\d+`, res)

		if !match {
			log.Fatal("version is not in proper format 0.0.0")
		}

		version := SplitTag(tag)

		res, err = exec.Command("git", "rev-list", tag+"..", "--count").Output()
		if err != nil {
			log.Fatal("git rev-list doesn't return proper count")
		}

		commits_count := strings.TrimSpace(string(res))

		version.CommitsSinceVersionSource, err = strconv.Atoi(commits_count)

		CheckIfError(err)

		version.BranchName =  GetGitBranch()

		version.PreReleaseLabel = version.BranchName // remove feature/

		jsVersion,_ := json.Marshal(version)
		fmt.Println(string(jsVersion))
	default:
		log.Fatal("wrong source for version, should be gradle or node")
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
