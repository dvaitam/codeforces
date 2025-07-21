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

type caseA struct {
	n   int
	arr []int
}

func generateCase(rng *rand.Rand) caseA {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(1000) + 1
	}
	return caseA{n, arr}
}

func solveCase(tc caseA) (int, int) {
	counts := make(map[int]int)
	maxH := 0
	for _, v := range tc.arr {
		counts[v]++
		if counts[v] > maxH {
			maxH = counts[v]
		}
	}
	return maxH, len(counts)
}

func runCase(bin string, tc caseA) error {
	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", tc.n)
	for i, v := range tc.arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprintf(&input, "%d", v)
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) < 2 {
		return fmt.Errorf("output too short: %s", out.String())
	}
	got1, err1 := strconv.Atoi(fields[0])
	got2, err2 := strconv.Atoi(fields[1])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("invalid integers in output: %s", out.String())
	}
	exp1, exp2 := solveCase(tc)
	if got1 != exp1 || got2 != exp2 {
		return fmt.Errorf("expected %d %d got %d %d", exp1, exp2, got1, got2)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput n=%d arr=%v\n", i+1, err, tc.n, tc.arr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
