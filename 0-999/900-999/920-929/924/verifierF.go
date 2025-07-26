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
	err := cmd.Run()
	if err != nil {
		return out.String() + errb.String(), err
	}
	return out.String(), nil
}

type Test struct {
	queries []query
	input   string
}

type query struct {
	l int64
	r int64
	k int
}

func genTest(rng *rand.Rand) Test {
	q := rng.Intn(5) + 1
	qs := make([]query, q)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		l := rng.Int63n(100) + 1
		r := l + rng.Int63n(20)
		k := rng.Intn(10)
		qs[i] = query{l: l, r: r, k: k}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", l, r, k))
	}
	return Test{queries: qs, input: sb.String()}
}

func isKBeautiful(x int64, k int) bool {
	if x == 0 {
		return k >= 0
	}
	digits := []int{}
	for x > 0 {
		digits = append(digits, int(x%10))
		x /= 10
	}
	sum := 0
	for _, d := range digits {
		sum += d
	}
	dp := make([]bool, sum+1)
	dp[0] = true
	for _, d := range digits {
		next := make([]bool, sum+1)
		for s := 0; s <= sum; s++ {
			if dp[s] {
				next[s] = true
				if s+d <= sum {
					next[s+d] = true
				}
			}
		}
		dp = next
	}
	for s := 0; s <= sum; s++ {
		if dp[s] {
			diff := sum - 2*s
			if diff < 0 {
				diff = -diff
			}
			if diff <= k {
				return true
			}
		}
	}
	return false
}

func solve(t Test) string {
	var sb strings.Builder
	for _, qu := range t.queries {
		count := int64(0)
		for x := qu.l; x <= qu.r; x++ {
			if isKBeautiful(x, qu.k) {
				count++
			}
		}
		sb.WriteString(fmt.Sprintf("%d\n", count))
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s got:%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
