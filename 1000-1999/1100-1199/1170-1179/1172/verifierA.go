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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveA(n int, a, b []int) int {
	const INF = int(1e9)
	avail := make([]int, n+1)
	for i := range avail {
		avail[i] = INF
	}
	for i := 0; i < n; i++ {
		if a[i] != 0 {
			avail[a[i]] = 0
		}
	}
	for i := 0; i < n; i++ {
		if b[i] != 0 && avail[b[i]] > i+1 {
			avail[b[i]] = i + 1
		}
	}
	tail := b[n-1]
	if tail != 0 {
		ok := true
		if tail <= n {
			for i := 0; i < tail; i++ {
				if b[n-tail+i] != i+1 {
					ok = false
					break
				}
			}
		} else {
			ok = false
		}
		if ok {
			valid := true
			for v := tail + 1; v <= n; v++ {
				if avail[v] > v-tail-1 {
					valid = false
					break
				}
			}
			if valid {
				return n - tail
			}
		}
	}
	maxDiff := 0
	for v := 1; v <= n; v++ {
		diff := avail[v] - v + 1
		if diff > maxDiff {
			maxDiff = diff
		}
	}
	if maxDiff < 0 {
		maxDiff = 0
	}
	return maxDiff + n
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	a := make([]int, n)
	b := make([]int, n)
	idxA := rng.Perm(n)
	idxB := rng.Perm(n)
	pa, pb := 0, 0
	for v := 1; v <= n; v++ {
		if rng.Intn(2) == 0 {
			if pa < n {
				a[idxA[pa]] = v
				pa++
			} else {
				b[idxB[pb]] = v
				pb++
			}
		} else {
			if pb < n {
				b[idxB[pb]] = v
				pb++
			} else {
				a[idxA[pa]] = v
				pa++
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	ans := solveA(n, a, b)
	return input, fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
