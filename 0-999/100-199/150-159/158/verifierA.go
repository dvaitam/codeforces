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

func expectedA(n, k int, scores []int) int {
	threshold := scores[k-1]
	cnt := 0
	for _, s := range scores {
		if s >= threshold && s > 0 {
			cnt++
		}
	}
	return cnt
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(50) + 1
		k := rng.Intn(n) + 1
		scores := make([]int, n)
		scores[0] = rng.Intn(101)
		for i := 1; i < n; i++ {
			scores[i] = rng.Intn(scores[i-1] + 1)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for i, s := range scores {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", s)
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := expectedA(n, k, scores)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", t+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != fmt.Sprint(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", t+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
