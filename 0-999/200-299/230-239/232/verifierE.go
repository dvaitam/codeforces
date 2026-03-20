package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases to keep the verifier self-contained.
const embeddedTestcases = `7 7 .#.###. #....## .#....# ....... ...##.# ##..#.# #.#.#.. 1 1 3 5 6
2 5 .#.#. #.... 4 2 5 2 5 2 4 2 4 2 5 2 5 1 5 2 5
8 8 .###.... #.#..#.# .....#.. ......#. .#...##. ........ ..#..### ..#..#.. 5 7 5 8 8 5 3 6 6 6 5 8 7 3 5 6 8 6 3 6 4
8 7 #.#.#.# #..#.## ....#.# #..###. .#...## ....... ..#..## ...#... 1 7 5 8 5
3 3 ... .#. ... 3 2 3 2 3 1 2 3 2 1 2 1 2
5 7 #.....# #...... ...##.# ......# ##...#. 2 5 7 5 7 5 4 5 7
7 8 .###.#.# .##...#. ...#.#.# ...#.### ......#. ...#.#.# ##..#.## 3 5 5 6 5 6 3 6 5 4 3 4 3
3 6 ....#. #..... ..#.## 1 2 3 2 6
6 2 ## .. .. .# .# #. 5 6 2 6 2 2 2 6 2 2 1 3 2 2 1 2 2 3 1 4 1
3 8 ..#.#... ##..#... ###..#.. 1 2 7 2 8
5 6 ...#.. .#.... ##...# ..#... #..... 1 4 5 5 5
4 5 .#... #..#. ..... ..##. 3 3 3 3 5 1 1 1 3 3 5 4 5
4 3 .#. ... .#. ##. 4 3 3 3 3 1 3 1 3 4 3 4 3 2 2 3 3
8 1 . . . # . . . # 5 2 1 7 1 1 1 1 1 1 1 5 1 3 1 6 1 6 1 6 1
8 6 #.#... #..... ...... #..#.. ..#### ..#... ....#. #.#... 1 6 2 6 4
5 5 ...#. ..... #.... .#.#. .#.#. 2 2 4 3 4 1 5 3 5
5 8 ......#. ..##.... ..#..#.. ##...... .....##. 2 5 2 5 8 2 8 3 8
3 4 #... #.#. .... 2 3 4 3 4 3 1 3 2
7 4 ##.. .##. .#.. .... #... .... ..#. 2 4 4 5 4 7 4 7 4
7 5 .#### ..#.. ..#.. .#.## ###.# .#..# ...#. 1 4 1 6 3
3 8 ..#..... #.##.... #...#... 3 3 4 3 7 1 6 3 6 2 5 3 7
8 7 .#....# #.####. ..#.... ..#.#.. .###..# .####.. #..#... #...#.. 4 1 6 4 6 8 2 8 6 2 2 3 5 1 5 8 6
3 5 ..... ..... ..### 3 2 5 2 5 1 2 2 2 1 5 1 5
4 5 .##.. ##... ..#.# .#.#. 3 1 1 4 5 4 5 4 5 2 3 2 3
3 1 . . . 1 3 1 3 1
3 6 ##.... ..#... .#.... 4 3 6 3 6 2 5 3 6 2 1 2 1 1 5 3 5
3 3 #.. ##. .#. 4 3 1 3 3 1 3 2 3 3 1 3 1 1 3 3 3
3 7 #...##. ......# ..#..#. 3 2 5 3 7 1 3 1 7 3 7 3 7
7 1 # # # . . . . 2 5 1 7 1 7 1 7 1
5 3 #.. ... ... ... ... 1 4 3 4 3
3 4 .... #..# .### 3 1 1 3 1 1 1 1 3 3 1 3 1
3 7 ..#.... #.##... ....#.. 4 2 2 2 6 2 2 3 6 1 4 1 4 2 2 3 2
3 1 # . . 1 3 1 3 1
6 2 .. .# .. .# #. .. 3 6 2 6 2 3 2 5 2 6 2 6 2
8 4 .... ...# .##. .... .#.. .... .#.# ..## 1 4 4 4 4
1 6 ##.... 1 1 3 1 4
3 7 .#..... ....... ..##..# 4 2 2 3 5 2 2 2 5 2 7 2 7 2 6 2 7
8 4 ..## .##. ...# .... .#.. ...# ...# ...# 5 4 1 5 4 1 1 7 2 4 1 5 4 5 4 5 4 6 3 6 3
4 1 . # . . 1 4 1 4 1
7 6 ..#... .#...# ..#... #..... ..#... .##... #..##. 1 5 1 7 2
8 1 . # . . # # . # 4 4 1 7 1 7 1 7 1 1 1 3 1 1 1 7 1
4 5 ##.## ...#. ..... .#... 4 3 3 4 4 4 1 4 5 3 2 3 3 2 3 2 3
5 1 . # # # . 2 5 1 5 1 5 1 5 1
2 1 . . 1 1 1 2 1
7 6 ...#.# ##.#.. ##..#. ...#.. ..###. ...### .#...# 3 2 5 3 6 4 6 4 6 3 4 3 4
4 2 ## ## .# .# 3 3 1 3 1 4 1 4 1 4 1 4 1
3 8 ....###. #..#...# #.#..#.. 4 1 8 1 8 3 5 3 8 2 3 3 4 3 2 3 8
3 8 .#.....# #..###.# ....#.#. 2 3 1 3 3 1 5 1 6
3 2 #. ## .# 1 3 1 3 1
6 7 #...#.. ...##.# #..#.## ###..## ....... #.....# 2 5 4 6 5 6 2 6 3
3 5 ..#.# ....# ..... 3 2 3 2 3 1 2 3 3 3 4 3 5
5 8 ...##.#. .##...## ...###.# #..#.... ###....# 5 5 4 5 6 2 4 4 7 3 3 5 7 4 3 4 7 5 5 5 5
3 4 #... ..#. ...# 3 3 2 3 3 1 2 2 2 1 4 2 4
3 2 .# .. #. 4 2 2 3 2 2 1 2 1 3 2 3 2 3 2 3 2
6 4 #... .#.. .... ###. ..#. ..#. 2 1 4 3 4 5 4 6 4
1 3 ... 5 1 3 1 3 1 3 1 3 1 3 1 3 1 1 1 3 1 1 1 2
8 2 .. .. .. ## .# #. .. .. 5 3 2 8 2 1 1 3 1 3 1 3 1 5 1 7 2 3 2 3 2
6 5 ..... #..#. .##.. ###.# #.... .#... 4 6 4 6 5 5 5 6 5 5 5 5 5 1 3 1 4
5 4 .... .#.. .##. ...# .#.. 2 3 1 3 4 2 1 2 4
8 4 .#.. .... ..#. ..#. #... .... .... .... 1 3 4 4 4
8 5 #..#. .#..# .#... #.... #.... #.##. #.... ..#.. 1 8 5 8 5
7 6 ....#. ...#.# .#.... .##..# ##.##. .....# #.#... 4 7 5 7 6 1 4 4 5 1 1 2 5 1 2 6 3
1 7 ....... 1 1 4 1 6
5 7 .#....# .#...#. .#..... ##...#. #..#.## 4 4 5 4 7 4 5 4 5 2 5 2 5 3 5 4 5
1 1 . 2 1 1 1 1 1 1 1 1
3 1 . . # 3 1 1 2 1 1 1 1 1 2 1 2 1
4 3 .## ... ... #.. 1 2 2 2 3
4 4 #..# .#.. ##.. #..# 2 3 4 3 4 1 3 2 4
7 2 .. .. #. ## #. #. .. 2 7 1 7 2 2 2 6 2
5 2 .# ## .. .. .. 5 3 1 3 2 4 2 4 2 4 2 5 2 4 2 4 2 1 1 4 2
6 5 #.... ..... #...# .#... .#... .#.#. 5 4 5 4 5 1 5 5 5 2 3 5 3 5 4 6 5 2 2 2 3
6 8 #...##.. ..#.#..# .....#.. ........ ..#..... ...#.#.. 4 1 2 5 8 1 3 1 3 3 5 4 5 6 1 6 3
6 7 ##.##.# ....... ..##### ...#..# ..####. ###.### 2 5 7 5 7 2 4 2 6
3 3 .#. .#. ### 3 1 3 2 3 2 3 2 3 1 3 1 3
7 7 ....... ...##.. ##..... .#..#.. ##.#... ......# #.....# 4 1 4 2 7 1 6 1 6 6 2 6 2 7 3 7 6
1 6 ..##.. 4 1 5 1 5 1 2 1 6 1 5 1 5 1 2 1 5
3 8 #....... ...#...# #....#.. 4 3 5 3 5 1 5 1 7 3 2 3 5 2 5 2 7
2 6 .#.#.. .#.#.. 5 1 5 2 5 2 6 2 6 1 1 2 6 2 6 2 6 2 3 2 3
3 8 #..#.#.. ...#..## ........ 1 1 3 3 8
2 2 .. .. 4 1 1 2 1 1 2 1 2 2 1 2 1 2 2 2 2
2 4 .#.# ...# 1 2 3 2 3
1 7 #..#.#. 1 1 2 1 3
1 5 ##.#. 2 1 3 1 3 1 5 1 5
8 1 . # # . . . . # 5 1 1 4 1 5 1 7 1 5 1 5 1 6 1 6 1 1 1 5 1
2 1 . . 5 2 1 2 1 2 1 2 1 1 1 1 1 1 1 2 1 2 1 2 1
3 8 .......# ........ ##..#### 2 2 4 2 4 2 1 2 8
5 4 ...# .... ..#. .... ...# 4 3 1 5 2 5 3 5 3 2 4 4 4 1 3 3 4
5 4 ..#. .#.# .... .... .... 1 2 1 4 3
6 2 #. .. .# .# .. .# 5 1 2 5 2 3 1 5 1 4 1 4 1 1 2 2 2 4 1 6 1
7 8 #......# ...#...# #.#..... ..#.#..# ...#..## .###.##. #.#.#.#. 4 7 8 7 8 2 5 7 8 1 5 3 6 2 6 3 8
3 1 . # . 3 3 1 3 1 3 1 3 1 3 1 3 1
4 2 .# .# .# .# 2 4 1 4 1 2 1 2 1
6 4 #... ..## .##. #... #.#. .#.# 3 2 1 6 1 3 4 3 4 3 1 3 4
8 3 #.. ... ..# .#. #.. #.# ... ... 3 7 1 7 3 8 1 8 3 1 2 8 2
7 4 .... ##.. ##.# .... ..## ..## ..## 5 6 1 6 2 1 2 1 2 5 2 7 2 3 3 4 3 4 1 5 1
5 5 ..... .##.. .#.#. ##.## ..##. 1 4 3 5 5
1 7 ..#..#. 3 1 7 1 7 1 1 1 7 1 4 1 5
2 2 .. ## 5 1 1 1 1 1 1 1 2 1 1 1 2 1 2 1 2 1 1 1 1
7 7 .#..#.. ....... ..#.... ..#.#.. .#.#... #..#..# ..#...# 5 3 2 6 5 3 4 4 7 3 1 5 6 3 4 3 4 1 1 4 4
4 2 .. .. .# .# 2 2 1 4 1 3 1 3 1`

