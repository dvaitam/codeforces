package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"math/rand"
)

// 1779E is an interactive problem.
// We simulate a judge with a concrete tournament graph and verify the candidate's answer.

// beats[i][j] = true means player i beats player j (0-indexed)
type tournament struct {
	n     int
	beats [][]bool
}

// canWin checks if player p can win a tournament on the set of players represented by 'alive' bitmask.
func (t *tournament) canWin(p int, alive int) bool {
	// Only p alive => p wins
	if alive == (1 << p) {
		return true
	}
	// Try removing some player q that p beats, and check if p can win the remaining
	// Actually, tournament elimination: pick two players, one is eliminated, repeat.
	// p can win if there exists some other player q in alive such that:
	//   either p beats q and p can win alive\{q}
	//   or there exists r != p in alive that beats q, and p can win alive\{q}
	// This is complex. Let's use the known characterization:
	// p can win iff for every non-empty subset S of alive\{p}, there exists
	// a player in alive\S that beats some player in S.
	// Actually simpler: p can win a single-elimination tournament on alive iff
	// p is a "king" of the sub-tournament on alive.
	// More precisely, p can win iff we can eliminate all other players one by one.

	// For small n, let's do BFS/recursion: p can win tournament on 'alive' iff
	// there exist two players a,b in alive where a beats b, and after removing b,
	// p can win on alive\{b}.
	// But this has exponential states. For n<=10, 2^10 * 10 = 10240 states, manageable with memoization.
	return false // placeholder, we use the memoized version below
}

func computeCandidateMasters(t *tournament) []bool {
	n := t.n
	full := (1 << n) - 1

	// memo[mask][p] = 0 (unknown), 1 (can win), 2 (cannot win)
	memo := make([][]int8, 1<<n)
	for i := range memo {
		memo[i] = make([]int8, n)
	}

	var canWinMemo func(p int, alive int) bool
	canWinMemo = func(p int, alive int) bool {
		if alive == (1 << p) {
			return true
		}
		if alive&(1<<p) == 0 {
			return false
		}
		if memo[alive][p] != 0 {
			return memo[alive][p] == 1
		}
		memo[alive][p] = 2 // assume false

		// p can win on 'alive' if there exists some player q != p in alive
		// such that someone in alive beats q (either p or someone else who can
		// be arranged to play q), and then p can win on alive\{q}.
		// More precisely: in a single-elimination tournament, two players are matched,
		// loser is eliminated. p can win iff we can find an elimination order.
		// Equivalent: there exists q in alive, q != p, such that:
		//   (1) there exists r in alive, r != q, with beats[r][q] = true
		//   (2) p can win on alive \ {q}
		for q := 0; q < n; q++ {
			if q == p || alive&(1<<q) == 0 {
				continue
			}
			// Check if someone in alive beats q
			canElimQ := false
			for r := 0; r < n; r++ {
				if r == q || alive&(1<<r) == 0 {
					continue
				}
				if t.beats[r][q] {
					canElimQ = true
					break
				}
			}
			if canElimQ && canWinMemo(p, alive&^(1<<q)) {
				memo[alive][p] = 1
				return true
			}
		}
		return false
	}

	result := make([]bool, n)
	for p := 0; p < n; p++ {
		result[p] = canWinMemo(p, full)
	}
	return result
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

func runInteractive(binary string, tourn *tournament) error {
	n := tourn.n

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

	// Send: 1 test case
	fmt.Fprintf(writer, "1\n")
	fmt.Fprintf(writer, "%d\n", n)
	writer.Flush()

	queryCount := 0

	for {
		line, err := readLine(reader)
		if err != nil {
			stdin.Close()
			cmd.Wait()
			return fmt.Errorf("reading from candidate: %v", err)
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if line[0] == '!' {
			// Answer line
			parts := strings.Fields(line)
			if len(parts) != 2 {
				stdin.Close()
				cmd.Wait()
				return fmt.Errorf("invalid answer format: %q", line)
			}
			ansStr := parts[1]

			expected := computeCandidateMasters(tourn)

			if n == 1 {
				// Special case: n=1, answer should be "1"
				if ansStr != "1" {
					stdin.Close()
					cmd.Wait()
					return fmt.Errorf("wrong answer for n=1: got %q want 1", ansStr)
				}
				stdin.Close()
				cmd.Wait()
				return nil
			}

			if len(ansStr) != n {
				stdin.Close()
				cmd.Wait()
				return fmt.Errorf("answer length %d != n=%d", len(ansStr), n)
			}

			for i := 0; i < n; i++ {
				got := ansStr[i] == '1'
				if got != expected[i] {
					stdin.Close()
					cmd.Wait()
					expStr := make([]byte, n)
					for j := 0; j < n; j++ {
						if expected[j] {
							expStr[j] = '1'
						} else {
							expStr[j] = '0'
						}
					}
					return fmt.Errorf("wrong answer: got %s want %s", ansStr, string(expStr))
				}
			}

			stdin.Close()
			cmd.Wait()
			return nil
		}

		if line[0] == '?' {
			// Query
			queryCount++
			if queryCount > 2*n {
				stdin.Close()
				cmd.Wait()
				return fmt.Errorf("too many queries: %d > 2*%d", queryCount, n)
			}
			parts := strings.Fields(line)
			if len(parts) != 3 {
				stdin.Close()
				cmd.Wait()
				return fmt.Errorf("invalid query format: %q", line)
			}
			player, err := strconv.Atoi(parts[1])
			if err != nil || player < 1 || player > n {
				stdin.Close()
				cmd.Wait()
				return fmt.Errorf("invalid player in query: %q", parts[1])
			}
			mask := parts[2]
			if len(mask) != n {
				stdin.Close()
				cmd.Wait()
				return fmt.Errorf("invalid mask length %d != %d", len(mask), n)
			}

			// Count wins: how many players j (0-indexed) with mask[j]='1' does player (1-indexed) beat?
			wins := 0
			pi := player - 1
			for j := 0; j < n; j++ {
				if mask[j] == '1' && j != pi && tourn.beats[pi][j] {
					wins++
				}
			}

			fmt.Fprintf(writer, "%d\n", wins)
			writer.Flush()
		} else {
			stdin.Close()
			cmd.Wait()
			return fmt.Errorf("unexpected line from candidate: %q", line)
		}
	}
}

func genTournament(rng *rand.Rand, n int) *tournament {
	beats := make([][]bool, n)
	for i := range beats {
		beats[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Intn(2) == 0 {
				beats[i][j] = true
			} else {
				beats[j][i] = true
			}
		}
	}
	return &tournament{n: n, beats: beats}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	rng := rand.New(rand.NewSource(42))

	// Test with small n values to keep brute-force feasible
	testSizes := []int{1, 2, 3, 4, 5, 6, 7, 8}
	numTests := 0

	for _, n := range testSizes {
		reps := 5
		if n <= 3 {
			reps = 3
		}
		for r := 0; r < reps; r++ {
			tourn := genTournament(rng, n)
			if err := runInteractive(binary, tourn); err != nil {
				fmt.Fprintf(os.Stderr, "test %d (n=%d) failed: %v\n", numTests+1, n, err)
				os.Exit(1)
			}
			numTests++
		}
	}
	fmt.Printf("All %d tests passed\n", numTests)
}
