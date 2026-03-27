package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var pref [3005][3005]int
var suff [3005][3005]int

// solveReference is the correct solver for 2085F1, embedded directly.
func solveReference(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var out bytes.Buffer
	writer := bufio.NewWriter(&out)

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return ""
	}

	type Item struct {
		CL int
		CR int
		D  int
	}
	items := make([]Item, 0, 3005)

	for tc := 0; tc < t; tc++ {
		var n, k int
		fmt.Fscan(reader, &n, &k)

		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		for i := 1; i <= n; i++ {
			for v := 1; v <= k; v++ {
				pref[i][v] = pref[i-1][v]
			}
			pref[i][a[i]] = i
		}

		for v := 1; v <= k; v++ {
			suff[n+1][v] = 1e9
		}
		for i := n; i >= 1; i-- {
			for v := 1; v <= k; v++ {
				suff[i][v] = suff[i+1][v]
			}
			suff[i][a[i]] = i
		}

		ans := int(1e15)

		if k%2 == 0 {
			h := k / 2
			C := h * h

			for M := 1; M < n; M++ {
				items = items[:0]
				for v := 1; v <= k; v++ {
					CL := int(1e9)
					if pref[M][v] > 0 {
						CL = -pref[M][v]
					}
					CR := int(1e9)
					if suff[M+1][v] <= n {
						CR = suff[M+1][v]
					}
					items = append(items, Item{CL, CR, CL - CR})
				}

				sort.Slice(items, func(i, j int) bool {
					return items[i].D < items[j].D
				})

				sum := 0
				for i := 0; i < h; i++ {
					sum += items[i].CL
				}
				for i := h; i < k; i++ {
					sum += items[i].CR
				}

				if sum < ans {
					ans = sum
				}
			}
			fmt.Fprintln(writer, ans-C)
		} else {
			h := (k - 1) / 2
			C := h * (h + 1)

			for i := 1; i <= n; i++ {
				vMid := a[i]
				items = items[:0]
				for v := 1; v <= k; v++ {
					if v == vMid {
						continue
					}
					CL := int(1e9)
					if pref[i-1][v] > 0 {
						CL = -pref[i-1][v]
					}
					CR := int(1e9)
					if suff[i+1][v] <= n {
						CR = suff[i+1][v]
					}
					items = append(items, Item{CL, CR, CL - CR})
				}

				sort.Slice(items, func(x, y int) bool {
					return items[x].D < items[y].D
				})

				sum := 0
				for x := 0; x < h; x++ {
					sum += items[x].CL
				}
				for x := h; x < k-1; x++ {
					sum += items[x].CR
				}

				if sum < ans {
					ans = sum
				}
			}
			fmt.Fprintln(writer, ans-C)
		}
	}
	writer.Flush()
	return strings.TrimSpace(out.String())
}

type testCase struct {
	input string
}

type testInstance struct {
	n, k int
	arr  []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}

	candidate := os.Args[1]
	tests := generateTests()

	for i, tc := range tests {
		expect := solveReference(tc.input)

		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}

		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		msg := stderr.String()
		if msg == "" {
			msg = stdout.String()
		}
		return "", fmt.Errorf("%v\n%s", err, msg)
	}
	return strings.TrimSpace(stdout.String()), nil
}

func equalTokens(a, b string) bool {
	ta := strings.Fields(a)
	tb := strings.Fields(b)
	if len(ta) != len(tb) {
		return false
	}
	for i := range ta {
		if ta[i] != tb[i] {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(20852085))
	var tests []testCase

	tests = append(tests, sampleTest())

	tests = append(tests, buildInput([]testInstance{
		{n: 2, k: 2, arr: []int{2, 1}},
	}))

	tests = append(tests, buildInput([]testInstance{
		{n: 5, k: 3, arr: []int{1, 1, 1, 2, 3}},
		{n: 6, k: 4, arr: []int{4, 4, 3, 2, 1, 1}},
	}))

	tests = append(tests, buildInput([]testInstance{
		{n: 8, k: 5, arr: []int{5, 5, 5, 1, 2, 3, 4, 4}},
	}))

	for i := 0; i < 30; i++ {
		tests = append(tests, randomTestCase(rng, rng.Intn(4)+1, 50))
	}

	tests = append(tests, randomTestCase(rng, 5, 400))
	tests = append(tests, worstCase())

	return tests
}

func sampleTest() testCase {
	return testCase{
		input: "6\n" +
			"3 2\n1 2 1\n" +
			"7 3\n2 1 1 3 1 1 2\n" +
			"6 3\n1 1 2 2 2 3\n" +
			"6 3\n1 2 2 2 2 3\n" +
			"10 5\n5 1 3 1 1 2 2 4 1 3\n" +
			"9 4\n1 2 3 3 3 3 3 2 4\n",
	}
}

func buildInput(instances []testInstance) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(instances))
	for _, inst := range instances {
		fmt.Fprintf(&b, "%d %d\n", inst.n, inst.k)
		for i, v := range inst.arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}

func randomTestCase(rng *rand.Rand, maxCases, maxN int) testCase {
	t := rng.Intn(maxCases) + 1
	var instances []testInstance
	for i := 0; i < t; i++ {
		n := rng.Intn(maxN-1) + 2
		k := rng.Intn(n-1) + 2
		arr := randomArray(rng, n, k)
		instances = append(instances, testInstance{n: n, k: k, arr: arr})
	}
	return buildInput(instances)
}

func randomArray(rng *rand.Rand, n, k int) []int {
	arr := make([]int, n)
	for i := 0; i < k; i++ {
		arr[i] = i + 1
	}
	for i := k; i < n; i++ {
		arr[i] = rng.Intn(k) + 1
	}
	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func worstCase() testCase {
	n := 3000
	k := 50
	arr := make([]int, n)
	idx := 0
	for val := 1; val <= k; val++ {
		arr[idx] = val
		idx++
	}
	for idx < n {
		arr[idx] = (idx % k) + 1
		idx++
	}
	return buildInput([]testInstance{{n: n, k: k, arr: arr}})
}
