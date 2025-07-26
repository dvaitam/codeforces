package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	n, m int
	a, b []int
}

func run(bin, input string) (string, error) {
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

func solveCase(tc testCase) int {
	a := append([]int(nil), tc.a...)
	b := append([]int(nil), tc.b...)
	sort.Ints(a)
	sort.Ints(b)
	best := tc.m
	for i := 0; i < tc.n; i++ {
		x := (b[i] - a[0]) % tc.m
		if x < 0 {
			x += tc.m
		}
		tmp := make([]int, tc.n)
		for j := 0; j < tc.n; j++ {
			val := (a[j] + x) % tc.m
			if val < 0 {
				val += tc.m
			}
			tmp[j] = val
		}
		sort.Ints(tmp)
		equal := true
		for j := 0; j < tc.n; j++ {
			if tmp[j] != b[j] {
				equal = false
				break
			}
		}
		if equal && x < best {
			best = x
		}
	}
	return best
}

func permute(rng *rand.Rand, arr []int) []int {
	perm := rng.Perm(len(arr))
	res := make([]int, len(arr))
	for i, p := range perm {
		res[i] = arr[p]
	}
	return res
}

func genRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 1
	if rng.Intn(10) == 0 {
		n = 2000
	}
	m := rng.Intn(1000000000) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(m)
	}
	x := rng.Intn(m)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		b[i] = (a[i] + x) % m
	}
	b = permute(rng, b)
	return testCase{n: n, m: m, a: a, b: b}
}

func generateTests(rng *rand.Rand) []testCase {
	tests := make([]testCase, 0, 100)
	tests = append(tests, testCase{n: 1, m: 2, a: []int{1}, b: []int{0}})
	tests = append(tests, testCase{n: 3, m: 5, a: []int{1, 2, 3}, b: []int{1, 2, 3}})
	tests = append(tests, testCase{n: 4, m: 3, a: []int{0, 0, 2, 1}, b: []int{2, 0, 1, 1}})
	tests = append(tests, testCase{n: 2, m: 5, a: []int{0, 0}, b: []int{3, 3}})
	tests = append(tests, testCase{n: 3, m: 1, a: []int{0, 0, 0}, b: []int{0, 0, 0}})
	tests = append(tests, testCase{n: 2, m: 1000000000, a: []int{999999999, 12345678}, b: []int{(999999999 + 500000000) % 1000000000, (12345678 + 500000000) % 1000000000}})

	for len(tests) < 100 {
		tests = append(tests, genRandomCase(rng))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := generateTests(rng)

	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.m)
		for j, val := range tc.a {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", val)
		}
		input += "\n"
		for j, val := range tc.b {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", val)
		}
		input += "\n"
		expected := fmt.Sprintf("%d", solveCase(tc))
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
