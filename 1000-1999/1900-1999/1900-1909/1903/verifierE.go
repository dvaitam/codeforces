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

type point struct {
	x, y int
}

func sqDist(a, b point) int64 {
	dx := int64(a.x - b.x)
	dy := int64(a.y - b.y)
	return dx*dx + dy*dy
}

type gameState struct {
	n      int
	start  point
	points []point
	used   []bool
	prev   point
	total  int64
	moves  int
}

func (g *gameState) applyMove(idx int) {
	g.used[idx] = true
	g.total += sqDist(g.prev, g.points[idx])
	g.prev = g.points[idx]
	g.moves++
}

func (g *gameState) availableIndices() []int {
	var res []int
	for i := 0; i < g.n; i++ {
		if !g.used[i] {
			res = append(res, i)
		}
	}
	return res
}

func genPoints(rng *rand.Rand, n int) (point, []point) {
	sx := rng.Intn(100) + 1
	sy := rng.Intn(100) + 1
	pts := make([]point, n)
	for i := 0; i < n; i++ {
		pts[i] = point{rng.Intn(100) + 1, rng.Intn(100) + 1}
	}
	return point{sx, sy}, pts
}

func runInteractiveGame(bin string, rng *rand.Rand, n int, start point, pts []point) error {
	cmd := exec.Command(bin)
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %v", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %v", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %v", err)
	}

	reader := bufio.NewReader(stdoutPipe)
	writer := bufio.NewWriter(stdinPipe)

	// Send the test case
	fmt.Fprintf(writer, "1\n")    // t = 1
	fmt.Fprintf(writer, "%d\n", n)
	fmt.Fprintf(writer, "%d %d\n", start.x, start.y)
	for _, p := range pts {
		fmt.Fprintf(writer, "%d %d\n", p.x, p.y)
	}
	writer.Flush()

	// Read player choice
	line, err := reader.ReadString('\n')
	if err != nil {
		stdinPipe.Close()
		cmd.Wait()
		return fmt.Errorf("read player choice: %v", err)
	}
	choice := strings.TrimSpace(line)
	choice = strings.ToLower(choice)

	var theofanisFirst bool
	if choice == "first" {
		theofanisFirst = true
	} else if choice == "second" {
		theofanisFirst = false
	} else {
		stdinPipe.Close()
		cmd.Wait()
		return fmt.Errorf("invalid player choice: %q", choice)
	}

	game := &gameState{
		n:      n,
		start:  start,
		points: pts,
		used:   make([]bool, n),
		prev:   start,
		total:  0,
		moves:  0,
	}

	for game.moves < n {
		moveNum := game.moves + 1 // 1-indexed
		isTheofanisTurn := (theofanisFirst && moveNum%2 == 1) || (!theofanisFirst && moveNum%2 == 0)

		if isTheofanisTurn {
			// Read Theofanis's move from candidate stdout
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					stdinPipe.Close()
					cmd.Wait()
					return fmt.Errorf("unexpected EOF from candidate on move %d", moveNum)
				}
				stdinPipe.Close()
				cmd.Wait()
				return fmt.Errorf("read move %d: %v", moveNum, err)
			}
			idx, err := strconv.Atoi(strings.TrimSpace(line))
			if err != nil {
				stdinPipe.Close()
				cmd.Wait()
				return fmt.Errorf("invalid move on turn %d: %q", moveNum, strings.TrimSpace(line))
			}
			idx-- // Convert to 0-indexed
			if idx < 0 || idx >= n || game.used[idx] {
				stdinPipe.Close()
				cmd.Wait()
				return fmt.Errorf("illegal move on turn %d: index %d (0-based)", moveNum, idx)
			}
			game.applyMove(idx)
		} else {
			// Opponent (judge) makes a random valid move
			avail := game.availableIndices()
			choice := avail[rng.Intn(len(avail))]
			game.applyMove(choice)
			// Send opponent's move to candidate (1-indexed)
			fmt.Fprintf(writer, "%d\n", choice+1)
			writer.Flush()
		}
	}

	// Game is over. First player wins if sum is even, second player wins if sum is odd.
	// Theofanis wins if the parity matches his role.
	sumEven := game.total%2 == 0
	theofanisWins := (theofanisFirst && sumEven) || (!theofanisFirst && !sumEven)

	stdinPipe.Close()
	cmd.Wait()

	if !theofanisWins {
		parity := "even"
		if game.total%2 != 0 {
			parity = "odd"
		}
		return fmt.Errorf("Theofanis lost: total sum = %d (%s), chose %s",
			game.total, parity, map[bool]string{true: "First", false: "Second"}[theofanisFirst])
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 50; i++ {
		n := rng.Intn(8) + 1
		start, pts := genPoints(rng, n)
		if err := runInteractiveGame(candidate, rng, n, start, pts); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (n=%d) failed: %v\n", i+1, n, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
