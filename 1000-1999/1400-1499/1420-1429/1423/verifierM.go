package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierM /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for tc := 1; tc <= 200; tc++ {
		n := rng.Intn(8) + 1
		m := rng.Intn(8) + 1
		mat, trueMin := generateMatrix(n, m, rng)

		got, queries, err := interact(bin, n, m, mat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d (n=%d m=%d): %v\n", tc, n, m, err)
			printMatrix(mat)
			os.Exit(1)
		}
		if got != trueMin {
			fmt.Fprintf(os.Stderr, "case %d (n=%d m=%d): wrong answer — expected %d got %d (queries=%d)\n", tc, n, m, trueMin, got, queries)
			printMatrix(mat)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

// generateMatrix builds an n×m totally monotone matrix using the parabola
// construction: A[i][j] = base[i] + (j - c[i])^2, where c is non-decreasing.
// The minimum of row i is base[i] (at column c[i]), so the global minimum is
// min(base[i]).
func generateMatrix(n, m int, rng *rand.Rand) (mat [][]int, trueMin int) {
	c := make([]int, n)
	c[0] = rng.Intn(m)
	for i := 1; i < n; i++ {
		c[i] = c[i-1] + rng.Intn(m-c[i-1]+1)
		if c[i] >= m {
			c[i] = m - 1
		}
	}

	mat = make([][]int, n)
	trueMin = 1 << 62
	for i := range mat {
		mat[i] = make([]int, m)
		base := rng.Intn(100) + 1
		for j := range mat[i] {
			diff := j - c[i]
			mat[i][j] = base + diff*diff
		}
		if base < trueMin {
			trueMin = base
		}
	}
	return
}

func interact(bin string, n, m int, mat [][]int) (answer, queries int, err error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}

	inPipe, err := cmd.StdinPipe()
	if err != nil {
		return 0, 0, fmt.Errorf("stdin pipe: %v", err)
	}
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return 0, 0, fmt.Errorf("stdout pipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		return 0, 0, fmt.Errorf("start: %v", err)
	}

	w := bufio.NewWriter(inPipe)
	r := bufio.NewReader(outPipe)

	fmt.Fprintf(w, "%d %d\n", n, m)
	w.Flush()

	maxQueries := (n + m) * 50
	for queries <= maxQueries {
		line, readErr := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			if readErr != nil {
				break
			}
			continue
		}

		if strings.HasPrefix(line, "! ") {
			val, convErr := strconv.Atoi(strings.TrimSpace(line[2:]))
			if convErr != nil {
				cmd.Process.Kill()
				return 0, queries, fmt.Errorf("bad answer line: %q", line)
			}
			inPipe.Close()
			cmd.Wait()
			return val, queries, nil
		}

		if strings.HasPrefix(line, "? ") {
			parts := strings.Fields(line[2:])
			if len(parts) != 2 {
				cmd.Process.Kill()
				return 0, queries, fmt.Errorf("bad query format: %q", line)
			}
			row, _ := strconv.Atoi(parts[0])
			col, _ := strconv.Atoi(parts[1])
			row--
			col--
			if row < 0 || row >= n || col < 0 || col >= m {
				cmd.Process.Kill()
				return 0, queries, fmt.Errorf("query out of bounds: r=%d c=%d (n=%d m=%d)", row+1, col+1, n, m)
			}
			queries++
			fmt.Fprintf(w, "%d\n", mat[row][col])
			w.Flush()
			continue
		}

		cmd.Process.Kill()
		return 0, queries, fmt.Errorf("unexpected output: %q", line)
	}

	cmd.Process.Kill()
	return 0, queries, fmt.Errorf("exceeded query limit or no answer produced")
}

func printMatrix(mat [][]int) {
	for _, row := range mat {
		fmt.Fprintln(os.Stderr, row)
	}
}
