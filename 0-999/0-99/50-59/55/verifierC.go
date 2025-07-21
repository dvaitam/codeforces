package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func expected(n, m int, pies [][2]int) string {
	INF := n + m + 5
	T := make([]int, 0, 2*(n+m))
	for j := 1; j <= m; j++ {
		best := INF
		for _, p := range pies {
			d := p[0] + abs(p[1]-j)
			if d < best {
				best = d
			}
		}
		T = append(T, best)
	}
	for j := 1; j <= m; j++ {
		best := INF
		for _, p := range pies {
			d := (n - p[0] + 1) + abs(p[1]-j)
			if d < best {
				best = d
			}
		}
		T = append(T, best)
	}
	for i := 1; i <= n; i++ {
		best := INF
		for _, p := range pies {
			d := p[1] + abs(p[0]-i)
			if d < best {
				best = d
			}
		}
		T = append(T, best)
	}
	for i := 1; i <= n; i++ {
		best := INF
		for _, p := range pies {
			d := (m - p[1] + 1) + abs(p[0]-i)
			if d < best {
				best = d
			}
		}
		T = append(T, best)
	}
	sort.Ints(T)
	ans := "NO"
	for idx, t := range T {
		if t <= idx+1 {
			ans = "YES"
			break
		}
	}
	return ans
}

func generateCase() (string, string) {
	n := rand.Intn(100) + 1
	m := rand.Intn(100) + 1
	k := rand.Intn(101)
	pies := make([][2]int, k)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i := 0; i < k; i++ {
		x := rand.Intn(n) + 1
		y := rand.Intn(m) + 1
		pies[i] = [2]int{x, y}
		fmt.Fprintf(&sb, "%d %d\n", x, y)
	}
	return sb.String(), expected(n, m, pies)
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input, exp := generateCase()
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
