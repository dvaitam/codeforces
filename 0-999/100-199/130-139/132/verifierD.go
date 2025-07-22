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

type operation struct {
	pow int
	op  byte
}

func computeOps(binStr string) []operation {
	n := len(binStr)
	s := make([]byte, n+2)
	s[0], s[1] = '0', '0'
	copy(s[2:], binStr)
	ops := make([]operation, 0)
	for i := n + 1; i >= 1; i-- {
		if s[i] == '1' {
			if s[i-1] == '0' {
				ops = append(ops, operation{pow: n + 1 - i, op: '+'})
			} else {
				j := i
				for j >= 0 && s[j] == '1' {
					s[j] = '0'
					j--
				}
				s[j] = '1'
				ops = append(ops, operation{pow: n + 1 - i, op: '-'})
			}
		}
	}
	return ops
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, string) {
	L := rng.Intn(50) + 1
	b := make([]byte, L)
	b[0] = '1'
	for i := 1; i < L; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	binStr := string(b)
	var sb strings.Builder
	sb.WriteString(binStr)
	sb.WriteByte('\n')
	return binStr, sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		binStr, input := genCase(rng)
		ops := computeOps(binStr)
		var expect strings.Builder
		fmt.Fprintf(&expect, "%d\n", len(ops))
		for _, o := range ops {
			fmt.Fprintf(&expect, "%c2^%d\n", o.op, o.pow)
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect.String()) {
			fmt.Fprintf(os.Stderr, "case %d mismatch:\nexpected:\n%s\ngot:\n%s", i+1, expect.String(), out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
