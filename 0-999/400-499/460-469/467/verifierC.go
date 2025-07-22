package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveC(r *bufio.Reader) string {
	var n, m, k int
	if _, err := fmt.Fscan(r, &n, &m, &k); err != nil {
		return ""
	}
	p := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(r, &p[i])
	}
	ps := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		ps[i] = ps[i-1] + p[i]
	}
	seg := make([]int64, n+1)
	for j := m; j <= n; j++ {
		seg[j] = ps[j] - ps[j-m]
	}
	const neg = int64(-1 << 63)
	dpPrev := make([]int64, n+1)
	dpCurr := make([]int64, n+1)
	for t := 1; t <= k; t++ {
		dpCurr[0] = neg
		for j := 1; j <= n; j++ {
			best := dpCurr[j-1]
			if j >= m {
				val := dpPrev[j-m] + seg[j]
				if val > best {
					best = val
				}
			}
			dpCurr[j] = best
		}
		dpPrev, dpCurr = dpCurr, dpPrev
	}
	return fmt.Sprint(dpPrev[n])
}

func generateCaseC(rng *rand.Rand) string {
	k := rng.Intn(3) + 1
	m := rng.Intn(5) + 1
	n := rng.Intn(20) + k*m
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i := 0; i < n; i++ {
		val := rng.Intn(20)
		fmt.Fprintf(&sb, "%d ", val)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseC(rng)
		expect := solveC(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
