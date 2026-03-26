package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ── embedded solver (CF-accepted 406E) ──────────────────────────────

func solve(input string) int64 {
	fields := strings.Fields(input)
	idx := 0
	nextInt := func() int64 {
		v, _ := strconv.ParseInt(fields[idx], 10, 64)
		idx++
		return v
	}

	n := nextInt()
	m := int(nextInt())

	Z := make([]int, 0, m)
	O := make([]int, 0, m)

	for i := 0; i < m; i++ {
		s := int(nextInt())
		f := int(nextInt())
		if s == 0 {
			Z = append(Z, f)
		} else {
			O = append(O, f)
		}
	}

	sort.Ints(Z)
	sort.Ints(O)

	Nz := len(Z)
	No := len(O)

	INF := int64(2e18)
	E_min := INF

	if Nz >= 3 {
		E0 := n - int64(Z[Nz-1]-Z[0])
		if E0 < E_min {
			E_min = E0
		}
	}
	if No >= 3 {
		E3 := n - int64(O[No-1]-O[0])
		if E3 < E_min {
			E_min = E3
		}
	}
	if Nz >= 2 && No >= 1 {
		for i := 0; i < No; i++ {
			o := O[i]
			var val int64
			if o < Z[0] {
				val = int64(Z[0] - o)
			} else if o > Z[Nz-1] {
				val = int64(o - Z[Nz-1])
			}
			if val < E_min {
				E_min = val
			}
		}
	}
	if Nz >= 1 && No >= 2 {
		for i := 0; i < Nz; i++ {
			z := Z[i]
			var val int64
			if z < O[0] {
				val = int64(O[0] - z)
			} else if z > O[No-1] {
				val = int64(z - O[No-1])
			}
			if val < E_min {
				E_min = val
			}
		}
	}

	var total_count int64

	if Nz >= 3 {
		E0 := n - int64(Z[Nz-1]-Z[0])
		if E0 == E_min {
			V_min := Z[0]
			V_max := Z[Nz-1]
			if V_min < V_max {
				C_min := int64(sort.Search(Nz, func(i int) bool { return Z[i] > V_min }))
				idx_max := sort.Search(Nz, func(i int) bool { return Z[i] >= V_max })
				C_max := int64(Nz - idx_max)
				C_mid := int64(Nz) - C_min - C_max
				cnt := C_max*(C_min*(C_min-1)/2) + C_min*(C_max*(C_max-1)/2) + C_min*C_max*C_mid
				total_count += cnt
			} else {
				nz := int64(Nz)
				total_count += nz * (nz - 1) * (nz - 2) / 6
			}
		}
	}
	if No >= 3 {
		E3 := n - int64(O[No-1]-O[0])
		if E3 == E_min {
			V_min := O[0]
			V_max := O[No-1]
			if V_min < V_max {
				C_min := int64(sort.Search(No, func(i int) bool { return O[i] > V_min }))
				idx_max := sort.Search(No, func(i int) bool { return O[i] >= V_max })
				C_max := int64(No - idx_max)
				C_mid := int64(No) - C_min - C_max
				cnt := C_max*(C_min*(C_min-1)/2) + C_min*(C_max*(C_max-1)/2) + C_min*C_max*C_mid
				total_count += cnt
			} else {
				no := int64(No)
				total_count += no * (no - 1) * (no - 2) / 6
			}
		}
	}
	if Nz >= 2 && No >= 1 {
		C_min := int64(sort.Search(Nz, func(i int) bool { return Z[i] > Z[0] }))
		idx_max := sort.Search(Nz, func(i int) bool { return Z[i] >= Z[Nz-1] })
		C_max := int64(Nz - idx_max)

		for i := 0; i < No; {
			o := O[i]
			j := i
			for j < No && O[j] == o {
				j++
			}
			cnt_o := int64(j - i)
			var val int64
			if o < Z[0] {
				val = int64(Z[0] - o)
			} else if o > Z[Nz-1] {
				val = int64(o - Z[Nz-1])
			}
			if val == E_min {
				var pairs_in_Z int64
				if val == 0 {
					idx_ge := sort.Search(Nz, func(k int) bool { return Z[k] >= o })
					C_less := int64(idx_ge)
					idx_gt := sort.Search(Nz, func(k int) bool { return Z[k] > o })
					C_equal := int64(idx_gt - idx_ge)
					C_greater := int64(Nz - idx_gt)
					pairs_in_Z = C_less*C_greater + C_less*C_equal + C_equal*C_greater + C_equal*(C_equal-1)/2
				} else if o < Z[0] {
					pairs_in_Z = C_min*(int64(Nz)-C_min) + C_min*(C_min-1)/2
				} else if o > Z[Nz-1] {
					pairs_in_Z = C_max*(int64(Nz)-C_max) + C_max*(C_max-1)/2
				}
				total_count += pairs_in_Z * cnt_o
			}
			i = j
		}
	}
	if Nz >= 1 && No >= 2 {
		C_min := int64(sort.Search(No, func(i int) bool { return O[i] > O[0] }))
		idx_max := sort.Search(No, func(i int) bool { return O[i] >= O[No-1] })
		C_max := int64(No - idx_max)

		for i := 0; i < Nz; {
			z := Z[i]
			j := i
			for j < Nz && Z[j] == z {
				j++
			}
			cnt_z := int64(j - i)
			var val int64
			if z < O[0] {
				val = int64(O[0] - z)
			} else if z > O[No-1] {
				val = int64(z - O[No-1])
			}
			if val == E_min {
				var pairs_in_O int64
				if val == 0 {
					idx_ge := sort.Search(No, func(k int) bool { return O[k] >= z })
					C_less := int64(idx_ge)
					idx_gt := sort.Search(No, func(k int) bool { return O[k] > z })
					C_equal := int64(idx_gt - idx_ge)
					C_greater := int64(No - idx_gt)
					pairs_in_O = C_less*C_greater + C_less*C_equal + C_equal*C_greater + C_equal*(C_equal-1)/2
				} else if z < O[0] {
					pairs_in_O = C_min*(int64(No)-C_min) + C_min*(C_min-1)/2
				} else if z > O[No-1] {
					pairs_in_O = C_max*(int64(No)-C_max) + C_max*(C_max-1)/2
				}
				total_count += pairs_in_O * cnt_z
			}
			i = j
		}
	}

	return total_count
}

