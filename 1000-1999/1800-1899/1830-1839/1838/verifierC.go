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

// Base64-encoded contents of testcasesC.txt.
const testcasesC = "MTAwCjEgMQoxIDMKMiAzCjMgNQoyIDUKMSA1CjIgNAo0IDUKMyA1CjQgNQozIDEKMSAzCjQgMwo0IDQKNSAyCjUgMgoyIDIKMSAyCjMgMgoyIDUKNSAzCjUgNQoyIDQKNCA1CjMgNQozIDMKNCAyCjQgNAo1IDIKNCAzCjQgNQo1IDMKNCA0CjMgNQo1IDQKNCAyCjMgMgo1IDMKNCAzCjMgNQo1IDUKNSA1CjUgNAozIDIKNCA1CjMgNQoxIDMKMSAyCjEgMQo1IDEKMyA1CjIgMQo1IDIKMyAyCjIgMQo0IDEKMSAzCjMgMgoyIDEKMSAxCjEgMQoxIDEKMyAzCjIgMgoyIDUKMSA0CjUgMQoyIDIKMSAxCjMgNQoxIDMKMyA0CjEgMwo0IDUKNSAxCjMgNAo1IDIKNCAyCjEgMwoxIDEKNCAyCjUgNQo0IDQKNSAzCjIgMwozIDMKNSA0CjEgNQoyIDEKMyAxCjIgMgoyIDEKNCAyCjUgMQoyIDIKNCAxCjMgMQo1IDIKNSA1CjMgMwo="

type testCase struct {
	n, m int
}

// Embedded solver logic from 1838C.go.
func solve(tc testCase) string {
	var sb strings.Builder
	for i := 2; i <= tc.n; i += 2 {
		base := (i - 1) * tc.m
		for j := 1; j <= tc.m; j++ {
			sb.WriteString(strconv.Itoa(base + j))
			if j < tc.m {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	for i := 1; i <= tc.n; i += 2 {
		base := (i - 1) * tc.m
		for j := 1; j <= tc.m; j++ {
			sb.WriteString(strconv.Itoa(base + j))
			if j < tc.m {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}

func parseTestcases() ([]testCase, error) {
	raw, err := base64.StdEncoding.DecodeString(testcasesC)
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
		if !sc.Scan() {
			return nil, fmt.Errorf("case %d missing m", i+1)
		}
		m, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d m: %v", i+1, err)
		}
		cases = append(cases, testCase{n: n, m: m})
	}
	if err := sc.Err(); err != nil {
		return nil, err
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
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
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
		input := fmt.Sprintf("1\n%d %d\n", tc.n, tc.m)
		want := solve(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
