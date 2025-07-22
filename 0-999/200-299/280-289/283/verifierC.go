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

const mod = 1000000007

func solveC(n, q, t int, arr []int, pairs [][2]int) int {
	indegree := make([]int, n)
	out := make([]int, n)
	nxt := make([]int, n)
	for i := range nxt {
		nxt[i] = -1
	}
	prev := make([]int, n)
	for i := range prev {
		prev[i] = -1
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
	visited := make([]bool, n)
	weights := make([]int, 0, n)
	C := 0
	for i := 0; i < n; i++ {
		if indegree[i] == 0 {
			path := []int{}
			cur := i
			for cur != -1 {
				path = append(path, cur)
				visited[cur] = true
				cur = nxt[cur]
			}
			k := len(path)
			prefix := make([]int, k)
			sum := 0
			for j, v := range path {
				sum += arr[v]
				prefix[j] = sum
			}
			weights = append(weights, prefix...)
			for j, v := range path {
				C += (k - 1 - j) * arr[v]
			}
		}
	}
	for i := 0; i < n; i++ {
		if !visited[i] {
			return 0
		}
	}
	T := t - C
	if T < 0 {
		return 0
	}
	dp := make([]int, T+1)
	dp[0] = 1
	for _, w := range weights {
		for s := w; s <= T; s++ {
			dp[s] += dp[s-w]
			if dp[s] >= mod {
				dp[s] -= mod
			}
		}
	}
	return dp[T]
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(6) + 1
	q := rng.Intn(n + 1)
	t := rng.Intn(50) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(10) + 1
	}
	usedB := make(map[int]bool)
	usedC := make(map[int]bool)
	pairs := make([][2]int, 0, q)
	for len(pairs) < q {
		b := rng.Intn(n)
		c := rng.Intn(n)
		if b == c || usedB[b] || usedC[c] {
			continue
		}
		usedB[b] = true
		usedC[c] = true
		pairs = append(pairs, [2]int{b, c})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, q, t))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	for _, p := range pairs {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0]+1, p[1]+1))
	}
	ans := solveC(n, q, t, arr, pairs)
	return sb.String(), ans
}

func runCase(bin, input string, exp int) error {
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
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
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
