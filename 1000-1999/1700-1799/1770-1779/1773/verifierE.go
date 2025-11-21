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
	refSource1773E = "1773E.go"
	refBinary1773E = "ref1773E.bin"
	maxTests       = 120
	maxTotalBlocks = 10000
)

type testCase struct {
	towers [][]int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := range tests {
		if expected[i][0] != got[i][0] || expected[i][1] != got[i][1] {
			fmt.Printf("Mismatch on test %d: expected (%d %d), got (%d %d)\n",
				i+1, expected[i][0], expected[i][1], got[i][0], got[i][1])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary1773E, refSource1773E)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary1773E), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, t int) ([][2]int64, error) {
	lines := strings.Fields(out)
	if len(lines) != 2*t {
		return nil, fmt.Errorf("expected %d numbers, got %d", 2*t, len(lines))
	}
	res := make([][2]int64, t)
	for i := 0; i < t; i++ {
		a, err := strconv.ParseInt(lines[2*i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", 2*i+1, err)
		}
		b, err := strconv.ParseInt(lines[2*i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", 2*i+2, err)
		}
		res[i] = [2]int64{a, b}
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", len(tc.towers))
		for _, tower := range tc.towers {
			fmt.Fprintf(&sb, "%d", len(tower))
			for _, v := range tower {
				fmt.Fprintf(&sb, " %d", v)
			}
			sb.WriteByte('\n')
		}
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(1773))
	var tests []testCase
	totalBlocks := 0

	add := func(tc testCase) {
		tests = append(tests, tc)
		for _, tower := range tc.towers {
			totalBlocks += len(tower)
		}
	}

	// Simple cases
	add(testCase{towers: [][]int{{1}, {2}, {3}}})
	add(testCase{towers: [][]int{{3, 1}, {5, 2}, {4}}})

	nextValue := 100 // ensure unique

	for len(tests) < maxTests && totalBlocks < maxTotalBlocks {
		remain := maxTotalBlocks - totalBlocks
		if remain <= 0 {
			break
		}
		maxN := 50
		n := rnd.Intn(maxN) + 1
		towers := make([][]int, n)
		blocksLeft := remain
		for i := 0; i < n; i++ {
			if blocksLeft == 0 {
				towers = towers[:i]
				break
			}
			maxHeight := 1 + blocksLeft - (n - i - 1)
			if maxHeight > 100 {
				maxHeight = 100
			}
			h := rnd.Intn(maxHeight) + 1
			blocksLeft -= h
			tower := make([]int, h)
			for j := 0; j < h; j++ {
				tower[j] = nextValue
				nextValue += rnd.Intn(3) + 1
			}
			towers[i] = tower
		}
		if len(towers) == 0 {
			break
		}
		add(testCase{towers: towers})
	}
	return tests
}
