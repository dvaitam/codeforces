package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func runProg(path, input string) (string, error) {
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
	return strings.TrimSpace(out.String()), err
}

func runRef(input string) (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "683I.go")
	return runProg(ref, input)
}

func genGrid(n, m int) []string {
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rand.Intn(5) == 0 {
				row[j] = 'X'
			} else {
				row[j] = '.'
			}
		}
		grid[i] = string(row)
	}
	// place Y, B, T distinct
	sx, sy := rand.Intn(n), rand.Intn(m)
	grid[sx] = grid[sx][:sy] + "Y" + grid[sx][sy+1:]
	bx, by := rand.Intn(n), rand.Intn(m)
	for bx == sx && by == sy {
		bx, by = rand.Intn(n), rand.Intn(m)
	}
	grid[bx] = grid[bx][:by] + "B" + grid[bx][by+1:]
	tx, ty := rand.Intn(n), rand.Intn(m)
	for (tx == sx && ty == sy) || (tx == bx && ty == by) {
		tx, ty = rand.Intn(n), rand.Intn(m)
	}
	grid[tx] = grid[tx][:ty] + "T" + grid[tx][ty+1:]
	return grid
}

func genCase() string {
	n := rand.Intn(4) + 3
	m := rand.Intn(4) + 3
	grid := genGrid(n, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(grid[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		in := genCase()
		expect, err := runRef(in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\nactual:%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
