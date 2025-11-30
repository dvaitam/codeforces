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

const testcasesB = `1 1 4 0
3 3 1 1 0 1 0 0 0 1 0 1
2 3 3 2 3 1 3 1 0
3 3 2 2 2 0 1 0 0 2 1 2
1 2 2 2 1
3 1 4 4 0 4
1 3 3 0 0 0
2 3 3 3 2 3 2 0 1
2 1 2 2 0
1 2 2 0 1
3 1 3 3 0 1
2 1 4 0 3
2 3 3 2 3 3 1 0 3
2 1 4 0 2
1 1 1 1
1 1 4 4
3 3 1 1 1 1 1 0 0 1 1 0
3 3 1 0 1 1 0 1 0 0 1 0
1 2 3 1 1
1 1 2 2
3 1 3 0 3 3
1 1 2 2
2 2 2 0 1 0 0
2 2 3 1 3 0 1
3 1 4 1 1 0
1 2 1 0 0
1 3 2 0 2 1
1 2 4 4 2
3 1 2 1 2 0
3 2 1 0 1 1 0 1 1
1 2 4 1 3
1 3 4 3 4 0
2 2 3 0 3 1 2
1 3 3 3 0 0
1 3 4 0 3 3
2 1 4 4 3
1 1 3 3
2 1 3 0 2
2 1 3 3 0
1 2 2 2 2
2 2 4 0 2 4 3
1 3 4 2 3 1
2 1 2 0 2
2 1 1 1 1
2 1 1 1 0
1 1 2 0
1 1 2 1
2 3 3 2 0 3 2 1 2
3 3 3 1 3 2 1 3 2 3 3 2
2 2 1 0 0 0 0
3 3 2 2 2 2 2 0 1 2 2 2
1 2 1 0 0
1 2 4 2 0
3 2 3 2 3 1 1 0 1
1 3 2 2 1 2
3 1 4 3 2 0
3 1 4 1 3 4
1 2 1 1 1
1 3 4 0 3 4
1 2 4 4 3
2 2 2 0 2 1 0
2 3 2 1 0 0 1 1 2
2 3 1 1 1 0 0 1 0
2 1 2 2 0
3 1 2 0 0 1
3 3 4 0 3 4 2 2 0 0 1 1
2 2 3 1 0 3 2
3 2 1 0 1 1 0 1 1
1 1 3 1
3 2 3 2 3 1 0 3 2
3 2 2 2 2 1 2 1 1
2 2 4 2 0 3 2
1 2 2 1 0
3 1 3 3 1 0
3 2 4 3 4 3 4 4 1
3 2 3 3 1 2 0 0 3
2 2 2 1 0 2 0
2 1 3 1 2
1 1 2 2
2 1 2 0 0
2 3 1 0 0 0 0 1 0
1 3 4 0 4 1
1 1 4 3
2 1 4 1 3
1 3 3 0 2 2
2 3 4 1 4 4 0 2 0
3 1 3 0 1 1
2 1 3 1 3
3 3 1 0 1 0 1 1 0 1 1 0
2 2 4 0 4 4 2
1 2 2 1 1
2 1 4 0 4
2 2 2 1 0 1 1
2 1 2 1 1
2 3 4 4 1 4 0 1 2
3 3 1 0 0 1 0 1 0 1 1 0
2 2 2 0 1 0 1
3 3 2 2 1 0 2 2 2 2 2 1
1 3 4 4 1 4
2 1 2 1 2`

const mod = 1000000007

func expected(n, m, k int, board [][]int) int {
	coords := make([][2]int, 0, n*m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			coords = append(coords, [2]int{i, j})
		}
	}
	N := len(coords)
	prevList := make([][]int, N)
	for i := 0; i < N; i++ {
		x1, y1 := coords[i][0], coords[i][1]
		for j := 0; j < i; j++ {
			x2, y2 := coords[j][0], coords[j][1]
			if (x1 <= x2 && y1 <= y2) || (x2 <= x1 && y2 <= y1) {
				prevList[i] = append(prevList[i], j)
			}
		}
	}
	color := make([]int, N)
	var dfs func(int) int
	dfs = func(pos int) int {
		if pos == N {
			return 1
		}
		x, y := coords[pos][0], coords[pos][1]
		if board[x][y] != 0 {
			c := board[x][y]
			for _, v := range prevList[pos] {
				if color[v] == c {
					return 0
				}
			}
			color[pos] = c
			res := dfs(pos + 1)
			color[pos] = 0
			return res
		}
		total := 0
		for c := 1; c <= k; c++ {
			ok := true
			for _, v := range prevList[pos] {
				if color[v] == c {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}
			color[pos] = c
			total += dfs(pos + 1)
			if total >= mod {
				total -= mod
			}
			color[pos] = 0
		}
		return total
	}

	if n > k || m > k || n+m-1 > k {
		return 0
	}
	return dfs(0) % mod
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseCase(line string) (int, int, int, [][]int, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return 0, 0, 0, nil, fmt.Errorf("invalid line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, 0, nil, err
	}
	m, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, 0, nil, err
	}
	k, err := strconv.Atoi(fields[2])
	if err != nil {
		return 0, 0, 0, nil, err
	}
	expectedCells := n * m
	if len(fields)-3 != expectedCells {
		return 0, 0, 0, nil, fmt.Errorf("expected %d cells got %d", expectedCells, len(fields)-3)
	}
	board := make([][]int, n)
	idx := 3
	for i := 0; i < n; i++ {
		board[i] = make([]int, m)
		for j := 0; j < m; j++ {
			val, convErr := strconv.Atoi(fields[idx])
			if convErr != nil {
				return 0, 0, 0, nil, convErr
			}
			board[i][j] = val
			idx++
		}
	}
	return n, m, k, board, nil
}

func buildInput(n, m, k int, board [][]int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(board[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesB))
	scanner.Buffer(make([]byte, 0, 1024), 1<<20)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, m, k, board, err := parseCase(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d parse error: %v\n", idx, err)
			os.Exit(1)
		}
		input := buildInput(n, m, k, board)
		expected := strconv.Itoa(expected(n, m, k, board))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed. Expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
