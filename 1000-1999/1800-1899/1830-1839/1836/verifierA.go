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
const testcasesA = "MyA5IDEgNAoyIDcgNwo4IDYgMyAxIDcgMCA2IDYgOQoxIDcKNSAzIDkgMSA1IDAKMSAwCjkgMCA2IDMgNiAwIDggMyA3IDcKOSAzIDUgMyAzIDcgNCAwIDYgOAoyIDIgNAoyIDUgOAo3IDggMyA0IDQgOSA3IDgKNyA5IDAgNyAzIDYgNiAyCjYgOCA1IDEgNyA4IDEKMyA4IDYgNQo4IDAgNyAwIDQgOSA5IDkgNgozIDIgOCAzCjEgMwo5IDggMyA2IDggNSA5IDUgNyA0CjkgOSAwIDYgOCAyIDggOCAzIDYKMSA3CjYgOSA4IDMgOCA2IDcKNiA2IDUgMCA4IDggOQoxMCA1IDcgOSAwIDMgMiA4IDkgMiAxCjkgNCAwIDEgMSAwIDcgMCA0IDMKNSAxIDkgMiA1IDQKMiAyIDIKNSA4IDIgNCA0IDcKNiA3IDcgMSAwIDQgNgo2IDYgMyA0IDEgNCA4CjQgOSA2IDAgMwoxIDYKMyAwIDIgNwo5IDYgOCAzIDggNyAzIDggMCA2CjEwIDUgNiAwIDQgMiAzIDAgNCAxIDEKNSA0IDIgNiA5IDQKMyAwIDggMAoxMCAzIDkgNyAyIDkgOCAwIDYgMyA1CjIgMyA5CjcgOSAzIDcgMSA2IDQgOAo4IDAgNSA5IDYgNCAwIDIgMwo2IDkgMiA1IDYgMyA0CjIgNiA4CjYgOCA3IDggMyAxIDAKMiAyIDIKMyA4IDMgNAo2IDkgOCA0IDUgNSA1CjIgNCAzCjEwIDcgMiA5IDggMSA1IDAgNiAxIDYKMyAyIDUgMQoxMCA5IDYgMSA5IDggMyA5IDEgNCA1CjUgOSA4IDEgNyA0CjIgMCA0CjEgOQoxIDEKNyAxIDAgMyAzIDkgNiAyCjIgNyAyCjQgMiAxIDYgNgo5IDQgOCA0IDcgNSAxIDMgNSAwCjEgMAo1IDkgNSA3IDYgNQo3IDEgMSA1IDkgNyAxIDQKNCA5IDggNyA1CjUgMiA4IDMgNCAzCjQgNSAxIDQgMQo4IDEgOSA1IDMgNiA0IDAgNQozIDUgOSA0CjQgNSAxIDggOQoxMCA5IDEgMyAzIDAgMyA2IDEgNCA4CjIgMSAwCjEgNAo2IDcgNyAyIDEgOCA1CjIgOCAyCjMgMiAyIDUKNSAxIDggOSA0IDIKNCAyIDggMCA1CjEwIDggMyAyIDQgNiA4IDIgMCAzIDQKMiA3IDYKOSA0IDggNyA4IDcgMCA2IDUgMgo1IDcgMCA2IDkgMAoxIDUKMTAgMiA5IDIgMiA0IDQgNiA5IDYgMgoxMCAxIDMgNyAwIDIgOCA1IDggNyAzCjQgNSA3IDcgMwo3IDUgOCA5IDQgMyAwIDEKOSA1IDIgOCAzIDQgNCA0IDggNQozIDcgOSAxCjIgOSA4CjEwIDYgMiAyIDQgNiAzIDkgMCA3IDYKNiA2IDggMiA4IDAgOAoyIDQgMQo1IDEgMiA5IDEgNwo0IDYgNiA2IDIKNiA3IDIgOSA3IDMgMQo3IDkgOCA2IDEgNCA0IDMKNyA4IDAgMyA4IDcgOSAwCjEgOQo0IDQgMyAyIDQKMyA4IDMgNAo1IDkgNCA3IDIgOAo2IDcgNiAxIDMgOSA2Cg=="

type testCase struct {
	arr []int
}

// Embedded solver logic from 1836A.go.
func solve(tc testCase) string {
	freq := make([]int, 101)
	maxVal := 0
	for _, x := range tc.arr {
		if x > maxVal {
			maxVal = x
		}
		if x >= len(freq) {
			tmp := make([]int, x+1)
			copy(tmp, freq)
			freq = tmp
		}
		freq[x]++
	}
	for i := 1; i <= maxVal; i++ {
		if freq[i] > freq[i-1] {
			return "NO"
		}
	}
	return "YES"
}

func parseTestcases() ([]testCase, error) {
	raw, err := base64.StdEncoding.DecodeString(testcasesA)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(bytes.NewReader(raw))
	sc.Split(bufio.ScanWords)
	cases := []testCase{}
	for sc.Scan() {
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("parse n: %v", err)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("case %d missing value %d", len(cases)+1, i)
			}
			arr[i], err = strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("case %d value %d: %v", len(cases)+1, i, err)
			}
		}
		cases = append(cases, testCase{arr: arr})
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
		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d\n", len(tc.arr))
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		want := solve(tc)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(got)) != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input.String(), want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
