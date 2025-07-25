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

type testCaseD struct {
	abbr string
}

func solveD(abbr string) int64 {
	digits := strings.TrimSpace(abbr[4:])
	pow10 := make([]int64, 10)
	pow10[0] = 1
	for i := 1; i < 10; i++ {
		pow10[i] = pow10[i-1] * 10
	}
	start := make([]int64, 10)
	start[1] = 1989
	for i := 2; i < 10; i++ {
		start[i] = start[i-1] + pow10[i-1]
	}
	k := len(digits)
	val, _ := strconv.ParseInt(digits, 10, 64)
	year := val
	mod := pow10[k]
	for year < start[k] {
		year += mod
	}
	return year
}

func runCaseD(bin string, tc testCaseD) error {
	input := fmt.Sprintf("1\n%s\n", tc.abbr)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := solveD(tc.abbr)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func genCaseD(rng *rand.Rand) testCaseD {
	k := rng.Intn(9) + 1
	maxVal := 1
	for i := 0; i < k; i++ {
		maxVal *= 10
	}
	val := rng.Intn(maxVal)
	abbr := fmt.Sprintf("IAO'%0*d", k, val)
	return testCaseD{abbr}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseD(rng)
		if err := runCaseD(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
