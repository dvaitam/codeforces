package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Interactive verifier for 1838F (Stuck Conveyor).
// We pick a random stuck conveyor and direction, then simulate queries.

func simulate(n int, stuckR, stuckC int, stuckDir byte, startR, startC int, grid [][]byte) (int, int) {
	// Simulate the box movement. The box starts at (startR, startC).
	// All conveyors are as specified in grid, except (stuckR, stuckC) which
	// always uses stuckDir.
	// If box moves to an empty square (outside [1..n]x[1..n]), it stops there.
	// If box enters a cycle, return (-1, -1).

	getDir := func(r, c int) byte {
		if r == stuckR && c == stuckC {
			return stuckDir
		}
		return grid[r-1][c-1]
	}

	// Use tortoise and hare to detect cycles
	r, c := startR, startC
	visited := make(map[[2]int]int)
	step := 0
	maxSteps := n*n + 10

	for {
		if r < 1 || r > n || c < 1 || c > n {
			return r, c
		}
		key := [2]int{r, c}
		if _, ok := visited[key]; ok {
			return -1, -1
		}
		visited[key] = step
		step++
		if step > maxSteps {
			return -1, -1
		}

		d := getDir(r, c)
		switch d {
		case '^':
			r--
		case 'v':
			r++
		case '<':
			c--
		case '>':
			c++
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	dirs := []byte{'^', 'v', '<', '>'}

	for tc := 0; tc < 50; tc++ {
		n := rng.Intn(5) + 2 // n in [2..6]
		stuckR := rng.Intn(n) + 1
		stuckC := rng.Intn(n) + 1
		stuckDir := dirs[rng.Intn(4)]

		cmd := exec.Command(bin)
		candIn, _ := cmd.StdinPipe()
		candOut, _ := cmd.StdoutPipe()
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			fmt.Printf("case %d: failed to start binary: %v\n", tc+1, err)
			os.Exit(1)
		}

		writer := bufio.NewWriter(candIn)
		reader := bufio.NewReader(candOut)

		// Send n
		fmt.Fprintf(writer, "%d\n", n)
		writer.Flush()

		queryCount := 0
		ok := true
		var failReason string

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					failReason = "unexpected EOF from candidate"
					ok = false
				} else {
					failReason = fmt.Sprintf("read error: %v", err)
					ok = false
				}
				break
			}
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}

			if line[0] == '!' {
				// Final answer: ! r c dir
				var ar, ac int
				var ad string
				_, err := fmt.Sscanf(line, "! %d %d %s", &ar, &ac, &ad)
				if err != nil {
					failReason = fmt.Sprintf("bad answer format: %s", line)
					ok = false
					break
				}
				if ar != stuckR || ac != stuckC || ad != string(stuckDir) {
					failReason = fmt.Sprintf("wrong answer: expected ! %d %d %c, got %s", stuckR, stuckC, stuckDir, line)
					ok = false
				}
				break
			}

			if line[0] == '?' {
				queryCount++
				if queryCount > 25 {
					failReason = "too many queries (>25)"
					ok = false
					// Send dummy response
					fmt.Fprintf(writer, "-1 -1\n")
					writer.Flush()
					break
				}

				// Parse: ? sr sc
				var sr, sc int
				_, err := fmt.Sscanf(line, "? %d %d", &sr, &sc)
				if err != nil {
					failReason = fmt.Sprintf("bad query format: %s", line)
					ok = false
					break
				}

				// Read n lines of grid
				grid := make([][]byte, n)
				for i := 0; i < n; i++ {
					gline, err := reader.ReadString('\n')
					if err != nil {
						failReason = fmt.Sprintf("failed to read grid line %d: %v", i+1, err)
						ok = false
						break
					}
					gline = strings.TrimSpace(gline)
					if len(gline) != n {
						failReason = fmt.Sprintf("grid line %d has length %d, expected %d: %q", i+1, len(gline), n, gline)
						ok = false
						break
					}
					grid[i] = []byte(gline)
				}
				if !ok {
					break
				}

				// Simulate
				resR, resC := simulate(n, stuckR, stuckC, stuckDir, sr, sc, grid)
				fmt.Fprintf(writer, "%d %d\n", resR, resC)
				writer.Flush()
			} else {
				failReason = fmt.Sprintf("unexpected output: %s", line)
				ok = false
				break
			}
		}

		candIn.Close()
		cmd.Wait()

		if !ok {
			fmt.Printf("case %d failed: %s (n=%d stuck=(%d,%d,'%c'))\n", tc+1, failReason, n, stuckR, stuckC, stuckDir)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
