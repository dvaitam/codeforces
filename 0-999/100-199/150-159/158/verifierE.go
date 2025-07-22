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

func expectedE(n, k int, tvals, dvals []int) int {
	if n == 0 || k >= n {
		return 86400
	}
	costPre := make([]int, n)
	for i := 0; i < n; i++ {
		cntGood := 0
		for j := 0; j < i; j++ {
			if tvals[j]+dvals[j] <= tvals[i] {
				cntGood++
			}
		}
		costPre[i] = i - cntGood
	}
	costPost := make([]int, n)
	for i := 0; i < n; i++ {
		cntBad := 0
		endI := tvals[i] + dvals[i]
		for j := i + 1; j < n; j++ {
			if tvals[j] < endI {
				cntBad++
			}
		}
		costPost[i] = cntBad
	}
	best := 0
	for e := 0; e < n; e++ {
		if e <= k {
			gap := tvals[e] - 1
			if gap > best {
				best = gap
			}
		}
	}
	for s := 0; s < n; s++ {
		if costPre[s]+costPost[s] <= k {
			gap := 86400 - (tvals[s] + dvals[s]) + 1
			if gap > best {
				best = gap
			}
		}
	}
	for s := 0; s < n; s++ {
		pre := costPre[s]
		if pre > k {
			continue
		}
		baseEnd := tvals[s] + dvals[s]
		maxMid := k - pre
		limit := s + 1 + maxMid
		if limit > n-1 {
			limit = n - 1
		}
		for e := s + 1; e <= limit; e++ {
			gap := tvals[e] - baseEnd
			if gap > best {
				best = gap
			}
		}
	}
	if best < 0 {
		best = 0
	}
	if best > 86400 {
		best = 86400
	}
	return best
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
	err := cmd.Run()
	return out.String(), err
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(10)
	k := 0
	if n > 0 {
		k = rng.Intn(n + 1)
	}
	tvals := make([]int, n)
	dvals := make([]int, n)
	cur := rng.Intn(86400-n+1) + 1
	for i := 0; i < n; i++ {
		tvals[i] = cur
		dvals[i] = rng.Intn(1000) + 1
		if i < n-1 {
			cur += rng.Intn((86400-cur)-(n-i-1)) + 1
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", tvals[i], dvals[i])
	}
	return sb.String(), expectedE(n, k, tvals, dvals)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, exp := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", t+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != fmt.Sprint(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", t+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
