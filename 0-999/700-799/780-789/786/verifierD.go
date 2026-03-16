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

const testcasesDRaw = `100
3 5
2 1 b
3 2 c
3 1
3 1
2 3
1 3
2 3
6 4
2 1 c
3 1 a
4 3 a
5 4 c
6 1 c
1 2
5 1
3 1
3 4
6 4
2 1 b
3 2 a
4 2 a
5 1 a
6 4 a
3 6
4 6
3 4
5 4
6 3
2 1 c
3 1 b
4 3 a
5 3 c
6 2 c
3 5
5 1
6 2
7 5
2 1 b
3 1 a
4 2 c
5 4 a
6 3 a
7 4 a
1 3
4 7
4 7
1 5
5 7
2 4
2 1 c
2 1
1 2
1 2
2 1
7 1
2 1 b
3 2 a
4 2 b
5 4 c
6 4 c
7 5 c
5 1
6 5
2 1 b
3 1 b
4 2 b
5 3 c
6 3 a
4 5
3 1
4 5
5 6
2 1
7 3
2 1 b
3 2 c
4 3 b
5 4 a
6 5 a
7 6 a
3 6
4 3
5 3
3 3
2 1 b
3 2 c
2 1
1 3
3 1
4 5
2 1 c
3 2 a
4 2 a
4 1
1 3
3 2
4 2
1 3
7 2
2 1 b
3 1 a
4 1 c
5 2 b
6 5 a
7 3 b
7 6
1 7
6 3
2 1 b
3 2 c
4 2 b
5 3 c
6 4 b
4 5
4 1
4 2
3 1
2 1 c
3 2 c
3 1
2 4
2 1 c
2 1
1 2
1 2
1 2
7 1
2 1 c
3 2 a
4 3 a
5 4 a
6 5 a
7 3 a
3 7
6 4
2 1 b
3 1 a
4 1 c
5 1 a
6 2 b
2 1
4 6
5 4
1 3
3 3
2 1 a
3 2 b
1 2
1 3
1 3
6 4
2 1 a
3 2 a
4 2 b
5 4 c
6 3 b
3 2
3 4
1 2
5 1
7 4
2 1 c
3 1 a
4 2 b
5 4 c
6 1 c
7 4 a
3 6
4 7
6 3
4 6
5 4
2 1 a
3 1 c
4 2 c
5 1 b
2 4
2 1
3 5
3 1
5 1
2 1 c
3 1 c
4 2 c
5 1 c
1 4
3 2
2 1 a
3 1 c
1 3
2 1
2 3
2 1 a
1 2
2 1
1 2
6 2
2 1 a
3 1 c
4 1 b
5 4 b
6 3 c
4 3
5 4
2 4
2 1 b
1 2
2 1
2 1
1 2
7 4
2 1 c
3 1 a
4 2 a
5 2 c
6 5 b
7 2 c
5 7
3 6
6 7
4 5
6 5
2 1 a
3 2 a
4 2 b
5 4 a
6 2 c
3 2
3 4
2 4
6 4
6 5
3 4
2 1 b
3 1 b
3 1
2 1
1 3
3 1
3 3
2 1 a
3 2 a
1 3
3 1
2 1
2 3
2 1 b
1 2
2 1
2 1
2 3
2 1 b
2 1
1 2
2 1
6 3
2 1 b
3 1 c
4 2 a
5 4 c
6 3 a
5 6
3 6
3 4
4 3
2 1 b
3 1 c
4 2 c
3 1
3 2
2 3
7 4
2 1 a
3 1 b
4 3 c
5 4 c
6 3 c
7 2 b
4 3
7 2
6 1
1 6
3 5
2 1 b
3 1 b
2 1
1 2
3 2
1 3
3 2
4 4
2 1 c
3 1 b
4 2 a
2 3
1 2
2 4
4 2
6 5
2 1 b
3 1 c
4 3 a
5 2 a
6 3 b
3 1
4 3
5 3
6 4
5 6
7 3
2 1 a
3 1 c
4 3 a
5 1 c
6 2 b
7 2 c
6 4
2 7
5 2
2 5
2 1 a
2 1
1 2
1 2
1 2
1 2
4 3
2 1 c
3 1 a
4 3 a
1 3
2 1
2 4
7 2
2 1 b
3 1 a
4 2 a
5 2 c
6 2 b
7 4 b
6 1
5 7
5 2
2 1 b
3 1 b
4 2 c
5 4 b
2 4
1 3
4 3
2 1 c
3 2 c
4 2 a
1 2
3 2
4 2
3 2
2 1 a
3 2 c
1 3
1 3
3 4
2 1 a
3 1 a
2 1
3 1
1 2
1 2
5 1
2 1 a
3 1 b
4 2 a
5 2 b
5 1
4 1
2 1 c
3 2 b
4 2 b
4 2
6 3
2 1 b
3 2 b
4 1 a
5 4 a
6 5 c
1 6
1 2
5 1
5 1
2 1 c
3 2 a
4 2 b
5 4 b
1 4
4 5
2 1 a
3 2 a
4 2 a
4 3
4 2
4 3
2 3
4 2
5 4
2 1 b
3 1 c
4 3 c
5 4 b
1 2
2 1
4 1
1 3
3 5
2 1 c
3 1 b
2 3
3 2
2 1
1 2
1 2
6 3
2 1 b
3 1 c
4 1 a
5 1 a
6 2 c
1 2
3 2
2 5
6 4
2 1 c
3 1 b
4 1 c
5 1 a
6 1 c
1 2
3 1
1 4
4 2
7 5
2 1 c
3 2 c
4 3 b
5 4 c
6 4 a
7 3 c
4 7
4 2
1 5
1 4
3 6
7 1
2 1 b
3 1 b
4 1 c
5 3 b
6 1 c
7 5 a
2 3
2 3
2 1 b
2 1
2 1
2 1
7 2
2 1 c
3 2 c
4 3 b
5 3 a
6 5 b
7 4 b
4 5
1 3
2 2
2 1 a
1 2
1 2
2 2
2 1 b
2 1
2 1
7 3
2 1 a
3 2 b
4 2 b
5 1 a
6 2 b
7 6 b
4 7
6 3
2 4
4 4
2 1 c
3 1 b
4 2 b
1 4
3 2
1 2
2 1
3 1
2 1 a
3 2 b
2 3
2 4
2 1 a
1 2
2 1
1 2
2 1
7 4
2 1 c
3 2 b
4 1 a
5 4 a
6 2 c
7 2 a
6 7
7 4
5 7
7 3
6 2
2 1 a
3 1 a
4 1 b
5 1 b
6 4 c
6 5
5 4
4 5
2 1 b
3 1 a
4 3 a
1 4
3 4
4 3
3 4
3 2
6 3
2 1 c
3 1 b
4 2 a
5 4 c
6 5 b
1 3
5 1
4 2
4 3
2 1 b
3 1 b
4 1 c
4 2
3 2
4 3
2 3
2 1 a
2 1
1 2
2 1
5 2
2 1 c
3 1 c
4 3 a
5 1 c
1 3
1 4
4 4
2 1 b
3 2 c
4 3 a
2 4
2 4
4 1
2 4
2 2
2 1 b
1 2
1 2
5 2
2 1 b
3 1 a
4 1 a
5 3 a
4 3
5 1
3 1
2 1 a
3 2 a
1 2
6 4
2 1 c
3 2 b
4 1 c
5 3 c
6 4 c
4 5
5 3
4 2
5 4
5 2
2 1 c
3 1 b
4 1 b
5 3 a
2 1
5 1
5 5
2 1 a
3 1 a
4 1 b
5 3 c
3 1
4 5
5 4
3 2
5 2
4 4
2 1 c
3 2 c
4 2 c
3 2
1 2
3 2
3 4
4 4
2 1 b
3 2 c
4 1 c
4 1
1 4
4 3
1 2
5 4
2 1 a
3 2 c
4 1 a
5 4 a
4 5
3 1
3 5
3 4
5 3
2 1 c
3 2 a
4 2 a
5 1 b
1 5
1 2
1 2
4 4
2 1 a
3 1 a
4 3 c
3 4
3 2
4 3
1 4
3 4
2 1 b
3 2 c
1 3
1 2
1 2
2 1
5 2
2 1 b
3 1 b
4 1 c
5 4 c
3 4
3 1
7 1
2 1 b
3 2 a
4 1 b
5 4 b
6 1 c
7 3 c
4 5
3 3
2 1 c
3 2 a
1 3
3 1
3 1
6 1
2 1 a
3 2 c
4 2 b
5 3 b
6 5 a
2 5
4 3
2 1 c
3 2 c
4 3 b
3 2
4 1
4 3
7 2
2 1 a
3 1 a
4 2 a
5 4 a
6 3 a
7 1 b
4 2
4 7
7 1
2 1 b
3 2 b
4 1 a
5 3 b
6 3 b
7 1 c
3 6
2 2
2 1 b
2 1
1 2
7 2
2 1 c
3 2 c
4 1 a
5 4 c
6 5 b
7 5 a
6 3
4 3
2 1
2 1 b
1 2
6 4
2 1 a
3 2 c
4 2 b
5 4 c
6 4 a
1 6
4 3
6 2
5 3
2 3
2 1 b
2 1
2 1
2 1
6 1
2 1 c
3 2 c
4 2 b
5 3 a
6 4 c
5 6
4 4
2 1 b
3 1 a
4 2 c
2 4
1 3
2 3
4 2
6 4
2 1 a
3 2 c
4 3 b
5 1 c
6 1 c
4 6
2 5
6 2
1 5
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
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "786D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func readTests() ([][]byte, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesDRaw))
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
		qStr := scan.Text()
		q, _ := strconv.Atoi(qStr)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for j := 0; j < n-1; j++ {
			scan.Scan()
			vStr := scan.Text()
			scan.Scan()
			uStr := scan.Text()
			scan.Scan()
			ch := scan.Text()
			sb.WriteString(fmt.Sprintf("%s %s %s\n", vStr, uStr, ch))
		}
		for j := 0; j < q; j++ {
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
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
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
