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

func check(a float64, arr []float64, k float64) bool {
	b := 0.0
	for _, v := range arr {
		if v > a {
			b += (v - a) / 100 * (100 - k)
		} else {
			b -= a - v
		}
	}
	return b > 0
}

func expectedValue(arr []float64, k float64) float64 {
	low, high := 0.0, 1000.0
	for it := 0; it < 100; it++ {
		mid := (low + high) / 2
		if check(mid, arr, k) {
			low = mid
		} else {
			high = mid
		}
	}
	return (low + high) / 2
}

func generateCase(rng *rand.Rand) (string, float64) {
	n := rng.Intn(20) + 1
	k := float64(rng.Intn(100))
	arr := make([]float64, n)
	for i := 0; i < n; i++ {
		arr[i] = float64(rng.Intn(1001))
	}
	exp := expectedValue(arr, k)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, int(k)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%.0f", v))
	}
	sb.WriteByte('\n')
	return sb.String(), exp
}

func runCase(bin, input string, exp float64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if diff := got - exp; diff < -1e-6 || diff > 1e-6 {
		return fmt.Errorf("expected %.6f got %.6f", exp, got)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
