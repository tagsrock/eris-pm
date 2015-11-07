package util

import (
	"regexp"

	"github.com/eris-ltd/eris-pm/definitions"
)

func BundleHttpPathCorrect(do *definitions.Do) {
	do.Chain = HttpPathCorrect(do.Chain, true)
	do.Signer = HttpPathCorrect(do.Signer, false)
	do.Compiler = HttpPathCorrect(do.Compiler, false)
}

func HttpPathCorrect(oldPath string, trailingSlash bool) string {
	var newPath string
	protoReg := regexp.MustCompile("https*://.*")
	trailer := regexp.MustCompile("/$")

	if !protoReg.MatchString(oldPath) {
		newPath = "http://" + oldPath
	} else {
		newPath = oldPath
	}

	if trailingSlash {
		if !trailer.MatchString(newPath) {
			newPath = newPath + "/"
		}
	}

	return newPath
}
