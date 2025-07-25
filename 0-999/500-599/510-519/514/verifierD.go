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

func expected(n, m int, k int64, a [][]int64) string {
	bestLen := 0
	best := make([]int64, m)
	for l := 0; l < n; l++ {
		maxCol := make([]int64, m)
		for r := l; r < n; r++ {
			for j := 0; j < m; j++ {
				if a[r][j] > maxCol[j] {
					maxCol[j] = a[r][j]
				}
			}
			sum := int64(0)
			for j := 0; j < m; j++ {
				sum += maxCol[j]
			}
			if sum <= k {
				if r-l+1 > bestLen {
					bestLen = r - l + 1
					copy(best, maxCol)
				}
			}
		}
	}
	var sb strings.Builder
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", best[i]))
	}
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	k := int64(rng.Intn(1000))
	a := make([][]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		a[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			v := int64(rng.Intn(1000))
			a[i][j] = v
			sb.WriteString(fmt.Sprintf("%d", v))
			if j+1 < m {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	expect := expected(n, m, k, a)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// simple edge case
	in := "1 1 5\n3\n"
	exp := "3"
	out, err := runCandidate(bin, in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "edge case failed: %v\n", err)
		os.Exit(1)
	}
	if out != exp {
		fmt.Fprintf(os.Stderr, "edge case failed: expected %s got %s\n", exp, out)
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCandidate(bin, in)
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
