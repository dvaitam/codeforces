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

func maxCost(n int) int {
	best := 0
	for r := 0; r < n; r++ {
		sum := 0
		maxProd := 0
		for i := 1; i <= r; i++ {
			prod := i * i
			sum += prod
			if prod > maxProd {
				maxProd = prod
			}
		}
		for j := 1; j <= n-r; j++ {
			val := n - j + 1
			idx := r + j
			prod := idx * val
			sum += prod
			if prod > maxProd {
				maxProd = prod
			}
		}
		cost := sum - maxProd
		if cost > best {
			best = cost
		}
	}
	return best
}

func generateCasesC() []int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]int, 0, 100)
	cases = append(cases, 2)
	for len(cases) < 100 {
		n := rng.Intn(50) + 2
		cases = append(cases, n)
	}
	return cases
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("1\n%d\n", n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("non-integer output %q", gotStr)
	}
	exp := maxCost(n)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCasesC()
	for i, n := range cases {
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
