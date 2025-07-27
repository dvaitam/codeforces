package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Test struct {
	arr []int64
}

func expected(arr []int64) int64 {
	a := make([]int64, len(arr))
	copy(a, arr)
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	var ans int64
	n := len(a)
	for i, v := range a {
		revenue := v * int64(n-i)
		if revenue > ans {
			ans = revenue
		}
	}
	return ans
}

func genCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(20) + 1
	arr := make([]int64, n)
	for i := range arr {
		// up to 1e12
		arr[i] = rng.Int63n(1_000_000_000_000) + 1
	}
	input := fmt.Sprintf("1\n%d\n", n)
	for i, v := range arr {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v)
	}
	input += "\n"
	return input, expected(arr)
}

func runCase(bin, input string, exp int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
