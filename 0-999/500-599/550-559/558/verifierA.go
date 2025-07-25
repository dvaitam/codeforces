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

type testCase struct {
	input    string
	expected string
}

func expectedAnswer(n int, trees [][2]int) int {
	left := make([][2]int, 0, n)
	right := make([][2]int, 0, n)
	for _, t := range trees {
		if t[0] < 0 {
			left = append(left, t)
		} else {
			right = append(right, t)
		}
	}
	sort.Slice(left, func(i, j int) bool { return left[i][0] > left[j][0] })
	sort.Slice(right, func(i, j int) bool { return right[i][0] < right[j][0] })
	k := len(left)
	if len(right) < k {
		k = len(right)
	}
	sum := 0
	for i := 0; i < k; i++ {
		sum += left[i][1] + right[i][1]
	}
	if len(left) > len(right) && len(left) > k {
		sum += left[k][1]
	} else if len(right) > len(left) && len(right) > k {
		sum += right[k][1]
	}
	return sum
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	trees := make([][2]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		x := rng.Intn(201) - 100
		if x == 0 {
			x = 1
		}
		a := rng.Intn(100) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", x, a))
		trees[i] = [2]int{x, a}
	}
	expect := expectedAnswer(n, trees)
	return testCase{input: sb.String(), expected: fmt.Sprint(expect)}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != tc.expected {
		return fmt.Errorf("expected %s got %s", tc.expected, got)
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
