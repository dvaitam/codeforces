package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Point struct {
	x, y int
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	// If the binary path ends in .go, run it with 'go run'
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func verify(input, output string) error {
	ininScan := bufio.NewScanner(strings.NewReader(input))
	ininScan.Split(bufio.ScanWords)
	
outScan := bufio.NewScanner(strings.NewReader(output))
	outScan.Split(bufio.ScanWords)

	if !ininScan.Scan() {
		return fmt.Errorf("empty input")
	}
	q := 0
	fmt.Sscanf(ininScan.Text(), "%d", &q)

	for k := 0; k < q; k++ {
		if !ininScan.Scan() {
			return fmt.Errorf("input too short (b)")
		}
		var b int
		fmt.Sscanf(ininScan.Text(), "%d", &b)
		
		if !ininScan.Scan() {
			return fmt.Errorf("input too short (w)")
		}
		var w int
		fmt.Sscanf(ininScan.Text(), "%d", &w)

		if !outScan.Scan() {
			return fmt.Errorf("output too short (YES/NO)")
		}
		ans := outScan.Text()

		canSolve := true
		if b > 3*w+1 || w > 3*b+1 {
			canSolve = false
		}

		if !canSolve {
			if ans != "NO" {
				return fmt.Errorf("case %d: expected NO, got %s for b=%d, w=%d", k+1, ans, b, w)
			}
		} else {
			if ans != "YES" {
				return fmt.Errorf("case %d: expected YES, got %s for b=%d, w=%d", k+1, ans, b, w)
			}
			
			// Read b+w coordinates
			points := make(map[Point]bool)
			var cells []Point
			cntB, cntW := 0, 0
			
total := b + w
			for i := 0; i < total; i++ {
				if !outScan.Scan() {
					return fmt.Errorf("case %d: output too short, expected coord x", k+1)
				}
				var x int
				fmt.Sscanf(outScan.Text(), "%d", &x)
				
				if !outScan.Scan() {
					return fmt.Errorf("case %d: output too short, expected coord y", k+1)
				}
				var y int
				fmt.Sscanf(outScan.Text(), "%d", &y)

				p := Point{x, y}
				if points[p] {
					return fmt.Errorf("case %d: duplicate point (%d, %d)", k+1, x, y)
				}
				points[p] = true
				cells = append(cells, p)

				// Check color. (1,1) is White.
				// (x, y) color depends on parity of x+y.
				// 1+1 = 2 (even) -> White.
				// So even sum -> White, odd sum -> Black.
				if (x+y)%2 == 0 {
					cntW++
				} else {
					cntB++
				}
			}

			if cntB != b {
				return fmt.Errorf("case %d: black count mismatch, expected %d got %d", k+1, b, cntB)
			}
			if cntW != w {
				return fmt.Errorf("case %d: white count mismatch, expected %d got %d", k+1, w, cntW)
			}

			// Check connectivity
			if len(cells) > 0 {
				visited := make(map[Point]bool)
				queue := []Point{cells[0]}
				visited[cells[0]] = true
				count := 0
				
				for len(queue) > 0 {
					curr := queue[0]
					queue = queue[1:]
					count++
					
					// neighbors
				dirs := []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
				for _, d := range dirs {
						next := Point{curr.x + d.x, curr.y + d.y}
						if points[next] && !visited[next] {
							visited[next] = true
							queue = append(queue, next)
						}
						}
				}
				
				if count != total {
					return fmt.Errorf("case %d: component not connected, visited %d of %d cells", k+1, count, total)
				}
			}
		}
	}
	return nil
}

func genCase(r *rand.Rand) string {
	// Constraints: 1 <= b, w <= 10^5 (sum <= 2*10^5)
	// For verification, small numbers are fine, but should cover edge cases.
	b := r.Intn(20) + 1
	w := r.Intn(20) + 1
	return fmt.Sprintf("1\n%d %d\n", b, w)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		output, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := verify(input, output); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s", i+1, err, input, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}