package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
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

func expected(n, m int, votes [][]int) int {
	wins := make([]int, n)
	for i := 0; i < m; i++ {
		bestCand := 0
		bestVotes := votes[i][0]
		for j := 1; j < n; j++ {
			if votes[i][j] > bestVotes {
				bestVotes = votes[i][j]
				bestCand = j
			}
		}
		wins[bestCand]++
	}
	winner := 0
	for i := 1; i < n; i++ {
		if wins[i] > wins[winner] {
			winner = i
		}
	}
	return winner + 1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		votes := make([][]int, m)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for r := 0; r < m; r++ {
			votes[r] = make([]int, n)
			for c := 0; c < n; c++ {
				v := rng.Intn(21)
				votes[r][c] = v
				fmt.Fprintf(&sb, "%d ", v)
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		want := expected(n, m, votes)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil || got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", i+1, want, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
