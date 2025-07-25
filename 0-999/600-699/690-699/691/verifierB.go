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
	in  string
	out string
}

var mirror = map[byte]byte{
	'A': 'A', 'H': 'H', 'I': 'I', 'M': 'M',
	'O': 'O', 'T': 'T', 'U': 'U', 'V': 'V',
	'W': 'W', 'X': 'X', 'Y': 'Y',
	'b': 'd', 'd': 'b', 'p': 'q', 'q': 'p',
	'o': 'o', 'v': 'v', 'w': 'w', 'x': 'x',
}

func solveCase(s string) string {
	n := len(s)
	for i := 0; i < n; i++ {
		m, ok := mirror[s[i]]
		if !ok || m != s[n-1-i] {
			return "NIE\n"
		}
	}
	return "TAK\n"
}

func buildCase(s string) testCase {
	return testCase{in: s + "\n", out: solveCase(s)}
}

func randomCase(rng *rand.Rand) testCase {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return buildCase(string(b))
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := out.String()
	if strings.TrimSpace(got) != strings.TrimSpace(tc.out) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(tc.out), strings.TrimSpace(got))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	cases = append(cases, buildCase("A"))
	cases = append(cases, buildCase("AB"))
	cases = append(cases, buildCase("oHo"))
	cases = append(cases, buildCase("aa"))
	cases = append(cases, buildCase("pqp"))
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
