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

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(200) + 1
	d := rng.Int63n(1000000000)
	heights := make([]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, d))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		heights[i] = rng.Int63n(1000000000)
		sb.WriteString(fmt.Sprintf("%d", heights[i]))
	}
	sb.WriteByte('\n')
	var count int64
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if absInt64(heights[i]-heights[j]) <= d {
				count += 2
			}
		}
	}
	return sb.String(), fmt.Sprintf("%d", count)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
