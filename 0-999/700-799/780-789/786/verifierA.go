package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesARaw = `100
5
4 1 2 3 4
3 4 2 1
6
2 3 2
1 5
4
3 3 1 2
1 3
2
1 1
1 1
5
3 2 3 4
4 3 1 4 2
7
4 6 1 4 2
2 6 3
7
1 2
5 2 6 5 3 4
2
1 1
1 1
4
3 2 1 3
3 1 2 3
2
1 1
1 1
3
1 1
1 2
5
1 1
2 2 1
2
1 1
1 1
7
5 4 5 3 2 6
6 6 3 1 2 4 5
6
3 2 5 1
3 1 2 4
3
2 2 1
1 1
7
2 1 5
6 5 6 1 4 3 2
6
5 1 4 5 2 3
1 5
2
1 1
1 1
3
1 1
2 1 2
2
1 1
1 1
5
2 1 3
4 1 3 4 2
3
2 2 1
1 1
2
1 1
1 1
2
1 1
1 1
7
4 5 6 3 4
3 4 3 2
6
1 4
1 3
7
1 5
3 2 6 4
4
3 2 3 1
3 2 3 1
7
1 1
5 2 3 6 1 5
7
4 4 5 6 1
4 6 5 4 3
7
1 2
4 1 3 2 4
6
4 5 1 4 2
3 3 4 1
5
2 1 3
2 1 2
7
4 3 1 2 5
6 1 5 6 4 3 2
4
2 3 1
1 2
5
1 1
3 4 1 2
3
2 1 2
2 1 2
2
1 1
1 1
6
1 5
4 4 5 2 1
3
2 2 1
1 1
4
2 2 3
3 1 2 3
2
1 1
1 1
4
2 3 1
2 1 2
7
2 1 3
2 5 1
4
2 2 3
2 1 3
6
4 4 3 2 1
4 1 4 2 5
4
2 3 1
1 3
6
4 1 5 4 3
2 1 4
2
1 1
1 1
5
2 4 1
3 2 4 1
5
2 3 1
1 1
6
4 2 1 5 3
3 2 1 4
6
2 1 2
4 4 3 5 1
2
1 1
1 1
6
4 3 4 2 1
5 5 2 1 4 3
4
1 3
1 2
6
2 4 3
4 1 5 3 4
2
1 1
1 1
2
1 1
1 1
7
4 3 5 4 6
5 1 5 6 3 4
7
4 2 3 4 6
4 5 6 2 1
7
1 6
2 3 5
6
3 3 1 2
3 3 4 2
5
1 2
2 2 4
7
3 1 6 4
4 2 4 1 3
7
2 3 4
1 5
5
4 4 1 3 2
2 1 4
6
5 2 3 1 5 4
4 4 1 5 2
6
3 1 3 5
4 3 1 5 4
2
1 1
1 1
5
3 1 4 3
4 3 1 4 2
2
1 1
1 1
5
4 1 4 2 3
1 4
3
1 2
2 2 1
5
3 3 1 2
1 1
4
3 1 2 3
1 1
3
2 1 2
1 1
3
1 1
2 2 1
6
2 2 5
4 1 4 2 3
3
1 1
2 2 1
4
2 3 2
1 3
2
1 1
1 1
3
1 1
2 1 2
3
2 2 1
1 1
6
1 4
1 1
3
2 2 1
1 2
4
3 2 1 3
2 3 2
7
5 1 4 5 6 2
6 3 4 2 5 6 1
2
1 1
1 1
6
5 2 4 5 1 3
4 2 5 4 3
6
2 3 4
1 4
3
2 1 2
2 1 2
4
1 3
2 3 1
5
4 1 3 2 4
4 3 4 1 2
4
3 1 3 2
2 3 2
5
1 3
3 3 4 1
6
1 1
1 3
4
3 2 1 3
1 1
5
3 3 4 1
2 1 2
`

func runExe(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "786A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func readTests() ([][]byte, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesARaw))
	if !scan.Scan() {
		return nil, fmt.Errorf("empty test file")
	}
	t, err := strconv.Atoi(strings.TrimSpace(scan.Text()))
	if err != nil {
		return nil, err
	}
	tests := make([][]byte, 0, t)
	for i := 0; i < t; i++ {
		var lines []string
		for j := 0; j < 3; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("not enough lines for test %d", i+1)
			}
			lines = append(lines, scan.Text())
		}
		tests = append(tests, []byte(strings.Join(lines, "\n")+"\n"))
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests, err := readTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for i, input := range tests {
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, string(input), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
