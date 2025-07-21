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

var kVal int64
var memo map[int64]int

func grundy(x int64) int {
	if x < 2*kVal {
		return 0
	}
	if v, ok := memo[x]; ok {
		return v
	}
	lim := x / kVal
	for i := int64(1); i*i <= x; i++ {
		if x%i != 0 {
			continue
		}
		d1 := i
		if d1 >= 2 && d1 <= lim {
			if d1%2 == 0 || grundy(x/d1) == 0 {
				memo[x] = 1
				return 1
			}
		}
		d2 := x / i
		if d2 != d1 && d2 >= 2 && d2 <= lim {
			if d2%2 == 0 || grundy(x/d2) == 0 {
				memo[x] = 1
				return 1
			}
		}
	}
	memo[x] = 0
	return 0
}

func solve(n, m, k int64) string {
	if n%2 == 0 {
		return "Marsel"
	}
	kVal = k
	memo = make(map[int64]int)
	if grundy(m) != 0 {
		return "Timur"
	}
	return "Marsel"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := int64(rng.Intn(5) + 1)
	m := int64(rng.Intn(200) + 1)
	k := int64(rng.Intn(int(m)) + 1)
	input := fmt.Sprintf("%d %d %d\n", n, m, k)
	return input, solve(n, m, k)
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
		return fmt.Errorf("expected %q got %q", expected, got)
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