// Embedded correct oracle source for 232E.
const oracleSource = `package main

import (
	"bufio"
	"io"
	"os"
)

type Query struct {
	r1, c1, r2, c2 int
}

var (
	buffer  []byte
	cursor  int
	queries []Query
	answers []bool
	grid    []string
	mask1   [505][505][8]uint64
	mask2   [505][505][8]uint64
)

func nextInt() int {
	for cursor < len(buffer) && (buffer[cursor] < '0' || buffer[cursor] > '9') {
		cursor++
	}
	if cursor >= len(buffer) {
		return 0
	}
	res := 0
	for cursor < len(buffer) && buffer[cursor] >= '0' && buffer[cursor] <= '9' {
		res = res*10 + int(buffer[cursor]-'0')
		cursor++
	}
	return res
}

func nextString() string {
	for cursor < len(buffer) && buffer[cursor] <= ' ' {
		cursor++
	}
	if cursor >= len(buffer) {
		return ""
	}
	start := cursor
	for cursor < len(buffer) && buffer[cursor] > ' ' {
		cursor++
	}
	return string(buffer[start:cursor])
}

func solve(r1, r2, c1, c2 int, qIdx []int) {
	if len(qIdx) == 0 {
		return
	}
	if r1 == r2 && c1 == c2 {
		for _, id := range qIdx {
			answers[id] = true
		}
		return
	}

	if r2-r1 > c2-c1 {
		mid := (r1 + r2) / 2
		var qCross, qTop, qBot []int
		for _, id := range qIdx {
			q := queries[id]
			if q.r1 <= mid && q.r2 >= mid {
				qCross = append(qCross, id)
			} else if q.r2 < mid {
				qTop = append(qTop, id)
			} else {
				qBot = append(qBot, id)
			}
		}

		if len(qCross) > 0 {
			numWords := (c2 - c1 + 1 + 63) / 64

			for r := mid; r >= r1; r-- {
				for c := c2; c >= c1; c-- {
					for w := 0; w < numWords; w++ {
						mask1[r][c][w] = 0
					}
					if grid[r][c] == '#' {
						continue
					}
					if r == mid {
						bitIdx := c - c1
						mask1[r][c][bitIdx/64] |= 1 << (bitIdx % 64)
					}
					if r < mid && grid[r+1][c] == '.' {
						for w := 0; w < numWords; w++ {
							mask1[r][c][w] |= mask1[r+1][c][w]
						}
					}
					if c < c2 && grid[r][c+1] == '.' {
						for w := 0; w < numWords; w++ {
							mask1[r][c][w] |= mask1[r][c+1][w]
						}
					}
				}
			}

			for r := mid; r <= r2; r++ {
				for c := c1; c <= c2; c++ {
					for w := 0; w < numWords; w++ {
						mask2[r][c][w] = 0
					}
					if grid[r][c] == '#' {
						continue
					}
					if r == mid {
						bitIdx := c - c1
						mask2[r][c][bitIdx/64] |= 1 << (bitIdx % 64)
					}
					if r > mid && grid[r-1][c] == '.' {
						for w := 0; w < numWords; w++ {
							mask2[r][c][w] |= mask2[r-1][c][w]
						}
					}
					if c > c1 && grid[r][c-1] == '.' {
						for w := 0; w < numWords; w++ {
							mask2[r][c][w] |= mask2[r][c-1][w]
						}
					}
				}
			}

			for _, id := range qCross {
				q := queries[id]
				ans := false
				for w := 0; w < numWords; w++ {
					if (mask1[q.r1][q.c1][w] & mask2[q.r2][q.c2][w]) != 0 {
						ans = true
						break
					}
				}
				answers[id] = ans
			}
		}

		solve(r1, mid-1, c1, c2, qTop)
		solve(mid+1, r2, c1, c2, qBot)

	} else {
		mid := (c1 + c2) / 2
		var qCross, qLeft, qRight []int
		for _, id := range qIdx {
			q := queries[id]
			if q.c1 <= mid && q.c2 >= mid {
				qCross = append(qCross, id)
			} else if q.c2 < mid {
				qLeft = append(qLeft, id)
			} else {
				qRight = append(qRight, id)
			}
		}

		if len(qCross) > 0 {
			numWords := (r2 - r1 + 1 + 63) / 64

			for c := mid; c >= c1; c-- {
				for r := r2; r >= r1; r-- {
					for w := 0; w < numWords; w++ {
						mask1[r][c][w] = 0
					}
					if grid[r][c] == '#' {
						continue
					}
					if c == mid {
						bitIdx := r - r1
						mask1[r][c][bitIdx/64] |= 1 << (bitIdx % 64)
					}
					if c < mid && grid[r][c+1] == '.' {
						for w := 0; w < numWords; w++ {
							mask1[r][c][w] |= mask1[r][c+1][w]
						}
					}
					if r < r2 && grid[r+1][c] == '.' {
						for w := 0; w < numWords; w++ {
							mask1[r][c][w] |= mask1[r+1][c][w]
						}
					}
				}
			}

			for c := mid; c <= c2; c++ {
				for r := r1; r <= r2; r++ {
					for w := 0; w < numWords; w++ {
						mask2[r][c][w] = 0
					}
					if grid[r][c] == '#' {
						continue
					}
					if c == mid {
						bitIdx := r - r1
						mask2[r][c][bitIdx/64] |= 1 << (bitIdx % 64)
					}
					if c > mid && grid[r][c-1] == '.' {
						for w := 0; w < numWords; w++ {
							mask2[r][c][w] |= mask2[r][c-1][w]
						}
					}
					if r > r1 && grid[r-1][c] == '.' {
						for w := 0; w < numWords; w++ {
							mask2[r][c][w] |= mask2[r-1][c][w]
						}
					}
				}
			}

			for _, id := range qCross {
				q := queries[id]
				ans := false
				for w := 0; w < numWords; w++ {
					if (mask1[q.r1][q.c1][w] & mask2[q.r2][q.c2][w]) != 0 {
						ans = true
						break
					}
				}
				answers[id] = ans
			}
		}

		solve(r1, r2, c1, mid-1, qLeft)
		solve(r1, r2, mid+1, c2, qRight)
	}
}

func main() {
	buffer, _ = io.ReadAll(os.Stdin)

	n := nextInt()
	m := nextInt()
	if n == 0 || m == 0 {
		return
	}

	grid = make([]string, n+1)
	for i := 1; i <= n; i++ {
		grid[i] = " " + nextString()
	}

	q := nextInt()
	queries = make([]Query, q)
	answers = make([]bool, q)
	qIdx := make([]int, q)

	for i := 0; i < q; i++ {
		queries[i].r1 = nextInt()
		queries[i].c1 = nextInt()
		queries[i].r2 = nextInt()
		queries[i].c2 = nextInt()
		qIdx[i] = i
	}

	solve(1, n, 1, m, qIdx)

	out := bufio.NewWriter(os.Stdout)
	for i := 0; i < q; i++ {
		if answers[i] {
			out.WriteString("Yes\n")
		} else {
			out.WriteString("No\n")
		}
	}
	out.Flush()
}
`

