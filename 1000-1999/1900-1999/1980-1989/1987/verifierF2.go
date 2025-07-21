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

func solve(arr []int) int {
	memo := make(map[string]int)
	var dfs func([]int) int
	dfs = func(a []int) int {
		if len(a) < 2 {
			return 0
		}
		key := fmt.Sprint(a)
		if v, ok := memo[key]; ok {
			return v
		}
		best := 0
		for i := 0; i < len(a)-1; i++ {
			if a[i] == i+1 {
				b := append([]int{}, a[:i]...)
				b = append(b, a[i+2:]...)
				if val := 1 + dfs(b); val > best {
					best = val
				}
			}
		}
		memo[key] = best
		return best
	}
	return dfs(arr)
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(8) + 2
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	cp := append([]int(nil), arr...)
	return sb.String(), solve(cp)
}

func runCase(bin string, input string, expected int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
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
