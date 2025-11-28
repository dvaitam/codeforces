package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const mod = 1000000007

func solveC(n, q, t int, arr []int, pairs [][2]int) int {
	indegree := make([]int, n+1)
	out := make([]int, n+1)
	nxt := make([]int, n+1)
	for i := range nxt {
		nxt[i] = 0
	}
	prev := make([]int, n+1)
	for i := range prev {
		prev[i] = 0
	}
	for _, p := range pairs {
		b := p[0]
		c := p[1]
		out[b]++
		indegree[c]++
		if out[b] > 1 || indegree[c] > 1 {
			return 0
		}
		nxt[b] = c
		prev[c] = b
	}
	visited := make([]bool, n+1)
	var weights []int
	var totalOffset int64
	
	for i := 1; i <= n; i++ {
		if indegree[i] == 0 && !visited[i] {
			path := []int{}
			curr := i
			for curr != 0 && curr <= n && !visited[curr] {
				visited[curr] = true
				path = append(path, curr)
				curr = nxt[curr]
			}
            
			k := len(path)
			
			for j := 0; j < k; j++ {
				idx := path[j]
				totalOffset += int64(arr[idx]) * int64(j)
			}

			suffixSum := 0
			for j := k - 1; j >= 0; j-- {
				suffixSum += arr[path[j]]
				weights = append(weights, suffixSum)
			}
		}
	}

	for i := 1; i <= n; i++ {
		if !visited[i] {
			return 0
		}
	}

	target := t - int(totalOffset)
	if target < 0 {
		return 0
	}
	
	dp := make([]int, target+1)
	dp[0] = 1

	for _, w := range weights {
		for s := w; s <= target; s++ {
			dp[s] = (dp[s] + dp[s-w]) % mod
			if dp[s] >= mod {
				dp[s] -= mod
			}
		}
	}
	return dp[target]
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(298) + 2 // N between 2 and 300
	q := rng.Intn(n + 1)
	t := rng.Intn(100000) + 1 // T between 1 and 100,000
	arr := make([]int, n+1) // 1-based indexing for coins
	for i := 1; i <= n; i++ {
		arr[i] = rng.Intn(100000) + 1 // coin values up to 10^5
	}
	usedB := make(map[int]bool)
	usedC := make(map[int]bool)
	pairs := make([][2]int, 0, q)
	for len(pairs) < q {
		b := rng.Intn(n) + 1 // 1-based node ID
		c := rng.Intn(n) + 1 // 1-based node ID
		if b == c || usedB[b] || usedC[c] {
			continue
		}
		usedB[b] = true
		usedC[c] = true
		pairs = append(pairs, [2]int{b, c})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, q, t))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteString(" ")
		}
		sb.WriteString(strconv.Itoa(arr[i]))
	}
	sb.WriteString("\n")
	for _, p := range pairs {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	ans := solveC(n, q, t, arr, pairs)
	return sb.String(), ans
}

func runCase(bin, input string, exp int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second) // Increased timeout
	defer cancel()

	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout (1s): potential infinite loop or too slow")
		}
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var val int
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &val); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if val != exp {
		return fmt.Errorf("expected %d got %d", exp, val)
	}
	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}