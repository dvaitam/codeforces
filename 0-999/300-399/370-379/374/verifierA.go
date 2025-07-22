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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func expected(n, m, i, j, a, b int) string {
	if (i == 1 || i == n) && (j == 1 || j == m) {
		return "0"
	}
	if a > n-1 || b > m-1 {
		return "Poor Inna and pony!"
	}
	inf := int(1e9)
	ans := inf
	corners := [][2]int{{1, 1}, {1, m}, {n, 1}, {n, m}}
	for _, c := range corners {
		x, y := c[0], c[1]
		dx := abs(x - i)
		dy := abs(y - j)
		if dx%a != 0 || dy%b != 0 {
			continue
		}
		u := dx / a
		v := dy / b
		if u%2 != v%2 {
			continue
		}
		if k := max(u, v); k < ans {
			ans = k
		}
	}
	if ans == inf {
		return "Poor Inna and pony!"
	}
	return fmt.Sprintf("%d", ans)
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for caseNum := 0; caseNum < 100; caseNum++ {
		n := rng.Intn(20) + 1
		m := rng.Intn(20) + 1
		i := rng.Intn(n) + 1
		j := rng.Intn(m) + 1
		a := rng.Intn(20) + 1
		b := rng.Intn(20) + 1
		input := fmt.Sprintf("%d %d %d %d %d %d\n", n, m, i, j, a, b)
		exp := expected(n, m, i, j, a, b)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
