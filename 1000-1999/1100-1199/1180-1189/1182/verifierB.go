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

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func buildRef() (string, error) {
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1182B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genPlus(h, w int) []string {
	grid := make([][]byte, h)
	for i := range grid {
		row := make([]byte, w)
		for j := range row {
			row[j] = '.'
		}
		grid[i] = row
	}
	cx := rand.Intn(h-2) + 1
	cy := rand.Intn(w-2) + 1
	up := rand.Intn(cx) + 1
	down := rand.Intn(h-cx-1) + 1
	left := rand.Intn(cy) + 1
	right := rand.Intn(w-cy-1) + 1
	for d := 0; d <= up; d++ {
		grid[cx-d][cy] = '*'
	}
	for d := 1; d <= down; d++ {
		grid[cx+d][cy] = '*'
	}
	for d := 1; d <= left; d++ {
		grid[cx][cy-d] = '*'
	}
	for d := 1; d <= right; d++ {
		grid[cx][cy+d] = '*'
	}
	res := make([]string, h)
	for i := range grid {
		res[i] = string(grid[i])
	}
	return res
}

func genRandomGrid(h, w int) []string {
	grid := make([]string, h)
	for i := 0; i < h; i++ {
		var sb strings.Builder
		for j := 0; j < w; j++ {
			if rand.Intn(2) == 0 {
				sb.WriteByte('.')
			} else {
				sb.WriteByte('*')
			}
		}
		grid[i] = sb.String()
	}
	return grid
}

func genTest(i int) []byte {
	h := rand.Intn(8) + 3
	w := rand.Intn(8) + 3
	var grid []string
	if i%5 == 0 {
		grid = genPlus(h, w)
	} else {
		grid = genRandomGrid(h, w)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", h, w))
	for _, row := range grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genTest(i)
		want, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
			os.Exit(1)
		}
		got, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, string(input), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
