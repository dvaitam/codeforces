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

func expected(n, m int, grid []string) string {
	ans := make([]int, m)
	for i := 1; i <= n; i++ {
		line := grid[i-1]
		for j := 1; j <= m; j++ {
			switch line[j-1] {
			case 'L':
				target := j - i + 1
				if target >= 1 && target <= m {
					ans[target-1]++
				}
			case 'R':
				target := j + i - 1
				if target >= 1 && target <= m {
					ans[target-1]++
				}
			case 'U':
				if i%2 == 1 {
					ans[j-1]++
				}
			}
		}
	}
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func runCase(bin, input, expected string) error {
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
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		grid := make([]string, n)
		k := 0
		for r := 0; r < n; r++ {
			lineBytes := make([]byte, m)
			for c := 0; c < m; c++ {
				if r == 0 {
					lineBytes[c] = '.'
					continue
				}
				ch := ".LRU"[rng.Intn(4)]
				lineBytes[c] = ch
				if ch != '.' {
					k++
				}
			}
			grid[r] = string(lineBytes)
		}
		input := fmt.Sprintf("%d %d %d\n%s\n", n, m, k, strings.Join(grid, "\n"))
		exp := expected(n, m, grid)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
