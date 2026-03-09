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

// solveC returns the optimal d (as float64) for the reference.
func solveC(a []int64) float64 {
	n := len(a)
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	// test(x): returns true if 3 segments of span x are NOT sufficient.
	test := func(x int64) bool {
		b := 0
		for i := 0; i < 3; i++ {
			if b >= n {
				return false
			}
			key := a[b] + x
			pos := sort.Search(n, func(j int) bool { return a[j] > key })
			b = pos
			if b == n {
				return false
			}
		}
		return true
	}
	var l, r int64 = 0, 1000000000
	for l < r {
		mid := l + (r-l)/2
		if test(mid) {
			l = mid + 1
		} else {
			r = mid
		}
	}
	return float64(l) / 2.0
}

func generateCase(rng *rand.Rand) (string, []int64) {
	n := rng.Intn(10) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = int64(rng.Intn(1000)) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), append([]int64(nil), arr...)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// validate checks that the candidate output is a valid answer for the given houses.
// optD is the reference optimal d.
func validate(output string, houses []int64, optD float64) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("expected 2 lines of output, got %d", len(lines))
	}

	dCand, err := strconv.ParseFloat(strings.TrimSpace(lines[0]), 64)
	if err != nil {
		return fmt.Errorf("cannot parse d: %v", err)
	}

	fields := strings.Fields(lines[1])
	if len(fields) != 3 {
		return fmt.Errorf("expected 3 center coordinates, got %d", len(fields))
	}
	centers := make([]float64, 3)
	for i, f := range fields {
		centers[i], err = strconv.ParseFloat(f, 64)
		if err != nil {
			return fmt.Errorf("cannot parse center %d: %v", i, err)
		}
	}

	const eps = 1e-4

	// Check optimality: candidate d must not exceed reference d by more than eps.
	if dCand > optD+eps {
		return fmt.Errorf("d=%.6f is worse than optimal %.6f", dCand, optD)
	}

	// Check feasibility: every house must be within dCand of some center.
	for _, h := range houses {
		covered := false
		for _, c := range centers {
			if abs(float64(h)-c) <= dCand+eps {
				covered = true
				break
			}
		}
		if !covered {
			return fmt.Errorf("house %d not covered by any center (d=%.6f, centers=%v)", h, dCand, centers)
		}
	}

	return nil
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Fixed edge cases
	type fixedCase struct {
		input  string
		houses []int64
	}
	fixed := []fixedCase{
		{"1\n500\n", []int64{500}},
		{"3\n846 687 287\n", []int64{846, 687, 287}},
		{"3\n1 1 1\n", []int64{1, 1, 1}},
		{"3\n1 500 1000\n", []int64{1, 500, 1000}},
	}

	caseNum := 1
	for _, fc := range fixed {
		optD := solveC(append([]int64(nil), fc.houses...))
		got, err := run(bin, fc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum, err, fc.input)
			os.Exit(1)
		}
		if err := validate(got, fc.houses, optD); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\noutput:\n%s\ninput:\n%s", caseNum, err, got, fc.input)
			os.Exit(1)
		}
		caseNum++
	}

	for caseNum <= 100 {
		in, houses := generateCase(rng)
		optD := solveC(append([]int64(nil), houses...))
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum, err, in)
			os.Exit(1)
		}
		if err := validate(got, houses, optD); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\noutput:\n%s\ninput:\n%s", caseNum, err, got, in)
			os.Exit(1)
		}
		caseNum++
	}
	fmt.Println("All tests passed")
}
