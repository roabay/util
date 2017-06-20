
### First install aapt (on Debian or derivatives install using apt-get)

### Install using go get
```
$ go get github.com/betacraft/aapt-parser
```

#### How to use
```go
func main() {
   apk := Parse("to_be_parsed.apk")
   if apk == nil {
      // some error, install aapt or give correct path
      return
	}
  // Use all data in Apk(see below Apk struct)
}

```

#### Apk Struct

```go
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
	NativeCode          string
}

```
