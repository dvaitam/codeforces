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

func generateCase(rng *rand.Rand) (string, []int, []int) {
	n := rng.Intn(20) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(100)
	}
	b := make([]int, n)
	copy(b, a)
	// shuffle b
	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		b[i], b[j] = b[j], b[i]
	}
	var sb strings.Builder
	w := bufio.NewWriter(&sb)
	fmt.Fprintln(w, n)
	for i, v := range a {
		if i > 0 {
			fmt.Fprint(w, " ")
		}
		fmt.Fprint(w, v)
	}
	fmt.Fprintln(w)
	for i, v := range b {
		if i > 0 {
			fmt.Fprint(w, " ")
		}
		fmt.Fprint(w, v)
	}
	fmt.Fprintln(w)
	w.Flush()
	return sb.String(), a, b
}

func applySwaps(b []int, ops [][2]int) {
	for _, op := range ops {
		i := op[0]
		j := op[1]
		b[i], b[j] = b[j], b[i]
	}
}

func runCase(exe, input string, aOrig, bOrig []int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return fmt.Errorf("invalid k: %v", err)
	}
	if k < 0 || k > 1_000_000 {
		return fmt.Errorf("k out of range")
	}
	ops := make([][2]int, k)
	for i := 0; i < k; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("not enough numbers")
		}
		x, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("bad op index: %v", err)
		}
		if !scanner.Scan() {
			return fmt.Errorf("not enough numbers")
		}
		y, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("bad op index: %v", err)
		}
		ops[i] = [2]int{x - 1, y - 1}
	}
	b := make([]int, len(bOrig))
	copy(b, bOrig)
	for _, op := range ops {
		if op[0] < 0 || op[1] < 0 || op[0] >= len(b) || op[1] >= len(b) || op[1] != op[0]+1 {
			return fmt.Errorf("invalid swap indices")
		}
		b[op[0]], b[op[1]] = b[op[1]], b[op[0]]
	}
	for i := range aOrig {
		if aOrig[i] != b[i] {
			return fmt.Errorf("final array incorrect")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, a, b := generateCase(rng)
		if err := runCase(exe, in, a, b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
