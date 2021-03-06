package monit

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/monkeyherder/salus/checks"
)

type MonitFile struct {
	Checks []checks.Check
}

func ReadMonitFile(filepath string) (MonitFile, error) {
	bytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		// Do something
	}

	lines := strings.Split(string(bytes), "\n")

	checks := []checks.Check{}

	i := 0
	for _, line := range lines {
		processMatch, err := regexp.Match("check process", []byte(line))
		fileMatch, err := regexp.Match("check file", []byte(line))

		if err != nil {
			// Do something
		}

		if processMatch {
			check := createProcessCheck(lines, i)
			checks = append(checks, check)
		} else if fileMatch {
			check := createFileCheck(lines, i)
			checks = append(checks, check)
		}

		i++
	}

	monitFile := MonitFile{checks}

	return monitFile, nil
}

func createProcessCheck(lines []string, startingIndex int) checks.ProcessCheck {
	name, lines := captureWithRegex(lines, `check process ([\w"\.]+)`, true)

	totalMemChecks, lines := parseAllTotalMemChecks(lines)

	pidfile, lines := captureWithRegex(lines, `with pidfile ([\w"/\.]+)`, true)
	startProgram, lines := captureWithRegex(lines, `start program (.*)$`, true)
	stopProgram, lines := captureWithRegex(lines, `stop program (.*)$`, true)
	group, lines := captureWithRegex(lines, `group (\w+)`, true)
	dependsOn, lines := captureWithRegex(lines, `depends on (\w+)`, true)

	failedSocket, lines := parseFailedUnixSocket(lines)
	failedHost, lines := parseFailedHost(lines)

	check := checks.ProcessCheck{
		Name:           name,
		Pidfile:        pidfile,
		StartProgram:   startProgram,
		StopProgram:    stopProgram,
		FailedSocket:   failedSocket,
		FailedHost:     failedHost,
		TotalMemChecks: totalMemChecks,
		Group:          group,
		DependsOn:      dependsOn,
	}

	return check
}

func createFileCheck(lines []string, startingIndex int) checks.FileCheck {
	name, lines := captureWithRegex(lines, `check file ([\w"\.]+)`, true)

	path, lines := captureWithRegex(lines, `with path ([\w"/\.]+)`, true)
	ifChanged, lines := captureWithRegex(lines, `if changed timestamp then exec (.*)$`, true)
	group, lines := captureWithRegex(lines, `group (\w+)`, true)
	dependsOn, lines := captureWithRegex(lines, `depends on (\w+)`, true)

	check := checks.FileCheck{
		Name:      name,
		Path:      path,
		IfChanged: ifChanged,
		Group:     group,
		DependsOn: dependsOn,
	}

	return check
}

func parseFailedUnixSocket(lines []string) (checks.FailedSocket, []string) {
	values, lines := parseGroupBlock(
		lines,
		"socketFile",
		map[string]string{
			"socketFile": `if failed unixsocket (["/\w\.]+)`,
			"timeout":    `with timeout (\d+) seconds`,
			"numCycles":  `for (\d+) cycles`,
			"action":     `then ([a-z]+)`,
		},
	)

	socketFile := values["socketFile"]
	timeout := values["timeout"]
	numCycles := values["numCycles"]
	action := values["action"]

	timeoutInt, err := strconv.Atoi(timeout)
	if err != nil {
		// Do something
	}

	numCyclesInt, err := strconv.Atoi(numCycles)
	if err != nil {
		// Do something
	}

	fs := checks.FailedSocket{
		SocketFile: socketFile,
		Timeout:    timeoutInt,
		NumCycles:  numCyclesInt,
		Action:     action,
	}

	return fs, lines
}

func parseFailedHost(lines []string) (checks.FailedHost, []string) {
	values, lines := parseGroupBlock(
		lines,
		"host",
		map[string]string{
			"host":      `if failed host ([\w\.]+)`,
			"port":      `port (\d+)`,
			"protocol":  `protocol (\w+)`,
			"timeout":   `with timeout (\d+) seconds`,
			"numCycles": `for (\d+) cycles`,
			"action":    `then ([a-z]+)`,
		},
	)

	host := values["host"]
	port := values["port"]
	protocol := values["protocol"]
	timeout := values["timeout"]
	numCycles := values["numCycles"]
	action := values["action"]

	timeoutInt, err := strconv.Atoi(timeout)
	if err != nil {
		// Do something
	}

	numCyclesInt, err := strconv.Atoi(numCycles)
	if err != nil {
		// Do something
	}

	fh := checks.FailedHost{
		Host:      host,
		Port:      port,
		Protocol:  protocol,
		Timeout:   timeoutInt,
		NumCycles: numCyclesInt,
		Action:    action,
	}

	return fh, lines
}

