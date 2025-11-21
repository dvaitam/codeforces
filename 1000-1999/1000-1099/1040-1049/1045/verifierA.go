package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type weaponDef struct {
	typ     int
	list    []int
	allowed map[int]struct{}
	l, r    int
	triple  [3]int
}

type testCase struct {
	name    string
	n, m    int
	weapons []weaponDef
	input   string
}

const refSource = "1045A.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate_binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expected, err := parseReference(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		obtained, err := validateCandidate(tc, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if obtained != expected {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d destroyed ships, got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
				idx+1, tc.name, expected, obtained, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1045A-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleA")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseReference(out string) (int, error) {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.Atoi(tokens[0])
	if err != nil {
		return 0, fmt.Errorf("invalid first line: %v", err)
	}
	if val < 0 {
		return 0, fmt.Errorf("negative destroyed count %d", val)
	}
	return val, nil
}

func validateCandidate(tc testCase, out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	x, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid destroyed count: %v", err)
	}
	if x < 0 {
		return 0, fmt.Errorf("negative number of destroyed ships %d", x)
	}
	expectedTokens := 1 + 2*x
	if len(fields) != expectedTokens {
		return 0, fmt.Errorf("expected %d integers, got %d", expectedTokens, len(fields))
	}
	if x == 0 {
		return 0, nil
	}
	shipUsed := make([]bool, tc.m+1)
	weaponUsage := make([]int, tc.n+1)
	for i := 0; i < x; i++ {
		wIdx, err := strconv.Atoi(fields[1+2*i])
		if err != nil {
			return 0, fmt.Errorf("invalid weapon index: %v", err)
		}
		sIdx, err := strconv.Atoi(fields[2+2*i])
		if err != nil {
			return 0, fmt.Errorf("invalid spaceship index: %v", err)
		}
		if wIdx < 1 || wIdx > tc.n {
			return 0, fmt.Errorf("weapon index %d out of range", wIdx)
		}
		if sIdx < 1 || sIdx > tc.m {
			return 0, fmt.Errorf("spaceship index %d out of range", sIdx)
		}
		if shipUsed[sIdx] {
			return 0, fmt.Errorf("spaceship %d destroyed multiple times", sIdx)
		}
		shipUsed[sIdx] = true
		weapon := tc.weapons[wIdx-1]
		switch weapon.typ {
		case 0:
			if _, ok := weapon.allowed[sIdx]; !ok {
				return 0, fmt.Errorf("weapon %d cannot target spaceship %d", wIdx, sIdx)
			}
		case 1:
			if sIdx < weapon.l || sIdx > weapon.r {
				return 0, fmt.Errorf("weapon %d targets outside [%d,%d]", wIdx, weapon.l, weapon.r)
			}
		case 2:
			if sIdx != weapon.triple[0] && sIdx != weapon.triple[1] && sIdx != weapon.triple[2] {
				return 0, fmt.Errorf("bazooka %d cannot target spaceship %d", wIdx, sIdx)
			}
		default:
			return 0, fmt.Errorf("unknown weapon type %d", weapon.typ)
		}
		weaponUsage[wIdx]++
	}
	for idx, weapon := range tc.weapons {
		count := weaponUsage[idx+1]
		switch weapon.typ {
		case 0, 1:
			if count > 1 {
				return 0, fmt.Errorf("weapon %d used %d times", idx+1, count)
			}
		case 2:
			if count != 0 && count != 2 {
				return 0, fmt.Errorf("bazooka %d must be used 0 or 2 times, got %d", idx+1, count)
			}
		}
	}
	return x, nil
}