func buildOracle() (string, func(), error) {
	tmpSrc, err := os.CreateTemp("", "oracle-232E-*.go")
	if err != nil {
		return "", nil, err
	}
	if _, err := tmpSrc.WriteString(oracleSource); err != nil {
		tmpSrc.Close()
		os.Remove(tmpSrc.Name())
		return "", nil, err
	}
	tmpSrc.Close()

	tmpBin, err := os.CreateTemp("", "oracle-232E-bin-*")
	if err != nil {
		os.Remove(tmpSrc.Name())
		return "", nil, err
	}
	tmpBin.Close()

	if out, err := exec.Command("go", "build", "-o", tmpBin.Name(), tmpSrc.Name()).CombinedOutput(); err != nil {
		os.Remove(tmpSrc.Name())
		os.Remove(tmpBin.Name())
		return "", nil, fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	os.Remove(tmpSrc.Name())
	return tmpBin.Name(), func() { os.Remove(tmpBin.Name()) }, nil
}

func parseInput(fields []string) (string, error) {
	if len(fields) < 3 {
		return "", fmt.Errorf("invalid testcase line")
	}
	nVal, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	mVal, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", err
	}
	if len(fields) < 2+nVal+1 {
		return "", fmt.Errorf("too few fields")
	}
	idx := 2
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", nVal, mVal))
	for i := 0; i < nVal; i++ {
		input.WriteString(fields[idx])
		input.WriteByte('\n')
		idx++
	}
	qVal, err := strconv.Atoi(fields[idx])
	if err != nil {
		return "", err
	}
	idx++
	if len(fields) != 2+nVal+1+4*qVal {
		return "", fmt.Errorf("field count mismatch")
	}
	input.WriteString(fmt.Sprintf("%d\n", qVal))
	for i := 0; i < qVal; i++ {
		input.WriteString(fmt.Sprintf("%s %s %s %s\n", fields[idx], fields[idx+1], fields[idx+2], fields[idx+3]))
		idx += 4
	}
	return input.String(), nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	lines := strings.Split(strings.TrimSpace(embeddedTestcases), "\n")
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 0 {
			continue
		}
		input, err := parseInput(fields)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad testcase %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		expect, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\nexpected:\n%s\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
