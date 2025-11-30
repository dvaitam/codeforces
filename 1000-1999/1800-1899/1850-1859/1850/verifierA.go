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
const testcasesA = "MCAwIDAKMCAwIDEKMCAwIDIKMCAwIDMKMCAwIDQKMCAwIDUKMCAwIDYKMCAwIDcKMCAwIDgKMCAwIDkKMCAxIDAKMCAxIDEKMCAxIDIKMCAxIDMKMCAxIDQKMCAxIDUKMCAxIDYKMCAxIDcKMCAxIDgKMCAxIDkKMCAyIDAKMCAyIDEKMCAyIDIKMCAyIDMKMCAyIDQKMCAyIDUKMCAyIDYKMCAyIDcKMCAyIDgKMCAyIDkKMCAzIDAKMCAzIDEKMCAzIDIKMCAzIDMKMCAzIDQKMCAzIDUKMCAzIDYKMCAzIDcKMCAzIDgKMCAzIDkKMCA0IDAKMCA0IDEKMCA0IDIKMCA0IDMKMCA0IDUKMCA0IDYKMCA0IDcKMCA0IDgKMCA0IDkKMCA1IDAKMCA1IDEKMCA1IDIKMCA1IDMKMCA1IDQKMCA1IDUKMCA1IDYKMCA1IDcKMCA1IDgKMCA1IDkKMCA2IDAKMCA2IDEKMCA2IDIKMCA2IDMKMCA2IDQKMCA2IDUKMCA2IDYKMCA2IDcKMCA2IDgKMCA2IDkKMCA3IDAKMCA3IDEKMCA3IDIKMCA3IDMKMCA3IDQKMCA3IDUKMCA3IDYKMCA3IDcKMCA3IDgKMCA3IDkKMCA4IDAKMCA4IDEKMCA4IDIKMCA4IDMKMCA4IDQKMCA4IDUKMCA4IDYKMCA4IDcKMCA4IDgKMCA4IDkKMCA5IDAKMCA5IDEKMCA5IDIKMCA5IDMKMCA5IDQKMCA5IDUKMCA5IDYKMCA5IDcKMCA5IDgKMCA5IDkKMSAwIDAKMSAwIDEKMSAwIDIKMSAwIDMKMSAwIDQKMSAwIDUKMSAwIDYKMSAwIDcKMSAwIDgKMSAwIDkKMSAxIDAKMSAxIDEKMSAxIDIKMSAxIDMKMSAxIDQKMSAxIDUKMSAxIDYKMSAxIDcKMSAxIDgKMSAxIDkKMSAyIDAKMSAyIDEKMSAyIDIKMSAyIDMKMSAyIDQKMSAyIDUKMSAyIDYKMSAyIDcKMSAyIDgKMSAyIDkKMSAzIDAKMSAzIDEKMSAzIDIKMSAzIDMKMSAzIDQKMSAzIDUKMSAzIDYKMSAzIDcKMSAzIDgKMSAzIDkKMSA0IDAKMSA0IDEKMSA0IDIKMSA0IDMKMSA0IDQKMSA0IDUKMSA0IDYKMSA0IDcKMSA0IDgKMSA0IDkKMSA1IDAKMSA1IDEKMSA1IDIKMSA1IDMKMSA1IDQKMSA1IDUKMSA1IDYKMSA1IDcKMSA1IDgKMSA1IDkKMSA2IDAKMSA2IDEKMSA2IDIKMSA2IDMKMSA2IDQKMSA2IDUKMSA2IDYKMSA2IDcKMSA2IDgKMSA2IDkKMSA3IDAKMSA3IDEKMSA3IDIKMSA3IDMKMSA3IDQKMSA3IDUKMSA3IDYKMSA3IDcKMSA3IDgKMSA3IDkKMSA4IDAKMSA4IDEKMSA4IDIKMSA4IDMKMSA4IDQKMSA4IDUKMSA4IDYKMSA4IDcKMSA4IDgKMSA4IDkKMSA5IDAKMSA5IDEKMSA5IDIKMSA5IDMKMSA5IDQKMSA5IDUKMSA5IDYKMSA5IDcKMSA5IDgKMSA5IDkK"

type testCase struct {
	a, b, c int
}

// Embedded solver logic from 1850A.go.
func solve(tc testCase) string {
	if tc.a+tc.b >= 10 || tc.a+tc.c >= 10 || tc.b+tc.c >= 10 {
		return "YES"
	}
	return "NO"
}

func parseTestcases() ([]testCase, error) {
	raw, err := base64.StdEncoding.DecodeString(testcasesA)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(bytes.NewReader(raw))
	sc.Split(bufio.ScanWords)
	cases := []testCase{}
	for {
		if !sc.Scan() {
			break
		}
		a, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("parse a: %v", err)
		}
		if !sc.Scan() {
			return nil, fmt.Errorf("missing b")
		}
		b, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("parse b: %v", err)
		}
		if !sc.Scan() {
			return nil, fmt.Errorf("missing c")
		}
		c, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("parse c: %v", err)
		}
		cases = append(cases, testCase{a: a, b: b, c: c})
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

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&input, "%d %d %d\n", tc.a, tc.b, tc.c)
	}

	output, err := runCandidate(bin, input.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "execution failed: %v\n", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(strings.NewReader(output))
	outScan.Split(bufio.ScanWords)
	for idx, tc := range cases {
		if !outScan.Scan() {
			fmt.Fprintf(os.Stderr, "missing output for case %d\n", idx+1)
			os.Exit(1)
		}
		got := strings.ToUpper(strings.TrimSpace(outScan.Text()))
		want := solve(tc)
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected: %s\ngot: %s\n", idx+1, input.String(), want, got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
