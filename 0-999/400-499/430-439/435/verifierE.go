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

func solveE(n, m int, grid []string) []string {
	s := make([][]byte, n)
	for i := 0; i < n; i++ {
		s[i] = []byte(grid[i])
	}
	work := func() bool {
		e := [5]int{-1, -1, -1, -1, -1}
		for i := 0; i < n; i++ {
			g := [5]int{-1, -1, -1, -1, -1}
			for j := 0; j < m; j++ {
				if s[i][j] != '0' {
					c := int(s[i][j] - '0')
					pr := (i + 1) & 1
					pc := (j + 1) & 1
					if g[c] != -1 && g[c] != pc {
						return false
					}
					if e[c] != -1 && e[c] != pr {
						return false
					}
					g[c] = pc
					e[c] = pr
				}
			}
		}
		cnt0, cnt1 := 0, 0
		for c := 1; c <= 4; c++ {
			if e[c] == 0 {
				cnt0++
			}
			if e[c] == 1 {
				cnt1++
			}
		}
		if cnt0 > 2 || cnt1 > 2 {
			return false
		}
		s1, s2, s3, s4 := 0, 0, 0, 0
		for c := 1; c <= 4; c++ {
			if e[c] != -1 {
				if e[c] == 1 {
					if s1 != 0 {
						s2 = c
					} else {
						s1 = c
					}
				} else {
					if s3 != 0 {
						s4 = c
					} else {
						s3 = c
					}
				}
			}
		}
		for c := 1; c <= 4; c++ {
			if e[c] == -1 {
				if s1 == 0 {
					s1 = c
				} else if s2 == 0 {
					s2 = c
				} else if s3 == 0 {
					s3 = c
				} else {
					s4 = c
				}
			}
		}
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if s[i][j] != '0' {
					c := int(s[i][j] - '0')
					pr := (i + 1) & 1
					pc := (j + 1) & 1
					if pr == 1 && pc == 1 && s1 != c {
						s1, s2 = s2, s1
					} else if pr == 1 && pc == 0 && s2 != c {
						s1, s2 = s2, s1
					} else if pr == 0 && pc == 1 && s3 != c {
						s3, s4 = s4, s3
					} else if pr == 0 && pc == 0 && s4 != c {
						s3, s4 = s4, s3
					}
					break
				}
			}
			for j := 0; j < m; j++ {
				pr := (i + 1) & 1
				pc := (j + 1) & 1
				var nc int
				if pr == 1 && pc == 1 {
					nc = s1
				} else if pr == 1 && pc == 0 {
					nc = s2
				} else if pr == 0 && pc == 1 {
					nc = s3
				} else {
					nc = s4
				}
				s[i][j] = byte('0' + nc)
			}
		}
		return true
	}
	if work() {
		res := make([]string, n)
		for i := 0; i < n; i++ {
			res[i] = string(s[i])
		}
		return res
	}
	// transpose
	t := make([][]byte, m)
	for i := 0; i < m; i++ {
		t[i] = make([]byte, n)
		for j := 0; j < n; j++ {
			t[i][j] = s[j][i]
		}
	}
	s = t
	n, m = m, n
	if work() {
		res := make([]string, n)
		for i := 0; i < n; i++ {
			b := make([]byte, m)
			for j := 0; j < m; j++ {
				b[j] = s[j][i]
			}
			res[i] = string(b)
		}
		return res
	}
	return []string{"0"}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	var pending string
	for {
		var line string
		if pending != "" {
			line = pending
			pending = ""
		} else {
			if !scanner.Scan() {
				break
			}
			line = strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "bad header on test %d\n", idx+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])

		grid := make([]string, n)
		for i := 0; i < n; i++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "unexpected EOF in test %d\n", idx+1)
				os.Exit(1)
			}
			grid[i] = strings.TrimSpace(scanner.Text())
		}

		expect := solveE(n, m, grid)
		input := fmt.Sprintf("%d %d\n", n, m)
		for i := 0; i < n; i++ {
			input += grid[i] + "\n"
		}

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}

		outputLines := strings.Split(strings.TrimSpace(out.String()), "\n")
		if len(outputLines) != len(expect) {
			fmt.Printf("test %d failed: expected %d lines got %d\n", idx+1, len(expect), len(outputLines))
			os.Exit(1)
		}
		mismatch := false
		for i := range expect {
			if strings.TrimSpace(outputLines[i]) != expect[i] {
				mismatch = true
				break
			}
		}
		if mismatch {
			fmt.Printf("test %d failed: expected\n%s\ngot\n%s\n", idx+1, strings.Join(expect, "\n"), strings.Join(outputLines, "\n"))
			os.Exit(1)
		}

		idx++
		if scanner.Scan() {
			pending = strings.TrimSpace(scanner.Text())
		} else {
			pending = ""
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
