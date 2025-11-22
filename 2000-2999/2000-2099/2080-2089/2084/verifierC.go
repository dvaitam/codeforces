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
)

const (
	refSource = "2084C.go"
	refBinary = "ref2084C.bin"
)

type pair struct {
	x, y int
}

type caseData struct {
	n int
	a []int
	b []int
}

type testCase struct {
	name  string
	cases []caseData
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		if err := validateSolution(tc, runProgramMust(refPath, input)); err != nil {
			fmt.Printf("reference failed on test %d (%s): %v\n", idx+1, tc.name, err)
			printInput(input)
			os.Exit(1)
		}

		if err := validateSolution(tc, runProgramMust(candidate, input)); err != nil {
			fmt.Printf("candidate failed on test %d (%s): %v\n", idx+1, tc.name, err)
			printInput(input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary, refSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary), nil
}

func runProgramMust(path string, input []byte) string {
	out, err := runProgram(path, input)
	if err != nil {
		fmt.Printf("runtime error while running %s: %v\noutput:\n%s\n", path, err, out)
		os.Exit(1)
	}
	return out
}

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.cases))
	for _, cs := range tc.cases {
		fmt.Fprintf(&sb, "%d\n", cs.n)
		writeIntSlice(&sb, cs.a)
		writeIntSlice(&sb, cs.b)
	}
	return []byte(sb.String())
}

func writeIntSlice(sb *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(sb, "%d", v)
	}
	sb.WriteByte('\n')
}

func validateSolution(tc testCase, out string) error {
	tokens := strings.Fields(out)
	ptr := 0
	for ci, cs := range tc.cases {
		if ptr >= len(tokens) {
			return fmt.Errorf("case %d: insufficient output", ci+1)
		}
		token := tokens[ptr]
		ptr++

		origA := append([]int(nil), cs.a...)
		origB := append([]int(nil), cs.b...)
		n := cs.n
		possible := isPossible(origA, origB)

		if token == "-1" {
			if possible {
				return fmt.Errorf("case %d: solution exists but got -1", ci+1)
			}
			continue
		}

		m, err := strconv.Atoi(token)
		if err != nil {
			return fmt.Errorf("case %d: invalid m %q", ci+1, token)
		}
		if m < 0 || m > n {
			return fmt.Errorf("case %d: m out of range", ci+1)
		}
		if ptr+2*m > len(tokens) {
			return fmt.Errorf("case %d: insufficient operations lines", ci+1)
		}

		a := origA
		b := origB
		for op := 0; op < m; op++ {
			i, err1 := strconv.Atoi(tokens[ptr])
			j, err2 := strconv.Atoi(tokens[ptr+1])
			ptr += 2
			if err1 != nil || err2 != nil {
				return fmt.Errorf("case %d: invalid indices", ci+1)
			}
			if i < 1 || i > n || j < 1 || j > n || i == j {
				return fmt.Errorf("case %d: indices out of range or equal", ci+1)
			}
			i--
			j--
			a[i], a[j] = a[j], a[i]
			b[i], b[j] = b[j], b[i]
		}

		if !isReverse(a, b) {
			return fmt.Errorf("case %d: final arrays not reverse", ci+1)
		}
	}
	if ptr != len(tokens) {
		return fmt.Errorf("extra output tokens")
	}
	return nil
}

func isReverse(a, b []int) bool {
	n := len(a)
	for i := 0; i < n; i++ {
		if a[i] != b[n-1-i] {
			return false
		}
	}
	return true
}

func isPossible(a, b []int) bool {
	n := len(a)
	cnt := make(map[pair]int)
	for i := 0; i < n; i++ {
		cnt[pair{a[i], b[i]}]++
	}

	oddSame := 0
	visited := make(map[pair]bool)
	for k, c := range cnt {
		if visited[k] {
			continue
		}
		if k.x == k.y {
			if c%2 == 1 {
				oddSame++
			}
			visited[k] = true
			continue
		}
		rev := pair{k.y, k.x}
		if cnt[rev] != c {
			return false
		}
		visited[k] = true
		visited[rev] = true
	}
	if oddSame > 1 || (oddSame == 1 && n%2 == 0) {
		return false
	}
	return true
}

func buildTests() []testCase {
	tests := []testCase{}

	tests = append(tests, testCase{
		name: "already-reverse",
		cases: []caseData{{
			n: 2,
			a: []int{1, 2},
			b: []int{2, 1},
		}},
	})
	tests = append(tests, testCase{
		name: "simple-impossible",
		cases: []caseData{{
			n: 2,
			a: []int{1, 2},
			b: []int{1, 2},
		}},
	})

	rnd := rand.New(rand.NewSource(20240601))

	// A possible case that is not already correct.
	for found := false; !found; {
		n := 6
		a := randomPerm(n, rnd)
		b := randomPerm(n, rnd)
		if isPossible(a, b) && !isReverse(a, b) {
			tests = append(tests, testCase{
				name: "needs-ops",
				cases: []caseData{{
					n: n,
					a: a,
					b: b,
				}},
			})
			found = true
		}
	}

	for t := 0; t < 12; t++ {
		caseCount := rnd.Intn(4) + 1
		cs := make([]caseData, 0, caseCount)
		for i := 0; i < caseCount; i++ {
			n := rnd.Intn(16) + 2
			a := randomPerm(n, rnd)
			b := randomPerm(n, rnd)
			cs = append(cs, caseData{n: n, a: a, b: b})
		}
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random-%d", t+1),
			cases: cs,
		})
	}

	// Larger but still reasonable size.
	largeCases := []caseData{}
	for i := 0; i < 3; i++ {
		n := 50 + 10*i
		a := randomPerm(n, rnd)
		b := randomPerm(n, rnd)
		largeCases = append(largeCases, caseData{n: n, a: a, b: b})
	}
	tests = append(tests, testCase{name: "medium-batch", cases: largeCases})

	return tests
}

func randomPerm(n int, rnd *rand.Rand) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	for i := n - 1; i > 0; i-- {
		j := rnd.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
	return p
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Print(string(in))
}
