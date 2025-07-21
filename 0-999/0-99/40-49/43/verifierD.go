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
	"time"
)

type coord struct{ x, y int }

type tele struct{ from, to coord }

func minTeleporters(n, m int) int {
	if (n == 1 && m == 1) || (n == 1 && m == 2) || (n == 2 && m == 1) {
		return 0
	}
	if n == 1 || m == 1 {
		return 1
	}
	if (n*m)%2 == 0 {
		return 0
	}
	return 1
}

func generateCases() [][2]int {
	cases := [][2]int{{1, 1}, {1, 2}, {2, 1}, {2, 2}, {3, 3}, {2, 3}, {3, 2}, {4, 5}}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 100 {
		n := rng.Intn(7) + 1
		m := rng.Intn(7) + 1
		if n*m < 2 {
			continue
		}
		cases = append(cases, [2]int{n, m})
	}
	return cases
}

func verify(n, m int, exe string) error {
	input := fmt.Sprintf("%d %d\n", n, m)
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return fmt.Errorf("invalid teleporter count: %v", err)
	}
	if k != minTeleporters(n, m) {
		return fmt.Errorf("expected %d teleporters got %d", minTeleporters(n, m), k)
	}
	teleports := make(map[coord]coord)
	for i := 0; i < k; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("missing teleporter line %d", i+1)
		}
		var x1, y1, x2, y2 int
		if _, err := fmt.Sscan(scanner.Text(), &x1, &y1, &x2, &y2); err != nil {
			return fmt.Errorf("invalid teleporter line: %v", err)
		}
		if x1 < 1 || x1 > n || y1 < 1 || y1 > m || x2 < 1 || x2 > n || y2 < 1 || y2 > m {
			return fmt.Errorf("teleporter coords out of range")
		}
		from := coord{x1, y1}
		if _, ok := teleports[from]; ok {
			return fmt.Errorf("duplicate teleporter from %v", from)
		}
		teleports[from] = coord{x2, y2}
	}
	path := make([]coord, 0, n*m+1)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var x, y int
		if _, err := fmt.Sscan(line, &x, &y); err != nil {
			return fmt.Errorf("bad path line: %v", err)
		}
		path = append(path, coord{x, y})
	}
	if len(path) != n*m+1 {
		return fmt.Errorf("path length %d expected %d", len(path), n*m+1)
	}
	if path[0] != (coord{1, 1}) || path[len(path)-1] != (coord{1, 1}) {
		return fmt.Errorf("path must start and end at 1 1")
	}
	counts := make(map[coord]int)
	counts[path[0]]++
	cur := path[0]
	for i := 1; i < len(path); i++ {
		nxt := path[i]
		if nxt.x < 1 || nxt.x > n || nxt.y < 1 || nxt.y > m {
			return fmt.Errorf("path coordinate out of range")
		}
		if t, ok := teleports[cur]; ok && t == nxt {
			// teleport
		} else {
			dx := cur.x - nxt.x
			if dx < 0 {
				dx = -dx
			}
			dy := cur.y - nxt.y
			if dy < 0 {
				dy = -dy
			}
			if dx+dy != 1 {
				return fmt.Errorf("invalid move from %v to %v", cur, nxt)
			}
		}
		counts[nxt]++
		cur = nxt
	}
	for x := 1; x <= n; x++ {
		for y := 1; y <= m; y++ {
			c := counts[coord{x, y}]
			if x == 1 && y == 1 {
				if c != 2 {
					return fmt.Errorf("(1,1) visited %d times", c)
				}
			} else if c != 1 {
				return fmt.Errorf("cell %d %d visited %d times", x, y, c)
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	for i, tc := range generateCases() {
		if err := verify(tc[0], tc[1], exe); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (%d %d) failed: %v\n", i+1, tc[0], tc[1], err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
