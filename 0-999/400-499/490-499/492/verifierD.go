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

func runCandidate(bin, input string) (string, error) {
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

func solveCase(n int, x, y int64, a []int64) string {
	var sb strings.Builder
	for _, hits := range a {
		low, high := int64(0), hits*min(x, y)
		for low < high {
			mid := (low + high) / 2
			if mid/x+mid/y < hits {
				low = mid + 1
			} else {
				high = mid
			}
		}
		p := low
		vanya := p%y == 0
		vova := p%x == 0
		if vanya && vova {
			sb.WriteString("Both\n")
		} else if vanya {
			sb.WriteString("Vanya\n")
		} else {
			sb.WriteString("Vova\n")
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	x := int64(rng.Intn(10) + 1)
	y := int64(rng.Intn(10) + 1)
	a := make([]int64, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, x, y)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(20) + 1)
		fmt.Fprintf(&sb, "%d\n", a[i])
	}
	expect := solveCase(n, x, y, a)
	return sb.String(), expect
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
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
