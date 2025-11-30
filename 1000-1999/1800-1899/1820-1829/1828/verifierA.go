package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Base64-encoded contents of testcasesA.txt.
const testcasesA = "MTAwCjk5CjE5NQoxMDgKMTEKNjcKMTMxCjEyNQoxMDQKNzgKMTIzCjkyCjE1MAo1NgoxMzAKMzYKNzMKMzYKMTk0CjI1CjE1OQo2NQoxMzcKMTgxCjE1NQozOAo4MAoyNgoxODcKMTkKMTc2Cjg1CjEyMQoxNDQKMjYKOTEKMTEyCjgxCjE1NwoxNjQKNTMKMTQyCjEyMwoxMTQKMTM0CjY3CjE2CjE0MQo0CjI0CjE4NQoxMDMKMTgyCjE3MgoxNjEKMQoxNTcKMTI3Cjg2CjYzCjE4Nwo4NAoxODEKMTcKNDkKMTQ2CjU3CjYyCjM3CjE0MAoxMTUKMjQKMjEKODIKMTMxCjEyNgoyOAo3OAoxNDIKNzUKMTgxCjMyCjE0MQo4NgoxMzkKNTMKMTU1CjE0MQoxNTEKNzQKMTE0CjI0CjE1Mwo5OQo4MgoxNDgKNjIKNzUKNDgKNDkKNDgK"

type testCase struct {
	n int
}

// Embedded solver logic from 1828A.go.
func solve(tc testCase) string {
	var sb strings.Builder
	for i := 1; i <= tc.n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(2 * i))
	}
	return sb.String()
}

func parseTestcases() ([]testCase, error) {
	raw, err := base64.StdEncoding.DecodeString(testcasesA)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(bytes.NewReader(raw))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return nil, fmt.Errorf("invalid test data")
	}
	t, err := strconv.Atoi(sc.Text())
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("case %d missing n", i+1)
		}
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d n: %v", i+1, err)
		}
		cases = append(cases, testCase{n: n})
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "bin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		input := fmt.Sprintf("1\n%d\n", tc.n)
		want := solve(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
