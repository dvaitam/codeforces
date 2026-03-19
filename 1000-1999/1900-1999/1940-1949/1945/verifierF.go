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

type refPair struct {
	v   int
	pos int
}

// refSolve is the correct embedded reference solver for 1945F.
// Ported directly from the accepted CF solution.
func refSolve(input string) string {
	tokens := strings.Fields(input)
	idx := 0
	nextInt := func() int {
		v := 0
		s := tokens[idx]
		idx++
		neg := false
		start := 0
		if s[0] == '-' {
			neg = true
			start = 1
		}
		for i := start; i < len(s); i++ {
			v = v*10 + int(s[i]-'0')
		}
		if neg {
			v = -v
		}
		return v
	}

	t := nextInt()
	var results []string

	for ; t > 0; t-- {
		n := nextInt()

		values := make([]int, n+1)
		for i := 1; i <= n; i++ {
			values[i] = nextInt()
		}

		pairs := make([]refPair, n)
		for pos := 1; pos <= n; pos++ {
			pidx := nextInt()
			pairs[pos-1] = refPair{v: values[pidx], pos: pos}
		}

		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].v > pairs[j].v
		})

		active := make([]bool, n+2)
		total, h, small := 0, 0, 0
		bestProd := int64(-1)
		bestK := 0

		for i := 0; i < n; {
			val := pairs[i].v
			j := i
			for j < n && pairs[j].v == val {
				pos := pairs[j].pos
				active[pos] = true
				total++
				if pos <= h {
					small++
				}
				for total-small >= h+1 {
					h++
					if active[h] {
						small++
					}
				}
				j++
			}

			prod := int64(h) * int64(val)
			if prod > bestProd || (prod == bestProd && h < bestK) {
				bestProd = prod
				bestK = h
			}
			i = j
		}

		results = append(results, fmt.Sprintf("%d %d", bestProd, bestK))
	}
	return strings.Join(results, "\n")
}

func generatePerm(rng *rand.Rand, n int) []int {
	perm := rng.Perm(n)
	for i := range perm {
		perm[i]++
	}
	return perm
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	v := make([]int, n)
	for i := range v {
		v[i] = rng.Intn(100) + 1
	}
	p := generatePerm(rng, n)
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d\n", n)
	for i, val := range v {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", val)
	}
	sb.WriteByte('\n')
	for i, val := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", val)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(exe, input string) error {
	expected := strings.TrimSpace(refSolve(input))

	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		if err := runCase(exe, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
