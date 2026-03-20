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
   "fmt"
   "os"
)

type Query struct {
   x1, y1, x2, y2 int
   idx            int
   Pbits          []uint64
   ans            bool
}

var (
   n, m int
   grid [][]byte
   rowWalls [][]int
   answers []bool
   mWords   int
)

func solve(l, r int, qs []*Query) {
   if len(qs) == 0 {
       return
   }
   if l == r {
       for _, q := range qs {
           if rowWalls[l][q.y2] - rowWalls[l][q.y1-1] == 0 {
               answers[q.idx] = true
           } else {
               answers[q.idx] = false
           }
       }
       return
   }
   mid := (l + r) >> 1
   startBuckets := make([][]*Query, mid-l+1)
   endBuckets := make([][]*Query, r-mid+1)
   var leftQ, rightQ []*Query
   for _, q := range qs {
       if q.x2 < mid {
           leftQ = append(leftQ, q)
       } else if q.x1 > mid {
           rightQ = append(rightQ, q)
       } else {
           startBuckets[q.x1-l] = append(startBuckets[q.x1-l], q)
           endBuckets[q.x2-mid] = append(endBuckets[q.x2-mid], q)
       }
   }
   dpNext := make([][]uint64, m+2)
   dpRow := make([][]uint64, m+2)
   for j := 0; j <= m+1; j++ {
       dpNext[j] = make([]uint64, mWords)
       dpRow[j] = make([]uint64, mWords)
   }
   for j := 1; j <= m; j++ {
       if grid[mid][j] == '.' {
           bit := uint(j - 1)
           dpNext[j][bit/64] |= 1 << (bit % 64)
       }
   }
   for i := mid; i >= l; i-- {
       if i < mid {
           for j := m; j >= 1; j-- {
               if grid[i][j] == '.' {
                   for d := 0; d < mWords; d++ {
                       dpRow[j][d] = dpNext[j][d] | dpRow[j+1][d]
                   }
               } else {
                   for d := 0; d < mWords; d++ {
                       dpRow[j][d] = 0
                   }
               }
           }
           for j := 1; j <= m; j++ {
               copy(dpNext[j], dpRow[j])
           }
       }
       for _, q := range startBuckets[i-l] {
           q.Pbits = make([]uint64, mWords)
           copy(q.Pbits, dpNext[q.y1])
       }
   }
   for j := 0; j <= m+1; j++ {
       for d := 0; d < mWords; d++ {
           dpNext[j][d] = 0
           dpRow[j][d] = 0
       }
   }
   for j := 1; j <= m; j++ {
       if grid[mid][j] == '.' {
           bit := uint(j - 1)
           dpNext[j][bit/64] |= 1 << (bit % 64)
       }
   }
   for i := mid; i <= r; i++ {
       if i > mid {
           for j := 1; j <= m; j++ {
               if grid[i][j] == '.' {
                   for d := 0; d < mWords; d++ {
                       dpRow[j][d] = dpNext[j][d] | dpRow[j-1][d]
                   }
               } else {
                   for d := 0; d < mWords; d++ {
                       dpRow[j][d] = 0
                   }
               }
           }
           for j := 1; j <= m; j++ {
               copy(dpNext[j], dpRow[j])
           }
       }
       for _, q := range endBuckets[i-mid] {
           lbit := q.y1 - 1
           rbit := q.y2 - 1
           w1 := lbit / 64
           w2 := rbit / 64
           ok := false
           for w := w1; w <= w2; w++ {
               mask := ^uint64(0)
               if w == w1 {
                   mask &= ^((1 << (lbit % 64)) - 1)
               }
               if w == w2 {
                   mask &= (1 << ((rbit % 64) + 1)) - 1
               }
               if q.Pbits[w]&dpNext[q.y2][w]&mask != 0 {
                   ok = true
                   break
               }
           }
           answers[q.idx] = ok
       }
   }
   solve(l, mid-1, leftQ)
   solve(mid+1, r, rightQ)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &m)
   grid = make([][]byte, n+1)
   var line string
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &line)
       grid[i] = []byte(" " + line)
   }
   rowWalls = make([][]int, n+1)
   for i := 1; i <= n; i++ {
       rowWalls[i] = make([]int, m+1)
       for j := 1; j <= m; j++ {
           rowWalls[i][j] = rowWalls[i][j-1]
           if grid[i][j] == '#' {
               rowWalls[i][j]++
           }
       }
   }
   var q int
   fmt.Fscan(in, &q)
   queries := make([]*Query, q)
   answers = make([]bool, q)
   mWords = (m + 63) / 64
   for i := 0; i < q; i++ {
       var x1, y1, x2, y2 int
       fmt.Fscan(in, &x1, &y1, &x2, &y2)
       queries[i] = &Query{x1: x1, y1: y1, x2: x2, y2: y2, idx: i}
   }
   solve(1, n, queries)
   for i := 0; i < q; i++ {
       if answers[i] {
           fmt.Fprintln(out, "Yes")
       } else {
           fmt.Fprintln(out, "No")
       }
   }
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
