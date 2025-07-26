package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type cell struct{ r, c int }

type Test struct {
	n    int
	path []cell
}

func generatePath(n int, rng *rand.Rand) []cell {
	length := 1 // at least 1 step
	visited := make(map[cell]bool)
	dirs := []struct{ dr, dc int }{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	x := rng.Intn(n) + 1
	y := rng.Intn(n) + 1
	path := []cell{{x, y}}
	visited[cell{x, y}] = true
	for len(path)-1 < length {
		cur := path[len(path)-1]
		candidates := make([]cell, 0, 4)
		for _, d := range dirs {
			nx, ny := cur.r+d.dr, cur.c+d.dc
			if nx >= 1 && nx <= n && ny >= 1 && ny <= n {
				c := cell{nx, ny}
				if !visited[c] {
					candidates = append(candidates, c)
				}
			}
		}
		if len(candidates) == 0 {
			break
		}
		next := candidates[rng.Intn(len(candidates))]
		path = append(path, next)
		visited[next] = true
	}
	if len(path) == 1 {
		// pick random adjacent cell
		d := dirs[rng.Intn(4)]
		nx, ny := x+d.dr, y+d.dc
		if nx < 1 || nx > n || ny < 1 || ny > n {
			nx, ny = x-d.dr, y-d.dc
		}
		path = append(path, cell{nx, ny})
	}
	return path
}

func generateTests() []Test {
	rng := rand.New(rand.NewSource(46))
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 3 // n in [3,6]
		path := generatePath(n, rng)
		tests = append(tests, Test{n: n, path: path})
	}
	return tests
}

func inside(c cell, x1, y1, x2, y2 int) bool {
	return c.r >= x1 && c.r <= x2 && c.c >= y1 && c.c <= y2
}

func queryAnswer(path []cell, x1, y1, x2, y2 int) int {
	cnt := 0
	for i := 0; i < len(path)-1; i++ {
		a := inside(path[i], x1, y1, x2, y2)
		b := inside(path[i+1], x1, y1, x2, y2)
		if a != b {
			cnt++
		}
	}
	return cnt
}

func runCase(bin string, t Test) error {
	var cmd *exec.Cmd
	cmd = exec.Command(bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	inw := bufio.NewWriter(stdin)
	outr := bufio.NewReader(stdout)
	fmt.Fprintln(inw, t.n)
	inw.Flush()
	queries := 0
	for {
		line, err := outr.ReadString('\n')
		if err != nil {
			cmd.Process.Kill()
			return fmt.Errorf("failed read: %v", err)
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "?") {
			queries++
			if queries > 2019 {
				cmd.Process.Kill()
				return fmt.Errorf("too many queries")
			}
			var x1, y1, x2, y2 int
			fmt.Sscanf(line, "? %d %d %d %d", &x1, &y1, &x2, &y2)
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			if y1 > y2 {
				y1, y2 = y2, y1
			}
			ans := queryAnswer(t.path, x1, y1, x2, y2)
			fmt.Fprintln(inw, ans)
			inw.Flush()
		} else if strings.HasPrefix(line, "!") {
			var x1, y1, x2, y2 int
			fmt.Sscanf(line, "! %d %d %d %d", &x1, &y1, &x2, &y2)
			head := t.path[0]
			tail := t.path[len(t.path)-1]
			ok := (x1 == head.r && y1 == head.c && x2 == tail.r && y2 == tail.c) ||
				(x2 == head.r && y2 == head.c && x1 == tail.r && y1 == tail.c)
			if !ok {
				cmd.Process.Kill()
				return fmt.Errorf("wrong answer")
			}
			cmd.Wait()
			return nil
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		if err := runCase(bin, t); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
