package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func probability(order []byte, R, C int) float64 {
	n := R + C
	L := 2 * n
	dp := make([][]float64, R+1)
	for i := range dp {
		dp[i] = make([]float64, R+1)
	}
	dp[0][0] = 1
	var winA float64
	for step := 0; step < L; step++ {
		next := make([][]float64, R+1)
		for i := range next {
			next[i] = make([]float64, R+1)
		}
		player := order[step]
		remaining := L - step
		for a := 0; a < R; a++ {
			for b := 0; b < R; b++ {
				prob := dp[a][b]
				if prob == 0 {
					continue
				}
				rawUsed := a + b
				rawLeft := 2*R - rawUsed
				pRaw := float64(rawLeft) / float64(remaining)
				pCook := 1 - pRaw
				if player == 'A' {
					if a+1 >= R {
						winA += 0
					} else {
						next[a+1][b] += prob * pRaw
					}
					next[a][b] += prob * pCook
				} else {
					if b+1 >= R {
						winA += prob * pRaw
					} else {
						next[a][b+1] += prob * pRaw
					}
					next[a][b] += prob * pCook
				}
			}
		}
		dp = next
	}
	return winA
}

func solve(R, C int, S string) int64 {
	n := R + C
	L := 2 * n

	// First pass: find global minimum unfairness across ALL valid orderings
	globalBest := 1e9
	allQ := make([]byte, L)
	for i := range allQ {
		allQ[i] = '?'
	}
	allQStr := string(allQ)
	seqAll := make([]byte, L)
	var dfsAll func(int, int, int)
	dfsAll = func(pos, a, b int) {
		if pos == L {
			if a == n && b == n {
				pA := probability(seqAll, R, C)
				diff := math.Abs(pA - (1 - pA))
				if diff < globalBest-1e-12 {
					globalBest = diff
				}
			}
			return
		}
		if (allQStr[pos] == 'A' || allQStr[pos] == '?') && a < n {
			seqAll[pos] = 'A'
			dfsAll(pos+1, a+1, b)
		}
		if (allQStr[pos] == 'B' || allQStr[pos] == '?') && b < n {
			seqAll[pos] = 'B'
			dfsAll(pos+1, a, b+1)
		}
	}
	dfsAll(0, 0, 0)

	// Second pass: count orderings matching S that achieve globalBest
	count := int64(0)
	seq := make([]byte, L)
	var dfs func(int, int, int)
	dfs = func(pos, a, b int) {
		if pos == L {
			if a == n && b == n {
				pA := probability(seq, R, C)
				diff := math.Abs(pA - (1 - pA))
				if math.Abs(diff-globalBest) <= 1e-12 {
					count++
				}
			}
			return
		}
		if (S[pos] == 'A' || S[pos] == '?') && a < n {
			seq[pos] = 'A'
			dfs(pos+1, a+1, b)
		}
		if (S[pos] == 'B' || S[pos] == '?') && b < n {
			seq[pos] = 'B'
			dfs(pos+1, a, b+1)
		}
	}
	dfs(0, 0, 0)
	return count
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		R := rng.Intn(3) + 1
		C := rng.Intn(3) + 1
		n := R + C
		length := 2 * n
		chars := []byte{'A', 'B', '?'}
		b := make([]byte, length)
		for i := 0; i < length; i++ {
			b[i] = chars[rng.Intn(len(chars))]
		}
		S := string(b)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n%s\n", R, C, S)
		expected := fmt.Sprintf("%d", solve(R, C, S))
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected:\n%s\n---\ngot:\n%s\n", t+1, sb.String(), expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
