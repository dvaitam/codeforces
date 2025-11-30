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

// Base64-encoded contents of testcasesE.txt.
const testcasesE = "MTAwCjQKMSA0IDUgNAo1IDMgMyA1CjEgMiA0IDQKMiAzIDUgMQozIDQgNCAxCjIKMSAzCjQgNQoxIDIgMSAzCjUKMiAyIDUgNSAzCjQgMiAxIDUgMwoxIDIgMyA1CjIgMyAzIDUKMyA0IDUgMgo0IDUgMiAxCjMKMyAxIDMKNSA0IDUKMSAyIDUgMwoyIDMgMSAyCjQKMiA0IDEgNAo1IDMgMyAyCjEgMiA1IDMKMiAzIDMgNQozIDQgNSAxCjUKMyA1IDQgMyA1CjIgNSAxIDIgNAoxIDIgNCAzCjIgMyA1IDUKMyA0IDIgNAo0IDUgMSAyCjIKMSAzCjUgNQoxIDIgNSA0CjUKMyAxIDUgMyA1CjUgNCAzIDEgNQoxIDIgMiAyCjIgMyAzIDQKMyA0IDIgMgo0IDUgNCAyCjMKMSAyIDUKMyA1IDQKMSAyIDEgMwoyIDMgNCAxCjIKNCAxCjMgMQoxIDIgMyAxCjMKMiA1IDQKMiA0IDUKMSAyIDQgNQoyIDMgNSAxCjIKMiAyCjUgMwoxIDIgMSA1CjMKMSA0IDEKMSAxIDMKMSAyIDEgNQoyIDMgMyAyCjMKNSA1IDMKNSA1IDEKMSAyIDUgNAoyIDMgMyA1CjIKMSA0CjQgMwoxIDIgMyA1CjUKNSAyIDEgNSAxCjQgNCAxIDQgNQoxIDIgNSAxCjIgMyAzIDQKMyA0IDUgMQo0IDUgMiAzCjQKMiA0IDEgNQoxIDQgNCAzCjEgMiA1IDEKMiAzIDUgMgozIDQgMyAxCjQKNSA0IDQgMQoxIDUgNSA1CjEgMiA0IDUKMiAzIDMgNQozIDQgMyA0CjMKMSAxIDUKMyA0IDMKMSAyIDUgMwoyIDMgMSAxCjUKMyAxIDMgNSAxCjQgMyAzIDIgMwoxIDIgMiA0CjIgMyAzIDQKMyA0IDQgMwo0IDUgMSAyCjUKMSAxIDQgNCA1CjQgNSAzIDMgMQoxIDIgMyAxCjIgMyAyIDQKMyA0IDMgNAo0IDUgNCAyCjQKMyA1IDEgNQo0IDUgNSA1CjEgMiA1IDQKMiAzIDUgNQozIDQgMiA0CjIKMiA1CjIgNQoxIDIgMSA0CjUKMSAyIDQgMiA1CjQgMSA1IDIgMQoxIDIgNSA0CjIgMyAxIDQKMyA0IDQgNAo0IDUgMyAxCjUKMSAxIDQgNCA1CjEgNSAyIDQgNAoxIDIgNCA1CjIgMyAxIDQKMyA0IDUgNQo0IDUgMyA0CjMKMyA0IDEKMiA0IDMKMSAyIDEgMwoyIDMgNCAxCjQKMSA0IDUgMQo1IDIgNSAyCjEgMiA1IDMKMiAzIDMgNQozIDQgNSAxCjIKMyAyCjQgMgoxIDIgNCA1CjUKMSA1IDQgNSAyCjUgMyAxIDMgMQoxIDIgMiA1CjIgMyAxIDIKMyA0IDIgMQo0IDUgMiA1CjMKMiAxIDEKMiA1IDQKMSAyIDMgMQoyIDMgNCAzCjMKMyAxIDMKMSAxIDIKMSAyIDUgMgoyIDMgMyAxCjQKNCA1IDUgNQozIDUgNSA0CjEgMiA0IDQKMiAzIDQgNAozIDQgMSAyCjQKMSA1IDQgMQo0IDUgMSAxCjEgMiAzIDQKMiAzIDEgMwozIDQgMiA0CjUKNSA0IDUgMSAyCjIgMSA1IDQgMwoxIDIgMiAyCjIgMyAzIDEKMyA0IDEgNQo0IDUgMyAxCjIKMyAzCjIgNQoxIDIgMiAzCjQKMiAyIDQgNQoxIDMgNSAzCjEgMiA0IDEKMiAzIDQgNAozIDQgMyA0CjQKMSA1IDQgMgozIDQgMSAxCjEgMiA0IDIKMiAzIDQgMwozIDQgNSAxCjQKNSAzIDQgMQoxIDQgMyAxCjEgMiAyIDIKMiAzIDIgNAozIDQgNCAyCjQKMSAxIDUgNQoxIDMgNSAyCjEgMiA1IDMKMiAzIDQgMgozIDQgMyA0CjIKNSA0CjMgNQoxIDIgMyAxCjIKMSAzCjEgNQoxIDIgMiA1CjIKMSA1CjUgMQoxIDIgNCA0CjMKMSAzIDQKNSAzIDEKMSAyIDEgMwoyIDMgNSAzCjIKNSA0CjQgNAoxIDIgMyAxCjQKNCAzIDEgNAo0IDIgMyAxCjEgMiAxIDUKMiAzIDUgMgozIDQgMiA1CjMKMSAxIDQKNCA0IDQKMSAyIDUgMgoyIDMgMiA0CjUKMyAyIDEgMSA1CjIgNSAyIDEgMQoxIDIgMiAxCjIgMyA1IDIKMyA0IDUgMwo0IDUgNSA1CjUKNSA0IDUgNSAxCjIgNSA1IDUgMgoxIDIgMiAxCjIgMyAyIDMKMyA0IDUgNAo0IDUgNCA0CjQKNCA0IDUgMgoxIDQgNCA0CjEgMiA0IDMKMiAzIDMgNQozIDQgMyAzCjUKNCAxIDUgNSAzCjIgNCA0IDEgMQoxIDIgNSAxCjIgMyAzIDMKMyA0IDUgMgo0IDUgMSAzCjQKMiA1IDEgNAo0IDMgMiAzCjEgMiAzIDMKMiAzIDQgMQozIDQgMiAyCjMKMyAyIDQKNSAzIDIKMSAyIDUgMwoyIDMgMSAyCjMKNSAyIDUKNSA1IDMKMSAyIDUgMgoyIDMgMiAxCjQKNCAyIDQgNQoyIDIgNCAxCjEgMiA0IDMKMiAzIDUgNQozIDQgMSAzCjIKMiA1CjUgMgoxIDIgMyA1CjUKMSAzIDUgMyAxCjQgMyA1IDIgMwoxIDIgMyA1CjIgMyAzIDUKMyA0IDEgNQo0IDUgMiAzCjUKMyAxIDIgNSAyCjMgMyAzIDEgMwoxIDIgMiAxCjIgMyAzIDUKMyA0IDIgMgo0IDUgNCA0CjIKNCAxCjUgMgoxIDIgMSA0CjQKMiA0IDIgMQoyIDIgNSAxCjEgMiA1IDQKMiAzIDMgMwozIDQgMyAxCjIKMyA0CjIgMQoxIDIgMSA1CjQKNSA0IDIgMQozIDMgNCAzCjEgMiA1IDEKMiAzIDEgMwozIDQgNCA1CjMKNSA1IDQKMyAzIDQKMSAyIDUgMgoyIDMgNSAzCjMKNSAyIDIKMiAxIDEKMSAyIDUgNQoyIDMgNSAzCjQKMiAzIDUgNQo1IDMgNSAyCjEgMiAzIDMKMiAzIDIgMgozIDQgNCAxCjIKNCA0CjQgMwoxIDIgNSAxCjUKMyAzIDMgMiA0CjEgMyA1IDQgMwoxIDIgNCAzCjIgMyAyIDIKMyA0IDIgMwo0IDUgMiAzCjIKMiA1CjQgMQoxIDIgMyAxCjQKMyA1IDQgNQoyIDQgNCAxCjEgMiAxIDIKMiAzIDQgNAozIDQgMiAyCjIKNCA0CjUgNQoxIDIgMSA1CjIKNSA0CjUgNAoxIDIgNSAyCjMKMiAzIDMKMiA1IDQKMSAyIDUgMwoyIDMgNCAzCjQKNSAxIDIgMwoyIDEgMiAyCjEgMiA0IDEKMiAzIDMgNQozIDQgMiAzCjIKMyA1CjIgMwoxIDIgMyAyCjMKMSAyIDEKNSAyIDUKMSAyIDUgNQoyIDMgNSAxCjQKMSAzIDIgMQo0IDMgMSAzCjEgMiAyIDUKMiAzIDQgNQozIDQgMiAyCjMKNSAyIDUKMyAzIDQKMSAyIDEgMwoyIDMgMiAxCjUKMiAzIDQgMiA0CjEgMyAyIDQgNQoxIDIgMiAyCjIgMyAyIDUKMyA0IDEgMgo0IDUgMyA1CjUKMyAzIDEgMyAyCjEgNSA0IDMgMgoxIDIgMyA0CjIgMyAyIDIKMyA0IDUgMwo0IDUgNCAzCjQKNSAxIDMgMQoxIDUgNSA0CjEgMiA1IDUKMiAzIDQgMgozIDQgNCA0CjMKNSA1IDQKNCA1IDIKMSAyIDMgNAoyIDMgMyAzCjIKNCA0CjIgMQoxIDIgNCAyCjMKMSAzIDUKMSAzIDQKMSAyIDUgMwoyIDMgMiAzCjIKMSAyCjIgNQoxIDIgMyAzCjIKNCAzCjMgNQoxIDIgMyAyCjIKMyA0CjIgMQoxIDIgNSAzCjMKNCAzIDUKMiA1IDMKMSAyIDIgMgoyIDMgMyAyCjQKNCAxIDMgMwoyIDQgMiAyCjEgMiAyIDEKMiAzIDMgNQozIDQgMiA0CjQKMSAxIDUgMgoyIDEgMiA0CjEgMiA0IDMKMiAzIDMgMQozIDQgMyAzCjUKMSAxIDMgMyA1CjIgNSAyIDIgNAoxIDIgNSAyCjIgMyAxIDIKMyA0IDEgNQo0IDUgMyAyCjQKMiAyIDQgMwoyIDQgMiA0CjEgMiA1IDUKMiAzIDIgMQozIDQgNCA0CjMKNCAxIDQKNSA1IDEKMSAyIDQgMgoyIDMgNCAyCjIKMiA0CjUgMQoxIDIgNSAzCjUKMyAzIDUgNSA1CjQgNSAyIDEgMgoxIDIgMiA0CjIgMyA1IDQKMyA0IDQgNAo0IDUgNSA0CjMKNCAzIDEKNSAzIDUKMSAyIDEgNQoyIDMgMyA1CjUKMyAyIDEgMiAxCjUgMyAzIDIgMwoxIDIgNSAzCjIgMyAyIDUKMyA0IDUgNQo0IDUgMiAyCjIKMiAzCjEgMgoxIDIgMiA0CjQKNSAxIDUgNQoyIDIgMyA0CjEgMiA1IDMKMiAzIDEgMwozIDQgMyAxCjQKNSAzIDQgNQo1IDQgMiAzCjEgMiAyIDIKMiAzIDEgMQozIDQgMSA0CjUKMiA1IDUgMiAzCjEgNSAyIDUgMgoxIDIgNCA0CjIgMyA0IDEKMyA0IDQgMgo0IDUgNSA1CjQKNSAyIDQgNAoyIDIgMiAyCjEgMiAzIDEKMiAzIDUgNQozIDQgMiA1Cg=="

