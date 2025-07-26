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

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func expectedD(n int, val, left, right []int) int {
	child := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		if left[i] != -1 {
			child[left[i]] = true
		}
		if right[i] != -1 {
			child[right[i]] = true
		}
	}
	root := 1
	for i := 1; i <= n; i++ {
		if !child[i] {
			root = i
			break
		}
	}
	type item struct {
		idx       int
		low, high int64
	}
	stack := []item{{root, -1 << 62, 1 << 62}}
	found := make(map[int]bool)
	freq := make(map[int]int)
	for i := 1; i <= n; i++ {
		freq[val[i]]++
	}
	for len(stack) > 0 {
		it := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		v := val[it.idx]
		if int64(v) > it.low && int64(v) < it.high {
			found[v] = true
		}
		if left[it.idx] != -1 {
			nh := minInt64(it.high, int64(v))
			stack = append(stack, item{left[it.idx], it.low, nh})
		}
		if right[it.idx] != -1 {
			nl := maxInt64(it.low, int64(v))
			stack = append(stack, item{right[it.idx], nl, it.high})
		}
	}
	success := 0
	for v := range found {
		success += freq[v]
	}
	return n - success
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(7) + 1
	val := make([]int, n+1)
	left := make([]int, n+1)
	right := make([]int, n+1)
	for i := 1; i <= n; i++ {
		val[i] = rng.Intn(21)
		left[i] = -1
		right[i] = -1
	}
	for i := 2; i <= n; i++ {
		for {
			p := rng.Intn(i-1) + 1
			if rng.Intn(2) == 0 {
				if left[p] == -1 {
					left[p] = i
					break
				}
			} else {
				if right[p] == -1 {
					right[p] = i
					break
				}
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", val[i], left[i], right[i]))
	}
	expect := expectedD(n, val, left, right)
	return sb.String(), expect
}

func runCase(bin string, input string, expect int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	resStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(resStr)
	if err != nil {
		return fmt.Errorf("bad output %q", resStr)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
