package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Test struct {
	input string
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// parse input into n, m, a, segments
func parseInput(input string) (int, int, []int64, [][2]int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) < 2 {
		return 0, 0, nil, nil, fmt.Errorf("malformed input")
	}
	hdr := strings.Fields(lines[0])
	if len(hdr) < 2 {
		return 0, 0, nil, nil, fmt.Errorf("bad header")
	}
	n, _ := strconv.Atoi(hdr[0])
	m, _ := strconv.Atoi(hdr[1])
	arrFields := strings.Fields(lines[1])
	if len(arrFields) != n {
		return 0, 0, nil, nil, fmt.Errorf("bad array length")
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		v, _ := strconv.ParseInt(arrFields[i], 10, 64)
		a[i] = v
	}
	segs := make([][2]int, 0, m)
	for i := 0; i < m; i++ {
		if 2+i >= len(lines) {
			return 0, 0, nil, nil, fmt.Errorf("missing segments")
		}
		f := strings.Fields(lines[2+i])
		if len(f) < 2 {
			return 0, 0, nil, nil, fmt.Errorf("bad segment line")
		}
		l, _ := strconv.Atoi(f[0])
		r, _ := strconv.Atoi(f[1])
		segs = append(segs, [2]int{l - 1, r - 1})
	}
	return n, m, a, segs, nil
}

func applySubset(a []int64, segs [][2]int, mask int) ([]int64, int) {
	n := len(a)
	b := make([]int64, n)
	copy(b, a)
	m := len(segs)
	count := 0
	for i := 0; i < m; i++ {
		if (mask>>i)&1 == 1 {
			count++
			l := segs[i][0]
			r := segs[i][1]
			for j := l; j <= r; j++ {
				b[j]--
			}
		}
	}
	return b, count
}

func score(a []int64) int64 {
	mn := int64(math.MaxInt64)
	mx := int64(math.MinInt64)
	for _, v := range a {
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	return mx - mn
}

func genTests() []Test {
	rand.Seed(5)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(6) + 1
		m := rand.Intn(6)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j := 0; j < n; j++ {
			val := rand.Intn(11) - 5
			sb.WriteString(strconv.Itoa(val))
			if j+1 < n {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
		for j := 0; j < m; j++ {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
		}
		tests = append(tests, Test{sb.String()})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE2.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	tests := genTests()
	for i, tc := range tests {
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		// Validate candidate output
		n, m, a, segs, perr := parseInput(tc.input)
		if perr != nil {
			fmt.Printf("internal parse error on test %d: %v\n", i+1, perr)
			os.Exit(1)
		}
		// Compute optimal score by brute force over all subsets (m is small in generator)
		best := int64(math.MinInt64)
		for mask := 0; mask < (1 << m); mask++ {
			b, _ := applySubset(a, segs, mask)
			sc := score(b)
			if sc > best {
				best = sc
			}
		}
		out := strings.Fields(strings.TrimSpace(got))
		if len(out) < 2 {
			fmt.Printf("Test %d failed: malformed output\nInput:\n%sGot:\n%s\n", i+1, tc.input, got)
			os.Exit(1)
		}
		claimed, e1 := strconv.ParseInt(out[0], 10, 64)
		k, e2 := strconv.Atoi(out[1])
		if e1 != nil || e2 != nil {
			fmt.Printf("Test %d failed: non-integer header\nInput:\n%sGot:\n%s\n", i+1, tc.input, got)
			os.Exit(1)
		}
		idxs := []int{}
		for t := 0; t < k && 2+t < len(out); t++ {
			v, er := strconv.Atoi(out[2+t])
			if er != nil {
				break
			}
			idxs = append(idxs, v)
		}
		// If not enough indices on the same line, try to parse across entire output
		if len(idxs) < k {
			toks := strings.Fields(got)
			idxs = []int{}
			for p := 2; p < len(toks) && len(idxs) < k; p++ {
				v, er := strconv.Atoi(toks[p])
				if er == nil {
					idxs = append(idxs, v)
				}
			}
		}
		if claimed != best {
			fmt.Printf("Test %d failed: wrong best value (expected %d got %d)\nInput:\n%sGot:\n%s\n", i+1, best, claimed, tc.input, got)
			os.Exit(1)
		}
		// Validate chosen subset yields the claimed score
		mask := 0
		used := make(map[int]bool)
		ok := true
		for _, id := range idxs {
			if id < 1 || id > m || used[id] {
				ok = false
				break
			}
			used[id] = true
			mask |= 1 << (id - 1)
		}
		if !ok || len(idxs) != k {
			fmt.Printf("Test %d failed: invalid indices list\nInput:\n%sGot:\n%s\n", i+1, tc.input, got)
			os.Exit(1)
		}
		b, _ := applySubset(a, segs, mask)
		sc := score(b)
		if sc != claimed {
			fmt.Printf("Test %d failed: chosen subset score mismatch (expected %d got %d)\nInput:\n%sGot:\n%s\n", i+1, claimed, sc, tc.input, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
