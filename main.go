package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	Epoch    int
	Major    int
	Minor    int
	Patch    int
	Hotfix   int
	Revision string
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

	outputHandle, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer outputHandle.Close()

	if strings.HasPrefix(refName, "debug_v") || strings.HasPrefix(refName, "prod_v") {
		version := "nil"
		parts := strings.Split(refName, "_v")
		if len(parts) > 1 {
			version = parts[1]
		} else {
			fmt.Println("NO SECOND PART FOUND AFTER _v, VERSION WILL BE 'nil'")
		}
		fmt.Printf("tag=%s\nversionnr=%s\nfilenameversion=%s\n", parts[0], version, version)

		fmt.Fprintf(outputHandle, `tag=%s\n`, parts[0])
		fmt.Fprintf(outputHandle, `versionnr=%s\n`, version)
		fmt.Fprintf(outputHandle, `filenameversion=%s\n`, version)
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
		fmt.Fprintln(outputHandle, "tag=debug")
		fmt.Println("tag=debug")
	default:
		fmt.Fprintln(outputHandle, "tag=prod")
		fmt.Println("tag=prod")
	}
	fmt.Fprintln(outputHandle, "versionnr="+version.String())
	fmt.Fprintln(outputHandle, "filenameversion="+version.FilenameVersion())
	fmt.Printf("versionnr=%s\nfilenameversion=%s\n", version.String(), version.FilenameVersion())
}

func (v Version) String() string {
	result := strings.Builder{}
	result.WriteString(fmt.Sprintf("%d:", v.Epoch))
	result.WriteString(fmt.Sprintf("%d.%d.%d.%d", v.Major, v.Minor, v.Patch, v.Hotfix))
	result.WriteString(fmt.Sprintf("-e%d", v.Epoch))
	if v.Revision != "" {
		result.WriteString(fmt.Sprintf("+%s", v.Revision))
	}
	return result.String()
}

func (v Version) FilenameVersion() string {
	result := strings.Builder{}
	result.WriteString(fmt.Sprintf("%d.%d.%d.%d", v.Major, v.Minor, v.Patch, v.Hotfix))
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
		v.Hotfix == 0 &&
		v.Revision == "" {
		return false
	}
	return true
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
		case 3:
			result.Hotfix, _ = strconv.Atoi(parts[i])
		}
	}

	result.Revision = matches[revisionIndex]
	return result
}
