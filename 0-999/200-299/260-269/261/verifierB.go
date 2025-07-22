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

const LMT = 55

func factorials() [LMT]float64 {
	var fi [LMT]float64
	fi[0] = 1
	for i := 1; i < LMT; i++ {
		fi[i] = float64(i) * fi[i-1]
	}
	return fi
}

func expectedB(n int, arr []int, p int) float64 {
	var dp [LMT][LMT]float64
	dp[0][0] = 1
	for i := 0; i < n; i++ {
		ai := arr[i]
		for j := n; j >= 1; j-- {
			for t := p; t >= ai; t-- {
				dp[j][t] += dp[j-1][t-ai]
			}
		}
	}
	fi := factorials()
	var ans float64
	for i := 1; i <= n; i++ {
		for t := 1; t <= p; t++ {
			ans += dp[i][t] * fi[i] * fi[n-i]
		}
	}
	return ans / fi[n]
}

func genCaseB(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(10) + 1
	}
	p := rng.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", p)
	exp := fmt.Sprintf("%.6f\n", expectedB(n, arr, p))
	return sb.String(), exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	// we compare as floats with tolerance
	gv, err := strconv.ParseFloat(got, 64)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	ev, _ := strconv.ParseFloat(strings.TrimSpace(expected), 64)
	if diff := gv - ev; diff < -1e-4 || diff > 1e-4 {
		return fmt.Errorf("expected %.6f got %.6f", ev, gv)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseB(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
