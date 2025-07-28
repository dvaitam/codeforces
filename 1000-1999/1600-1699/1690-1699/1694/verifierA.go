package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runBinary(bin string, input []byte) ([]byte, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return append(out.Bytes(), errb.Bytes()...), fmt.Errorf("%v", err)
	}
	return out.Bytes(), nil
}

func minCreepiness(a, b int) int {
	if a == 0 || b == 0 {
		if a > b {
			return a
		}
		return b
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	if diff == 0 {
		return 1
	}
	return diff
}

func creepiness(s string) (int, bool) {
	zeros, ones, maxDiff := 0, 0, 0
	for i := 0; i < len(s); i++ {
		if s[i] == '0' {
			zeros++
		} else if s[i] == '1' {
			ones++
		} else {
			return 0, false
		}
		diff := zeros - ones
		if diff < 0 {
			diff = -diff
		}
		if diff > maxDiff {
			maxDiff = diff
		}
	}
	return maxDiff, true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	data, err := ioutil.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read testcasesA.txt: %v\n", err)
		os.Exit(1)
	}

	outBytes, err := runBinary(bin, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "binary execution failed: %v\n%s", err, outBytes)
		os.Exit(1)
	}
	outBuf := bytes.NewBuffer(outBytes)

	inScanner := bufio.NewScanner(bytes.NewReader(data))
	if !inScanner.Scan() {
		fmt.Fprintln(os.Stderr, "testcase file empty")
		os.Exit(1)
	}
	t, err := strconv.Atoi(strings.TrimSpace(inScanner.Text()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "bad test count")
		os.Exit(1)
	}

	outScanner := bufio.NewScanner(outBuf)
	outScanner.Split(bufio.ScanWords)
	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		if !inScanner.Scan() {
			fmt.Fprintf(os.Stderr, "not enough cases in file at case %d\n", caseIdx)
			os.Exit(1)
		}
		parts := strings.Fields(inScanner.Text())
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "bad line in testcase %d\n", caseIdx)
			os.Exit(1)
		}
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])

		if !outScanner.Scan() {
			fmt.Fprintf(os.Stderr, "not enough output for case %d\n", caseIdx)
			os.Exit(1)
		}
		ans := outScanner.Text()
		if len(ans) != a+b {
			fmt.Fprintf(os.Stderr, "case %d wrong length\n", caseIdx)
			os.Exit(1)
		}
		zeros := strings.Count(ans, "0")
		ones := len(ans) - zeros
		if zeros != a || ones != b {
			fmt.Fprintf(os.Stderr, "case %d wrong number of zeros/ones\n", caseIdx)
			os.Exit(1)
		}
		c, ok := creepiness(ans)
		if !ok {
			fmt.Fprintf(os.Stderr, "case %d output has invalid characters\n", caseIdx)
			os.Exit(1)
		}
		exp := minCreepiness(a, b)
		if c != exp {
			fmt.Fprintf(os.Stderr, "case %d wrong creepiness got %d expected %d\n", caseIdx, c, exp)
			os.Exit(1)
		}
	}
	if outScanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
