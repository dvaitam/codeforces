package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 2
	m := rng.Intn(3) + 2
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			fmt.Fprintf(&sb, "%d ", rng.Intn(5))
		}
		sb.WriteByte('\n')
	}
	has := false
	grid := make([][]int, n)
	for i := range grid {
		grid[i] = make([]int, m)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				grid[i][j] = -1
			} else {
				grid[i][j] = rng.Intn(5)
				has = true
			}
			fmt.Fprintf(&sb, "%d ", grid[i][j])
		}
		sb.WriteByte('\n')
	}
	if !has {
		grid[0][0] = 1
	}
	x := rng.Intn(n)
	y := rng.Intn(m)
	if grid[x][y] == -1 {
		grid[x][y] = 1
	}
	fmt.Fprintf(&sb, "%d %d\n", x+1, y+1)
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	ref := "./105D.go"
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		want, err := runBinary(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal solver failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
