package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSourceD3 = "331D3.go"
	refBinaryD3 = "ref331D3.bin"
	totalTests  = 25
	maxArrows   = 30
	maxQueries  = 40
)

type arrow struct {
	x0, y0, x1, y1 int64
}

type query struct {
	x, y int64
	dir  rune
	t    int64
}

type testCase struct {
	b       int64
	arrs    []arrow
	queries []query
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD3.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(refPath)

	tests := generateTests()

	for i, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			return
		}
		refAns, err := parseOutput(refOut, len(tc.queries))
		if err != nil {
			fmt.Printf("reference output parse error on test %d: %v\n", i+1, err)
			return
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("Input used:")
			fmt.Println(input)
			return
		}
		candAns, err := parseOutput(candOut, len(tc.queries))
		if err != nil {
			fmt.Printf("candidate output parse error on test %d: %v\noutput:\n%s", i+1, err, candOut)
			fmt.Println("Input used:")
			fmt.Println(input)
			return
		}

		for j := range refAns {
			if refAns[j][0] != candAns[j][0] || refAns[j][1] != candAns[j][1] {
				fmt.Printf("Mismatch on test %d query %d: expected (%d %d), got (%d %d)\n",
					i+1, j+1, refAns[j][0], refAns[j][1], candAns[j][0], candAns[j][1])
				fmt.Println("Input used:")
				fmt.Println(input)
				return
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryD3, refSourceD3)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryD3), nil
}

func runProgram(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewBufferString(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, q int) ([][2]int64, error) {
	lines := strings.Fields(out)
	if len(lines) != 2*q {
		return nil, fmt.Errorf("expected %d numbers, got %d", 2*q, len(lines))
	}
	ans := make([][2]int64, q)
	for i := 0; i < q; i++ {
		x, err := strconv.ParseInt(lines[2*i], 10, 64)
		if err != nil {
			return nil, err
		}
		y, err := strconv.ParseInt(lines[2*i+1], 10, 64)
		if err != nil {
			return nil, err
		}
		ans[i] = [2]int64{x, y}
	}
	return ans, nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", len(tc.arrs), tc.b)
	for _, ar := range tc.arrs {
		fmt.Fprintf(&sb, "%d %d %d %d\n", ar.x0, ar.y0, ar.x1, ar.y1)
	}
	fmt.Fprintf(&sb, "%d\n", len(tc.queries))
	for _, q := range tc.queries {
		fmt.Fprintf(&sb, "%d %d %c %d\n", q.x, q.y, q.dir, q.t)
	}
	return sb.String()
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, totalTests)
	for len(tests) < totalTests {
		b := int64(rnd.Intn(900) + 100)
		n := rnd.Intn(maxArrows + 1)
		arrs := make([]arrow, 0, n)
		usedY := make(map[int64]struct{})
		for i := 0; i < n; i++ {
			y := int64(rnd.Intn(int(b-1))) + 1
			for {
				if _, ok := usedY[y]; !ok {
					break
				}
				y = int64(rnd.Intn(int(b-1))) + 1
			}
			usedY[y] = struct{}{}
			x0 := int64(rnd.Intn(int(b)))
			x1 := int64(rnd.Intn(int(b)))
			if x0 == x1 {
				if x1 < b {
					x1++
				} else {
					x0--
				}
			}
			arrs = append(arrs, arrow{x0: x0, y0: y, x1: x1, y1: y})
		}
		qn := rnd.Intn(maxQueries) + 1
		queries := make([]query, qn)
		dirs := []rune{'U', 'D', 'L', 'R'}
		for i := 0; i < qn; i++ {
			queries[i] = query{
				x:   int64(rnd.Intn(int(b + 1))),
				y:   int64(rnd.Intn(int(b + 1))),
				dir: dirs[rnd.Intn(len(dirs))],
				t:   int64(rnd.Intn(1_000_000)),
			}
		}
		tests = append(tests, testCase{b: b, arrs: arrs, queries: queries})
	}
	return tests
}
