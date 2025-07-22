package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const S = 1000000

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	used := make(map[int]bool)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		for {
			x := rng.Intn(S) + 1
			if !used[x] {
				used[x] = true
				arr[i] = x
				break
			}
		}
	}
	sort.Ints(arr)
	var in strings.Builder
	fmt.Fprintf(&in, "%d\n", n)
	for i, v := range arr {
		if i > 0 {
			in.WriteByte(' ')
		}
		fmt.Fprintf(&in, "%d", v)
	}
	in.WriteByte('\n')
	return in.String()
}

func checkAnswer(input, output string) error {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return fmt.Errorf("input parse: %v", err)
	}
	X := make(map[int]bool)
	sumX := int64(0)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		X[x] = true
		sumX += int64(x)
	}
	sumL := sumX - int64(n)

	out := bufio.NewReader(strings.NewReader(output))
	var m int
	if _, err := fmt.Fscan(out, &m); err != nil {
		return fmt.Errorf("output parse m: %v", err)
	}
	if m < 1 || m > S-n {
		return fmt.Errorf("invalid m %d", m)
	}
	Ys := make(map[int]bool)
	sumR := int64(0)
	for i := 0; i < m; i++ {
		var y int
		if _, err := fmt.Fscan(out, &y); err != nil {
			return fmt.Errorf("read y: %v", err)
		}
		if y < 1 || y > S {
			return fmt.Errorf("y out of range")
		}
		if X[y] {
			return fmt.Errorf("y intersects X")
		}
		if Ys[y] {
			return fmt.Errorf("duplicate y")
		}
		Ys[y] = true
		sumR += int64(S - y)
	}
	if sumR != sumL {
		return fmt.Errorf("sums mismatch")
	}
	return nil
}

func runCase(bin, input string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	return checkAnswer(input, buf.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
