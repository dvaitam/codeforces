package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveCase(n, k int, s string) (int, string) {
	bytesArr := []byte(s)
	bestCost := math.MaxInt64
	bestResult := ""
	for target := byte('0'); target <= '9'; target++ {
		cnt := 0
		for i := 0; i < n; i++ {
			if bytesArr[i] == target {
				cnt++
			}
		}
		need := k - cnt
		cost := 0
		b := make([]byte, n)
		copy(b, bytesArr)
		if need > 0 {
			for d := 1; d <= 9 && need > 0; d++ {
				for i := 0; i < n && need > 0; i++ {
					if b[i] > target && int(b[i]-target) == d {
						cost += d
						b[i] = target
						need--
					}
				}
				for i := n - 1; i >= 0 && need > 0; i-- {
					if b[i] < target && int(target-b[i]) == d {
						cost += d
						b[i] = target
						need--
					}
				}
			}
		}
		result := string(b)
		if cost < bestCost || (cost == bestCost && result < bestResult) {
			bestCost = cost
			bestResult = result
		}
	}
	return bestCost, bestResult
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	digits := make([]byte, n)
	for i := 0; i < n; i++ {
		digits[i] = byte('0' + rng.Intn(10))
	}
	s := string(digits)
	sb.WriteString(s)
	sb.WriteByte('\n')
	cost, res := solveCase(n, k, s)
	input := sb.String()
	expected := fmt.Sprintf("%d\n%s", cost, res)
	return input, expected
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
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
