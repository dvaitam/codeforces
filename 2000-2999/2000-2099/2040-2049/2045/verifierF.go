package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type cell struct {
	r, c, a int64
}

type gameCase struct {
	N   int64
	M   int
	K   int64
	pts []cell
}

type testCase struct {
	name     string
	input    string
	expected string
}

// solveF implements the game theory directly:
// group cells by r%(K+1), XOR stone counts mod (K+1) per group;
// if any group has a non-zero XOR, Anda (first player) wins.
func solveF(cs gameCase) string {
	k := cs.K + 1
	mp := make(map[int64]int64)
	for _, p := range cs.pts {
		r := p.r % k
		x := p.a % k
		mp[r] ^= x
	}
	for _, v := range mp {
		if v != 0 {
			return "Anda"
		}
	}
	return "Kamu"
}

func packCases(name string, cases []gameCase) testCase {
	var inp, exp strings.Builder
	fmt.Fprintf(&inp, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&inp, "%d %d %d\n", cs.N, cs.M, cs.K)
		for _, p := range cs.pts {
			fmt.Fprintf(&inp, "%d %d %d\n", p.r, p.c, p.a)
		}
		exp.WriteString(solveF(cs))
		exp.WriteByte('\n')
	}
	return testCase{name: name, input: inp.String(), expected: strings.TrimSpace(exp.String())}
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func sampleTest() testCase {
	cases := []gameCase{
		{N: 2, M: 2, K: 4, pts: []cell{{1, 1, 3}, {2, 1, 2}}},
		{N: 100, M: 2, K: 1, pts: []cell{{4, 1, 10}, {4, 4, 10}}},
		{N: 10, M: 5, K: 2, pts: []cell{{1, 1, 4}, {3, 1, 2}, {4, 2, 5}, {2, 2, 1}, {5, 3, 4}}},
	}
	return packCases("sample", cases)
}

func singlePileTests() testCase {
	cases := []gameCase{
		{N: 1, M: 1, K: 1, pts: []cell{{1, 1, 1}}},
		{N: 5, M: 1, K: 2, pts: []cell{{3, 2, 7}}},
		{N: 1_000_000_000, M: 1, K: 200000, pts: []cell{{1, 1, 1}}},
		{N: 1_000_000_000, M: 1, K: 3, pts: []cell{{1_000_000_000, 1, 5}}},
	}
	return packCases("single-pile", cases)
}

func layeredRowsTest() testCase {
	cases := []gameCase{
		{
			N: 15, M: 4, K: 3,
			pts: []cell{{15, 1, 3}, {12, 5, 4}, {9, 2, 6}, {6, 3, 8}},
		},
		{
			N: 50, M: 6, K: 7,
			pts: []cell{{5, 1, 10}, {10, 5, 12}, {20, 7, 15}, {30, 10, 9}, {40, 15, 20}, {45, 25, 11}},
		},
	}
	return packCases("layered-rows", cases)
}

func randomCase(rng *rand.Rand, maxM int, maxN int64) gameCase {
	N := int64(1 + rng.Intn(int(maxN)))
	// A triangular grid of size N has N*(N+1)/2 distinct cells.
	// Cap M so the unique-cell sampling loop always terminates.
	available := N * (N + 1) / 2
	if available > int64(maxM) {
		available = int64(maxM)
	}
	M := int(1 + rng.Intn(int(available)))
	K := int64(1 + rng.Intn(200000))
	pts := make([]cell, 0, M)
	used := make(map[int64]struct{})
	for len(pts) < M {
		r := int64(1 + rng.Intn(int(N)))
		c := int64(1 + rng.Intn(int(r)))
		key := (r << 32) ^ c
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		a := int64(1 + rng.Intn(1_000_000_000))
		pts = append(pts, cell{r, c, a})
	}
	return gameCase{N: N, M: M, K: K, pts: pts}
}

func randomPack(name string, rng *rand.Rand, t int, maxM int, maxN int64) testCase {
	cases := make([]gameCase, 0, t)
	for i := 0; i < t; i++ {
		cases = append(cases, randomCase(rng, maxM, maxN))
	}
	return packCases(name, cases)
}

func skewedHeavyCase(rng *rand.Rand) testCase {
	N := int64(1_000_000_000)
	K := int64(200000)
	M := 2000
	pts := make([]cell, 0, M)
	used := make(map[int64]struct{})
	base := N - 5000
	for len(pts) < M {
		r := base + int64(rng.Intn(5000))
		c := int64(1 + rng.Intn(int(r)))
		key := (r << 32) ^ c
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		a := int64(1 + rng.Intn(1_000_000_000))
		pts = append(pts, cell{r, c, a})
	}
	return packCases("skewed-heavy", []gameCase{{N: N, M: M, K: K, pts: pts}})
}

func buildTests(rng *rand.Rand) []testCase {
	return []testCase{
		sampleTest(),
		singlePileTests(),
		layeredRowsTest(),
		randomPack("random-small", rng, 10, 10, 10),
		randomPack("random-mid", rng, 8, 2000, 8000),
		randomPack("random-largeK", rng, 5, 50000, 120000),
		skewedHeavyCase(rng),
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	rng := rand.New(rand.NewSource(2045007))
	tests := buildTests(rng)

	for i, tc := range tests {
		got, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on case %d (%s): %v\ninput:\n%s", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "wrong answer on case %d (%s)\nexpected:\n%s\ngot:\n%s\ninput:\n%s",
				i+1, tc.name, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("Accepted (%d tests)\n", len(tests))
}
