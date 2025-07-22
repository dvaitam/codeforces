package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveE(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return ""
	}
	floors := make([][]byte, n)
	for i := 0; i < n; i++ {
		var line string
		fmt.Fscan(reader, &line)
		floors[n-1-i] = []byte(line)
	}
	visited := make([][2]int, m)
	curVer := 1
	floor := n - 1
	pos := 0
	dir := 1
	var timeCount int64
	reset := func() { curVer++ }
	reset()
	for {
		if floor == 0 {
			return fmt.Sprintf("%d\n", timeCount)
		}
		if floors[floor-1][pos] == '.' {
			floor--
			timeCount++
			reset()
			continue
		}
		next := pos + dir
		var cell byte
		if next < 0 || next >= m {
			cell = '#'
		} else {
			cell = floors[floor][next]
		}
		if cell == '+' {
			floors[floor][next] = '.'
			dir = -dir
			timeCount++
			reset()
			continue
		}
		if cell == '.' {
			pos = next
		} else if cell == '#' {
			dir = -dir
		}
		timeCount++
		di := 1
		if dir < 0 {
			di = 0
		}
		if visited[pos][di] == curVer {
			return "Never\n"
		}
		visited[pos][di] = curVer
	}
}

func generateCaseE(rng *rand.Rand) string {
	n := rng.Intn(5) + 2
	m := rng.Intn(10) + 1
	floors := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			ch := rng.Intn(3)
			if ch == 0 {
				row[j] = '.'
			} else if ch == 1 {
				row[j] = '+'
			} else {
				row[j] = '#'
			}
		}
		floors[n-1-i] = string(row)
	}
	// ensure first cell of top floor is '.'
	top := []byte(floors[n-1])
	top[0] = '.'
	floors[n-1] = string(top)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(floors[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseE(rng)
	}
	for i, tc := range cases {
		expect := solveE(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
