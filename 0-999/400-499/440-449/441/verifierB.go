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

type pair struct{ a, b int }

type testCase struct {
	n     int
	v     int
	pairs []pair
}

func (tc testCase) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.v))
	for _, p := range tc.pairs {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.a, p.b))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	v := rng.Intn(10) + 1
	pairs := make([]pair, n)
	for i := 0; i < n; i++ {
		a := rng.Intn(10) + 1
		b := rng.Intn(10)
		pairs[i] = pair{a, b}
	}
	return testCase{n, v, pairs}
}

func expected(tc testCase) int {
	maxA := 0
	for _, p := range tc.pairs {
		if p.a > maxA {
			maxA = p.a
		}
	}
	fruits := make([]int, maxA+2)
	for _, p := range tc.pairs {
		fruits[p.a] += p.b
	}
	prev := 0
	total := 0
	for day := 1; day <= maxA+1; day++ {
		takePrev := prev
		if takePrev > tc.v {
			takePrev = tc.v
		}
		total += takePrev
		cap := tc.v - takePrev
		takeCur := fruits[day]
		if takeCur > cap {
			takeCur = cap
		}
		total += takeCur
		prev = fruits[day] - takeCur
	}
	return total
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	input := tc.Input()
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expect := expected(tc)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{
		{n: 1, v: 1, pairs: []pair{{1, 1}}},
		{n: 2, v: 2, pairs: []pair{{1, 1}, {1, 2}}},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.Input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
