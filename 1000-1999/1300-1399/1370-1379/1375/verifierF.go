package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// 1375F is an interactive problem.
// First player chooses y, Second player adds y to a pile (cannot repeat same pile consecutively).
// Second player loses if any two piles become equal.
// First player loses after 1000 turns.
// Candidate (Harris) chooses First or Second to guarantee a win.
//
// This verifier plays as the adversary (Second player) using various strategies
// and checks that the candidate (First player) can force two piles equal.

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	numTests := 50

	for t := 0; t < numTests; t++ {
		// Generate 3 distinct piles
		a := [3]int64{
			rng.Int63n(1000) + 1,
			rng.Int63n(1000) + 1,
			rng.Int63n(1000) + 1,
		}
		// Ensure distinct
		for a[0] == a[1] || a[1] == a[2] || a[0] == a[2] {
			a[0] = rng.Int63n(1000) + 1
			a[1] = rng.Int63n(1000) + 1
			a[2] = rng.Int63n(1000) + 1
		}

		// Generate all valid adversary strategies (no consecutive repeats) up to depth 3
		// Each strategy is a sequence of pile choices; we prune invalid (consecutive repeat) ones.
		// We try multiple strategies. The candidate should win against all of them.
		strategies := generateStrategies(3)
		for _, strat := range strategies {
			err := simulate(binary, a, strat)
			if err != nil {
				fmt.Fprintf(os.Stderr, "test %d failed (piles=%v strat=%v): %v\n", t+1, a, strat, err)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", numTests)
}

// generateStrategies generates all valid adversary strategies of length maxTurns
// where no two consecutive choices are the same pile.
func generateStrategies(maxTurns int) [][]int {
	var result [][]int
	var dfs func(cur []int, last int)
	dfs = func(cur []int, last int) {
		if len(cur) == maxTurns {
			cp := make([]int, len(cur))
			copy(cp, cur)
			result = append(result, cp)
			return
		}
		for p := 0; p < 3; p++ {
			if p == last {
				continue // cannot repeat same pile consecutively
			}
			dfs(append(cur, p), p)
		}
	}
	dfs(nil, -1)
	return result
}

func simulate(binary string, piles [3]int64, strategy []int) error {
	cmd := exec.Command(binary)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %v", err)
	}
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %v", err)
	}

	reader := bufio.NewReader(stdout)
	writer := bufio.NewWriter(stdin)

	// Send piles
	fmt.Fprintf(writer, "%d %d %d\n", piles[0], piles[1], piles[2])
	writer.Flush()

	a := piles

	// Read "First" or "Second"
	line, err := readLine(reader)
	if err != nil {
		stdin.Close()
		cmd.Wait()
		return fmt.Errorf("reading First/Second: %v", err)
	}
	role := strings.TrimSpace(line)
	if role != "First" && role != "Second" {
		stdin.Close()
		cmd.Wait()
		return fmt.Errorf("expected First or Second, got %q", role)
	}

	if role != "First" {
		// Most accepted solutions play as First. If candidate says Second,
		// just skip this strategy (we only test First player strategies).
		stdin.Close()
		cmd.Wait()
		return nil
	}

	// First player makes moves; we play as adversary (Second player)
	won := false
	for turn := 0; turn < len(strategy); turn++ {
		// Read y from candidate
		line, err := readLine(reader)
		if err != nil {
			stdin.Close()
			cmd.Wait()
			if won {
				return nil
			}
			return fmt.Errorf("reading y on turn %d: %v", turn+1, err)
		}
		y, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
		if err != nil {
			stdin.Close()
			cmd.Wait()
			return fmt.Errorf("invalid y on turn %d: %v", turn+1, err)
		}
		if y <= 0 {
			stdin.Close()
			cmd.Wait()
			return fmt.Errorf("y must be positive, got %d on turn %d", y, turn+1)
		}

		// Adversary adds y to pile strategy[turn]
		pile := strategy[turn]
		a[pile] += y

		// Check if two piles are now equal => Second player (adversary) loses, candidate wins
		if a[0] == a[1] || a[1] == a[2] || a[0] == a[2] {
			won = true
			// Send pile choice so candidate can read it, then close
			fmt.Fprintf(writer, "%d\n", pile+1)
			writer.Flush()
			stdin.Close()
			cmd.Wait()
			return nil
		}

		// Send pile choice to candidate (1-indexed)
		fmt.Fprintf(writer, "%d\n", pile+1)
		writer.Flush()
	}

	stdin.Close()
	cmd.Wait()

	if !won {
		return fmt.Errorf("piles not equal after %d turns: %v", len(strategy), a)
	}
	return nil
}

func readLine(r *bufio.Reader) (string, error) {
	var sb strings.Builder
	for {
		line, isPrefix, err := r.ReadLine()
		if err != nil {
			if err == io.EOF && sb.Len() > 0 {
				return sb.String(), nil
			}
			return "", err
		}
		sb.Write(line)
		if !isPrefix {
			return sb.String(), nil
		}
	}
}
