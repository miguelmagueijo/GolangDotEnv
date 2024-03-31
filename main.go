package golangDotEnv

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

var removeQuotes = true
var removeApostrophe = true
var injectToEnv = false

func SetRemoveQuotes(ignore bool) {
	removeQuotes = ignore
}

func SetRemoveApostrophe(ignore bool) {
	removeApostrophe = ignore
}

func SetInjectToEnv(inject bool) {
	injectToEnv = inject
}

func panicIfError(e error) {
	if e != nil {
		panic(e)
	}
}

func logText(text string) string {
	return "[GolangDotEnv#WARNING] " + text
}

func Load() map[string]string {
	envData := make(map[string]string)

	fileData, err := os.Open("./.env")

	panicIfError(err)

	scanner := bufio.NewScanner(fileData)

	scanner.Split(bufio.ScanLines)

	addedData := false
	lineNumber := 0
	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		lineNumber++

		if bytes.HasPrefix(lineBytes, []byte("#")) || len(lineBytes) == 0 {
			continue
		}

		arr := bytes.SplitN(lineBytes, []byte("="), 2)

		if len(arr) != 2 {
			fmt.Printf(logText("Bad file line on line %d, ignoring: \"%s\"\n"), lineNumber, scanner.Text())
			continue
		}

		keyBytes, valueBytes := arr[0], arr[1]

		if len(keyBytes) == 0 || len(valueBytes) == 0 {
			fmt.Printf(logText("Bad variable composition on line %d, key: \"%s\", value: \"%s\"\n"), lineNumber, string(keyBytes), string(valueBytes))
			continue
		}

		bytes.ReplaceAll(keyBytes, []byte(" "), []byte("_"))

		if valueBytes[0] == '"' && removeQuotes {
			valueBytes = bytes.Trim(valueBytes, "\"")
		} else if valueBytes[0] == '\'' && removeApostrophe {
			valueBytes = bytes.Trim(valueBytes, "'")
		}

		if injectToEnv {
			err = os.Setenv(string(keyBytes), string(valueBytes))
			panicIfError(err)
			addedData = true
		} else {
			envData[string(keyBytes)] = string(valueBytes)
			addedData = true
		}
	}

	if !addedData {
		fmt.Println(logText("No data was loaded. Please check printed logs for warnings."))
		return nil
	}

	if injectToEnv {
		fmt.Println(logText("Injecting to environment. Returned map is nil."))
		return nil
	}

	return envData
}

func LoadWithPath(path string) string {
	panic("LoadEnvWithPath not yet implemented")
}
