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

func bestScore(n int, k int64, start int, p []int, a []int64) int64 {
	pos := start - 1
	best := int64(k) * a[pos]
	prefix := int64(0)
	limit := int64(n)
	if k < limit {
		limit = k
	}
	for i := int64(1); i < limit; i++ {
		prefix += a[pos]
		pos = p[pos] - 1
		cand := prefix + (int64(k)-i)*a[pos]
		if cand > best {
			best = cand
		}
	}
	return best
}

func solve(n int, k int64, pb, ps int, p []int, a []int64) string {
	sb := bestScore(n, k, pb, p, a)
	ss := bestScore(n, k, ps, p, a)
	if sb > ss {
		return "Bodya"
	} else if ss > sb {
		return "Sasha"
	}
	return "Draw"
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	k := int64(rng.Intn(5) + 1)
	pb := rng.Intn(n) + 1
	ps := rng.Intn(n) + 1
	p := rng.Perm(n)
	for i := range p {
		p[i]++
	}
	aVals := make([]int64, n)
	for i := range aVals {
		aVals[i] = int64(rng.Intn(10) + 1)
	}
	var sbuilder strings.Builder
	sbuilder.WriteString("1\n")
	fmt.Fprintf(&sbuilder, "%d %d %d %d\n", n, k, pb, ps)
	for i := 0; i < n; i++ {
		if i > 0 {
			sbuilder.WriteByte(' ')
		}
		sbuilder.WriteString(fmt.Sprint(p[i]))
	}
	sbuilder.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sbuilder.WriteByte(' ')
		}
		sbuilder.WriteString(fmt.Sprint(aVals[i]))
	}
	sbuilder.WriteByte('\n')
	expect := solve(n, k, pb, ps, p, aVals)
	return sbuilder.String(), expect
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		inp, exp := genCase(rng)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\ninput:\n%s", i+1, exp, got, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
