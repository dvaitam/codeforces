package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generateCase(rng *rand.Rand) (string, []int64) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	colors := make([]int, n)
	for i := 0; i < n; i++ {
		colors[i] = i + 1
	}
	colorful := make([]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	answers := make([]int64, 0)
	for i := 0; i < m; i++ {
		typ := rng.Intn(2) + 1
		if typ == 1 {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			x := rng.Intn(10) + 1
			sb.WriteString(fmt.Sprintf("1 %d %d %d\n", l, r, x))
			for j := l - 1; j < r; j++ {
				diff := abs(x - colors[j])
				colorful[j] += int64(diff)
				colors[j] = x
			}
		} else {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			sb.WriteString(fmt.Sprintf("2 %d %d\n", l, r))
			sum := int64(0)
			for j := l - 1; j < r; j++ {
				sum += colorful[j]
			}
			answers = append(answers, sum)
		}
	}
	return sb.String(), answers
}

func runCase(bin string, input string, answers []int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	scanner.Split(bufio.ScanWords)
	var got []int64
	for scanner.Scan() {
		v, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		got = append(got, v)
	}
	if len(got) != len(answers) {
		return fmt.Errorf("expected %d numbers got %d", len(answers), len(got))
	}
	for i, v := range answers {
		if got[i] != v {
			return fmt.Errorf("query %d expected %d got %d", i, v, got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
