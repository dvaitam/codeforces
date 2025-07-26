package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func solveB(n, m int, a []int) int {
	if n > m {
		return 0
	}
	freq := make(map[int]int)
	for _, v := range a {
		freq[v]++
	}
	counts := make([]int, 0, len(freq))
	for _, v := range freq {
		counts = append(counts, v)
	}
	sort.Slice(counts, func(i, j int) bool { return counts[i] > counts[j] })
	ans := 0
	for days := 1; days <= 100; days++ {
		total := 0
		for _, c := range counts {
			total += c / days
		}
		if total < n {
			ans = days - 1
			break
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(100) + 1
	m := rng.Intn(100) + 1
	a := make([]int, m)
	for i := range a {
		a[i] = rng.Intn(100) + 1
	}
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range a {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
	input.WriteByte('\n')
	expected := solveB(n, m, a)
	return input.String(), expected
}

func runCase(exe string, input string, expected int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(outStr)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
