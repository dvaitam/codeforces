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

type Arrow struct {
	x0, y0, x1, y1 int
}

type Query struct {
	x, y int
	dir  byte
	t    int64
}

type testCase struct {
	n       int
	b       int
	arrows  []Arrow
	queries []Query
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := serializeInput(tests)

	expected, err := runAndParse(refBin, input, totalQueries(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	got, err := runAndParse(candidate, input, totalQueries(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	idx := 0
	for tcIdx, tc := range tests {
		for qIdx := range tc.queries {
			if expected[idx] != got[idx] {
				fmt.Fprintf(os.Stderr, "Mismatch in test %d query %d: expected %s got %s\n", tcIdx+1, qIdx+1, expected[idx], got[idx])
				fmt.Fprintf(os.Stderr, "n=%d b=%d arrows=%v query=%v\n", tc.n, tc.b, tc.arrows, tc.queries[qIdx])
				os.Exit(1)
			}
			idx++
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_331D1.bin"
	cmd := exec.Command("go", "build", "-o", refName, "331D1.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runAndParse(target, input string, count int) ([]string, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(out)
	if len(fields) != count*2 {
		return nil, fmt.Errorf("expected %d coordinate pairs, got %d tokens (output: %q)", count, len(fields), out)
	}
	results := make([]string, count)
	for i := 0; i < count; i++ {
		results[i] = fmt.Sprintf("%s %s", fields[2*i], fields[2*i+1])
	}
	return results, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	sumQueries := totalQueries(tests)

	for len(tests) < 50 && sumQueries < 2000 {
		n := rng.Intn(5) + 1
		b := rng.Intn(30) + 1
		arrows := make([]Arrow, n)
		occupied := make(map[[4]int]struct{})
		for i := 0; i < n; i++ {
			for {
				if rng.Intn(2) == 0 {
					x := rng.Intn(b + 1)
					y0 := rng.Intn(b)
					y1 := rng.Intn(b)
					if y0 == y1 {
						continue
					}
					if y0 > y1 {
						y0, y1 = y1, y0
					}
					a := Arrow{x0: x, y0: y0, x1: x, y1: y1}
					key := [4]int{a.x0, a.y0, a.x1, a.y1}
					if _, ok := occupied[key]; !ok {
						occupied[key] = struct{}{}
						arrows[i] = a
						break
					}
				} else {
					y := rng.Intn(b + 1)
					x0 := rng.Intn(b)
					x1 := rng.Intn(b)
					if x0 == x1 {
						continue
					}
					if x0 > x1 {
						x0, x1 = x1, x0
					}
					a := Arrow{x0: x0, y0: y, x1: x1, y1: y}
					key := [4]int{a.x0, a.y0, a.x1, a.y1}
					if _, ok := occupied[key]; !ok {
						occupied[key] = struct{}{}
						arrows[i] = a
						break
					}
				}
			}
		}
		q := rng.Intn(20) + 1
		queries := make([]Query, q)
		dirs := []byte{'U', 'D', 'L', 'R'}
		for i := 0; i < q; i++ {
			queries[i] = Query{
				x:   rng.Intn(b + 1),
				y:   rng.Intn(b + 1),
				dir: dirs[rng.Intn(len(dirs))],
				t:   rng.Int63n(1000),
			}
		}
		tests = append(tests, testCase{n: n, b: b, arrows: arrows, queries: queries})
		sumQueries += q
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n:      0,
			b:      5,
			arrows: nil,
			queries: []Query{
				{x: 0, y: 0, dir: 'R', t: 3},
				{x: 5, y: 5, dir: 'L', t: 10},
			},
		},
		{
			n: 1,
			b: 10,
			arrows: []Arrow{
				{x0: 0, y0: 0, x1: 0, y1: 5},
			},
			queries: []Query{
				{x: 0, y: 0, dir: 'U', t: 6},
				{x: 3, y: 3, dir: 'R', t: 5},
			},
		},
	}
}

func totalQueries(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += len(tc.queries)
	}
	return total
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.b))
		for _, a := range tc.arrows {
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", a.x0, a.y0, a.x1, a.y1))
		}
		sb.WriteString(strconv.Itoa(len(tc.queries)))
		sb.WriteByte('\n')
		for _, q := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d %c %d\n", q.x, q.y, q.dir, q.t))
		}
	}
	return sb.String()
}
