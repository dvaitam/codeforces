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

func maxSubarray(a []int64) int64 {
	best := a[0]
	for i := 0; i < len(a); i++ {
		sum := int64(0)
		for j := i; j < len(a); j++ {
			sum += a[j]
			if sum > best {
				best = sum
			}
		}
	}
	if best < 0 {
		return 0
	}
	return best
}

func solve(n, m, c int, x []int64, p []int64, queries [][2]int) string {
	if n == 1 {
		return "0.000000000000"
	}
	arr := make([]int64, n)
	for i := 1; i < n; i++ {
		diff := x[i+1] - x[i]
		arr[i] = (diff-int64(2*c))*p[i] + diff*(100-p[i])
	}
	var ans int64
	for _, q := range queries {
		l, r := q[0], q[1]
		if l < r {
			sub := arr[l:r]
			ans += maxSubarray(sub)
		}
	}
	return fmt.Sprintf("%.12f", float64(ans)/200.0)
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 2
		m := rng.Intn(5) + 1
		c := rng.Intn(10) + 1
		x := make([]int64, n+1)
		for j := 1; j <= n; j++ {
			if j == 1 {
				x[j] = 0
			} else {
				x[j] = x[j-1] + int64(rng.Intn(10)+1)
			}
		}
		p := make([]int64, n+1)
		for j := 1; j < n; j++ {
			p[j] = int64(rng.Intn(101))
		}
		queries := make([][2]int, m)
		for j := 0; j < m; j++ {
			l := rng.Intn(n-1) + 1
			r := rng.Intn(n-l) + l + 1
			queries[j] = [2]int{l, r}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, c))
		for j := 1; j <= n; j++ {
			sb.WriteString(fmt.Sprintf("%d ", x[j]))
		}
		sb.WriteString("\n")
		for j := 1; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d ", p[j]))
		}
		sb.WriteString("\n")
		for j := 0; j < m; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", queries[j][0], queries[j][1]))
		}
		input := sb.String()
		expected := solve(n, m, c, x, p, queries)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(expected) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
