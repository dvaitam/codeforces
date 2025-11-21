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

type caseOutput struct {
	k    int
	x, y []int64
}

type testCase struct {
	p, s int64
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func parseRef(out string, t int) ([]caseOutput, error) {
	reader := strings.NewReader(out)
	res := make([]caseOutput, 0, t)
	for i := 0; i < t; i++ {
		var k int
		if _, err := fmt.Fscan(reader, &k); err != nil {
			return nil, fmt.Errorf("reference output ended early at case %d: %v", i+1, err)
		}
		co := caseOutput{k: k}
		if k != -1 {
			for j := 0; j < k; j++ {
				var x, y int64
				if _, err := fmt.Fscan(reader, &x, &y); err != nil {
					return nil, fmt.Errorf("reference output ended early at case %d coord %d: %v", i+1, j+1, err)
				}
				co.x = append(co.x, x)
				co.y = append(co.y, y)
			}
		}
		res = append(res, co)
	}
	return res, nil
}

func parseCandidate(out string, t int) ([]caseOutput, error) {
	reader := strings.NewReader(out)
	res := make([]caseOutput, 0, t)
	for i := 0; i < t; i++ {
		var k int
		if _, err := fmt.Fscan(reader, &k); err != nil {
			return nil, fmt.Errorf("candidate output ended early at case %d: %v", i+1, err)
		}
		co := caseOutput{k: k}
		if k != -1 {
			if k < 1 || k > 50000 {
				return nil, fmt.Errorf("candidate k out of bounds at case %d: %d", i+1, k)
			}
			for j := 0; j < k; j++ {
				var x, y int64
				if _, err := fmt.Fscan(reader, &x, &y); err != nil {
					return nil, fmt.Errorf("candidate output ended early at case %d coord %d: %v", i+1, j+1, err)
				}
				if x < -1_000_000_000 || x > 1_000_000_000 || y < -1_000_000_000 || y > 1_000_000_000 {
					return nil, fmt.Errorf("candidate coordinate out of bounds at case %d coord %d", i+1, j+1)
				}
				co.x = append(co.x, x)
				co.y = append(co.y, y)
			}
		}
		res = append(res, co)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected after %d cases", t)
	}
	return res, nil
}

func verifySolution(tc testCase, sol caseOutput) error {
	k := sol.k
	if k == -1 {
		return fmt.Errorf("solution missing")
	}
	if len(sol.x) != k || len(sol.y) != k {
		return fmt.Errorf("coordinate count mismatch")
	}
	type key struct {
		x, y int64
	}
	cells := make(map[key]struct{}, k)
	for i := 0; i < k; i++ {
		kd := key{sol.x[i], sol.y[i]}
		if _, ok := cells[kd]; ok {
			return fmt.Errorf("overlapping pieces")
		}
		cells[kd] = struct{}{}
	}
	if len(cells) != k {
		return fmt.Errorf("duplicate pieces")
	}

	// Perimeter and connectivity.
	var area int64 = int64(k)
	perimeter := int64(0)
	dirs := [][2]int64{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for c := range cells {
		shared := 0
		for _, d := range dirs {
			_, ok := cells[key{c.x + d[0], c.y + d[1]}]
			if ok {
				shared++
			}
		}
		perimeter += int64(4 - shared)
	}

	if perimeter*tc.s != area*tc.p {
		return fmt.Errorf("ratio mismatch: per=%d area=%d p=%d s=%d", perimeter, area, tc.p, tc.s)
	}

	// Connectivity via DFS.
	var stack []key
	for c := range cells {
		stack = append(stack, c)
		break
	}
	visited := make(map[key]struct{}, k)
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if _, ok := visited[cur]; ok {
			continue
		}
		visited[cur] = struct{}{}
		for _, d := range dirs {
			nk := key{cur.x + d[0], cur.y + d[1]}
			if _, ok := cells[nk]; ok {
				if _, seen := visited[nk]; !seen {
					stack = append(stack, nk)
				}
			}
		}
	}
	if len(visited) != len(cells) {
		return fmt.Errorf("pieces not connected")
	}
	return nil
}

func buildTests() ([]byte, []testCase) {
	t := rand.Intn(4) + 1
	var sb strings.Builder
	tests := make([]testCase, 0, t)
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		var p, s int64
		switch rand.Intn(6) {
		case 0:
			p = 4
			s = 1
		case 1:
			p = rand.Int63n(50) + 1
			s = rand.Int63n(50) + 1
			if p < 4*s && p > 0 {
				p = 4 * s
			}
		case 2:
			p = rand.Int63n(50) + 1
			s = rand.Int63n(50) + 1
		case 3:
			// Impossible p>4s
			s = rand.Int63n(50) + 1
			p = 4*s + int64(rand.Intn(5)+1)
		default:
			// Small random
			p = int64(rand.Intn(10) + 1)
			s = int64(rand.Intn(10) + 1)
		}
		tests = append(tests, testCase{p: p, s: s})
		sb.WriteString(fmt.Sprintf("%d %d\n", p, s))
	}
	return []byte(sb.String()), tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refF.bin"
	if err := exec.Command("go", "build", "-o", ref, "2111F.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())

	for iter := 0; iter < 200; iter++ {
		input, tests := buildTests()
		refOut, err := run(ref, input)
		if err != nil {
			fmt.Println("reference failed on iteration", iter+1, ":", err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candOut, err := run(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on iteration %d: %v\n", iter+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		refParsed, err := parseRef(refOut, len(tests))
		if err != nil {
			fmt.Println("failed to parse reference output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", refOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candParsed, err := parseCandidate(candOut, len(tests))
		if err != nil {
			fmt.Println("failed to parse candidate output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", candOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		for i := range tests {
			refPossible := refParsed[i].k != -1
			if !refPossible {
				if candParsed[i].k != -1 {
					fmt.Printf("candidate provided solution when reference says impossible on iteration %d case %d\n", iter+1, i+1)
					fmt.Println("input case:", tests[i])
					os.Exit(1)
				}
				continue
			}
			if err := verifySolution(tests[i], candParsed[i]); err != nil {
				fmt.Printf("invalid solution on iteration %d case %d: %v\n", iter+1, i+1, err)
				fmt.Println("input case:", tests[i])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
