package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	x []int
	y []int
}

func expected(x, y []int) int {
	i, j := 0, 0
	sumX, sumY := 0, 0
	count := 0
	n, m := len(x), len(y)
	for i < n || j < m {
		if sumX == sumY {
			if sumX != 0 {
				count++
				sumX, sumY = 0, 0
				continue
			}
			if i < n {
				sumX += x[i]
				i++
			}
			if j < m {
				sumY += y[j]
				j++
			}
		} else if sumX < sumY {
			if i < n {
				sumX += x[i]
				i++
			} else {
				break
			}
		} else {
			if j < m {
				sumY += y[j]
				j++
			} else {
				break
			}
		}
	}
	if sumX == sumY && sumX != 0 {
		count++
	}
	return count
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d\n", len(tc.x), len(tc.y))
	for i, v := range tc.x {
		if i == len(tc.x)-1 {
			input += fmt.Sprintf("%d\n", v)
		} else {
			input += fmt.Sprintf("%d ", v)
		}
	}
	for i, v := range tc.y {
		if i == len(tc.y)-1 {
			input += fmt.Sprintf("%d\n", v)
		} else {
			input += fmt.Sprintf("%d ", v)
		}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expected(tc.x, tc.y)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func generateCases() []testCase {
	rng := rand.New(rand.NewSource(1))
	cases := make([]testCase, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		x := make([]int, n)
		y := make([]int, m)
		sumX := 0
		sumY := 0
		for j := range x {
			x[j] = rng.Intn(10) + 1
			sumX += x[j]
		}
		for j := range y {
			y[j] = rng.Intn(10) + 1
			sumY += y[j]
		}
		if sumX > sumY {
			y[m-1] += sumX - sumY
		} else if sumY > sumX {
			x[n-1] += sumY - sumX
		}
		cases = append(cases, testCase{x: x, y: y})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
