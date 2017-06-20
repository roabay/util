package aaptparse

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestParseApk(t *testing.T) {
	for {
		apk, err := Parse("138.apk")
		log.Println("error: ", err)
		if err != nil {
			log.Println("error: ", err)
		}
		if apk == nil {
			t.Fail()
		}
		fmt.Println("Apk = ", "\n",
			"Label:", apk.AppLabel, "\n",
			"FeaturesNotRequired:", apk.FeaturesNotRequired, "\n",
			"Gl Use: ", apk.GlUse, "\n",
			"LibsNotRequired", apk.LibsNotRequired, "\n",
			"PackageName: ", apk.PackageName, "\n",
			"VersionCode: ", apk.VersionCode, "\n",
			"VersionName:", apk.VersionName, "\n",
			"SdkVersion:", apk.SdkVersion, "\n",
			"TargetSdkVersion:", apk.TargetSdkVersion, "\n",
			"UsesPermissions:", apk.Permissions, "\n",
			"NativeCode:", apk.NativeCode, "\n",
			"FeaturesRequired:", apk.FeaturesRequired)
		fmt.Println("Waiting 2 secs...")
		time.Sleep(2 * time.Second)
	}
}
