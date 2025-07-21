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

func run(bin, input string) (string, error) {
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

func validAssignment(prices [3]int, ans [3]int) bool {
	used := [4]bool{}
	for i := 0; i < 3; i++ {
		if ans[i] < 1 || ans[i] > 3 || used[ans[i]] {
			return false
		}
		used[ans[i]] = true
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if prices[i] > prices[j] && ans[i] >= ans[j] {
				return false
			}
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		prices := [3]int{rng.Intn(100) + 1, rng.Intn(100) + 1, rng.Intn(100) + 1}
		input := fmt.Sprintf("%d %d %d\n", prices[0], prices[1], prices[2])
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		var ans [3]int
		n, _ := fmt.Sscan(out, &ans[0], &ans[1], &ans[2])
		if n != 3 || !validAssignment(prices, ans) {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid answer %s\ninput:%s", i+1, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