// ── test generation ─────────────────────────────────────────────────

type testCase struct {
	n int64
	m int
	s []int
	f []int64
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i := 0; i < tc.m; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.s[i], tc.f[i]))
	}
	return sb.String()
}

func randInt64(rnd *rand.Rand, lo, hi int64) int64 {
	if lo == hi {
		return lo
	}
	return lo + rnd.Int63n(hi-lo+1)
}

func generateTests() []testCase {
	tests := []testCase{
		{1, 3, []int{0, 1, 0}, []int64{1, 1, 1}},
		{5, 3, []int{0, 0, 1}, []int64{3, 1, 4}},
		{10, 4, []int{0, 1, 1, 0}, []int64{1, 10, 5, 7}},
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for len(tests) < 57 {
		m := rnd.Intn(500) + 3
		n := randInt64(rnd, 1, 1_000_000_000)
		s := make([]int, m)
		f := make([]int64, m)
		for i := 0; i < m; i++ {
			s[i] = rnd.Intn(2)
			f[i] = randInt64(rnd, 1, n)
		}
		tests = append(tests, testCase{n, m, s, f})
	}

	// large tests
	for i := 0; i < 3; i++ {
		m := 100000
		n := int64(1_000_000_000)
		s := make([]int, m)
		f := make([]int64, m)
		for j := 0; j < m; j++ {
			s[j] = rnd.Intn(2)
			f[j] = randInt64(rnd, 1, n)
		}
		tests = append(tests, testCase{n, m, s, f})
	}

	return tests
}

func runProgram(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string) (int64, error) {
	f := strings.Fields(out)
	if len(f) != 1 {
		return 0, fmt.Errorf("expected 1 integer, got %d fields", len(f))
	}
	return strconv.ParseInt(f[0], 10, 64)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := generateTests()
	for i, tc := range tests {
		input := formatInput(tc)
		expected := solve(input)

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: bad output: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %d\n", i+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
