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

func solveA(n, k int, s string) int {
	if k > n || k > 13 {
		return -1
	}
	b := []byte(s)
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	const inf = int(^uint(0) >> 1)
	mn := inf
	found := false
	for i := 0; i < n; i++ {
		cnt := 1
		prev := b[i]
		sum := int(b[i] - 'a' + 1)
		for j := i + 1; j < n && cnt < k; j++ {
			if b[j]-prev >= 2 {
				sum += int(b[j] - 'a' + 1)
				prev = b[j]
				cnt++
			}
		}
		if cnt == k {
			found = true
			if sum < mn {
				mn = sum
			}
		}
	}
	if !found {
		return -1
	}
	return mn
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(50) + 1
	k := rng.Intn(n) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	s := string(b)
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", n, k))
	input.WriteString(s)
	input.WriteByte('\n')
	expected := solveA(n, k, s)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
