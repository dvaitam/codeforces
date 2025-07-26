package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func feasible(b int, pos []int, k int) bool {
	maxDiff := b + 1
	prev := pos[0]
	used := 1
	i := 1
	last := pos[len(pos)-1]
	for {
		if last-prev <= maxDiff {
			used++
			return used <= k
		}
		j := i
		for j < len(pos) && pos[j]-prev <= maxDiff {
			j++
		}
		if j == i {
			return false
		}
		prev = pos[j-1]
		used++
		if used > k {
			return false
		}
		i = j
	}
}

func expected(n, k int, s string) int {
	pos := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			pos = append(pos, i+1)
		}
	}
	left, right := 0, n
	for left < right {
		mid := (left + right) / 2
		if feasible(mid, pos, k) {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return left
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

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(10) + 2
	if rng.Float64() < 0.1 {
		n = rng.Intn(50) + 2
	}
	k := rng.Intn(n-1) + 2
	sb := make([]byte, n)
	sb[0] = '0'
	sb[n-1] = '0'
	for i := 1; i < n-1; i++ {
		if rng.Intn(2) == 0 {
			sb[i] = '0'
		} else {
			sb[i] = '1'
		}
	}
	s := string(sb)
	input := fmt.Sprintf("%d %d\n%s\n", n, k, s)
	return input, expected(n, k, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierL.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
