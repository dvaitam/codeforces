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

func runCandidate(bin, input string) (string, error) {
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

func checkSlope(y []int64, i, j int) bool {
	slopeNumer := y[j-1] - y[i-1]
	slopeDenom := int64(j - i)
	intercepts := make(map[int64]struct{})
	n := len(y)
	for k := 1; k <= n; k++ {
		val := slopeDenom*y[k-1] - slopeNumer*int64(k)
		if _, ok := intercepts[val]; !ok {
			intercepts[val] = struct{}{}
			if len(intercepts) > 2 {
				return false
			}
		}
	}
	return len(intercepts) == 2
}

func expected(y []int64) string {
	n := len(y)
	if n < 3 {
		return "No"
	}
	pairs := [][2]int{{1, 2}, {1, 3}, {2, 3}}
	for _, p := range pairs {
		if checkSlope(y, p[0], p[1]) {
			return "Yes"
		}
	}
	return "No"
}

func generateYesCase(rng *rand.Rand) []int64 {
	n := rng.Intn(8) + 3 // 3..10
	slope := int64(rng.Intn(5) - 2)
	b1 := int64(rng.Intn(11) - 5)
	delta := int64(rng.Intn(5) + 1)
	b2 := b1 + delta
	y := make([]int64, n)
	for i := 1; i <= n; i++ {
		if rng.Intn(2) == 0 {
			y[i-1] = slope*int64(i) + b1
		} else {
			y[i-1] = slope*int64(i) + b2
		}
	}
	return y
}

func generateNoCase(rng *rand.Rand) []int64 {
	n := rng.Intn(8) + 3
	y := make([]int64, n)
	for i := 0; i < n; i++ {
		y[i] = int64(rng.Intn(21) - 10)
	}
	// ensure it's a no case
	if expected(y) == "Yes" {
		y[0]++
	}
	return y
}

func makeInput(y []int64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(y)))
	for i, v := range y {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type testcase struct{ in, out string }
	cases := []testcase{}

	// deterministic edge cases
	cases = append(cases, testcase{makeInput([]int64{0, 1, 2}), "No"})
	cases = append(cases, testcase{makeInput([]int64{7, 5, 8, 6, 9}), "Yes"})
	cases = append(cases, testcase{makeInput([]int64{1, 1, 1}), "No"})
	cases = append(cases, testcase{makeInput([]int64{1, 2, 3}), "No"})

	for i := 0; i < 50; i++ {
		y := generateYesCase(rng)
		cases = append(cases, testcase{makeInput(y), expected(y)})
	}
	for i := 0; i < 50; i++ {
		y := generateNoCase(rng)
		cases = append(cases, testcase{makeInput(y), expected(y)})
	}

	for i, tc := range cases {
		got, err := runCandidate(bin, tc.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(tc.out) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.out, got, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