func deterministicTests() []testCase {
	return []testCase{
		finalizeTestCase("single_sql", 5, []weaponDef{
			{typ: 0, list: []int{3}},
		}),
		finalizeTestCase("simple_beam", 6, []weaponDef{
			{typ: 1, l: 1, r: 6},
			{typ: 0, list: []int{2, 4}},
		}),
		finalizeTestCase("two_bazookas", 9, []weaponDef{
			{typ: 2, triple: [3]int{1, 2, 3}},
			{typ: 2, triple: [3]int{4, 5, 6}},
			{typ: 1, l: 3, r: 9},
		}),
		finalizeTestCase("mixed_choices", 12, []weaponDef{
			{typ: 0, list: []int{1, 4, 7}},
			{typ: 1, l: 2, r: 5},
			{typ: 2, triple: [3]int{8, 9, 10}},
			{typ: 0, list: []int{6, 11, 12}},
			{typ: 1, l: 7, r: 12},
		}),
		finalizeTestCase("dense_intervals", 40, func() []weaponDef {
			var defs []weaponDef
			for i := 0; i < 10; i++ {
				l := 1 + i*4
				r := l + 9
				if r > 40 {
					r = 40
				}
				defs = append(defs, weaponDef{typ: 1, l: l, r: r})
			}
			defs = append(defs, weaponDef{typ: 0, list: []int{5, 15, 25, 35}})
			defs = append(defs, weaponDef{typ: 2, triple: [3]int{31, 32, 33}})
			return defs
		}()),
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	m := rng.Intn(250) + 1
	n := rng.Intn(200) + 1
	defs := make([]weaponDef, 0, n)
	bazookaUsed := make([]bool, m+1)
	for i := 0; i < n; i++ {
		typ := rng.Intn(3)
		if typ == 2 {
			available := availableShips(bazookaUsed)
			if len(available) < 3 {
				typ = rng.Intn(2)
			}
		}
		w := weaponDef{typ: typ}
		switch typ {
		case 0:
			maxK := 6
			if m < maxK {
				maxK = m
			}
			k := rng.Intn(maxK) + 1
			w.list = sampleUnique(rng, m, k)
		case 1:
			l := rng.Intn(m) + 1
			r := rng.Intn(m) + 1
			if l > r {
				l, r = r, l
			}
			w.l = l
			w.r = r
		case 2:
			available := availableShips(bazookaUsed)
			rng.Shuffle(len(available), func(i, j int) {
				available[i], available[j] = available[j], available[i]
			})
			a, b, c := available[0], available[1], available[2]
			w.triple = [3]int{a, b, c}
			bazookaUsed[a] = true
			bazookaUsed[b] = true
			bazookaUsed[c] = true
		}
		defs = append(defs, w)
	}
	return finalizeTestCase(fmt.Sprintf("random_%d", idx+1), m, defs)
}

func finalizeTestCase(name string, m int, defs []weaponDef) testCase {
	weapons := make([]weaponDef, len(defs))
	for i, w := range defs {
		copyList := append([]int(nil), w.list...)
		var allowed map[int]struct{}
		if w.typ == 0 {
			allowed = make(map[int]struct{}, len(copyList))
			for _, v := range copyList {
				allowed[v] = struct{}{}
			}
		}
		weapons[i] = weaponDef{
			typ:     w.typ,
			list:    copyList,
			allowed: allowed,
			l:       w.l,
			r:       w.r,
			triple:  w.triple,
		}
	}
	tc := testCase{
		name:    name,
		n:       len(defs),
		m:       m,
		weapons: weapons,
	}
	tc.input = formatInput(tc)
	return tc
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	sb.Grow(tc.n*20 + tc.m*3 + 32)
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, w := range tc.weapons {
		switch w.typ {
		case 0:
			sb.WriteString(fmt.Sprintf("0 %d", len(w.list)))
			for _, v := range w.list {
				sb.WriteByte(' ')
				sb.WriteString(strconv.Itoa(v))
			}
			sb.WriteByte('\n')
		case 1:
			sb.WriteString(fmt.Sprintf("1 %d %d\n", w.l, w.r))
		case 2:
			sb.WriteString(fmt.Sprintf("2 %d %d %d\n", w.triple[0], w.triple[1], w.triple[2]))
		}
	}
	return sb.String()
}

func availableShips(used []bool) []int {
	res := make([]int, 0, len(used))
	for i := 1; i < len(used); i++ {
		if !used[i] {
			res = append(res, i)
		}
	}
	return res
}

func sampleUnique(rng *rand.Rand, m, k int) []int {
	if k > m {
		k = m
	}
	perm := rng.Perm(m)
	res := make([]int, k)
	for i := 0; i < k; i++ {
		res[i] = perm[i] + 1
	}
	sort.Ints(res)
	return res
}
