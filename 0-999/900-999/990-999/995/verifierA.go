package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// structure for a single test case

type caseA struct {
	n    int
	k    int
	spot [4][]int // rows with designated spots (rows 0 and 3) and start positions (rows1 and2?)
}

func genCaseA(rng *rand.Rand) caseA {
	n := rng.Intn(10) + 1    // keep small for speed
	k := rng.Intn(2*n-1) + 1 // ensure at least one empty in inner rows
	var spot [4][]int
	for i := 0; i < 4; i++ {
		spot[i] = make([]int, n)
	}
	// assign parking spots in rows 0 and 3
	usedCols1 := rng.Perm(n)
	usedCols2 := rng.Perm(n)
	for id := 1; id <= k; id++ {
		if rng.Intn(2) == 0 {
			c := usedCols1[id%n]
			for spot[0][c] != 0 {
				c = rng.Intn(n)
			}
			spot[0][c] = id
		} else {
			c := usedCols2[id%n]
			for spot[3][c] != 0 {
				c = rng.Intn(n)
			}
			spot[3][c] = id
		}
	}
	// assign starting positions rows 1 and 2 (index 1 and2?)
	free := make([][2]int, 0, 2*n)
	for r := 1; r <= 2; r++ {
		for c := 0; c < n; c++ {
			free = append(free, [2]int{r, c})
		}
	}
	rng.Shuffle(len(free), func(i, j int) { free[i], free[j] = free[j], free[i] })
	for id := 1; id <= k; id++ {
		pos := free[id-1]
		spot[pos[0]][pos[1]] = id
	}
	return caseA{n: n, k: k, spot: spot}
}

// simulate moves and verify
func verifyOutput(tc caseA, output string) error {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	if fields[0] == "-1" {
		return fmt.Errorf("solution reported impossible")
	}
	m, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("bad move count: %v", err)
	}
	if m < 0 || m > 20000 {
		return fmt.Errorf("invalid move count %d", m)
	}
	if len(fields) != 1+3*m {
		return fmt.Errorf("expected %d numbers got %d", 1+3*m, len(fields))
	}
	// grid occupancy: rows 0..3, cols 0..n-1
	grid := make([][]int, 4)
	for i := 0; i < 4; i++ {
		grid[i] = make([]int, tc.n)
		copy(grid[i], tc.spot[i])
	}
	// car positions map from id->(r,c) for rows 1 and2 (index 1 and2)
	pos := make(map[int][2]int)
	for r := 1; r <= 2; r++ {
		for c := 0; c < tc.n; c++ {
			id := grid[r][c]
			if id != 0 {
				pos[id] = [2]int{r, c}
			}
		}
	}
	// clear rows 0 and 3 since they are targets only
	for _, r := range []int{0, 3} {
		for c := 0; c < tc.n; c++ {
			grid[r][c] = 0
		}
	}
	idx := 1
	for step := 0; step < m; step++ {
		id, err1 := strconv.Atoi(fields[idx])
		r, err2 := strconv.Atoi(fields[idx+1])
		c, err3 := strconv.Atoi(fields[idx+2])
		if err1 != nil || err2 != nil || err3 != nil {
			return fmt.Errorf("bad move on line %d", step+1)
		}
		idx += 3
		r--
		c--
		if r < 0 || r >= 4 || c < 0 || c >= tc.n {
			return fmt.Errorf("move out of bounds")
		}
		p, ok := pos[id]
		if !ok {
			return fmt.Errorf("car %d not found", id)
		}
		if abs(p[0]-r)+abs(p[1]-c) != 1 {
			return fmt.Errorf("illegal move for car %d", id)
		}
		if grid[r][c] != 0 {
			return fmt.Errorf("destination occupied")
		}
		if r == 0 || r == 3 {
			if tc.spot[r][c] != id {
				return fmt.Errorf("car %d cannot park at %d,%d", id, r+1, c+1)
			}
		}
		// perform move
		grid[p[0]][p[1]] = 0
		grid[r][c] = id
		pos[id] = [2]int{r, c}
		if r == 0 || r == 3 {
			delete(pos, id)
		}
	}
	if len(pos) != 0 {
		return fmt.Errorf("not all cars parked")
	}
	return nil
}

func runCaseA(bin string, tc caseA) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
	for r := 0; r < 4; r++ {
		for c := 0; c < tc.n; c++ {
			if c > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", tc.spot[r][c])
		}
		sb.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return verifyOutput(tc, out.String())
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseA(rng)
		if err := runCaseA(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
