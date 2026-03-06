package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)


func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct{ input string }

func hasSolution(l, r, m int64) bool {
	diff := r - l
	for a := l; a <= r; a++ {
		rem := m % a
		n := m / a
		if n >= 1 && rem <= diff {
			return true
		}
		if a-rem <= diff {
			return true
		}
	}
	return false
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(2))
	cases := make([]Case, 0, 100)
	for len(cases) < 100 {
		l := int64(rng.Intn(10) + 1)
		r := l + int64(rng.Intn(10))
		m := rng.Int63n(1000) + 1
		if !hasSolution(l, r, m) {
			continue
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d %d\n", l, r, m)
		cases = append(cases, Case{sb.String()})
	}
	return cases
}

// parseInput extracts (l, r, m) from a single-testcase input string.
func parseInput(s string) (l, r, m int64) {
	fmt.Sscanf(strings.TrimSpace(s), "1\n%d %d %d", &l, &r, &m)
	return
}

func checkAnswer(input, got string) error {
	l, r, m := parseInput(input)
	var a, b, c int64
	if _, err := fmt.Sscan(strings.TrimSpace(got), &a, &b, &c); err != nil {
		return fmt.Errorf("could not parse output %q: %v", got, err)
	}
	if a < l || a > r {
		return fmt.Errorf("a=%d out of [%d,%d]", a, l, r)
	}
	if b < l || b > r {
		return fmt.Errorf("b=%d out of [%d,%d]", b, l, r)
	}
	if c < l || c > r {
		return fmt.Errorf("c=%d out of [%d,%d]", c, l, r)
	}
	// n*a + b - c = m  =>  n = (m - b + c) / a, must be a strictly positive integer
	num := m - b + c
	if num <= 0 || num%a != 0 {
		return fmt.Errorf("no valid n for a=%d b=%d c=%d m=%d (m-b+c=%d)", a, b, c, m, num)
	}
	return nil
}

func runCase(bin string, c Case) error {
	got, err := runBinary(bin, c.input)
	if err != nil {
		return err
	}
	if err := checkAnswer(c.input, got); err != nil {
		return fmt.Errorf("invalid answer %q: %v", got, err)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
