package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func expected(arr []int64) int64 {
	pref := int64(0)
	maxPref := int64(0)
	minPref := int64(0)
	for _, v := range arr {
		pref += v
		if pref > maxPref {
			maxPref = pref
		}
		if pref < minPref {
			minPref = pref
		}
	}
	return maxPref - minPref
}

func generateTests() [][]int64 {
	rand.Seed(4)
	t := 100
	res := make([][]int64, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(8) + 1
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = int64(rand.Intn(21) - 10)
		}
		res[i] = arr
	}
	return res
}

func verifyCase(bin string, arr []int64) error {
	exp := expected(arr)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("execution error: %v", err)
	}
	outStr := strings.TrimSpace(string(out))
	got, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid output %q", outStr)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, arr := range tests {
		if err := verifyCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
