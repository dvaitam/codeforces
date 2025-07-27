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

func spfSieve(maxN int) []int {
	spf := make([]int, maxN+1)
	for i := 0; i <= maxN; i++ {
		spf[i] = i
	}
	for i := 2; i*i <= maxN; i++ {
		for j := i * 2; j <= maxN; j += i {
			if spf[j] > i {
				spf[j] = i
			}
		}
	}
	return spf
}

func solveCase(arr []int, spf []int) (int, []int) {
	mp := map[int]int{}
	id := 1
	colors := make([]int, len(arr))
	for i, v := range arr {
		p := spf[v]
		if mp[p] == 0 {
			mp[p] = id
			id++
		}
		colors[i] = mp[p]
	}
	return len(mp), colors
}

func generateCase(rng *rand.Rand) []int {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(999) + 2
	}
	return arr
}

func runCase(bin string, arr []int, spf []int) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
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
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("expected at least 2 lines got %d", len(lines))
	}
	expectedCnt, expectedColors := solveCase(arr, spf)
	if strings.TrimSpace(lines[0]) != strconv.Itoa(expectedCnt) {
		return fmt.Errorf("expected %d got %s", expectedCnt, lines[0])
	}
	cols := strings.Fields(lines[1])
	if len(cols) != len(arr) {
		return fmt.Errorf("expected %d colors got %d", len(arr), len(cols))
	}
	for i, c := range cols {
		v, err := strconv.Atoi(c)
		if err != nil {
			return fmt.Errorf("bad color %q", c)
		}
		if v != expectedColors[i] {
			return fmt.Errorf("colors mismatch")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	spf := spfSieve(1000)
	for i := 0; i < 100; i++ {
		arr := generateCase(rng)
		if err := runCase(bin, arr, spf); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\narr: %v\n", i+1, err, arr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
