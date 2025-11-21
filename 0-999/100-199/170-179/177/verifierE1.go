package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const (
	maxInt64   = int64(^uint64(0) >> 1)
	maxSearchX = int64(1) << 60
)

type planet struct {
	a, b int64
}

type test struct {
	input    string
	expected string
}

func solveCase(n int, c int64, planets []planet) string {
	T := c - int64(n)
	if T < 0 {
		return "0"
	}
	allZero := true
	for _, p := range planets {
		if p.a != 0 {
			allZero = false
			break
		}
	}
	if T == 0 && allZero {
		return "-1"
	}
	capVal := T + 1
	sumFloor := func(x int64) int64 {
		if x <= 0 {
			return 0
		}
		var sum int64
		for _, p := range planets {
			if p.a == 0 {
				continue
			}
			if p.b <= 0 {
				continue
			}
			if x > maxInt64/p.a {
				return capVal
			}
			sum += p.a * x / p.b
			if sum >= capVal {
				return capVal
			}
		}
		return sum
	}
	lo, hi := int64(1), int64(1)
	for sumFloor(hi) < T && hi < maxSearchX {
		if hi > maxSearchX/2 {
			hi = maxSearchX
		} else {
			hi <<= 1
		}
	}
	if sumFloor(hi) < T {
		return "0"
	}
	for lo < hi {
		mid := lo + (hi-lo)/2
		if sumFloor(mid) < T {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	x0 := lo
	if sumFloor(x0) != T {
		return "0"
	}
	hi = x0
	for sumFloor(hi) == T && hi < maxSearchX {
		if hi > maxSearchX/2 {
			hi = maxSearchX
			break
		}
		hi <<= 1
	}
	if sumFloor(hi) == T {
		return "-1"
	}
	lo = x0
	for lo < hi {
		mid := lo + (hi-lo)/2
		if sumFloor(mid) <= T {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	x1 := lo
	return fmt.Sprintf("%d", x1-x0)
}

func formatInput(n int, c int64, planets []planet) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, c))
	for _, p := range planets {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.a, p.b))
	}
	return sb.String()
}

func fixedTests() []test {
	cases := []struct {
		n       int
		c       int64
		planets []planet
	}{
		{1, 1, []planet{{0, 5}}},
		{2, 5, []planet{{1, 5}, {2, 4}}},
		{1, 3, []planet{{5, 2}}},
		{3, 2, []planet{{1, 1}, {0, 4}, {2, 5}}},
		{2, 4, []planet{{0, 7}, {0, 3}}},
		{3, 12, []planet{{7, 3}, {0, 10}, {4, 5}}},
	}
	var tests []test
	for _, cs := range cases {
		inp := formatInput(cs.n, cs.c, cs.planets)
		exp := solveCase(cs.n, cs.c, cs.planets)
		tests = append(tests, test{inp, exp})
	}
	return tests
}

func randomSmallTests(rng *rand.Rand, need int) []test {
	var tests []test
	for len(tests) < need {
		n := rng.Intn(4) + 1
		c := int64(rng.Intn(200) + 1)
		planets := make([]planet, n)
		for i := 0; i < n; i++ {
			planets[i] = planet{
				a: int64(rng.Intn(6)),
				b: int64(rng.Intn(6) + 1),
			}
		}
		inp := formatInput(n, c, planets)
		exp := solveCase(n, c, planets)
		tests = append(tests, test{inp, exp})
	}
	return tests
}

func randomLargeTests(rng *rand.Rand, need int) []test {
	var tests []test
	for len(tests) < need {
		n := rng.Intn(40) + 60
		c := int64(rng.Int63n(1_000_000_000) + 1)
		planets := make([]planet, n)
		for i := 0; i < n; i++ {
			a := rng.Int63n(1_000_000_000)
			b := rng.Int63n(1_000_000_000) + 1
			planets[i] = planet{a: a, b: b}
		}
		inp := formatInput(n, c, planets)
		exp := solveCase(n, c, planets)
		tests = append(tests, test{inp, exp})
	}
	return tests
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(177))
	tests := fixedTests()
	tests = append(tests, randomSmallTests(rng, 80)...)
	tests = append(tests, randomLargeTests(rng, 20)...)
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\nInput:%s\n", i+1, err, t.input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
