package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type pair struct{ s, t int }

func check(seq []int, t int) (int, bool) {
	s1, s2 := 0, 0
	p1, p2 := 0, 0
	last := 0
	for _, v := range seq {
		if v == 1 {
			p1++
		} else {
			p2++
		}
		if p1 == t || p2 == t {
			if p1 == t && p2 == t {
				return 0, false
			}
			if p1 == t {
				s1++
				last = 1
			} else {
				s2++
				last = 2
			}
			p1, p2 = 0, 0
		}
	}
	if p1 != 0 || p2 != 0 {
		return 0, false
	}
	if s1 == s2 {
		return 0, false
	}
	if s1 > s2 && last == 1 {
		return s1, true
	}
	if s2 > s1 && last == 2 {
		return s2, true
	}
	return 0, false
}

func expectedB(seq []int) []pair {
	res := []pair{}
	for t := 1; t <= len(seq); t++ {
		if s, ok := check(seq, t); ok {
			res = append(res, pair{s, t})
		}
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].s == res[j].s {
			return res[i].t < res[j].t
		}
		return res[i].s < res[j].s
	})
	return res
}

func genCaseB(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	seq := make([]int, n)
	for i := 0; i < n; i++ {
		seq[i] = rng.Intn(2) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range seq {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	pairs := expectedB(seq)
	var exp strings.Builder
	exp.WriteString(fmt.Sprintf("%d\n", len(pairs)))
	for _, p := range pairs {
		fmt.Fprintf(&exp, "%d %d\n", p.s, p.t)
	}
	return sb.String(), exp.String()
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected \n%s\ngot \n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseB(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
