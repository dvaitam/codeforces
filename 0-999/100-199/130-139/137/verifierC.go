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

func computeIncluded(intervals [][2]int) int {
	sort.Slice(intervals, func(i, j int) bool { return intervals[i][0] < intervals[j][0] })
	count := 0
	maxEnd := -1 << 60
	for _, iv := range intervals {
		if iv[1] < maxEnd {
			count++
		} else if iv[1] > maxEnd {
			maxEnd = iv[1]
		}
	}
	return count
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	vals := rng.Perm(2 * n)
	intervals := make([][2]int, n)
	for i := 0; i < n; i++ {
		a := vals[2*i] + 1
		b := vals[2*i+1] + 1
		if a > b {
			a, b = b, a
		}
		intervals[i] = [2]int{a, b}
	}
	// shuffle order
	rng.Shuffle(n, func(i, j int) { intervals[i], intervals[j] = intervals[j], intervals[i] })
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, iv := range intervals {
		fmt.Fprintf(&sb, "%d %d\n", iv[0], iv[1])
	}
	exp := fmt.Sprintf("%d", computeIncluded(append([][2]int(nil), intervals...)))
	return sb.String(), exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
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
