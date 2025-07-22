package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func simulate(seq []byte) [10]int {
	cp := 0
	dir := 1
	var cnt [10]int
	for cp >= 0 && cp < len(seq) {
		c := seq[cp]
		if c >= '0' && c <= '9' {
			d := int(c - '0')
			cnt[d]++
			cp += dir
			d--
			idx := cp - dir
			if d < 0 {
				seq = append(seq[:idx], seq[idx+1:]...)
				if dir == 1 {
					cp--
				}
			} else {
				seq[idx] = byte('0' + d)
			}
		} else {
			if c == '<' {
				dir = -1
			} else {
				dir = 1
			}
			cp += dir
			prev := cp - dir
			if cp >= 0 && cp < len(seq) && (seq[cp] == '<' || seq[cp] == '>') {
				seq = append(seq[:prev], seq[prev+1:]...)
				if dir == 1 {
					cp--
				}
			}
		}
	}
	return cnt
}

func solveCase(s string, queries [][2]int) string {
	var sb strings.Builder
	for _, q := range queries {
		sub := []byte(s[q[0]-1 : q[1]])
		res := simulate(sub)
		for i := 0; i < 10; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", res[i]))
		}
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

type test struct{ input, expected string }

func generateTests() []test {
	rng := rand.New(rand.NewSource(45))
	var tests []test
	fixed := []struct {
		s string
		q [][2]int
	}{
		{"1", [][2]int{{1, 1}}},
		{"<>", [][2]int{{1, 2}}},
	}
	for _, f := range fixed {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n%s\n", len(f.s), len(f.q), f.s))
		for _, p := range f.q {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		inp := sb.String()
		tests = append(tests, test{inp, solveCase(f.s, f.q)})
	}
	for len(tests) < 100 {
		n := rng.Intn(8) + 1
		letters := "<>0123456789"
		b := make([]byte, n)
		for i := range b {
			b[i] = letters[rng.Intn(len(letters))]
		}
		s := string(b)
		qn := rng.Intn(3) + 1
		qs := make([][2]int, qn)
		for i := 0; i < qn; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			qs[i] = [2]int{l, r}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n%s\n", n, qn, s))
		for _, p := range qs {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		inp := sb.String()
		tests = append(tests, test{inp, solveCase(s, qs)})
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
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
