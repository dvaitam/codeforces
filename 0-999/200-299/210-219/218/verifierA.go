package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	n, k  int
	r     []int
	input string
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(100) + 1
	k := rng.Intn(n) + 1
	total := 2*n + 1
	y := make([]int, total)
	for i := 0; i < total; i += 2 {
		y[i] = rng.Intn(41)
	}
	for i := 1; i < total; i += 2 {
		left := y[i-1]
		right := y[i+1]
		base := max(left, right) + 1
		y[i] = base + rng.Intn(99-base+1)
	}
	r := append([]int(nil), y...)
	peaks := make([]int, 0, n)
	for i := 1; i < total; i += 2 {
		peaks = append(peaks, i)
	}
	rng.Shuffle(len(peaks), func(i, j int) { peaks[i], peaks[j] = peaks[j], peaks[i] })
	for i := 0; i < k; i++ {
		r[peaks[i]]++
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < total; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", r[i]))
	}
	sb.WriteByte('\n')
	return testCase{n: n, k: k, r: r, input: sb.String()}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func checkOutput(tc testCase, out string) error {
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) != len(tc.r) {
		return fmt.Errorf("expected %d numbers got %d", len(tc.r), len(fields))
	}
	y := make([]int, len(tc.r))
	for i, f := range fields {
		if _, err := fmt.Sscan(f, &y[i]); err != nil {
			return fmt.Errorf("bad int %q", f)
		}
		if y[i] < 0 || y[i] > 100 {
			return fmt.Errorf("value out of range %d", y[i])
		}
	}
	for i := 1; i < len(y)-1; i += 2 {
		if !(y[i] > y[i-1] && y[i] > y[i+1]) {
			return fmt.Errorf("peak %d not higher than neighbours", i+1)
		}
	}
	count := 0
	for i := 0; i < len(y); i++ {
		if i%2 == 1 {
			diff := tc.r[i] - y[i]
			if diff == 1 {
				count++
			} else if diff != 0 {
				return fmt.Errorf("peak %d difference %d not 0 or 1", i+1, diff)
			}
		} else if tc.r[i] != y[i] {
			return fmt.Errorf("index %d expected %d got %d", i+1, tc.r[i], y[i])
		}
	}
	if count != tc.k {
		return fmt.Errorf("expected %d changed peaks got %d", tc.k, count)
	}
	return nil
}

func runCase(exe string, tc testCase) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return checkOutput(tc, out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(exe, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
