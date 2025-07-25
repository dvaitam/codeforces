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
	input    string
	expected string
}

type friend struct {
	money  int64
	factor int64
}

func solveCase(friends []friend, d int64) int64 {
	sort.Slice(friends, func(i, j int) bool { return friends[i].money < friends[j].money })
	var best, cur int64
	l := 0
	for r := 0; r < len(friends); r++ {
		cur += friends[r].factor
		for l <= r && friends[r].money-friends[l].money >= d {
			cur -= friends[l].factor
			l++
		}
		if cur > best {
			best = cur
		}
	}
	return best
}

func buildCase(friends []friend, d int64) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", len(friends), d))
	for _, f := range friends {
		sb.WriteString(fmt.Sprintf("%d %d\n", f.money, f.factor))
	}
	ans := solveCase(append([]friend(nil), friends...), d)
	return testCase{input: sb.String(), expected: fmt.Sprintf("%d\n", ans)}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 1
	d := int64(rng.Intn(20) + 1)
	friends := make([]friend, n)
	for i := range friends {
		friends[i].money = int64(rng.Intn(100))
		friends[i].factor = int64(rng.Intn(100))
	}
	return buildCase(friends, d)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{
		buildCase([]friend{{10, 5}}, 1),
		buildCase([]friend{{1, 1}, {2, 2}, {3, 3}}, 5),
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
