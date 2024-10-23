package version

import (
	"fmt"

	"github.com/blang/semver/v4"
)

func Do() {

	v, _ := semver.Make("0.0.1-alpha.preview+123.github")
	fmt.Printf("Major: %d\n", v.Major)
	fmt.Printf("Minor: %d\n", v.Minor)
	fmt.Printf("Patch: %d\n", v.Patch)
	fmt.Printf("Pre: %s\n", v.Pre)
	fmt.Printf("Build: %s\n", v.Build)

	// Prerelease versions array
	if len(v.Pre) > 0 {
		fmt.Println("Prerelease versions:")
		for i, pre := range v.Pre {
			fmt.Printf("%d: %q\n", i, pre)
		}
	}

	// Build meta data array
	if len(v.Build) > 0 {
		fmt.Println("Build meta data:")
		for i, build := range v.Build {
			fmt.Printf("%d: %q\n", i, build)
		}
	}

	v001, err := semver.Make("0.0.1")
	// Compare using helpers: v.GT(v2), v.LT, v.GTE, v.LTE
	fmt.Println(
		v001.GT(v),
		v.LT(v001),
		v.GTE(v),
		v.LTE(v),
		v.Validate(),
	)

	// Or use v.Compare(v2) for comparisons (-1, 0, 1):
	fmt.Println(
		v001.Compare(v) == 1,
		v.Compare(v001) == -1,
		v.Compare(v) == 0,
	)

	// Manipulate Version in place:
	v.Pre[0], err = semver.NewPRVersion("beta")
	if err != nil {
		fmt.Printf("Error parsing pre release version: %q", err)
	}

	fmt.Println("\nValidate versions:")
	v.Build[0] = "?"

	err = v.Validate()
	if err != nil {
		fmt.Printf("Validation failed: %s\n", err)
	}

}
