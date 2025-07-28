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

type testCase struct {
	input  string
	output string
}

type query struct {
	typ int
	l   int
	r   int
	k   int
	p   int
	x   int
}

func solveCase(a []int, qs []query) string {
	var sb strings.Builder
	for _, q := range qs {
		if q.typ == 1 {
			l := q.l - 1
			r := q.r - 1
			cnt := make(map[int]int)
			for i := l; i <= r; i++ {
				cnt[a[i]]++
			}
			if len(cnt) < q.k {
				sb.WriteString("-1\n")
				continue
			}
			freq := make([]int, 0, len(cnt))
			for _, v := range cnt {
				freq = append(freq, v)
			}
			sort.Ints(freq)
			best := int(1e9)
			for i := 0; i+q.k-1 < len(freq); i++ {
				diff := freq[i+q.k-1] - freq[i]
				if diff < best {
					best = diff
				}
			}
			sb.WriteString(fmt.Sprintf("%d\n", best))
		} else {
			a[q.p-1] = q.x
		}
	}
	return sb.String()
}

func buildCase(a []int, qs []query) testCase {
	n := len(a)
	m := len(qs)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, q := range qs {
		if q.typ == 1 {
			sb.WriteString(fmt.Sprintf("1 %d %d %d\n", q.l, q.r, q.k))
		} else {
			sb.WriteString(fmt.Sprintf("2 %d %d\n", q.p, q.x))
		}
	}
	out := solveCase(append([]int(nil), a...), qs)
	return testCase{input: sb.String(), output: out}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(5) + 1
	}
	qs := make([]query, m)
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			k := rng.Intn(3) + 1
			qs[i] = query{typ: 1, l: l, r: r, k: k}
		} else {
			p := rng.Intn(n) + 1
			x := rng.Intn(5) + 1
			qs[i] = query{typ: 2, p: p, x: x}
		}
	}
	return buildCase(a, qs)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.output)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	a0 := []int{1, 2, 3}
	qs0 := []query{{typ: 1, l: 1, r: 3, k: 2}}
	cases = append(cases, buildCase(a0, qs0))
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
