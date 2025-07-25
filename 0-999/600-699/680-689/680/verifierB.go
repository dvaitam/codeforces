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

func runProg(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(n, a int, arr []int) int {
	ans := 0
	maxDist := a - 1
	if n-a > maxDist {
		maxDist = n - a
	}
	for d := 0; d <= maxDist; d++ {
		l := a - d
		r := a + d
		if l >= 1 && r <= n {
			if l == r {
				ans += arr[l-1]
			} else if arr[l-1] == 1 && arr[r-1] == 1 {
				ans += 2
			}
		} else if l >= 1 && l <= n && (r < 1 || r > n) {
			ans += arr[l-1]
		} else if r >= 1 && r <= n && (l < 1 || l > n) {
			ans += arr[r-1]
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(100) + 1
	a := rng.Intn(n) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, a)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(2))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseInput(input string) (int, int, []int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) < 2 {
		return 0, 0, nil, fmt.Errorf("invalid input")
	}
	var n, a int
	fmt.Sscanf(lines[0], "%d %d", &n, &a)
	arr := make([]int, n)
	fields := strings.Fields(lines[1])
	if len(fields) != n {
		return 0, 0, nil, fmt.Errorf("invalid input")
	}
	for i := 0; i < n; i++ {
		fmt.Sscan(fields[i], &arr[i])
	}
	return n, a, arr, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		input := generateCase(rng)
		n, a, arr, err := parseInput(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v", err)
			os.Exit(1)
		}
		expected := fmt.Sprint(solve(n, a, arr))
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
