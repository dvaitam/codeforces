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

func runBin(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

// can returns true if we can place k teleports in [0,x] each at distance >= R from every friend.
func can(R int64, friends []int64, x int64, k int64) bool {
	if R == 0 {
		return x+1 >= k
	}
	spread := R - 1
	prevEnd := int64(-1)
	available := int64(0)
	for _, a := range friends {
		left := a - spread
		right := a + spread
		if left < 0 {
			left = 0
		}
		if right > x {
			right = x
		}
		if left > prevEnd+1 {
			available += left - (prevEnd + 1)
			if available >= k {
				return true
			}
		}
		if right > prevEnd {
			prevEnd = right
		}
	}
	if prevEnd < x {
		available += x - prevEnd
	}
	return available >= k
}

func optimalR(friends []int64, x int64, k int64) int64 {
	lo, hi := int64(0), x+1
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if can(mid, friends, x, k) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

type testCase struct {
	n, k    int
	x       int64
	friends []int64
}

func (tc testCase) Input() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d\n", tc.n, tc.k, tc.x)
	for i, f := range tc.friends {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(f, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func verify(tc testCase, output string) error {
	sorted := make([]int64, len(tc.friends))
	copy(sorted, tc.friends)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	R := optimalR(sorted, tc.x, int64(tc.k))

	fields := strings.Fields(output)
	if len(fields) != tc.k {
		return fmt.Errorf("expected %d positions, got %d", tc.k, len(fields))
	}
	seen := make(map[int64]bool, tc.k)
	teleports := make([]int64, tc.k)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer %q", f)
		}
		if v < 0 || v > tc.x {
			return fmt.Errorf("position %d out of range [0, %d]", v, tc.x)
		}
		if seen[v] {
			return fmt.Errorf("duplicate position %d", v)
		}
		seen[v] = true
		teleports[i] = v
	}
	sort.Slice(teleports, func(i, j int) bool { return teleports[i] < teleports[j] })

	// Compute min over friends of (distance to nearest teleport).
	actualMin := int64(1) << 62
	for _, f := range tc.friends {
		idx := sort.Search(len(teleports), func(i int) bool { return teleports[i] >= f })
		if idx < len(teleports) {
			if d := teleports[idx] - f; d < actualMin {
				actualMin = d
			}
		}
		if idx > 0 {
			if d := f - teleports[idx-1]; d < actualMin {
				actualMin = d
			}
		}
	}

	if actualMin < R {
		return fmt.Errorf("achieved min dist %d < optimal %d", actualMin, R)
	}
	return nil
}

func sampleTests() []testCase {
	return []testCase{
		{4, 1, 4, []int64{1, 0, 2, 4}},
		{5, 5, 4, []int64{0, 1, 2, 3, 4}},
		{2, 1, 4, []int64{4, 0}},
		{3, 4, 6, []int64{2, 4, 3}},
		{3, 2, 12, []int64{6, 12, 0}},
		{4, 3, 12, []int64{8, 12, 0, 4}},
		{1, 1, 1000000000, []int64{1}},
		{1, 1, 1000000001, []int64{0}},
		{3, 4, 9, []int64{8, 7, 9}},
		{3, 4, 9, []int64{2, 0, 1}},
	}
}

func genTests(rng *rand.Rand, count int) []testCase {
	tests := make([]testCase, count)
	for i := range tests {
		x := int64(rng.Intn(30) + 1)
		k := int(rng.Intn(int(x)+1) + 1)
		n := rng.Intn(8) + 1
		friends := make([]int64, n)
		for j := range friends {
			friends[j] = int64(rng.Intn(int(x) + 1))
		}
		tests[i] = testCase{n: n, k: k, x: x, friends: friends}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	bin := os.Args[1]

	allTests := sampleTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	allTests = append(allTests, genTests(rng, 200)...)

	for i, tc := range allTests {
		input := tc.Input()
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := verify(tc, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(allTests))
}
