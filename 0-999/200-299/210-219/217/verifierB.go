package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func toggle(b byte) byte {
	if b == 'T' {
		return 'B'
	}
	return 'T'
}

func expected(line string) (string, string) {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return "", ""
	}
	n, _ := strconv.Atoi(parts[0])
	r, _ := strconv.Atoi(parts[1])
	if n == 1 {
		if r == 1 {
			return "0", "T"
		}
		return "IMPOSSIBLE", ""
	}
	if n == 2 {
		if r == 2 {
			return "0", "TB"
		}
		return "IMPOSSIBLE", ""
	}
	bi := -1
	const INF = int(1 << 60)
	bval := INF
	for i := 1; i+i <= r; i++ {
		x := i
		y := r - i
		ans := 0
		errc := 0
		for x > 0 && y > 0 {
			if x < y {
				x, y = y, x
			}
			dd := x / y
			ans += dd
			errc += dd - 1
			x %= y
		}
		x += y
		errc--
		if x == 1 && ans+1 == n {
			if errc < bval {
				bval = errc
				bi = i
			}
		}
	}
	if bi == -1 {
		return "IMPOSSIBLE", ""
	}
	u := bi
	d := r - bi
	moves := make([]byte, 0, n+2)
	for u > 1 || d > 1 {
		if u > d {
			moves = append(moves, 'T')
			u -= d
		} else {
			moves = append(moves, 'B')
			d -= u
		}
	}
	for i, j := 0, len(moves)-1; i < j; i, j = i+1, j-1 {
		moves[i], moves[j] = moves[j], moves[i]
	}
	if len(moves) == 0 {
		return "IMPOSSIBLE", ""
	}
	first := toggle(moves[0])
	last := toggle(moves[len(moves)-1])
	seq := make([]byte, 0, len(moves)+2)
	seq = append(seq, first)
	seq = append(seq, moves...)
	seq = append(seq, last)
	if seq[0] == 'B' {
		for i := range seq {
			seq[i] = toggle(seq[i])
		}
	}
	return fmt.Sprint(bval), string(seq)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		expVal, expSeq := expected(line)
		input := line + "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		parts := strings.Fields(strings.TrimSpace(out.String()))
		if expVal == "IMPOSSIBLE" {
			if len(parts) != 1 || parts[0] != "IMPOSSIBLE" {
				fmt.Printf("test %d failed: expected IMPOSSIBLE got %s\n", idx, strings.Join(parts, " "))
				os.Exit(1)
			}
			continue
		}
		if len(parts) != 2 {
			fmt.Printf("test %d: expected 2 outputs got %d\n", idx, len(parts))
			os.Exit(1)
		}
		if parts[0] != expVal || parts[1] != expSeq {
			fmt.Printf("test %d failed\nexpected: %s %s\n     got: %s %s\n", idx, expVal, expSeq, parts[0], parts[1])
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
