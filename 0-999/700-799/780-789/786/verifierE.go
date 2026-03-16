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

const testcasesERaw = `100
3 3
2 1
3 2
1 2
3 2
1 3
6 3
2 1
3 1
4 1
5 3
6 2
3 6
6 3
1 5
4 4
2 1
3 1
4 2
3 1
3 2
4 3
4 2
3 3
2 1
3 1
2 3
2 3
3 2
4 2
2 1
3 2
4 3
4 3
2 3
5 3
2 1
3 1
4 3
5 3
5 2
1 3
2 3
5 1
2 1
3 1
4 3
5 3
3 2
5 1
2 1
3 1
4 2
5 3
5 2
4 1
2 1
3 2
4 2
3 1
5 2
2 1
3 1
4 1
5 2
2 5
1 5
5 2
2 1
3 1
4 2
5 2
2 4
4 1
3 4
2 1
3 2
2 1
1 2
2 1
3 1
3 4
2 1
3 2
2 3
1 3
2 3
3 1
5 4
2 1
3 1
4 1
5 3
4 3
4 5
2 3
3 4
5 1
2 1
3 1
4 2
5 2
1 2
5 4
2 1
3 1
4 1
5 1
3 1
4 2
5 1
2 5
4 4
2 1
3 1
4 3
1 4
3 2
1 3
1 3
4 1
2 1
3 2
4 1
4 2
5 4
2 1
3 1
4 2
5 4
4 1
2 3
4 1
5 2
4 2
2 1
3 1
4 2
3 1
2 1
4 1
2 1
3 2
4 1
2 4
4 3
2 1
3 1
4 1
3 4
3 4
4 2
4 3
2 1
3 2
4 3
2 3
3 2
3 1
2 4
2 1
2 1
2 1
1 2
1 2
6 1
2 1
3 2
4 3
5 3
6 4
2 4
3 3
2 1
3 2
3 2
2 3
2 3
2 3
2 1
2 1
1 2
2 1
4 4
2 1
3 2
4 1
2 1
3 2
1 2
4 3
2 2
2 1
1 2
1 2
2 3
2 1
1 2
1 2
1 2
3 2
2 1
3 1
1 2
3 2
6 3
2 1
3 1
4 1
5 1
6 1
2 1
2 1
5 1
3 2
2 1
3 1
3 2
2 1
6 2
2 1
3 2
4 3
5 1
6 2
5 6
5 3
5 3
2 1
3 2
4 3
5 1
3 5
4 3
1 2
2 4
2 1
1 2
1 2
1 2
1 2
3 2
2 1
3 2
2 3
3 2
2 2
2 1
1 2
2 1
5 3
2 1
3 1
4 2
5 3
4 3
5 2
3 5
2 1
2 1
2 1
6 3
2 1
3 1
4 1
5 2
6 4
6 3
1 4
5 2
3 1
2 1
3 1
2 3
6 2
2 1
3 1
4 3
5 2
6 2
1 5
3 6
4 3
2 1
3 1
4 1
1 3
4 1
2 4
5 4
2 1
3 2
4 2
5 1
2 3
2 3
2 3
5 2
6 3
2 1
3 1
4 1
5 2
6 2
1 2
3 6
6 3
6 4
2 1
3 2
4 1
5 4
6 3
4 6
2 4
1 6
6 4
3 1
2 1
3 1
2 1
6 1
2 1
3 1
4 2
5 1
6 1
4 2
6 2
2 1
3 1
4 2
5 4
6 3
6 3
2 6
4 3
2 1
3 2
4 1
2 3
2 4
4 1
2 2
2 1
1 2
1 2
6 1
2 1
3 2
4 3
5 3
6 1
3 4
2 1
2 1
2 1
3 4
2 1
3 1
1 3
1 3
3 1
1 2
6 4
2 1
3 2
4 1
5 2
6 5
6 1
1 2
4 2
1 5
4 3
2 1
3 1
4 3
1 2
1 3
4 3
6 3
2 1
3 2
4 2
5 4
6 4
5 1
5 2
3 5
6 4
2 1
3 2
4 2
5 3
6 2
4 5
3 6
2 3
1 2
4 2
2 1
3 1
4 1
1 2
2 3
5 4
2 1
3 2
4 2
5 1
1 2
1 4
1 4
2 4
2 2
2 1
2 1
1 2
4 1
2 1
3 1
4 2
4 1
5 2
2 1
3 2
4 1
5 1
1 2
2 4
5 3
2 1
3 2
4 2
5 1
4 5
1 4
4 3
3 2
2 1
3 1
3 1
3 1
6 2
2 1
3 2
4 2
5 1
6 1
4 3
4 6
4 2
2 1
3 2
4 2
4 2
4 3
6 1
2 1
3 1
4 3
5 1
6 5
2 1
2 2
2 1
2 1
2 1
2 1
2 1
1 2
6 1
2 1
3 1
4 1
5 2
6 1
4 3
4 3
2 1
3 1
4 2
1 3
3 1
1 3
2 4
2 1
1 2
2 1
1 2
2 1
3 2
2 1
3 2
3 2
2 3
5 2
2 1
3 1
4 2
5 3
3 1
2 5
3 2
2 1
3 1
1 3
1 3
2 1
2 1
1 2
2 1
2 1
1 2
6 4
2 1
3 2
4 1
5 3
6 2
4 5
6 3
6 2
4 2
3 2
2 1
3 1
1 3
3 2
6 3
2 1
3 1
4 3
5 4
6 4
3 6
1 4
4 3
4 2
2 1
3 2
4 2
4 1
3 1
6 3
2 1
3 2
4 1
5 4
6 2
3 2
6 2
4 5
3 2
2 1
3 2
2 1
3 1
6 1
2 1
3 1
4 2
5 3
6 2
1 6
5 4
2 1
3 2
4 2
5 1
2 4
4 3
2 4
5 3
4 1
2 1
3 2
4 3
2 1
5 1
2 1
3 1
4 2
5 3
2 1
2 2
2 1
2 1
2 1
2 3
2 1
2 1
1 2
2 1
3 1
2 1
3 2
3 1
5 4
2 1
3 2
4 1
5 4
2 3
3 1
4 5
5 2
4 4
2 1
3 2
4 2
3 4
3 2
2 3
1 3
5 4
2 1
3 2
4 1
5 4
4 3
3 2
3 4
2 4
2 4
2 1
2 1
1 2
2 1
1 2
4 2
2 1
3 2
4 2
1 4
4 2
3 3
2 1
3 2
3 1
1 3
2 1
5 2
2 1
3 1
4 1
5 3
5 2
3 4
2 1
2 1
2 1
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
	ref := "./refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "786E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func readTests() ([][]byte, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesERaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("empty test file")
	}
	t, _ := strconv.Atoi(scan.Text())
	tests := make([][]byte, 0, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("bad test %d", i+1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		mStr := scan.Text()
		m, _ := strconv.Atoi(mStr)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j := 0; j < n-1; j++ {
			scan.Scan()
			vStr := scan.Text()
			scan.Scan()
			uStr := scan.Text()
			sb.WriteString(fmt.Sprintf("%s %s\n", vStr, uStr))
		}
		for j := 0; j < m; j++ {
			scan.Scan()
			xStr := scan.Text()
			scan.Scan()
			yStr := scan.Text()
			sb.WriteString(fmt.Sprintf("%s %s\n", xStr, yStr))
		}
		tests = append(tests, []byte(sb.String()))
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
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
