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

func solve(n, x, y int, arr []int) int {
	if x > y {
		return n
	}
	cnt := 0
	for _, v := range arr {
		if v <= x {
			cnt++
		}
	}
	if cnt%2 == 0 {
		return cnt / 2
	}
	return cnt/2 + 1
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	x := rng.Intn(100) + 1
	y := rng.Intn(100) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(100) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, x, y))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	ans := solve(n, x, y, arr)
	return sb.String(), fmt.Sprintf("%d", ans)
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
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
