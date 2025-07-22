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

func expected(arr []int, t, c int) int {
	cur := 0
	ans := 0
	for _, x := range arr {
		if x <= t {
			cur++
		} else {
			if cur >= c {
				ans += cur - c + 1
			}
			cur = 0
		}
	}
	if cur >= c {
		ans += cur - c + 1
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(50) + 1
	t := rng.Intn(20) + 1
	c := rng.Intn(n) + 1
	arr := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, t, c))
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(30)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), expected(arr, t, c)
}

func runCase(exe, input string, exp int) error {
	cmd := exec.Command(exe)
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
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
