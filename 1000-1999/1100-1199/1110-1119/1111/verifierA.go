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
	s string
	t string
}

func isVowel(c byte) bool {
	switch c {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	default:
		return false
	}
}

func solve(s, t string) string {
	if len(s) != len(t) {
		return "No"
	}
	for i := 0; i < len(s); i++ {
		if isVowel(s[i]) != isVowel(t[i]) {
			return "No"
		}
	}
	return "Yes"
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(42))
	vowels := "aeiou"
	cons := "bcdfghjklmnpqrstvwxyz"
	tests := make([]test, 0, 100)
	for i := 0; i < 50; i++ {
		n := rng.Intn(8) + 1
		sb := make([]byte, n)
		tb := make([]byte, n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				sb[j] = vowels[rng.Intn(len(vowels))]
			} else {
				sb[j] = cons[rng.Intn(len(cons))]
			}
			if rng.Intn(2) == 0 {
				if isVowel(sb[j]) {
					tb[j] = vowels[rng.Intn(len(vowels))]
				} else {
					tb[j] = cons[rng.Intn(len(cons))]
				}
			} else {
				if isVowel(sb[j]) {
					tb[j] = cons[rng.Intn(len(cons))]
				} else {
					tb[j] = vowels[rng.Intn(len(vowels))]
				}
			}
		}
		tests = append(tests, test{string(sb), string(tb)})
	}
	for i := 0; i < 50; i++ {
		n := rng.Intn(8) + 1
		m := rng.Intn(8) + 1
		if m == n {
			m = (m % 8) + 1
		}
		sb := make([]byte, n)
		tb := make([]byte, m)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				sb[j] = vowels[rng.Intn(len(vowels))]
			} else {
				sb[j] = cons[rng.Intn(len(cons))]
			}
		}
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				tb[j] = vowels[rng.Intn(len(vowels))]
			} else {
				tb[j] = cons[rng.Intn(len(cons))]
			}
		}
		tests = append(tests, test{string(sb), string(tb)})
	}
	return tests
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%s\n%s\n", tc.s, tc.t)
		want := solve(tc.s, tc.t)
		got, err := run(binary, input)
		if err != nil {
			fmt.Printf("Test %d: error running binary: %v\n", i+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(strings.ToLower(got))
		want = strings.TrimSpace(strings.ToLower(want))
		if got != want {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
