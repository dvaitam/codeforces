package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type matrix struct {
	n, m int
	data [][]int
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveB(a, b matrix) (int, int) {
	limitX := max(a.n, b.n)
	limitY := max(a.m, b.m)
	bestSum := 0
	bestX, bestY := 0, 0
	for x := -limitX; x < limitX; x++ {
		for y := -limitY; y < limitY; y++ {
			sum := 0
			rStart := max(0, -x)
			rEnd := min(a.n, b.n-x)
			cStart := max(0, -y)
			cEnd := min(a.m, b.m-y)
			for i := rStart; i < rEnd; i++ {
				for j := cStart; j < cEnd; j++ {
					sum += a.data[i][j] * b.data[i+x][j+y]
				}
			}
			if sum > bestSum {
				bestSum = sum
				bestX = x
				bestY = y
			}
		}
	}
	return bestX, bestY
}

func genTests() []struct {
	in     string
	expect string
} {
	tests := []struct{ in, expect string }{}
	for m1 := 2; len(tests) < 100 && m1 <= 2; m1++ {
		for mask1 := 0; len(tests) < 100 && mask1 < (1<<(2*m1)); mask1++ {
			for mask2 := 0; len(tests) < 100 && mask2 < (1<<(2*m1)); mask2++ {
				a := matrix{2, m1, make([][]int, 2)}
				b := matrix{2, m1, make([][]int, 2)}
				for i := 0; i < 2; i++ {
					a.data[i] = make([]int, m1)
					b.data[i] = make([]int, m1)
					for j := 0; j < m1; j++ {
						if mask1&(1<<(i*m1+j)) != 0 {
							a.data[i][j] = 1
						}
						if mask2&(1<<(i*m1+j)) != 0 {
							b.data[i][j] = 1
						}
					}
				}
				bx, by := solveB(a, b)
				input := fmt.Sprintf("2 %d\n", m1)
				for i := 0; i < 2; i++ {
					for j := 0; j < m1; j++ {
						input += fmt.Sprintf("%d", a.data[i][j])
					}
					input += "\n"
				}
				input += fmt.Sprintf("2 %d\n", m1)
				for i := 0; i < 2; i++ {
					for j := 0; j < m1; j++ {
						input += fmt.Sprintf("%d", b.data[i][j])
					}
					input += "\n"
				}
				exp := fmt.Sprintf("%d %d", bx, by)
				tests = append(tests, struct{ in, expect string }{input, exp})
			}
		}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != t.expect {
			fmt.Printf("test %d failed: expected=%s got=%s\ninput:\n%s\n", i+1, t.expect, out, t.in)
			os.Exit(1)
		}
	}
	fmt.Printf("ok %d tests\n", len(tests))
}
