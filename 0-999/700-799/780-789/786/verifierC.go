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

const testcasesCRaw = `100
1
1
2
2 1
5
3 5 2 5 1
3
2 3 2
6
5 4 5 3 1 1
6
4 3 4 4 5 2
3
1 1 1
3
2 1 1
6
5 6 5 2 4 4
6
5 3 3 4 2 4
8
4 8 5 8 6 8 8 6
8
8 4 6 3 5 8 5 5
7
3 6 2 4 5 3 6
2
2 1
4
1 1 1 3
4
1 2 3 2
4
1 4 1 1
6
3 2 2 6 1 1
2
1 1
1
1
6
3 2 2 6 2 5
1
1
1
1
3
1 1 2
2
2 2
8
1 5 8 1 5 7 3 8
4
1 3 1 1
8
3 7 8 6 3 6 5 5
7
6 1 6 5 2 6 1
5
1 2 2 2 1
8
4 1 4 4 8 2 5 2
4
3 3 4 3
1
1
1
1
7
2 1 5 6 1 2 1
2
1 1
4
1 2 1 4
8
5 7 4 4 7 7 1 1
7
5 5 2 1 6 7 4
6
1 5 1 5 3 3
6
3 1 6 4 1 1
5
2 1 4 1 4
8
8 4 2 1 5 1 6 5
2
1 2
4
1 3 4 4
3
2 2 1
5
1 1 1 5 3
7
2 6 1 1 5 6 4
1
1
5
3 4 2 3 3
8
8 7 8 5 7 4 3 8
5
5 4 1 5 5
2
1 2
3
3 1 2
2
1 1
3
2 2 1
6
4 2 5 3 1 2
7
1 3 5 2 6 5 3
3
1 2 3
4
4 3 2 4
8
1 7 3 7 1 8 5 7
5
4 4 3 5 3
2
1 1
7
7 6 4 6 1 3 4
8
3 2 1 7 4 7 4 2
7
7 5 7 7 2 3 6
4
4 2 1 4
8
5 3 8 4 2 6 1 8
2
2 2
8
5 8 1 2 6 3 7 5
3
1 1 2
7
4 6 3 2 1 3 5
8
1 6 1 7 8 4 5 8
3
2 3 3
5
1 3 3 3 3
5
4 5 1 5 2
7
5 5 7 7 2 7 5
2
2 1
4
4 2 3 1
2
1 2
6
2 3 3 1 3 4
6
2 4 4 3 4 2
8
4 5 6 3 2 4 8 4
6
2 3 2 2 2 3
7
4 7 7 6 3 3 6
6
6 4 6 6 6 6
5
5 5 1 3 3
7
4 2 3 3 4 4 1
3
2 2 1
1
1
6
2 3 1 6 6 4
1
1
4
4 3 4 2
6
3 2 4 1 2 2
6
3 2 4 3 3 1
6
2 2 6 2 6 5
1
1
6
6 5 1 2 2 1
7
4 7 3 2 3 5 5
2
2 2
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
	ref := "./refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "786C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func readTests() ([][]byte, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesCRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("empty test file")
	}
	t, _ := strconv.Atoi(scan.Text())
	tests := make([][]byte, 0, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("missing n for test %d", i+1)
		}
		n, _ := strconv.Atoi(scan.Text())
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("bad array value test %d", i+1)
			}
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(scan.Text())
		}
		sb.WriteByte('\n')
		tests = append(tests, []byte(sb.String()))
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
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
