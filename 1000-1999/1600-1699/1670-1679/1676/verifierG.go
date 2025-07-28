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

type Tree struct {
	children [][]int
	colors   []byte
}

func countBalanced(p []int, colors string) int {
	n := len(colors)
	tree := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		tree[p[i-2]] = append(tree[p[i-2]], i)
	}
	bytesColors := []byte(colors)
	var ans int
	var dfs func(int) (int, int)
	dfs = func(v int) (int, int) {
		white, black := 0, 0
		if bytesColors[v-1] == 'W' {
			white = 1
		} else {
			black = 1
		}
		for _, u := range tree[v] {
			w, b := dfs(u)
			white += w
			black += b
		}
		if white == black {
			ans++
		}
		return white, black
	}
	dfs(1)
	return ans
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	parents := make([]int, n-1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", n)
	for i := 2; i <= n; i++ {
		parents[i-2] = rng.Intn(i-1) + 1
		if i > 2 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", parents[i-2])
	}
	sb.WriteByte('\n')
	var colorsSB strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			colorsSB.WriteByte('B')
		} else {
			colorsSB.WriteByte('W')
		}
	}
	colors := colorsSB.String()
	sb.WriteString(colors)
	sb.WriteByte('\n')
	ans := countBalanced(parents, colors)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
