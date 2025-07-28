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

const MaxB = 1000

var dist [MaxB + 1]int

func init() {
	const INF = int(1e9)
	for i := range dist {
		dist[i] = INF
	}
	dist[1] = 0
	q := []int{1}
	for head := 0; head < len(q); head++ {
		v := q[head]
		for x := 1; x <= v; x++ {
			nxt := v + v/x
			if nxt <= MaxB && dist[nxt] > dist[v]+1 {
				dist[nxt] = dist[v] + 1
				q = append(q, nxt)
			}
		}
	}
}

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
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveD(n, k int, bVals, cVals []int) int {
	w := make([]int, n)
	for i := 0; i < n; i++ {
		w[i] = dist[bVals[i]]
	}
	cap := min(k, 12*n)
	dp := make([]int, cap+1)
	for i := 0; i < n; i++ {
		cost := w[i]
		val := cVals[i]
		if cost > cap {
			continue
		}
		for j := cap; j >= cost; j-- {
			if dp[j-cost]+val > dp[j] {
				dp[j] = dp[j-cost] + val
			}
		}
	}
	return dp[min(k, cap)]
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(5) + 1
	k := r.Intn(20) + 1
	b := make([]int, n)
	c := make([]int, n)
	for i := 0; i < n; i++ {
		b[i] = r.Intn(1000) + 1
		c[i] = r.Intn(100) + 1
	}
	expect := fmt.Sprintf("%d", solveD(n, k, b, c))
	input := fmt.Sprintf("1\n%d %d\n", n, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", b[i])
	}
	input += "\n"
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", c[i])
	}
	input += "\n"
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
