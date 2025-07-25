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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func subset(weights []int, target int) bool {
	if target < 0 {
		return false
	}
	dp := make([]bool, target+1)
	dp[0] = true
	for _, w := range weights {
		for j := target; j >= w; j-- {
			if dp[j-w] {
				dp[j] = true
			}
		}
	}
	return dp[target]
}

func solveE(n int, a, b, c, d int, parent []int) string {
	g := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := parent[i-2]
		g[i] = append(g[i], p)
		g[p] = append(g[p], i)
	}
	belong := make([]int, n+1)
	leafCnt := make([]int, n+1)
	var dfs func(int, int, int) int
	dfs = func(v, p, root int) int {
		belong[v] = root
		cnt := 0
		isLeaf := true
		for _, to := range g[v] {
			if to == p {
				continue
			}
			isLeaf = false
			cnt += dfs(to, v, root)
		}
		if isLeaf {
			cnt = 1
		}
		leafCnt[v] = cnt
		return cnt
	}
	total := 0
	for _, ch := range g[1] {
		total += dfs(ch, 1, ch)
	}
	if total%2 == 1 {
		return "NO"
	}
	half := total / 2
	childA := belong[a]
	childC := belong[c]
	childD := belong[d]
	subLeaves := map[int]int{}
	for _, ch := range g[1] {
		subLeaves[ch] = leafCnt[ch]
	}
	check := func(x, y int) bool {
		fixed := subLeaves[x] + subLeaves[y]
		target := half - fixed
		others := []int{}
		for _, ch := range g[1] {
			if ch == x || ch == y {
				continue
			}
			others = append(others, subLeaves[ch])
		}
		return subset(others, target)
	}
	if check(childA, childC) || check(childA, childD) {
		return "YES"
	}
	return "NO"
}

func genCase(rng *rand.Rand) (int, int, int, int, int, []int) {
	n := rng.Intn(8) + 5
	parent := make([]int, n-1)
	for i := 2; i <= n; i++ {
		parent[i-2] = rng.Intn(i-1) + 1
	}
	leaves := []int{}
	deg := make([]int, n+1)
	for i := 2; i <= n; i++ {
		deg[i]++
		deg[parent[i-2]]++
	}
	for i := 2; i <= n; i++ {
		if deg[i] == 1 {
			leaves = append(leaves, i)
		}
	}
	if len(leaves) < 4 {
		for len(leaves) < 4 {
			leaves = append(leaves, 2)
		}
	}
	shuffle := rng.Perm(len(leaves))
	a := leaves[shuffle[0]]
	b := leaves[shuffle[1]]
	c := leaves[shuffle[2]]
	d := leaves[shuffle[3]]
	return n, a, b, c, d, parent
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, a, b, c, d, parent := genCase(rng)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		fmt.Fprintf(&sb, "%d %d %d %d\n", a, b, c, d)
		for i := 2; i <= n; i++ {
			fmt.Fprintf(&sb, "%d ", parent[i-2])
		}
		sb.WriteByte('\n')
		expect := solveE(n, a, b, c, d, parent)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", i+1, sb.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
