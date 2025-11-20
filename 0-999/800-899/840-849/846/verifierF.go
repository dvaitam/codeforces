package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

func runTests(dir, binary string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.in"))
	if err != nil {
		return err
	}
	sort.Strings(files)
	for _, inFile := range files {
		outFile := inFile[:len(inFile)-3] + ".out"
		input, err := os.ReadFile(inFile)
		if err != nil {
			return err
		}
		expectedBytes, err := os.ReadFile(outFile)
		if err != nil {
			return err
		}
        expectedStr := strings.TrimSpace(string(expectedBytes))

		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewReader(input)
		outBytes, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: runtime error: %v\noutput: %s", filepath.Base(inFile), err, string(outBytes))
		}
        gotStr := strings.TrimSpace(string(outBytes))
        
        gotFloat, errGot := strconv.ParseFloat(gotStr, 64)
        expectedFloat, errExpected := strconv.ParseFloat(expectedStr, 64)

        if errGot != nil || errExpected != nil {
            // If parsing to float fails, fall back to string comparison
            if gotStr != expectedStr {
                return fmt.Errorf("%s: expected\n%s\nbut got\n%s", filepath.Base(inFile), expectedStr, gotStr)
            }
        } else {
            // Compare floats with a tolerance
            epsilon := 1e-6
            if math.Abs(gotFloat - expectedFloat) > epsilon {
                 return fmt.Errorf("%s: expected\n%s\nbut got\n%s", filepath.Base(inFile), expectedStr, gotStr)
            }
        }
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	base := filepath.Dir(file)
	testDir := filepath.Join(base, "tests", "F")
	if err := runTests(testDir, binary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}