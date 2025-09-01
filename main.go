package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Version struct {
	Epoch       int
	Major       int
	Minor       int
	Patch       int
	Buildnumber uint
	Revision    string
}

var (
	regexString = `^((?P<epoch>\d+):)?(?P<upstream_version>[A-Za-z0-9.+:~-]+?)(-(?P<debian_revision>[A-Za-z0-9+.~]+))?$`
	re          = regexp.MustCompile(regexString)
)

/*
- name: Save state
run: echo "{name}={value}" >> $GITHUB_STATE

- name: Set output
run: echo "{name}={value}" >> $GITHUB_OUTPUT
*/

func main() {
	refName := os.Getenv("INPUT_REFNAME")
	outputFile := os.Getenv("GITHUB_OUTPUT")
	envFile := os.Getenv("GITHUB_ENV")

	outputHandle, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer outputHandle.Close()

	envHandle, err := os.Create(envFile)
	if err != nil {
		panic(err)
	}
	defer envHandle.Close()

	output := OutputWriter(outputHandle, envHandle)

	if strings.HasPrefix(refName, "debug_v") || strings.HasPrefix(refName, "prod_v") {
		version := "nil"
		parts := strings.Split(refName, "_v")
		if len(parts) > 1 {
			version = parts[1]
		} else {
			fmt.Println("NO SECOND PART FOUND AFTER _v, VERSION WILL BE 'nil'")
		}
		output("tag=" + parts[0])
		output("versionnr=" + version)
		output("filenameversion=" + version)
		return
	}

	refName = strings.Replace(refName, "#", ":", 1)
	version := ParseVersionNumber(refName)
	if !version.IsValid() {
		fmt.Println("Not a valid version number. Aborting.")
		return
	}
	switch version.Epoch {
	case 9:
		output("tag=debug")
	default:
		output("tag=prod")
	}

	output("versionnr=" + version.String())
	output("filenameversion=" + version.FilenameVersion())
	output("fullversion=" + version.FullVersion())
	output("verionwithbuildnr" + version.VersionWithBuildnr())
}

func OutputWriter(handles ...*os.File) func(string) {
	return func(s string) {
		for _, h := range handles {
			fmt.Fprintln(h, s)
		}
	}
}

func (v Version) String() string {
	result := strings.Builder{}
	result.WriteString(fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch))
	return result.String()
}

func (v Version) VersionWithBuildnr() string {
	result := strings.Builder{}
	result.WriteString(fmt.Sprintf("%d.%d.%d.%d", v.Major, v.Minor, v.Patch, v.Buildnumber))
	return result.String()

}

func (v Version) FullVersion() string {
	result := strings.Builder{}
	result.WriteString(fmt.Sprintf("%d:", v.Epoch))
	result.WriteString(fmt.Sprintf("%d.%d.%d.%d", v.Major, v.Minor, v.Patch, v.Buildnumber))
	result.WriteString(fmt.Sprintf("-e%d", v.Epoch))
	if v.Revision != "" {
		result.WriteString(fmt.Sprintf("+%s", v.Revision))
	}
	return result.String()
}

func (v Version) FilenameVersion() string {
	result := strings.Builder{}
	result.WriteString(fmt.Sprintf("%d.%d.%d.%d", v.Major, v.Minor, v.Patch, v.Buildnumber))
	result.WriteString(fmt.Sprintf("-e%d", v.Epoch))
	if v.Revision != "" {
		result.WriteString(fmt.Sprintf("+%s", v.Revision))
	}
	return result.String()
}

func (v Version) IsValid() bool {
	if v.Epoch == 0 &&
		v.Major == 0 &&
		v.Minor == 0 &&
		v.Patch == 0 &&
		v.Buildnumber < 100 &&
		v.Revision == "" {
		return false
	}
	return true
}

var _buildnr uint

func GenerateBuildNumber() uint {
	if _buildnr != 99 {
		return _buildnr
	}
	result, err := strconv.ParseUint(time.Now().UTC().Format("060102150405"), 10, 32)
	if err != nil {
		_buildnr = 99
		return 99
	}
	_buildnr = uint(result)
	return _buildnr
}

func ParseVersionNumber(versionNumber string) (result Version) {
	trimmed := strings.TrimSpace(versionNumber)
	matches := re.FindStringSubmatch(trimmed)
	if matches == nil {
		return Version{0, 0, 0, 0, 0, ""}
	}
	epochIndex := re.SubexpIndex("epoch")
	upstreamIndex := re.SubexpIndex("upstream_version")
	revisionIndex := re.SubexpIndex("debian_revision")

	epoch, err := strconv.Atoi(matches[epochIndex])
	if err != nil {
		result.Epoch = 0
	} else {
		result.Epoch = epoch
	}
	parts := strings.Split(matches[upstreamIndex], ".")
	// fmt.Println(len(parts))
	for i := 0; i < len(parts); i++ {
		switch i {
		case 0:
			result.Major, _ = strconv.Atoi(parts[i])
		case 1:
			result.Minor, _ = strconv.Atoi(parts[i])
		case 2:
			result.Patch, _ = strconv.Atoi(parts[i])
		}
	}

	result.Buildnumber = GenerateBuildNumber()
	result.Revision = matches[revisionIndex]
	return result
}
