package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solveCase(d int) (bool, float64, float64) {
	df := float64(d)
	D := df*df - 4*df
	if D < 0 {
		return false, 0, 0
	}
	sqrtD := math.Sqrt(D)
	a := (df + sqrtD) / 2
	b := (df - sqrtD) / 2
	return true, a, b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Build test cases: include edge cases plus random
	var cases [][]int
	cases = append(cases, []int{0})
	cases = append(cases, []int{1})
	cases = append(cases, []int{2})
	cases = append(cases, []int{3})
	cases = append(cases, []int{4})
	cases = append(cases, []int{5})
	cases = append(cases, []int{1000})
	for i := 0; i < 30; i++ {
		t := rng.Intn(5) + 1
		ds := make([]int, t)
		for j := 0; j < t; j++ {
			ds[j] = rng.Intn(1001)
		}
		cases = append(cases, ds)
	}

	for idx, ds := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(ds)))
		for _, d := range ds {
			sb.WriteString(fmt.Sprintf("%d\n", d))
		}
		input := sb.String()

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", idx+1, err, errBuf.String())
			os.Exit(1)
		}

		lines := strings.Split(strings.TrimSpace(out.String()), "\n")
		if len(lines) != len(ds) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d lines, got %d\ninput:\n%soutput:\n%s\n",
				idx+1, len(ds), len(lines), input, out.String())
			os.Exit(1)
		}

		for li, d := range ds {
			line := strings.TrimSpace(lines[li])
			hasSol, expA, expB := solveCase(d)
			_ = expA
			_ = expB

			if !hasSol {
				if line != "N" {
					fmt.Fprintf(os.Stderr, "case %d line %d: expected N for d=%d, got %q\n", idx+1, li+1, d, line)
					os.Exit(1)
				}
				continue
			}

			// Should start with Y
			if !strings.HasPrefix(line, "Y") {
				fmt.Fprintf(os.Stderr, "case %d line %d: expected Y for d=%d, got %q\n", idx+1, li+1, d, line)
				os.Exit(1)
			}

			parts := strings.Fields(line)
			if len(parts) != 3 {
				fmt.Fprintf(os.Stderr, "case %d line %d: expected 'Y a b', got %q\n", idx+1, li+1, line)
				os.Exit(1)
			}

			a, err1 := strconv.ParseFloat(parts[1], 64)
			b, err2 := strconv.ParseFloat(parts[2], 64)
			if err1 != nil || err2 != nil {
				fmt.Fprintf(os.Stderr, "case %d line %d: cannot parse floats from %q\n", idx+1, li+1, line)
				os.Exit(1)
			}

			df := float64(d)
			eps := 1e-6
			if math.Abs((a+b)-df) > eps {
				fmt.Fprintf(os.Stderr, "case %d line %d: a+b=%.15f != d=%d (diff=%.15f)\n",
					idx+1, li+1, a+b, d, math.Abs((a+b)-df))
				os.Exit(1)
			}
			if math.Abs(a*b-df) > eps {
				fmt.Fprintf(os.Stderr, "case %d line %d: a*b=%.15f != d=%d (diff=%.15f)\n",
					idx+1, li+1, a*b, d, math.Abs(a*b-df))
				os.Exit(1)
			}
			if a < -eps || b < -eps {
				fmt.Fprintf(os.Stderr, "case %d line %d: a=%.15f or b=%.15f is negative\n",
					idx+1, li+1, a, b)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
