package main

import (
	"flag"
	"regexp"
	"fmt"
	"log"
	"io/ioutil"
)
func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func FindInFile(re *regexp.Regexp, file string) string {
	dat, err := ioutil.ReadFile(file)
	check(err)
	version := re.FindSubmatch(dat)

	if version == nil {
		log.Fatal("version is not found in gradle.properties")
	}

	return string(version[1])
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
	// case "git":
	// 	sh(returnStdout: true, script: 'git rev-list --tags --no-walk --max-count=1').trim()
	// 	sh(returnStdout: true, script: "git tag --points-at=${TAG_SHA}").trim()
	// 	sh(returnStdout: true, script: "git rev-list ${TAG_SHA}.. --count").trim()
	default:
		log.Fatal("wrong source for version, should be gradle or node")
	}
}
