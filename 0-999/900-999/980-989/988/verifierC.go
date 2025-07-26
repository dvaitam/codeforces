package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	seqs [][]int64
}

func solve(seqs [][]int64) string {
	type occ struct{ seq, idx int }
	occMap := make(map[int64]occ)
	for i, seq := range seqs {
		var sum int64
		for _, v := range seq {
			sum += v
		}
		for j, v := range seq {
			exc := sum - v
			if prev, ok := occMap[exc]; ok {
				if prev.seq != i {
					return fmt.Sprintf("YES\n%d %d\n%d %d\n", prev.seq+1, prev.idx+1, i+1, j+1)
				}
			} else {
				occMap[exc] = occ{i, j}
			}
		}
	}
	return "NO\n"
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.seqs))
	for _, seq := range tc.seqs {
		fmt.Fprintf(&sb, "%d\n", len(seq))
		for i, v := range seq {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expected := strings.TrimSpace(solve(tc.seqs))
	got := strings.TrimSpace(out.String())
	if expected != got {
		return fmt.Errorf("expected:\n%s\n-- got:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	var cases []testCase
	// simple deterministic cases
	cases = append(cases, testCase{seqs: [][]int64{{1}, {1}}})
	cases = append(cases, testCase{seqs: [][]int64{{2, 3, 1, 3, 2}, {1, 1, 2, 2, 2, 1}}})

	rng := rand.New(rand.NewSource(2))
	for i := 0; i < 100; i++ {
		k := rng.Intn(4) + 2 // at least 2 sequences
		seqs := make([][]int64, k)
		total := 0
		for j := 0; j < k; j++ {
			n := rng.Intn(4) + 1
			total += n
			seq := make([]int64, n)
			for t := 0; t < n; t++ {
				seq[t] = int64(rng.Intn(21) - 10)
			}
			seqs[j] = seq
		}
		if total > 20 {
			i--
			continue
		}
		cases = append(cases, testCase{seqs: seqs})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", i+1, err, buildInput(tc))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
