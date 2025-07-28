package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func canTransform(s, t string) bool {
	n := len(s)
	sNoB := make([]byte, 0, n)
	tNoB := make([]byte, 0, n)
	posAS := make([]int, 0)
	posAT := make([]int, 0)
	posCS := make([]int, 0)
	posCT := make([]int, 0)

	for i := 0; i < n; i++ {
		if s[i] != 'b' {
			sNoB = append(sNoB, s[i])
			if s[i] == 'a' {
				posAS = append(posAS, i)
			} else {
				posCS = append(posCS, i)
			}
		}
		if t[i] != 'b' {
			tNoB = append(tNoB, t[i])
			if t[i] == 'a' {
				posAT = append(posAT, i)
			} else {
				posCT = append(posCT, i)
			}
		}
	}
	if string(sNoB) != string(tNoB) {
		return false
	}
	for i := range posAS {
		if posAS[i] > posAT[i] {
			return false
		}
	}
	for i := range posCS {
		if posCS[i] < posCT[i] {
			return false
		}
	}
	return true
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	b := make([]byte, n)
	c := make([]byte, n)
	letters := []byte{'a', 'b', 'c'}
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(3)]
		c[i] = letters[rng.Intn(3)]
	}
	s := string(b)
	t := string(c)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	sb.WriteString(fmt.Sprintf("%s\n%s\n", s, t))
	expected := "NO"
	if canTransform(s, t) {
		expected = "YES"
	}
	return testCase{input: sb.String(), expected: expected}
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{}
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}

	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d error: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, strings.TrimSpace(out), tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
