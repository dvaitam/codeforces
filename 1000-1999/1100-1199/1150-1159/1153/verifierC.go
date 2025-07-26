package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	input    string
	n        int
	original string
}

func generateTests() []Test {
	rng := rand.New(rand.NewSource(44))
	tests := make([]Test, 0, 100)
	for len(tests) < 100 {
		n := (rng.Intn(20) + 2) / 2 * 2 // even between 2 and 20
		var sb strings.Builder
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			r := rng.Intn(3)
			if r == 0 {
				b[i] = '('
			} else if r == 1 {
				b[i] = ')'
			} else {
				b[i] = '?'
			}
		}
		if b[0] == ')' {
			b[0] = '?'
		}
		if b[n-1] == '(' {
			b[n-1] = '?'
		}
		s := string(b)
		sb.WriteString(fmt.Sprintf("%d\n%s\n", n, s))
		// ensure there is a solution using algorithm from 1153C
		res := solveC(n, s)
		if res == ":(" {
			continue
		}
		tests = append(tests, Test{input: sb.String(), n: n, original: s})
	}
	return tests
}

func solveC(n int, s string) string {
	if n%2 == 1 {
		return ":("
	}
	if s[0] == ')' || s[n-1] == '(' {
		return ":("
	}
	b := []byte(s)
	b[0] = '('
	b[n-1] = ')'
	open := 0
	close := 0
	for i := 0; i < n; i++ {
		if b[i] == '(' {
			open++
		} else if b[i] == ')' {
			close++
		}
	}
	need := n / 2
	remOpen := need - open
	remClose := need - close
	if remOpen < 0 || remClose < 0 {
		return ":("
	}
	for i := 1; i < n-1; i++ {
		if b[i] == '?' {
			if remOpen > 0 {
				b[i] = '('
				remOpen--
			} else {
				b[i] = ')'
				remClose--
			}
		}
	}
	bal := 0
	for i := 0; i < n; i++ {
		if b[i] == '(' {
			bal++
		} else {
			bal--
		}
		if i < n-1 && bal <= 0 {
			return ":("
		}
	}
	if bal != 0 {
		return ":("
	}
	return string(b)
}

func validate(n int, orig, out string) bool {
	if len(out) != n {
		return false
	}
	for i := 0; i < n; i++ {
		if orig[i] != '?' && orig[i] != out[i] {
			return false
		}
		if out[i] != '(' && out[i] != ')' {
			return false
		}
	}
	bal := 0
	for i := 0; i < n; i++ {
		if out[i] == '(' {
			bal++
		} else {
			bal--
		}
		if bal < 0 {
			return false
		}
		if i < n-1 && bal == 0 {
			return false
		}
	}
	return bal == 0
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if !validate(tc.n, tc.original, got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected valid output\ngot: %s\n", i+1, tc.input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
