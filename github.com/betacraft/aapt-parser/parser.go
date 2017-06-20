package aaptparse

import (
	"log"
	"os/exec"
	"strings"
)

var (
	permissions         = "uses-permission"
	featuresNotRequired = "uses-feature-not-required"
	featuresRequired    = "uses-feature"
	packageKey          = "package"
	usesGl              = "uses-gl-es"
	appLabel            = "application-label:"
	libsNotRequired     = "uses-library-not-required"
	targetSdkVersion    = "targetSdkVersion"
	sdkVersion          = "sdkVersion"
	nativeCode          = "native-code"
)

type Apk struct {
	Permissions         []string
	FeaturesNotRequired []string
	FeaturesRequired    []string
	LibsNotRequired     []string
	LibsRequired        []string // TODO
	AppLabel            string
	PackageName         string
	VersionCode         int
	VersionName         string
	TargetSdkVersion    string
	SdkVersion          string
	GlUse               string
	NativeCode          []string
}

func Parse(apkPath string) (*Apk, error) {
	apk := new(Apk)
	op, err := exec.Command("aapt", "dump", "badging", apkPath).Output()
	if err == exec.ErrNotFound {
		log.Println("Install aapt first")
		return nil, err
	}

	if err != nil {
		log.Println("Check if path is correct, use absolute path")
		return nil, err
	}

	data := strings.TrimSpace(string(op))

	lines := strings.Split(data, "\n") // get all lines
	for _, line := range lines {
		if strings.Contains(line, permissions) {
			getPermissionInfo(line, apk)
			continue
		}
		if strings.Contains(line, featuresNotRequired) {
			getFeatureNotRequired(line, apk)
			continue
		}
		if strings.Contains(line, packageKey) {
			getPackageInfo(line, apk)
			continue
		}
		if strings.Contains(line, usesGl) {
			getGlInfo(line, apk)
			continue
		}
		if strings.Contains(line, appLabel) {
			getAppLabel(line, apk)
			continue
		}
		if strings.Contains(line, libsNotRequired) {
			getLibsNotRequired(line, apk)
			continue
		}
		if strings.Contains(line, targetSdkVersion) {
			getTargetSdk(line, apk)
			continue
		}
		if strings.Contains(line, sdkVersion) {
			getSdk(line, apk)
			continue
		}
		if strings.Contains(line, nativeCode) {
			getNativeCode(line, apk)
			continue
		}
		if strings.Contains(line, featuresRequired) {
			getFeatureRequired(line, apk)
			continue
		}
	}

	return apk, nil
}