type testCase struct {
	n     int
	a, b  []int
	edges [][4]int
}

// Embedded solver logic from 1824E.go (placeholder that always returns 0).
func solve(tc testCase) string {
	return "0"
}

func parseTestcases() ([]testCase, error) {
	raw, err := base64.StdEncoding.DecodeString(testcasesE)
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
		a := make([]int, n)
		b := make([]int, n)
		for j := 0; j < n; j++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("case %d a[%d]: missing", i+1, j)
			}
			a[j], err = strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("case %d a[%d]: %v", i+1, j, err)
			}
		}
		for j := 0; j < n; j++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("case %d b[%d]: missing", i+1, j)
			}
			b[j], err = strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("case %d b[%d]: %v", i+1, j, err)
			}
		}
		edges := make([][4]int, n-1)
		for j := 0; j < n-1; j++ {
			for k := 0; k < 4; k++ {
				if !sc.Scan() {
					return nil, fmt.Errorf("case %d edge %d field %d missing", i+1, j, k)
				}
				edges[j][k], err = strconv.Atoi(sc.Text())
				if err != nil {
					return nil, fmt.Errorf("case %d edge %d field %d: %v", i+1, j, k, err)
				}
			}
		}
		cases = append(cases, testCase{n: n, a: a, b: b, edges: edges})
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
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
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
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", tc.n)
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d %d %d\n", e[0], e[1], e[2], e[3])
		}

		want := solve(tc)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input.String(), want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
