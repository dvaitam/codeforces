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

func compatible(a, b int) bool { return a&b == 0 }

func expected(arr []int) []int {
	n := len(arr)
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = -1
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			if compatible(arr[i], arr[j]) {
				res[i] = arr[j]
				break
			}
		}
	}
	return res
}

func genCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(256)
	}
	var b strings.Builder
	fmt.Fprintln(&b, n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "%d ", arr[i])
	}
	fmt.Fprintln(&b)
	return b.String(), expected(arr)
}

func runCase(bin, input string, exp []int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != len(exp) {
		return fmt.Errorf("expected %v got %s", exp, out.String())
	}
	for i, f := range fields {
		var x int
		fmt.Sscan(f, &x)
		if x != exp[i] {
			return fmt.Errorf("expected %v got %v", exp, fields)
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
		input, exp := genCase(rng)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
