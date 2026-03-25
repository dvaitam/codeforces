package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var dir = []byte{'<', '^', '>'}

type Pair struct {
	r, c int
}

// solveD is the embedded CF-accepted solver for 487D.
func solveD(input string) string {
	tokens := strings.Fields(input)
	pos := 0
	nextToken := func() string {
		t := tokens[pos]
		pos++
		return t
	}
	nextInt := func() int {
		v, _ := strconv.Atoi(nextToken())
		return v
	}

	n := nextInt()
	m := nextInt()
	q := nextInt()

	grid := make([][]byte, n+1)
	for i := 1; i <= n; i++ {
		grid[i] = make([]byte, m+2)
		s := nextToken()
		for j := 1; j <= m; j++ {
			grid[i][j] = s[j-1]
		}
	}

	dest := make([][]Pair, n+1)
	for i := 1; i <= n; i++ {
		dest[i] = make([]Pair, m+2)
	}

	B := 250
	numBlocks := (n + B - 1) / B
	L := make([]int, numBlocks)
	R := make([]int, numBlocks)
	for k := 0; k < numBlocks; k++ {
		L[k] = k*B + 1
		R[k] = (k + 1) * B
		if R[k] > n {
			R[k] = n
		}
	}

	block_dest := make([][]Pair, numBlocks)
	for k := 0; k < numBlocks; k++ {
		block_dest[k] = make([]Pair, m+2)
	}

	computeDest := func(i int) {
		for y := 1; y <= m; y++ {
			curY := y
			vis := 0
			for {
				if curY == 0 || curY == m+1 {
					dest[i][y] = Pair{i, curY}
					break
				}
				if vis&(1<<curY) != 0 {
					dest[i][y] = Pair{-1, -1}
					break
				}
				vis |= (1 << curY)
				if grid[i][curY] == '^' {
					dest[i][y] = Pair{i - 1, curY}
					break
				} else if grid[i][curY] == '<' {
					curY--
				} else {
					curY++
				}
			}
		}
	}

	computeBlock := func(k int) {
		for y := 1; y <= m; y++ {
			curX := R[k]
			curY := y
			for curX >= L[k] {
				d := dest[curX][curY]
				if d.r == -1 {
					curX, curY = -1, -1
					break
				}
				if d.c == 0 || d.c == m+1 {
					curX, curY = d.r, d.c
					break
				}
				curX = d.r
				curY = d.c
			}
			block_dest[k][y] = Pair{curX, curY}
		}
	}

	for i := 1; i <= n; i++ {
		computeDest(i)
	}
	for k := 0; k < numBlocks; k++ {
		computeBlock(k)
	}

	var out strings.Builder
	out.Grow(q * 15)

	for i := 0; i < q; i++ {
		typ := nextToken()
		if typ == "A" {
			x := nextInt()
			y := nextInt()
			curX, curY := x, y
			for curX > 0 {
				k := (curX - 1) / B
				if curX == R[k] {
					d := block_dest[k][curY]
					curX, curY = d.r, d.c
					if curX == -1 || curY == 0 || curY == m+1 {
						break
					}
				} else {
					d := dest[curX][curY]
					curX, curY = d.r, d.c
					if curX == -1 || curY == 0 || curY == m+1 {
						break
					}
				}
			}
			out.WriteString(strconv.Itoa(curX))
			out.WriteByte(' ')
			out.WriteString(strconv.Itoa(curY))
			out.WriteByte('\n')
		} else if typ == "C" {
			x := nextInt()
			y := nextInt()
			c := nextToken()
			grid[x][y] = c[0]
			computeDest(x)
			computeBlock((x - 1) / B)
		}
	}
	return strings.TrimSpace(out.String())
}

func run(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
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

func randBoard(rng *rand.Rand, n, m int) []string {
	board := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			b[j] = dir[rng.Intn(3)]
		}
		board[i] = string(b)
	}
	return board
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		q := rng.Intn(5) + 1
		board := randBoard(rng, n, m)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for i := 0; i < n; i++ {
			sb.WriteString(board[i])
			sb.WriteByte('\n')
		}
		for i := 0; i < q; i++ {
			if rng.Intn(2) == 0 {
				x := rng.Intn(n) + 1
				y := rng.Intn(m) + 1
				sb.WriteString(fmt.Sprintf("A %d %d\n", x, y))
			} else {
				x := rng.Intn(n) + 1
				y := rng.Intn(m) + 1
				c := string(dir[rng.Intn(3)])
				sb.WriteString(fmt.Sprintf("C %d %d %s\n", x, y, c))
			}
		}
		input := sb.String()
		expected := solveD(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\ngot:\n%s\ninput:\n%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
