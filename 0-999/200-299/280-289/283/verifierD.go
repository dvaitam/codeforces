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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solveD(a []int64) int {
	n := len(a)
	bad := make([]bool, n)
	for i := 0; i < n-1; i++ {
		x := a[i]
		y := a[i+1]
		var d int64
		if y%2 == 1 {
			d = 0
		} else {
			d = y / 2
		}
		if (x-d)%y != 0 {
			bad[i] = true
		}
	}
	dp0 := make([]int, n+1)
	dp1 := make([]int, n+1)
	dp0[1] = 0
	dp1[1] = 1
	for i := 2; i <= n; i++ {
		if bad[i-2] {
			dp0[i] = max(dp0[i-1], dp1[i-1])
			dp1[i] = 1 + dp0[i-1]
		} else {
			mx := max(dp0[i-1], dp1[i-1])
			dp0[i] = mx
			dp1[i] = 1 + mx
		}
	}
	best := max(dp0[n], dp1[n])
	return n - best
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(20) + 2
	a := make([]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(30) + 1)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteString("\n")
	ans := solveD(a)
	return sb.String(), ans
}

func runCase(bin, input string, exp int) error {
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
	var val int
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &val); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if val != exp {
		return fmt.Errorf("expected %d got %d", exp, val)
	}
	return nil
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
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
