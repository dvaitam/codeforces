package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Op struct {
	l, r  int64
	start int64
	end   int64
}

func resolveChar(s string, ops []Op, k int64) byte {
	origLen := int64(len(s))
	for k > origLen {
		for i := len(ops) - 1; i >= 0; i-- {
			op := ops[i]
			if k >= op.start && k <= op.end {
				k = op.l + (k - op.start)
				break
			}
		}
	}
	return s[k-1]
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	c := rng.Intn(3) + 1
	q := rng.Intn(5) + 1
	var sb strings.Builder
	s := make([]byte, n)
	for i := range s {
		s[i] = byte('a' + rng.Intn(3))
	}
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d %d\n", n, c, q)
	sb.WriteString(string(s))
	sb.WriteByte('\n')
	ops := make([]Op, c)
	currLen := int64(n)
	for i := 0; i < c; i++ {
		l := rng.Int63n(currLen) + 1
		r := l + rng.Int63n(currLen-l+1)
		ops[i] = Op{l: l, r: r, start: currLen + 1, end: currLen + (r - l + 1)}
		currLen += r - l + 1
		fmt.Fprintf(&sb, "%d %d\n", l, r)
	}
	queries := make([]int64, q)
	for i := 0; i < q; i++ {
		queries[i] = rng.Int63n(currLen) + 1
	}
	for _, k := range queries {
		fmt.Fprintf(&sb, "%d\n", k)
	}
	var expSB strings.Builder
	for _, k := range queries {
		ch := resolveChar(string(s), ops, k)
		expSB.WriteByte(ch)
		expSB.WriteByte('\n')
	}
	return sb.String(), strings.TrimSpace(expSB.String())
}

func runCase(bin, input, expected string) error {
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
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1705))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
