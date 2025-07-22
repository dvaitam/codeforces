package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type pos struct{ r, c int }

func solve(n, m, x int, board []string, q int, text string) string {
	letters := make([][]pos, 26)
	var shifts []pos
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			ch := board[i][j]
			if ch == 'S' {
				shifts = append(shifts, pos{i, j})
			} else {
				letters[ch-'a'] = append(letters[ch-'a'], pos{i, j})
			}
		}
	}
	good := make([]bool, 26)
	if len(shifts) > 0 {
		maxd2 := x * x
		for i := 0; i < 26; i++ {
			if len(letters[i]) == 0 {
				continue
			}
			ok := false
			for _, p := range letters[i] {
				for _, s := range shifts {
					dr := p.r - s.r
					dc := p.c - s.c
					if dr*dr+dc*dc <= maxd2 {
						ok = true
						break
					}
				}
				if ok {
					break
				}
			}
			good[i] = ok
		}
	}
	count := 0
	for i := 0; i < q; i++ {
		ch := text[i]
		if ch >= 'a' && ch <= 'z' {
			idx := ch - 'a'
			if len(letters[idx]) == 0 {
				return "-1"
			}
		} else {
			idx := ch - 'A'
			if len(letters[idx]) == 0 || len(shifts) == 0 {
				return "-1"
			}
			if !good[idx] {
				count++
			}
		}
	}
	return fmt.Sprint(count)
}

func genCase() (string, string) {
	n := rand.Intn(4) + 1
	m := rand.Intn(4) + 1
	x := rand.Intn(5) + 1
	board := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rand.Float64() < 0.2 {
				row[j] = 'S'
			} else {
				row[j] = byte('a' + rand.Intn(26))
			}
		}
		board[i] = string(row)
	}
	q := rand.Intn(10) + 1
	textb := make([]byte, q)
	for i := 0; i < q; i++ {
		if rand.Float64() < 0.5 {
			textb[i] = byte('a' + rand.Intn(26))
		} else {
			textb[i] = byte('A' + rand.Intn(26))
		}
	}
	text := string(textb)
	in := fmt.Sprintf("%d %d %d\n", n, m, x)
	for _, row := range board {
		in += row + "\n"
	}
	in += fmt.Sprintf("%d\n%s\n", q, text)
	out := solve(n, m, x, board, q, text)
	return in, out
}

func runCase(bin, input, expected string, idx int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error on case %d: %v", idx, err)
	}
	out := strings.TrimSpace(string(outBytes))
	if out != expected {
		return fmt.Errorf("wrong answer on case %d: expected %q got %q input:\n%s", idx, expected, out, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]
	tests := 150
	for i := 1; i <= tests; i++ {
		in, exp := genCase()
		if err := runCase(bin, in, exp, i); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	fmt.Printf("ok %d\n", tests)
}
