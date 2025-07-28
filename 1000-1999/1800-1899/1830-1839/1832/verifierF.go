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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func zombies(n, k int, x, m int64, l, r []int64) int64 {
	maxStart := int(x - m)
	starts := make([]int, k)
	assign := make([]int, n)
	best := int64(-1)

	var calc func() // compute zombies for current starts and assignment
	calc = func() {
		var total int64
		for i := 0; i < n; i++ {
			s := int64(starts[assign[i]])
			for t := int64(0); t < x; t++ {
				defended := false
				if t >= l[i] && t < r[i] {
					defended = true
				}
				if !defended && t >= s && t < s+m {
					defended = true
				}
				if !defended {
					total++
				}
			}
		}
		if total > best {
			best = total
		}
	}

	var assignDFS func(int)
	assignDFS = func(i int) {
		if i == n {
			calc()
			return
		}
		for g := 0; g < k; g++ {
			assign[i] = g
			assignDFS(i + 1)
		}
	}

	var startDFS func(int)
	startDFS = func(pos int) {
		if pos == k {
			assignDFS(0)
			return
		}
		for t := 0; t <= maxStart; t++ {
			starts[pos] = t
			startDFS(pos + 1)
		}
	}

	startDFS(0)
	return best
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	k := rng.Intn(n) + 1
	x := int64(rng.Intn(5) + 1)
	m := int64(rng.Intn(int(x)) + 1)
	l := make([]int64, n)
	r := make([]int64, n)
	for i := 0; i < n; i++ {
		L := int64(rng.Intn(int(x)))
		R := int64(rng.Intn(int(x-L))) + L + 1
		l[i] = L
		r[i] = R
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, k, x, m))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", l[i], r[i]))
	}
	exp := zombies(n, k, x, m, l, r)
	return sb.String(), fmt.Sprintf("%d", exp)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
