package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	val int64
	id  int
}

type testCase struct {
	name      string
	n         int
	a, b      int64
	dist      []int64
	expectYes bool
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()

	for idx, tc := range tests {
		input := buildInputStr(tc)

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		if _, err := parseSolution(tc, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func runProgram(bin, input string) (string, error) {
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
	return out.String(), err
}

func parseSolution(tc testCase, output string) (bool, error) {
	lines := strings.Fields(output)
	if len(lines) == 0 {
		return false, errors.New("empty output")
	}
	ans := strings.ToUpper(lines[0])
	if ans != "YES" && ans != "NO" {
		return false, fmt.Errorf("first token must be YES or NO, got %q", lines[0])
	}
	if ans == "NO" {
		if tc.expectYes {
			return false, errors.New("reported NO but a solution is expected")
		}
		if len(lines) > 1 {
			return false, errors.New("extra tokens after NO")
		}
		return false, nil
	}

	if !tc.expectYes {
		return false, errors.New("reported YES but no solution is expected")
	}
	expectedTokens := 1 + 2*tc.n
	if len(lines) != expectedTokens {
		return false, fmt.Errorf("expected %d coordinate tokens, got %d", expectedTokens, len(lines))
	}
	coords := make([][2]int64, tc.n)
	for i := 0; i < tc.n; i++ {
		x, err := parseInt64(lines[1+2*i])
		if err != nil {
			return false, fmt.Errorf("invalid x_%d: %v", i+1, err)
		}
		y, err := parseInt64(lines[2+2*i])
		if err != nil {
			return false, fmt.Errorf("invalid y_%d: %v", i+1, err)
		}
		coords[i] = [2]int64{x, y}
	}

	if coords[0][0] != 0 || coords[0][1] != 0 {
		return false, fmt.Errorf("first statue must be at (0,0), got (%d,%d)", coords[0][0], coords[0][1])
	}
	last := coords[tc.n-1]
	if last[0] != tc.a || last[1] != tc.b {
		return false, fmt.Errorf("last statue must be at (%d,%d), got (%d,%d)", tc.a, tc.b, last[0], last[1])
	}
	for i := 0; i < tc.n-1; i++ {
		got := manhattan(coords[i], coords[i+1])
		if got != tc.dist[i] {
			return false, fmt.Errorf("distance between statue %d and %d expected %d, got %d", i+1, i+2, tc.dist[i], got)
		}
	}
	return true, nil
}

func parseInt64(s string) (int64, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func manhattan(a, b [2]int64) int64 {
	return abs64(a[0]-b[0]) + abs64(a[1]-b[1])
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func buildInputStr(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.a, tc.b))
	for i, v := range tc.dist {
		sb.WriteString(fmt.Sprintf("%d", v))
		if i+1 < len(tc.dist) {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

// solveRef is the embedded correct reference solver for 2068H.
func solveRef(tc testCase) (bool, [][2]int64) {
	n := tc.n
	a := tc.a
	b := tc.b
	d := tc.dist

	U := a + b
	V := a - b

	D := int64(0)
	for _, v := range d {
		D += v
	}

	if (D-U)%2 != 0 {
		return false, nil
	}

	Ru := (D - U) / 2
	Rv := (D - V) / 2

	if Ru < 0 || Ru > D || Rv < 0 || Rv > D {
		return false, nil
	}

	items := make([]Item, n-1)
	for i := 0; i < n-1; i++ {
		items[i] = Item{val: d[i], id: i}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].val > items[j].val
	})

	ansArr := make([]int, n-1)

	var dfs func(idx int, w1, w2, w3, w4 int64) bool
	dfs = func(idx int, w1, w2, w3, w4 int64) bool {
		if idx == len(items) {
			return true
		}
		val := items[idx].val
		origId := items[idx].id

		if w1 >= val {
			ansArr[origId] = 1
			if dfs(idx+1, w1-val, w2, w3, w4) {
				return true
			}
		}
		if w2 >= val && w2 != w1 {
			ansArr[origId] = 2
			if dfs(idx+1, w1, w2-val, w3, w4) {
				return true
			}
		}
		if w3 >= val && w3 != w1 && w3 != w2 {
			ansArr[origId] = 3
			if dfs(idx+1, w1, w2, w3-val, w4) {
				return true
			}
		}
		if w4 >= val && w4 != w1 && w4 != w2 && w4 != w3 {
			ansArr[origId] = 4
			if dfs(idx+1, w1, w2, w3, w4-val) {
				return true
			}
		}
		return false
	}

	if !dfs(0, Ru, D-Ru, Rv, D-Rv) {
		return false, nil
	}

	sumA1 := int64(0)
	sumB1 := int64(0)

	k := make([]int64, len(d))
	j := make([]int64, len(d))

	for i := 0; i < len(d); i++ {
		switch ansArr[i] {
		case 1:
			k[i] = d[i]
			sumA1 += d[i]
		case 2:
			k[i] = 0
		case 3:
			j[i] = d[i]
			sumB1 += d[i]
		case 4:
			j[i] = 0
		}
	}

	remU := Ru - sumA1
	for i := 0; i < len(d); i++ {
		if ansArr[i] == 3 || ansArr[i] == 4 {
			take := d[i]
			if remU < take {
				take = remU
			}
			k[i] = take
			remU -= take
		}
	}

	remV := Rv - sumB1
	for i := 0; i < len(d); i++ {
		if ansArr[i] == 1 || ansArr[i] == 2 {
			take := d[i]
			if remV < take {
				take = remV
			}
			j[i] = take
			remV -= take
		}
	}

	coords := make([][2]int64, n)
	x, y := int64(0), int64(0)
	coords[0] = [2]int64{x, y}
	for i := 0; i < len(d); i++ {
		du := d[i] - 2*k[i]
		dv := d[i] - 2*j[i]
		dx := (du + dv) / 2
		dy := (du - dv) / 2
		x += dx
		y += dy
		coords[i+1] = [2]int64{x, y}
	}

	return true, coords
}

func buildTests() []testCase {
	tests := []testCase{
		{
			name:      "sample-no",
			n:         3,
			a:         5,
			b:         8,
			dist:      []int64{9, 0},
			expectYes: false,
		},
		{
			name:      "sample-yes",
			n:         4,
			a:         10,
			b:         6,
			dist:      []int64{7, 8, 5},
			expectYes: true,
		},
		{
			name:      "stay-origin",
			n:         3,
			a:         0,
			b:         0,
			dist:      []int64{0, 0},
			expectYes: true,
		},
		{
			name:      "parity-impossible",
			n:         3,
			a:         1,
			b:         0,
			dist:      []int64{2, 2},
			expectYes: false,
		},
		{
			name:      "long-line",
			n:         5,
			a:         15,
			b:         0,
			dist:      []int64{3, 4, 5, 3},
			expectYes: true,
		},
	}

	// Determine expectYes using the embedded solver for deterministic tests
	for i := range tests {
		yes, _ := solveRef(tests[i])
		tests[i].expectYes = yes
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 5; i++ {
		tc := randomFeasible(rng, 5+i)
		tests = append(tests, tc)
	}
	for i := 0; i < 3; i++ {
		tc := randomImpossible(rng, 6+i)
		tests = append(tests, tc)
	}

	return tests
}

func randomFeasible(rng *rand.Rand, n int) testCase {
	if n < 3 {
		n = 3
	}
	coords := make([][2]int64, n)
	dist := make([]int64, n-1)
	for i := 1; i < n; i++ {
		l := int64(rng.Intn(6))
		var dx, dy int64
		dir := rng.Intn(4)
		switch dir {
		case 0:
			dx, dy = l, 0
		case 1:
			dx, dy = -l, 0
		case 2:
			dx, dy = 0, l
		default:
			dx, dy = 0, -l
		}
		coords[i][0] = coords[i-1][0] + dx
		coords[i][1] = coords[i-1][1] + dy
		dist[i-1] = abs64(dx) + abs64(dy)
	}
	return testCase{
		name:      fmt.Sprintf("feasible-%d", rng.Int()),
		n:         n,
		a:         coords[n-1][0],
		b:         coords[n-1][1],
		dist:      dist,
		expectYes: true,
	}
}

func randomImpossible(rng *rand.Rand, n int) testCase {
	if n < 3 {
		n = 3
	}
	dist := make([]int64, n-1)
	var sum int64
	for i := 0; i < n-1; i++ {
		di := int64(rng.Intn(10))
		dist[i] = di
		sum += di
	}
	a := sum + 5
	return testCase{
		name:      fmt.Sprintf("impossible-%d", rng.Int()),
		n:         n,
		a:         a,
		b:         0,
		dist:      dist,
		expectYes: false,
	}
}
