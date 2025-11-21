package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type photoScenario struct {
	h, w, k int
	pieceH  int
	pieces  [][][]int
}

type photoTestCase struct {
	input    string
	cases    []photoScenario
	bestCost []int
	costs    [][][]int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		lines, err := readOutputLines(out)
		if err != nil {
			fmt.Printf("test %d output error: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, tc.input, out)
			os.Exit(1)
		}
		if len(lines) != len(tc.cases) {
			fmt.Printf("test %d failed: expected %d lines got %d\ninput:\n%s\noutput:\n%s\n", idx+1, len(tc.cases), len(lines), tc.input, out)
			os.Exit(1)
		}
		for i, line := range lines {
			ans, err := parsePermutation(line, tc.cases[i].k)
			if err != nil {
				fmt.Printf("test %d scenario %d invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, i+1, err, tc.input, out)
				os.Exit(1)
			}
			if err := checkScenario(tc.cases[i], tc.costs[i], tc.bestCost[i], ans); err != nil {
				fmt.Printf("test %d scenario %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, i+1, err, tc.input, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}

func readOutputLines(out string) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines, nil
}

func parsePermutation(line string, k int) ([]int, error) {
	fields := strings.Fields(line)
	if len(fields) != k {
		return nil, fmt.Errorf("expected %d numbers got %d", k, len(fields))
	}
	perm := make([]int, k)
	seen := make([]bool, k+1)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		if val < 1 || val > k {
			return nil, fmt.Errorf("value %d out of range [1,%d]", val, k)
		}
		if seen[val] {
			return nil, fmt.Errorf("value %d appears multiple times", val)
		}
		seen[val] = true
		perm[i] = val
	}
	return perm, nil
}

func checkScenario(sc photoScenario, cost [][]int, best int, perm []int) error {
	k := sc.k
	order := make([]int, k)
	posUsed := make([]bool, k)
	for piece, position := range perm {
		pos := position - 1
		if posUsed[pos] {
			return fmt.Errorf("position %d assigned multiple times", pos+1)
		}
		posUsed[pos] = true
		order[pos] = piece
	}
	total := 0
	for i := 0; i < k-1; i++ {
		total += cost[order[i]][order[i+1]]
	}
	if total != best {
		return fmt.Errorf("expected minimal cost %d but got %d", best, total)
	}
	return nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []photoTestCase {
	rand.Seed(42)
	var tests []photoTestCase
	// simple deterministic cases
	tests = append(tests, buildTest([]photoScenario{
		makeScenario(2, 2, 2, 1),
	}))
	for t := 0; t < 100; t++ {
		numCases := rand.Intn(3) + 1
		scs := make([]photoScenario, numCases)
		for i := 0; i < numCases; i++ {
			k := rand.Intn(4) + 2
			pieceH := rand.Intn(3) + 1
			w := rand.Intn(4) + 1
			h := pieceH * k
			scs[i] = randomScenario(h, w, k, pieceH)
		}
		tests = append(tests, buildTest(scs))
	}
	return tests
}

func makeScenario(h, w, k, pieceH int) photoScenario {
	pieces := make([][][]int, k)
	for i := 0; i < k; i++ {
		pieces[i] = make([][]int, pieceH)
		for r := 0; r < pieceH; r++ {
			row := make([]int, w)
			for c := 0; c < w; c++ {
				row[c] = (i + r + c) % 10
			}
			pieces[i][r] = row
		}
	}
	return photoScenario{h: h, w: w, k: k, pieceH: pieceH, pieces: pieces}
}

func randomScenario(h, w, k, pieceH int) photoScenario {
	pieces := make([][][]int, k)
	for i := 0; i < k; i++ {
		pieces[i] = make([][]int, pieceH)
		for r := 0; r < pieceH; r++ {
			row := make([]int, w)
			for c := 0; c < w; c++ {
				row[c] = rand.Intn(256)
			}
			pieces[i][r] = row
		}
	}
	return photoScenario{h: h, w: w, k: k, pieceH: pieceH, pieces: pieces}
}

func buildTest(cases []photoScenario) photoTestCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	best := make([]int, len(cases))
	costs := make([][][]int, len(cases))
	for idx, sc := range cases {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", sc.h, sc.w, sc.k))
		for i := 0; i < sc.k; i++ {
			for r := 0; r < sc.pieceH; r++ {
				for c := 0; c < sc.w; c++ {
					if c > 0 {
						sb.WriteByte(' ')
					}
					sb.WriteString(strconv.Itoa(sc.pieces[i][r][c]))
				}
				sb.WriteByte('\n')
			}
		}
		_, bestCost, costMatrix := solveScenario(sc)
		best[idx] = bestCost
		costs[idx] = costMatrix
	}
	return photoTestCase{
		input:    sb.String(),
		cases:    cases,
		bestCost: best,
		costs:    costs,
	}
}

func solveScenario(sc photoScenario) ([]int, int, [][]int) {
	k := sc.k
	w := sc.w
	pieceH := sc.pieceH
	top := make([][]int, k)
	bottom := make([][]int, k)
	for i := 0; i < k; i++ {
		top[i] = make([]int, w)
		bottom[i] = make([]int, w)
		for r := 0; r < pieceH; r++ {
			row := sc.pieces[i][r]
			for c := 0; c < w; c++ {
				val := row[c]
				if r == 0 {
					top[i][c] = val
				}
				if r == pieceH-1 {
					bottom[i][c] = val
				}
			}
		}
	}
	cost := make([][]int, k)
	for i := 0; i < k; i++ {
		cost[i] = make([]int, k)
		for j := 0; j < k; j++ {
			if i == j {
				continue
			}
			diff := 0
			for c := 0; c < w; c++ {
				d := bottom[i][c] - top[j][c]
				if d < 0 {
					d = -d
				}
				diff += d
			}
			cost[i][j] = diff
		}
	}
	nMask := 1 << k
	const inf = int(1e9)
	dp := make([][]int, nMask)
	prev := make([][]int, nMask)
	for mask := 0; mask < nMask; mask++ {
		dp[mask] = make([]int, k)
		prev[mask] = make([]int, k)
		for j := 0; j < k; j++ {
			dp[mask][j] = inf
			prev[mask][j] = -1
		}
	}
	for i := 0; i < k; i++ {
		dp[1<<i][i] = 0
	}
	for mask := 1; mask < nMask; mask++ {
		for last := 0; last < k; last++ {
			if dp[mask][last] == inf {
				continue
			}
			for nxt := 0; nxt < k; nxt++ {
				if mask&(1<<nxt) != 0 {
					continue
				}
				nm := mask | (1 << nxt)
				cand := dp[mask][last] + cost[last][nxt]
				if cand < dp[nm][nxt] {
					dp[nm][nxt] = cand
					prev[nm][nxt] = last
				}
			}
		}
	}
	full := nMask - 1
	best := inf
	bestLast := 0
	for last := 0; last < k; last++ {
		if dp[full][last] < best {
			best = dp[full][last]
			bestLast = last
		}
	}
	order := make([]int, k)
	mask := full
	last := bestLast
	for idx := k - 1; idx >= 0; idx-- {
		order[idx] = last
		pl := prev[mask][last]
		mask ^= 1 << last
		last = pl
		if mask == 0 {
			break
		}
	}
	ans := make([]int, k)
	for pos, piece := range order {
		ans[piece] = pos + 1
	}
	return ans, best, cost
}
