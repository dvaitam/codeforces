package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func solveCase(n, m, k int64) string {
	var sb strings.Builder
	stBig := int64(1)
	for i := int64(0); i < k; i++ {
		temp := m
		big := n % m
		stSmall := stBig
		for j := int64(0); j < temp; j++ {
			if big > 0 {
				size := (n + m - 1) / m
				fmt.Fprintf(&sb, "%d ", size)
				eles := size
				for eles > 0 {
					fmt.Fprintf(&sb, "%d ", stBig)
					stBig++
					if stBig > n {
						stBig = 1
					}
					stSmall = stBig
					eles--
				}
				sb.WriteByte('\n')
				big--
			} else {
				size := n / m
				fmt.Fprintf(&sb, "%d ", size)
				eles := size
				for eles > 0 {
					fmt.Fprintf(&sb, "%d ", stSmall)
					stSmall++
					if stSmall > n {
						stSmall = 1
					}
					eles--
				}
				sb.WriteByte('\n')
			}
		}
	}
	sb.WriteByte('\n')
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) (string, string) {
	n := int64(rng.Intn(20) + 2)
	m := int64(rng.Intn(int(n/2)) + 1)
	if m == 0 {
		m = 1
	}
	if n < 2*m {
		n = 2 * m
	}
	k := int64(rng.Intn(5) + 1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d\n", n, m, k)
	return sb.String(), solveCase(n, m, k)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
