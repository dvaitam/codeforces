package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func solve(t string) string {
	n := len(t)
	for k := 1; k < n; k++ {
		if (n+k)%2 != 0 {
			continue
		}
		L := (n + k) / 2
		if k >= L || L > n {
			continue
		}
		s := t[:L]
		if t[n-L:] != s {
			continue
		}
		if L < len(t) {
			if t[L:] != s[k:] {
				continue
			}
		} else {
			if k != L {
				continue
			}
		}
		return fmt.Sprintf("YES\n%s", s)
	}
	return "NO"
}

func randString(rng *rand.Rand, n int) string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(43))
	var tests []test
	fixed := []string{"abrakadabrabrakadabra", "abrakadabrakadabra", "abcabc", "abcd"}
	for _, f := range fixed {
		tests = append(tests, test{f + "\n", solve(f)})
	}
	for len(tests) < 100 {
		if rng.Intn(2) == 0 { // valid case
			L := rng.Intn(10) + 2
			k := rng.Intn(L-1) + 1
			s := randString(rng, L)
			t := s + s[k:]
			tests = append(tests, test{t + "\n", solve(t)})
		} else { // invalid case by corruption
			L := rng.Intn(10) + 2
			k := rng.Intn(L-1) + 1
			s := randString(rng, L)
			t := s + s[k:]
			pos := rng.Intn(len(t))
			tb := []byte(t)
			orig := tb[pos]
			for tb[pos] == orig {
				tb[pos] = byte('a' + rng.Intn(26))
			}
			t = string(tb)
			tests = append(tests, test{t + "\n", solve(t)})
		}
	}
	return tests
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" {
		if len(os.Args) < 3 {
			fmt.Println("usage: go run verifierB.go /path/to/binary")
			os.Exit(1)
		}
		bin = os.Args[2]
	}
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
