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

func solveC(m, t, r int, w []int) int {
	const offset = 300
	const size = 1000
	used := make([]bool, size)
	total := 0
	for _, wi := range w {
		cnt := 0
		for s := wi - t; s <= wi-1; s++ {
			idx := s + offset
			if idx >= 0 && idx < size && used[idx] {
				cnt++
			}
		}
		if cnt >= r {
			continue
		}
		need := r - cnt
		for s := wi - 1; s >= wi-t && need > 0; s-- {
			idx := s + offset
			if idx >= 0 && idx < size && !used[idx] {
				used[idx] = true
				total++
				need--
			}
		}
		if need > 0 {
			return -1
		}
	}
	return total
}

func generateCase(rng *rand.Rand) (string, int) {
	m := rng.Intn(10) + 1
	t := rng.Intn(10) + 1
	r := rng.Intn(t) + 1
	w := make([]int, m)
	cur := rng.Intn(50) + 1
	for i := 0; i < m; i++ {
		cur += rng.Intn(5) + 1
		w[i] = cur
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", m, t, r))
	for i, v := range w {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	expect := solveC(m, t, r, w)
	return sb.String(), expect
}

func runCase(bin, input string, expected int) error {
	cmd := exec.Command(bin)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		if err := runCase(bin, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
