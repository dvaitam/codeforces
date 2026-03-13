package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// bruteForce computes the maximum beauty by enumerating all partitions.
// For each partition into contiguous segments, the beauty of a segment
// is b[i] where i is the index of the element with minimum h in that segment.
func bruteForce(n int, h []int, b []int64) int64 {
	if n == 0 {
		return 0
	}
	best := int64(-1e18)
	// Enumerate all 2^(n-1) partitions via bitmask on the n-1 gaps
	for mask := 0; mask < (1 << uint(n-1)); mask++ {
		total := int64(0)
		start := 0
		for start < n {
			end := start
			for end < n-1 {
				if mask&(1<<uint(end)) != 0 {
					break
				}
				end++
			}
			// Segment is [start, end]
			// Find the index of min h in this segment
			minIdx := start
			for j := start + 1; j <= end; j++ {
				if h[j] < h[minIdx] {
					minIdx = j
				}
			}
			total += b[minIdx]
			start = end + 1
		}
		if total > best {
			best = total
		}
	}
	return best
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("run error: %v\nstderr: %s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rand.Seed(5)
	const T = 200
	for i := 0; i < T; i++ {
		n := rand.Intn(8) + 1
		perm := rand.Perm(n)
		h := make([]int, n)
		for j := 0; j < n; j++ {
			h[j] = perm[j] + 1
		}
		b := make([]int64, n)

		var input bytes.Buffer
		fmt.Fprintln(&input, n)
		for j := 0; j < n; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", h[j])
		}
		input.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			bv := int64(rand.Intn(11) - 5)
			b[j] = bv
			fmt.Fprintf(&input, "%d", bv)
		}
		input.WriteByte('\n')

		expected := bruteForce(n, h, b)
		inp := input.String()

		got, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s\n", i+1, err, inp)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output %q\n", i+1, got)
			os.Exit(1)
		}
		if gotVal != expected {
			fmt.Fprintf(os.Stderr, "test %d mismatch\nexpected: %d\ngot: %d\ninput:\n%s\n", i+1, expected, gotVal, inp)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
