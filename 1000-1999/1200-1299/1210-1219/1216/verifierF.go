package main

import (
	"bytes"
	"container/heap"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type interval struct {
	r   int
	val int64
}

type pq []interval

func (h pq) Len() int { return len(h) }
func (h pq) Less(i, j int) bool {
	if h[i].val == h[j].val {
		return h[i].r < h[j].r
	}
	return h[i].val < h[j].val
}
func (h pq) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *pq) Push(x interface{}) { *h = append(*h, x.(interval)) }
func (h *pq) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func solve(n, k int, s string) int64 {
	starts := make([][]int, n+2)
	for i := 1; i <= n; i++ {
		if s[i-1] == '1' {
			L := i - k
			if L < 1 {
				L = 1
			}
			starts[L] = append(starts[L], i)
		}
	}
	dp := make([]int64, n+1)
	active := &pq{}
	heap.Init(active)
	for i := 1; i <= n; i++ {
		for _, j := range starts[i] {
			r := j + k
			if r > n {
				r = n
			}
			val := dp[i-1] + int64(j)
			heap.Push(active, interval{r, val})
		}
		for active.Len() > 0 && (*active)[0].r < i {
			heap.Pop(active)
		}
		best := int64(1 << 60)
		if active.Len() > 0 {
			best = (*active)[0].val
		}
		direct := dp[i-1] + int64(i)
		if direct < best {
			dp[i] = direct
		} else {
			dp[i] = best
		}
	}
	return dp[n]
}

func genCase(rng *rand.Rand) (string, int, int, string) {
	n := rng.Intn(30) + 1
	k := rng.Intn(n) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(3) == 0 {
			b[i] = '1'
		} else {
			b[i] = '0'
		}
	}
	s := string(b)
	input := fmt.Sprintf("%d %d\n%s\n", n, k, s)
	return input, n, k, s
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, n, k, s := genCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		var ans int64
		if _, err := fmt.Sscan(out, &ans); err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad output: %v\ninput:\n%s\noutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
		exp := solve(n, k, s)
		if ans != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, ans, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
