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

func runCandidate(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCase struct {
	n int
	k int
}

func expectedSets(tc testCase) (int, []string) {
	m := (6*tc.n - 1) * tc.k
	sets := make([]string, tc.n)
	for i := 0; i < tc.n; i++ {
		x := 1 + 6*i
		nums := []int{x * tc.k, (x + 2) * tc.k, (x + 4) * tc.k}
		if (x+1)%3 != 0 {
			nums = append(nums, (x+1)*tc.k)
		} else {
			nums = append(nums, (x+3)*tc.k)
		}
		sort.Ints(nums)
		sets[i] = fmt.Sprintf("%d %d %d %d", nums[0], nums[1], nums[2], nums[3])
	}
	sort.Strings(sets)
	return m, sets
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	k := rng.Intn(20) + 1
	return testCase{n, k}
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d\n", tc.n, tc.k)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != tc.n+1 {
		return fmt.Errorf("expected %d lines, got %d", tc.n+1, len(lines))
	}

	var mGot int
	if _, err := fmt.Sscan(lines[0], &mGot); err != nil {
		return fmt.Errorf("cannot parse m: %v", err)
	}
	mExp, expSets := expectedSets(tc)
	if mGot != mExp {
		return fmt.Errorf("expected m %d, got %d", mExp, mGot)
	}

	candSets := make([]string, tc.n)
	for i := 0; i < tc.n; i++ {
		var a, b, c, d int
		if _, err := fmt.Sscan(lines[i+1], &a, &b, &c, &d); err != nil {
			return fmt.Errorf("cannot parse set %d: %v", i+1, err)
		}
		nums := []int{a, b, c, d}
		sort.Ints(nums)
		candSets[i] = fmt.Sprintf("%d %d %d %d", nums[0], nums[1], nums[2], nums[3])
	}
	sort.Strings(candSets)

	for i := range expSets {
		if expSets[i] != candSets[i] {
			return fmt.Errorf("expected set %s but got %s", expSets[i], candSets[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{{1, 1}, {2, 3}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
