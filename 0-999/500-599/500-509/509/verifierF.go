package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const MOD = 1000000007

func solveTree(b []int) string {
	n := len(b) - 1
	if n <= 1 {
		return "1"
	}
	m := n - 1
	breakChild := make([]bool, m)
	for i := 1; i < m; i++ {
		if b[i+2] < b[i+1] {
			breakChild[i] = true
		}
	}
	dpPrev := make([]int, m+2)
	dpPrev[1] = 1
	for i := 1; i < m; i++ {
		pre := make([]int, i+2)
		for k := 1; k <= i; k++ {
			pre[k] = pre[k-1] + dpPrev[k]
			if pre[k] >= MOD {
				pre[k] -= MOD
			}
		}
		dpCurr := make([]int, m+2)
		for k := 1; k <= i+1; k++ {
			if breakChild[i] {
				if k-1 >= 1 {
					dpCurr[k] = pre[k-1]
				}
			} else {
				if k > i {
					dpCurr[k] = pre[i]
				} else {
					dpCurr[k] = pre[k]
				}
			}
		}
		dpPrev = dpCurr
	}
	ans := 0
	for k := 1; k <= m; k++ {
		ans += dpPrev[k]
		if ans >= MOD {
			ans -= MOD
		}
	}
	return strconv.Itoa(ans)
}

func generateCase(rng *rand.Rand) []int {
	n := rng.Intn(6) + 1
	b := make([]int, n+1)
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		b[i+1] = perm[i] + 1
	}
	b[0] = 0
	if b[1] != 1 {
		idx := 1
		for i := 1; i <= n; i++ {
			if b[i] == 1 {
				idx = i
				break
			}
		}
		b[1], b[idx] = b[idx], b[1]
	}
	return b
}

func runCase(bin string, b []int) error {
	var sb strings.Builder
	n := len(b) - 1
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(b[i]))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := solveTree(b)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		b := generateCase(rng)
		if err := runCase(bin, b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