func parseTotalMem(lines []string) (checks.MemUsage, []string) {
	totalMemLineEnding, lines := captureWithRegex(lines, `if totalmem > (.*$)`, true)

	tmpLines := []string{totalMemLineEnding}
	memLimit, _ := captureWithRegex(tmpLines, `(\d+) Mb`, false)
	numCycles, _ := captureWithRegex(tmpLines, `for (\d+) cycles`, false)
	action, _ := captureWithRegex(tmpLines, `then (\w+)`, false)

	memLimitInt, err := strconv.Atoi(memLimit)
	if err != nil {
		// Do something
	}

	numCyclesInt, err := strconv.Atoi(numCycles)
	if err != nil {
		// Do something
	}

	mu := checks.MemUsage{
		MemLimit:  memLimitInt,
		NumCycles: numCyclesInt,
		Action:    action,
	}

	return mu, lines
}

func parseAllTotalMemChecks(lines []string) ([]checks.MemUsage, []string) {
	var memChecks []checks.MemUsage
	var memCheck checks.MemUsage

	for _, line := range lines {
		newProcessCheck, err := regexp.Match("check process", []byte(line))
		if err != nil {
			// Do something
		}

		newFileCheck, err := regexp.Match("check file", []byte(line))
		if err != nil {
			// Do something
		}

		if newProcessCheck || newFileCheck {
			break
		}

		memCheckMatch, err := regexp.Match("if totalmem", []byte(line))
		if err != nil {
			// Do something
		}

		if memCheckMatch {
			memCheck, lines = parseTotalMem(lines)
			memChecks = append(memChecks, memCheck)
		}
	}

	return memChecks, lines
}

func parseGroupBlock(lines []string, keyRegex string, regexes map[string]string) (map[string]string, []string) {
	var startingIndex int
	//	var endingIndex int
	var newLines []string
	values := map[string]string{}

	startingRegex := regexp.MustCompile(regexes[keyRegex])

	for i, line := range lines {
		newProcessCheck, err := regexp.Match("check process", []byte(line))
		if err != nil {
			// Do something
		}

		newFileCheck, err := regexp.Match("check file", []byte(line))
		if err != nil {
			// Do something
		}

		if newProcessCheck || newFileCheck {
			break
		}

		match := startingRegex.Match([]byte(line))

		if match {
			startingIndex = i

			newLines = append([]string{}, lines[i:]...)

			for key, regex := range regexes {
				values[key], lines = captureWithRegex(newLines, regex, false)
			}

			//			for j, newLine := range newLines {
			//				thenMatch, err := regexp.Match("then ", []byte(newLine))
			//
			//				if err != nil {
			//					// Do something
			//				}
			//
			//				if thenMatch {
			//					endingIndex = i + j
			//				}
			//			}
		}
	}

	if len(values) > 0 {
		removeElementsFromSlice(lines, startingIndex, startingIndex+1)
	}

	return values, lines
}

func captureWithRegex(lines []string, reg string, removeLine bool) (string, []string) {
	var myString string

	for i, line := range lines {
		regex := regexp.MustCompile(reg)
		values := regex.FindStringSubmatch(line)

		newProcessCheck, err := regexp.Match("check process", []byte(line))
		if err != nil {
			// Do something
		}

		newFileCheck, err := regexp.Match("check file", []byte(line))
		if err != nil {
			// Do something
		}

		if len(values) > 1 {
			myString = strings.TrimSpace(values[1])

			if removeLine {
				lines = removeElementsFromSlice(lines, i, i+1)
			}

			break
		} else if newProcessCheck || newFileCheck {
			break
		}
	}

	stripReg := regexp.MustCompile(`"([^"]*)"`)
	return stripReg.ReplaceAllString(myString, "${1}"), lines
}

func removeElementsFromSlice(slice []string, startingIndex int, endingIndex int) []string {
	return append(slice[:startingIndex], slice[endingIndex:]...)
}
