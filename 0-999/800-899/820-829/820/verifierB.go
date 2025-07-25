package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected [3]int
}

func expectedTriple(n, a int) [3]int {
	centerStep := 360.0 / float64(n)
	center := centerStep
	minDiff := math.Abs(center/2.0 - float64(a))
	vC, vS := int64(3), int64(2)
	ver := int64(3)
	center += centerStep
	for center <= 180.0 {
		diff1 := math.Abs(center/2.0 - float64(a))
		if diff1 < minDiff {
			vC = ver + 1
			vS = ver
			minDiff = diff1
		}
		diff2 := math.Abs(180.0 - center/2.0 - float64(a))
		if diff2 < minDiff {
			vC = ver - 1
			vS = ver
			minDiff = diff2
		}
		ver++
		center += centerStep
	}
	return [3]int{1, int(vC), int(vS)}
}

func buildCase(n, a int) testCase {
	input := fmt.Sprintf("%d %d\n", n, a)
	return testCase{input: input, expected: expectedTriple(n, a)}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(100000-3+1) + 3
	a := rng.Intn(180) + 1
	return buildCase(n, a)
}

func runCase(bin string, tc testCase) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var v1, v2, v3 int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &v1, &v2, &v3); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if v1 != tc.expected[0] || v2 != tc.expected[1] || v3 != tc.expected[2] {
		return fmt.Errorf("expected %v got %d %d %d", tc.expected, v1, v2, v3)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, buildCase(3, 60))
	cases = append(cases, buildCase(4, 67))
	cases = append(cases, buildCase(50, 90))

	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
