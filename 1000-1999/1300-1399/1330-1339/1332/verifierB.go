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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func isComposite(n int) bool {
	if n < 4 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return true
		}
	}
	return false
}

func buildComposites() []int {
	var cs []int
	for i := 4; i <= 1000; i++ {
		if isComposite(i) {
			cs = append(cs, i)
		}
	}
	return cs
}

func generateCase(rng *rand.Rand, composites []int) []int {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = composites[rng.Intn(len(composites))]
	}
	return arr
}

func validate(arr []int, output string) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("expected 2 lines, got %d", len(lines))
	}
	m, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("cannot parse m: %v", err)
	}
	if m < 1 || m > 11 {
		return fmt.Errorf("m=%d out of [1,11]", m)
	}
	fields := strings.Fields(lines[1])
	n := len(arr)
	if len(fields) != n {
		return fmt.Errorf("expected %d colors, got %d", n, len(fields))
	}
	colors := make([]int, n)
	used := make([]bool, m+1)
	for i, f := range fields {
		c, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("bad color value %q", f)
		}
		if c < 1 || c > m {
			return fmt.Errorf("color %d out of [1,%d]", c, m)
		}
		colors[i] = c
		used[c] = true
	}
	for c := 1; c <= m; c++ {
		if !used[c] {
			return fmt.Errorf("color %d never used", c)
		}
	}
	// check same-color pairs have gcd > 1
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if colors[i] == colors[j] && gcd(arr[i], arr[j]) == 1 {
				return fmt.Errorf("arr[%d]=%d and arr[%d]=%d have same color %d but gcd=1",
					i, arr[i], j, arr[j], colors[i])
			}
		}
	}
	return nil
}

func runCase(bin string, arr []int) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
	input.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return validate(arr, out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	composites := buildComposites()
	for i := 0; i < 200; i++ {
		arr := generateCase(rng, composites)
		if err := runCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\narr: %v\n", i+1, err, arr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
