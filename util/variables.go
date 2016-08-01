package util

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/eris-ltd/eris-pm/definitions"

	log "github.com/eris-ltd/eris-logger"
)

func PreProcess(toProcess string, do *definitions.Do) (string, error) {
	// $block.... $account.... etc. should be caught. hell$$o should not
	// :$libAddr needs to be caught
	catchEr := regexp.MustCompile("(^|\\s|:)\\$([a-zA-Z0-9_.]+)")
	// If there's a match then run through the replacement process
	if catchEr.MatchString(toProcess) {
		log.WithField("match", toProcess).Debug("Replacement Match Found")

		// find what we need to catch.
		processedString := toProcess

		for _, jobMatch := range catchEr.FindAllStringSubmatch(toProcess, -1) {
			jobName := jobMatch[2]
			varName := "$" + jobName
			var innerVarName string
			var wantsInnerValues bool = false
			/*
				log.WithFields(log.Fields{
				 	"var": varName,
				 	"job": jobName,
				}).Debugf("Correcting match %d", i+1)
			*/
			// first parse the reserved words.
			if strings.Contains(jobName, "block") {
				block, err := replaceBlockVariable(toProcess, do)
				if err != nil {
					log.WithField("err", err).Error("Error replacing block variable.")
					return "", err
				}
				/*log.WithFields(log.Fields{
					"var": toProcess,
					"res": block,
				}).Debug("Fixing Variables =>")*/
				processedString = strings.Replace(processedString, toProcess, block, 1)
			}

			if strings.Contains(jobName, ".") { //for functions with multiple returns
				wantsInnerValues = true
				var splitStr = strings.Split(jobName, ".")
				jobName = splitStr[0]
				innerVarName = splitStr[1]
			}

			// second we loop through the jobNames to do a result replace
			for _, job := range do.Package.Jobs {
				if string(jobName) == job.JobName {
					if wantsInnerValues {
						for _, innerVal := range job.JobVars {
							if innerVal.Name == innerVarName { //find the value we want from the bunch
								processedString = strings.Replace(processedString, varName, innerVal.Value, 1)
								log.WithFields(log.Fields{
									"job":     string(jobName),
									"varName": innerVarName,
									"result":  innerVal.Value,
								}).Debug("Fixing Inner Vars =>")
							}
						}
					} else {
						log.WithFields(log.Fields{
							"var": string(jobName),
							"res": job.JobResult,
						}).Debug("Fixing Variables =>")
						processedString = strings.Replace(processedString, varName, job.JobResult, 1)
					}

				}
			}
		}
		return processedString, nil
	}

	// if no matches, return original
	return toProcess, nil
}

func replaceBlockVariable(toReplace string, do *definitions.Do) (string, error) {
	log.WithFields(log.Fields{
		"chain": do.Chain,
		"var":   toReplace,
	}).Debug("Correcting $block variable")
	block, err := ChainStatus("latest_block_height", do)
	log.WithField("=>", block).Debug("Current height is")
	if err != nil {
		return "", err
	}

	if toReplace == "$block" {
		log.WithField("=>", block).Debug("Replacement (=)")
		return block, nil
	}

	catchEr := regexp.MustCompile("\\$block\\+(\\d*)")
	if catchEr.MatchString(toReplace) {
		height := catchEr.FindStringSubmatch(toReplace)[1]
		h1, err := strconv.Atoi(height)
		if err != nil {
			return "", err
		}
		h2, err := strconv.Atoi(block)
		if err != nil {
			return "", err
		}
		height = strconv.Itoa(h1 + h2)
		log.WithField("=>", height).Debug("Replacement (+)")
		return height, nil
	}

	catchEr = regexp.MustCompile("\\$block\\-(\\d*)")
	if catchEr.MatchString(toReplace) {
		height := catchEr.FindStringSubmatch(toReplace)[1]
		h1, err := strconv.Atoi(height)
		if err != nil {
			return "", err
		}
		h2, err := strconv.Atoi(block)
		if err != nil {
			return "", err
		}
		height = strconv.Itoa(h1 - h2)
		log.WithField("=>", height).Debug("Replacement (-)")
		return height, nil
	}

	log.WithField("=>", toReplace).Debug("Replacement (unknown)")
	return toReplace, nil
}

func PreProcessLibs(libs string, do *definitions.Do) (string, error) {
	libraries, _ := PreProcess(libs, do)
	if libraries != "" {
		pairs := strings.Split(libraries, ",")
		for _, pair := range pairs {
			libAndAddr := strings.Split(pair, ":")
			libAndAddr[1] = strings.ToLower(libAndAddr[1])
			pair = strings.Join(libAndAddr, ":")
		}
		libraries = strings.Join(pairs, " ")
	}
	log.WithField("=>", libraries).Debug("Library String")
	return libraries, nil
}

func GetReturnValue(vars []*definitions.Variable) string {
	var result []string
	//log.WithField("=>", vars).Debug("GetReturnValue")
	for _, value := range vars {
		result = append(result, value.Value)
	}
	if len(vars) > 1 {
		return "(" + strings.Join(result, ", ") + ")"
	} else {
		return result[0]
	}
}
