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

const mod = 100000000

func solveCase(n1, n2, k1, k2 int) int {
	dpF := make([][]int, n1+1)
	dpH := make([][]int, n1+1)
	for i := 0; i <= n1; i++ {
		dpF[i] = make([]int, n2+1)
		dpH[i] = make([]int, n2+1)
	}
	dpF[0][0] = 1
	dpH[0][0] = 1
	for i := 0; i <= n1; i++ {
		for j := 0; j <= n2; j++ {
			if i == 0 && j == 0 {
				continue
			}
			if i > 0 {
				sum := 0
				for x := 1; x <= k1 && x <= i; x++ {
					sum += dpH[i-x][j]
					if sum >= mod {
						sum -= mod
					}
				}
				dpF[i][j] = sum
			}
			if j > 0 {
				sum := 0
				for y := 1; y <= k2 && y <= j; y++ {
					sum += dpF[i][j-y]
					if sum >= mod {
						sum -= mod
					}
				}
				dpH[i][j] = sum
			}
		}
	}
	ans := dpF[n1][n2] + dpH[n1][n2]
	ans %= mod
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n1 := rng.Intn(10) + 1
	n2 := rng.Intn(10) + 1
	k1 := rng.Intn(10) + 1
	k2 := rng.Intn(10) + 1
	input := fmt.Sprintf("%d %d %d %d\n", n1, n2, k1, k2)
	ans := solveCase(n1, n2, k1, k2)
	return input, fmt.Sprintf("%d", ans)
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
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
