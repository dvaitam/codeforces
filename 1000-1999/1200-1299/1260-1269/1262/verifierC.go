package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type op struct{ l, r int }

func targetString(n, k int) string {
	var b strings.Builder
	for i := 0; i < k-1; i++ {
		b.WriteString("()")
	}
	rem := n - 2*(k-1)
	b.WriteString(strings.Repeat("(", rem/2))
	b.WriteString(strings.Repeat(")", rem/2))
	return b.String()
}

func applyOps(s []byte, ops []op) []byte {
	for _, o := range ops {
		l, r := o.l-1, o.r-1
		for l < r {
			s[l], s[r] = s[r], s[l]
			l++
			r--
		}
	}
	return s
}

func parseOutput(out string) ([]op, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty output")
	}
	var m int
	if _, err := fmt.Sscan(scanner.Text(), &m); err != nil {
		return nil, fmt.Errorf("bad m: %v", err)
	}
	var ops []op
	for i := 0; i < m; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("missing op line %d", i+1)
		}
		var l, r int
		if _, err := fmt.Sscan(scanner.Text(), &l, &r); err != nil {
			return nil, fmt.Errorf("bad op %d: %v", i+1, err)
		}
		ops = append(ops, op{l, r})
	}
	return ops, nil
}

func run(bin, input string) (string, error) {
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
	return out.String(), nil
}

func runCase(bin string, n, k int, s string) error {
	input := fmt.Sprintf("1\n%d %d\n%s\n", n, k, s)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	ops, err := parseOutput(out)
	if err != nil {
		return err
	}
	if len(ops) > n {
		return fmt.Errorf("too many operations")
	}
	b := []byte(s)
	for _, o := range ops {
		if o.l < 1 || o.r > n || o.l > o.r {
			return fmt.Errorf("invalid op")
		}
	}
	b = applyOps(b, ops)
	if string(b) != targetString(n, k) {
		return fmt.Errorf("wrong final string")
	}
	return nil
}

func genCase(rng *rand.Rand) (int, int, string) {
	n := (rng.Intn(10) + 1) * 2
	k := rng.Intn(n/2) + 1
	open := n / 2
	close := n / 2
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if open > 0 && close > 0 {
			if rng.Intn(2) == 0 {
				b[i] = '('
				open--
			} else {
				b[i] = ')'
				close--
			}
		} else if open > 0 {
			b[i] = '('
			open--
		} else {
			b[i] = ')'
			close--
		}
	}
	return n, k, string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []struct {
		n, k int
		s    string
	}{
		{8, 2, "(()())()"},
	}
	for idx, c := range cases {
		if err := runCase(bin, c.n, c.k, c.s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	for i := len(cases); i < 100; i++ {
		n, k, s := genCase(rng)
		if err := runCase(bin, n, k, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
