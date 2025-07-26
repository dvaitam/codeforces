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

func expectedB(grades []int) int {
	n := len(grades)
	sum := 0
	for _, g := range grades {
		sum += g
	}
	cur := sum * 10
	target := 45 * n
	if cur >= target {
		return 0
	}
	sort.Ints(grades)
	cnt := 0
	for _, g := range grades {
		diff := 5 - g
		cur += diff * 10
		cnt++
		if cur >= target {
			return cnt
		}
	}
	return cnt
}

func generateCaseB(rng *rand.Rand) (string, int) {
	n := rng.Intn(100) + 1
	grades := make([]int, n)
	for i := range grades {
		grades[i] = rng.Intn(5) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, g := range grades {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(g))
	}
	sb.WriteByte('\n')
	expected := expectedB(append([]int(nil), grades...))
	return sb.String(), expected
}

func runCaseB(bin, input string, expected int) error {
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseB(rng)
		if err := runCaseB(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
