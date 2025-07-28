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

type testCase struct {
	input    string
	expected string
}

func expectedA(n int, row1, row2 string) string {
	grid := []string{row1, row2}
	dirs := [8][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	vis := make([][]bool, 2)
	for i := range vis {
		vis[i] = make([]bool, n)
	}
	type node struct{ r, c int }
	q := []node{{0, 0}}
	vis[0][0] = true
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur.r == 1 && cur.c == n-1 {
			return "YES"
		}
		for _, d := range dirs {
			nr, nc := cur.r+d[0], cur.c+d[1]
			if nr < 0 || nr >= 2 || nc < 0 || nc >= n {
				continue
			}
			if grid[nr][nc] == '1' || vis[nr][nc] {
				continue
			}
			vis[nr][nc] = true
			q = append(q, node{nr, nc})
		}
	}
	if vis[1][n-1] {
		return "YES"
	}
	return "NO"
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(98) + 3 // 3..100
		row1 := make([]byte, n)
		row2 := make([]byte, n)
		for j := 0; j < n; j++ {
			if j == 0 {
				row1[j] = '0'
			} else {
				if rng.Intn(2) == 0 {
					row1[j] = '0'
				} else {
					row1[j] = '1'
				}
			}
			if j == n-1 {
				row2[j] = '0'
			} else {
				if rng.Intn(2) == 0 {
					row2[j] = '0'
				} else {
					row2[j] = '1'
				}
			}
		}
		input := fmt.Sprintf("1\n%d\n%s\n%s\n", n, string(row1), string(row2))
		exp := expectedA(n, string(row1), string(row2))
		cases[i] = testCase{input: input, expected: exp}
	}
	return cases
}

func run(bin, stringInput string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(stringInput)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if out != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
