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

func expected(owner []int) int {
	const inf = int(1e9)
	dp1, dp2, dp3 := 0, inf, inf
	for _, o := range owner {
		cost1, cost2, cost3 := 0, 0, 0
		if o != 1 {
			cost1 = 1
		}
		if o != 2 {
			cost2 = 1
		}
		if o != 3 {
			cost3 = 1
		}
		newDp1 := dp1 + cost1
		newDp2 := min(dp1, dp2) + cost2
		newDp3 := min(dp2, dp3) + cost3
		dp1, dp2, dp3 = newDp1, newDp2, newDp3
	}
	return dp3
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func generateCase(rng *rand.Rand) (string, []string) {
	n := rng.Intn(10) + 1
	k1 := rng.Intn(n + 1)
	k2 := rng.Intn(n - k1 + 1)
	k3 := n - k1 - k2
	owner := make([]int, n+1)
	nums := rng.Perm(n)
	idx := 0
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", k1, k2, k3))
	for i := 0; i < k1; i++ {
		owner[nums[idx]+1] = 1
		sb.WriteString(fmt.Sprintf("%d ", nums[idx]+1))
		idx++
	}
	if k1 > 0 {
		sb.WriteByte('\n')
	} else {
		sb.WriteByte('\n')
	}
	for i := 0; i < k2; i++ {
		owner[nums[idx]+1] = 2
		sb.WriteString(fmt.Sprintf("%d ", nums[idx]+1))
		idx++
	}
	if k2 > 0 {
		sb.WriteByte('\n')
	} else {
		sb.WriteByte('\n')
	}
	for i := 0; i < k3; i++ {
		owner[nums[idx]+1] = 3
		sb.WriteString(fmt.Sprintf("%d ", nums[idx]+1))
		idx++
	}
	if k3 > 0 {
		sb.WriteByte('\n')
	} else {
		sb.WriteByte('\n')
	}
	// trim trailing spaces and compute expected
	lines := strings.Split(sb.String(), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	input := strings.Join(lines, "\n") + "\n"
	exp := []string{fmt.Sprintf("%d", expected(owner[1:]))}
	return input, exp
}

func runCase(bin, input string, exp []string) error {
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
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(lines))
	}
	for i, line := range lines {
		if strings.TrimSpace(line) != exp[i] {
			return fmt.Errorf("line %d expected %s got %s", i+1, exp[i], strings.TrimSpace(line))
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
