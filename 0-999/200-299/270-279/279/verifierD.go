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

func solveD(a []int) string {
	n := len(a) - 1
	P := make([][][2]int, n+1)
	for k := 2; k <= n; k++ {
		for i := 1; i < k; i++ {
			for j := i; j < k; j++ {
				if a[i]+a[j] == a[k] {
					P[k] = append(P[k], [2]int{i, j})
				}
			}
		}
		if len(P[k]) == 0 {
			return "-1"
		}
	}
	var dfs func(int, []int, int) bool
	dfs = func(k int, lastUse []int, M int) bool {
		if k > n {
			return true
		}
		for _, pair := range P[k] {
			i, j := pair[0], pair[1]
			oldI, oldJ := lastUse[i], lastUse[j]
			lastUse[i], lastUse[j] = k, k
			cnt := 0
			for t := 1; t < k; t++ {
				if lastUse[t] >= k {
					cnt++
				}
			}
			if cnt <= M {
				if dfs(k+1, lastUse, M) {
					return true
				}
			}
			lastUse[i], lastUse[j] = oldI, oldJ
		}
		return false
	}
	for m := 1; m <= n; m++ {
		last := make([]int, n+1)
		if dfs(2, last, m) {
			return fmt.Sprintf("%d", m)
		}
	}
	return "-1"
}

func generateCaseD(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = rng.Intn(20) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	input := sb.String()
	expected := solveD(a)
	return input, expected
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
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseD(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
