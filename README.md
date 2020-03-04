# getversion
small version utility for CI/CD sysyems like Jenkins written in Go.
Designed to fetch base project version from gradle.properties, package.json or git tags if project follows GitFlow branching strategy.

Mimics [GitVersion functionality](https://gitversion.net/docs/), but w/o dependency on .NET libraries

## Usage

### Jenkins

in Jenkinsfile you can call it like 

`def dockerTag = sh(script: "getversion -source=git-tag -build-id=${BUILD_ID} -showvariable=AssemblySemVer", returnStdout: true)`

### Git tags
Tries to find most recent tag and determines a number of commits since then.

Imagine you have latest reachable tag 0.2.0 and you're on develop branch

`getversion -source=git-tag`  returns
```json
{
  "Major": 0,
  "Minor": 3,
  "Patch": 0,
  "CommitsSinceVersionSource": 14,
  "CommitsSinceVersionSourcePadded": "0014",
  "PreReleaseTag": "alpha.14",
  "PreReleaseTagWithDash": "-alpha.14",
  "PreReleaseLabel": "alpha",
  "PreReleaseNumber": 0,
  "SemVer": "0.3.0-alpha.14",
  "BranchName": "develop",
  "AssemblySemVer": "0.3.0-alpha.14.0",
  "BuildMetaData": "0"
}

```

### Gradle 
in `gradle.properties` you need to have

`version=0.1.2-SNAPSHOT`

`getversion -source=gradle` returns

```json
{
  "Major": 0,
  "Minor": 2,
  "Patch": 2,
  "CommitsSinceVersionSource": 0,
  "CommitsSinceVersionSourcePadded": "",
  "PreReleaseTag": "alpha.0",
  "PreReleaseTagWithDash": "-alpha.0",
  "PreReleaseLabel": "alpha",
  "PreReleaseNumber": 0,
  "SemVer": "0.2.2-alpha.0",
  "BranchName": "develop",
  "AssemblySemVer": "0.2.2-alpha.0.0",
  "BuildMetaData": "0"
}

```

### Node JS

in `package.json` you need to have 

`"version": "1.0.0"`

`getversion -source=node` returns 

```json
{
  "Major": 1,
  "Minor": 1,
  "Patch": 0,
  "CommitsSinceVersionSource": 0,
  "CommitsSinceVersionSourcePadded": "",
  "PreReleaseTag": "alpha.0",
  "PreReleaseTagWithDash": "-alpha.0",
  "PreReleaseLabel": "alpha",
  "PreReleaseNumber": 0,
  "SemVer": "1.1.0-alpha.0",
  "BranchName": "develop",
  "AssemblySemVer": "1.1.0-alpha.0.0",
  "BuildMetaData": "0"
}

```


