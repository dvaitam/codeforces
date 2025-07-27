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

func runProg(bin, input string) (string, error) {
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
	return out.String(), nil
}

func genCase(rng *rand.Rand) (int, int, int, int, string) {
	n := rng.Intn(7) + 3
	m := rng.Intn(7) + 3
	x := rng.Intn(n-2) + 2
	y := rng.Intn(m-2) + 2
	input := fmt.Sprintf("%d %d %d %d\n", n, m, x, y)
	return n, m, x, y, input
}

func checkOutput(n, m, x, y int, output string) error {
	fields := strings.Fields(output)
	if len(fields) != 2*n*m {
		return fmt.Errorf("expected %d pairs, got %d", n*m, len(fields)/2)
	}
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}
	px, py := x, y
	for i := 0; i < n*m; i++ {
		xi, err1 := strconv.Atoi(fields[2*i])
		yi, err2 := strconv.Atoi(fields[2*i+1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid integer in output")
		}
		if xi < 1 || xi > n || yi < 1 || yi > m {
			return fmt.Errorf("cell out of range: %d %d", xi, yi)
		}
		if visited[xi-1][yi-1] {
			return fmt.Errorf("cell repeated: %d %d", xi, yi)
		}
		if i == 0 {
			if xi != x || yi != y {
				return fmt.Errorf("first cell should be %d %d", x, y)
			}
		} else {
			if xi != px && yi != py {
				return fmt.Errorf("illegal move from %d %d to %d %d", px, py, xi, yi)
			}
		}
		visited[xi-1][yi-1] = true
		px, py = xi, yi
	}
	// ensure all cells visited
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if !visited[i][j] {
				return fmt.Errorf("cell %d %d not visited", i+1, j+1)
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, x, y, input := genCase(rng)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			return
		}
		if err := checkOutput(n, m, x, y, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, input, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
