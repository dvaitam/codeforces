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
const testcasesA = "MTAwCjMgMyAyCjEzIDQgMgozIDMgMgozIDIgMQoxMSAxMCA5CjEwIDggNwoxOSAzIDEKMTEgMSAxCjIwIDE2IDEyCjYgMyAyCjYgNCAzCjIgMiAyCjEwIDUgNQoxNyA3IDQKMTkgNyA2CjE0IDE0IDQKMTYgMyAzCjE1IDEgMQo3IDYgMQoxNSAxMCA0CjE4IDE0IDEyCjE5IDE5IDE1CjYgMSAxCjEzIDQgMwozIDIgMQoxNSAxMyA0CjE2IDcgMwozIDEgMQoyIDEgMQoxNiAxMCA5CjE2IDEgMQo2IDIgMQoyMCAzIDMKMjMgMzIgMQo5IDQgMQozIDIgMQo2IDYgNAozIDMgMQo5IDE2IDEKMiAyIDEKOSA5IDEKMTAgOCA1CjggMyAxCjE1IDIgMgo3IDggMgoxIDEgMQozIDMgMQo2IDMgMwo0IDEgMQoxNCA5IDkKNiA1IDMKMTAgOCA0CjE2IDE1IDYKMyAxIDEKMTQgOCA0CjE5IDE3IDExCjEyIDEwIDIKNyA0IDEKMyAyIDEKNiA2IDQKNyA0IDMKMTAgMSAxCjExIDcgNgoxOSAyIDIKMSAxIDEKMiAyIDEKMTkgMTQgOAozIDEgMQoxOSAxNyAxNgoxMyAyIDIKNyAzIDEKMTYgMTIgMwo1IDEgMQozIDMgMwoxNiAyIDEKMTggOCAzCjEzIDcgMwoxOCAxMyA1CjEwIDMgMwoxOCAxIDEKMTMgMiAyCjE1IDEgMQoxOSAxNSAxMgo0IDQgNAoxNyAxMCA3CjEwIDggNgo5IDIgMgoyMCAxMyAzCjkgMSAxCjE0IDEgMQoxOCA5IDgKMTQgMiAyCjIwIDUgMgo3IDQgNAoxOCA5IDcKMTQgMTEgMgoxNiAxNSAxMwoxMCA0IDQKMTYgNCA0CjIgMSAxCjQgMyAxCjEzIDggNgo4IDcgNQoxMyA1IDIKMSAxIDEKMTEgOSA1CjYgMiAyCjE1IDYgNQo="

type testCase struct {
	n, k, x int
}

// Embedded solver logic from 1845A.go.
func solve(tc testCase) string {
	n, k, x := tc.n, tc.k, tc.x
	var sb strings.Builder
	if x != 1 {
		sb.WriteString("YES\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteByte('1')
		}
		return sb.String()
	}
	if k == 1 {
		return "NO"
	}
	if n%2 == 0 {
		sb.WriteString("YES\n")
		sb.WriteString(strconv.Itoa(n / 2))
		sb.WriteByte('\n')
		for i := 0; i < n/2; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteByte('2')
		}
		return sb.String()
	}
	if k == 2 || n == 1 {
		return "NO"
	}
	// odd n, k>2, n>1
	sb.WriteString("YES\n")
	sb.WriteString(strconv.Itoa(n / 2))
	sb.WriteByte('\n')
	for i := 0; i < n/2-1; i++ {
		sb.WriteString("2 ")
	}
	sb.WriteByte('3')
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
		if !sc.Scan() {
			return nil, fmt.Errorf("case %d missing k", i+1)
		}
		k, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d k: %v", i+1, err)
		}
		if !sc.Scan() {
			return nil, fmt.Errorf("case %d missing x", i+1)
		}
		x, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d x: %v", i+1, err)
		}
		cases = append(cases, testCase{n: n, k: k, x: x})
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
		input := fmt.Sprintf("1\n%d %d %d\n", tc.n, tc.k, tc.x)
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
