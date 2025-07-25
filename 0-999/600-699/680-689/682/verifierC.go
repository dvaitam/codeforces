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

type edge struct {
	to int
	w  int64
}

type testCase struct {
	n int
	a []int64
	p []int
	c []int64
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return out.String(), nil
}

func solve(tc testCase) int {
	g := make([][]edge, tc.n+1)
	for i := 2; i <= tc.n; i++ {
		g[tc.p[i-2]] = append(g[tc.p[i-2]], edge{to: i, w: tc.c[i-2]})
	}
	type node struct {
		v       int
		dist    int64
		minPref int64
		rem     bool
	}
	stack := []node{{v: 1, dist: 0, minPref: 0, rem: false}}
	removed := 0
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if cur.rem {
			removed++
			for _, e := range g[cur.v] {
				nd := cur.dist + e.w
				mp := cur.minPref
				if nd < mp {
					mp = nd
				}
				stack = append(stack, node{v: e.to, dist: nd, minPref: mp, rem: true})
			}
			continue
		}
		if cur.dist-cur.minPref > tc.a[cur.v-1] {
			removed++
			for _, e := range g[cur.v] {
				nd := cur.dist + e.w
				mp := cur.minPref
				if nd < mp {
					mp = nd
				}
				stack = append(stack, node{v: e.to, dist: nd, minPref: mp, rem: true})
			}
			continue
		}
		for _, e := range g[cur.v] {
			nd := cur.dist + e.w
			mp := cur.minPref
			if nd < mp {
				mp = nd
			}
			stack = append(stack, node{v: e.to, dist: nd, minPref: mp, rem: false})
		}
	}
	return removed
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 0; i < tc.n; i++ {
		sb.WriteString(fmt.Sprintf("%d", tc.a[i]))
		if i+1 < tc.n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	for i := 2; i <= tc.n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.p[i-2], tc.c[i-2]))
	}
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	out = strings.TrimSpace(out)
	var got int
	if _, err := fmt.Sscan(out, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expected := solve(tc)
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	a := make([]int64, n)
	for i := range a {
		a[i] = rng.Int63n(1000) + 1
	}
	p := make([]int, n-1)
	c := make([]int64, n-1)
	for i := 2; i <= n; i++ {
		p[i-2] = rng.Intn(i-1) + 1
		c[i-2] = rng.Int63n(2001) - 1000
	}
	return testCase{n: n, a: a, p: p, c: c}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, testCase{n: 1, a: []int64{1}})
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, buildInput(tc))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
