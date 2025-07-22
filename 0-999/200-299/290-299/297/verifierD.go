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

func run(bin, input string) (string, error) {
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
	return strings.TrimSpace(out.String()), nil
}

func bruteForce(h, w, k int, row []string, col []string) (int, bool) {
	total := h*(w-1) + w*(h-1)
	best := 0
	cells := h * w
	assign := make([]int, cells)
	var dfs func(pos int)
	dfs = func(pos int) {
		if pos == cells {
			sat := 0
			for i := 0; i < h; i++ {
				for j := 0; j < w-1; j++ {
					a := assign[i*w+j]
					b := assign[i*w+j+1]
					if row[i][j] == 'E' && a == b || row[i][j] == 'N' && a != b {
						sat++
					}
				}
			}
			for i := 0; i < h-1; i++ {
				for j := 0; j < w; j++ {
					a := assign[i*w+j]
					b := assign[(i+1)*w+j]
					if col[i][j] == 'E' && a == b || col[i][j] == 'N' && a != b {
						sat++
					}
				}
			}
			if sat > best {
				best = sat
			}
			return
		}
		for c := 1; c <= k; c++ {
			assign[pos] = c
			dfs(pos + 1)
		}
	}
	dfs(0)
	return best, best*4 >= 3*total
}

func satisfaction(h, w int, row []string, col []string, board []int) int {
	sat := 0
	for i := 0; i < h; i++ {
		for j := 0; j < w-1; j++ {
			a := board[i*w+j]
			b := board[i*w+j+1]
			if row[i][j] == 'E' && a == b || row[i][j] == 'N' && a != b {
				sat++
			}
		}
	}
	for i := 0; i < h-1; i++ {
		for j := 0; j < w; j++ {
			a := board[i*w+j]
			b := board[(i+1)*w+j]
			if col[i][j] == 'E' && a == b || col[i][j] == 'N' && a != b {
				sat++
			}
		}
	}
	return sat
}

func parseBoard(tokens []string, idx int, h, w, k int) ([]int, error) {
	if len(tokens) < idx+h*w {
		return nil, fmt.Errorf("not enough numbers")
	}
	board := make([]int, h*w)
	for i := 0; i < h*w; i++ {
		var v int
		fmt.Sscan(tokens[idx+i], &v)
		if v < 1 || v > k {
			return nil, fmt.Errorf("color out of range")
		}
		board[i] = v
	}
	return board, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for caseNum := 0; caseNum < 200; caseNum++ {
		h := rng.Intn(2) + 2
		w := rng.Intn(2) + 2
		k := rng.Intn(3) + 1
		row := make([]string, h)
		col := make([]string, h-1)
		for i := 0; i < h; i++ {
			if w-1 > 0 {
				b := make([]byte, w-1)
				for j := range b {
					if rng.Intn(2) == 0 {
						b[j] = 'E'
					} else {
						b[j] = 'N'
					}
				}
				row[i] = string(b)
			} else {
				row[i] = ""
			}
			if i < h-1 {
				b := make([]byte, w)
				for j := range b {
					if rng.Intn(2) == 0 {
						b[j] = 'E'
					} else {
						b[j] = 'N'
					}
				}
				col[i] = string(b)
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", h, w, k))
		for i := 0; i < 2*h-1; i++ {
			if i%2 == 0 {
				sb.WriteString(row[i/2])
			} else {
				sb.WriteString(col[i/2])
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		_, possible := bruteForce(h, w, k, row, col)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		tokens := strings.Fields(out)
		if len(tokens) == 0 {
			fmt.Fprintf(os.Stderr, "case %d failed: no output\n", caseNum+1)
			os.Exit(1)
		}
		ans := strings.ToUpper(tokens[0])
		if ans == "NO" {
			if possible {
				fmt.Fprintf(os.Stderr, "case %d failed: answer should be YES\ninput:\n%s", caseNum+1, input)
				os.Exit(1)
			}
			continue
		}
		if ans != "YES" {
			fmt.Fprintf(os.Stderr, "case %d failed: first token should be YES or NO\n", caseNum+1)
			os.Exit(1)
		}
		board, err := parseBoard(tokens, 1, h, w, k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		sat := satisfaction(h, w, row, col, board)
		total := h*(w-1) + w*(h-1)
		if sat*4 < 3*total {
			fmt.Fprintf(os.Stderr, "case %d failed: satisfaction too low\ninput:\n%soutput:\n%s", caseNum+1, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
