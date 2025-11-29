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

// Embedded source for the reference solution (was 1063B.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

type cell struct {
   r, c, lsteps, rsteps int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   var sr, sc int
   var x, y int
   fmt.Fscan(reader, &n, &m)
   fmt.Fscan(reader, &sr, &sc)
   fmt.Fscan(reader, &x, &y)
   maze := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       maze[i] = []byte(s)
   }
   visited := make([][]bool, n)
   for i := 0; i < n; i++ {
       visited[i] = make([]bool, m)
   }
   q := make([]cell, 0, n*m)
   head := 0
   // 0-based start
   sr--
   sc--
   q = append(q, cell{sr, sc, 0, 0})
   visited[sr][sc] = true
   ans := 1
   for head < len(q) {
       cur := q[head]
       head++
       r, c := cur.r, cur.c
       ls, rs := cur.lsteps, cur.rsteps
       // move left
       if c > 0 && ls < x && !visited[r][c-1] && maze[r][c-1] == '.' {
           visited[r][c-1] = true
           ans++
           q = append(q, cell{r, c - 1, ls + 1, rs})
       }
       // move right
       if c+1 < m && rs < y && !visited[r][c+1] && maze[r][c+1] == '.' {
           visited[r][c+1] = true
           ans++
           q = append(q, cell{r, c + 1, ls, rs + 1})
       }
       // move up
       for i := r - 1; i >= 0; i-- {
           if visited[i][c] || maze[i][c] != '.' {
               break
           }
           visited[i][c] = true
           ans++
           // left from new cell
           if c > 0 && ls < x && !visited[i][c-1] && maze[i][c-1] == '.' {
               visited[i][c-1] = true
               ans++
               q = append(q, cell{i, c - 1, ls + 1, rs})
           }
           // right from new cell
           if c+1 < m && rs < y && !visited[i][c+1] && maze[i][c+1] == '.' {
               visited[i][c+1] = true
               ans++
               q = append(q, cell{i, c + 1, ls, rs + 1})
           }
       }
       // move down
       for i := r + 1; i < n; i++ {
           if visited[i][c] || maze[i][c] != '.' {
               break
           }
           visited[i][c] = true
           ans++
           // left from new cell
           if c > 0 && ls < x && !visited[i][c-1] && maze[i][c-1] == '.' {
               visited[i][c-1] = true
               ans++
               q = append(q, cell{i, c - 1, ls + 1, rs})
           }
           // right from new cell
           if c+1 < m && rs < y && !visited[i][c+1] && maze[i][c+1] == '.' {
               visited[i][c+1] = true
               ans++
               q = append(q, cell{i, c + 1, ls, rs + 1})
           }
       }
   }
   fmt.Println(ans)
}
`

const testcasesRaw = `100
2 5 1 3 0 1
*....
*..*.
3 2 3 1 1 0
..
.*
.*
1 5 1 4 1 2
.....
1 4 1 2 2 2
*..*
5 4 5 2 1 1
.**.
...*
...*
....
*...
1 4 1 3 2 2
....
1 2 1 2 2 1
*.
4 3 1 2 2 2
*..
..*
...
...
5 5 5 5 1 1
.*...
..**.
*.*.*
..*..
..*..
3 4 3 3 1 1
...*
..**
*...
4 2 1 1 1 2
..
**
..
..
1 4 1 4 0 2
.**.
1 3 1 3 0 1
...
1 5 1 5 1 0
***..
2 3 1 1 2 2
...
*..
1 3 1 2 0 0
..*
2 3 2 1 1 2
...
.*.
5 2 1 1 0 0
.*
.*
.*
..
.*
5 4 2 1 1 0
..*.
...*
...*
....
*...
5 1 1 1 0 0
.
*
.
.
.
2 2 1 2 1 2
*.
..
1 2 1 1 0 0
.*
5 3 4 2 1 1
.**
*..
.*.
...
...
1 3 1 2 0 2
...
3 1 2 1 1 2
*
.
.
5 5 5 5 0 0
.*..*
*..*.
.*.*.
.....
*.**.
2 2 1 1 1 2
.*
**
2 2 2 2 2 0
.*
..
4 4 3 4 2 1
.*..
***.
....
*...
1 2 1 1 0 2
.*
4 2 2 2 1 2
..
..
**
..
1 1 1 1 2 0
.
3 5 2 2 2 2
*.***
.....
*..*.
5 2 5 1 2 0
*.
.*
*.
**
.*
4 2 4 2 1 0
*.
..
..
..
3 3 1 2 2 2
...
...
*..
5 2 3 2 2 1
*.
**
..
..
.*
2 3 1 1 0 2
..*
**.
5 3 5 2 1 2
.*.
.*.
..*
..*
..*
5 5 2 4 2 2
.**..
.....
..*.*
**.*.
.....
4 3 2 1 2 0
..*
...
*..
..*
4 1 3 1 2 1
.
.
.
.
4 3 3 1 2 0
...
.*.
...
.**
5 3 3 3 2 0
...
..*
...
.**
.**
2 3 1 1 2 2
.*.
**.
3 4 3 2 0 0
*...
....
..**
1 1 1 1 2 0
.
5 5 4 3 2 1
.....
...*.
..*..
**.**
...*.
2 1 1 1 0 2
.
.
3 2 1 2 2 1
*.
.*
..
1 4 1 2 0 1
...*
4 4 3 3 1 2
...*
....
**..
*.*.
3 5 3 1 1 0
*.***
..*..
.....
4 4 1 3 0 1
*...
*..*
**..
....
4 5 4 1 0 2
*..*.
..***
.....
..*..
3 3 2 2 2 2
...
..*
...
4 3 4 1 2 0
...
***
*.*
...
3 5 3 2 0 0
....*
*.**.
*....
1 2 1 2 0 2
..
3 2 1 2 0 1
..
*.
..
3 4 1 1 2 2
....
**..
*...
1 5 1 3 2 2
.*.*.
4 5 4 5 2 2
.....
*..**
...**
.***.
4 3 3 3 2 1
*.*
...
**.
.**
4 3 3 3 1 2
..*
.**
...
..*
1 3 1 2 0 1
...
4 5 1 5 0 2
...*.
.***.
.**..
.*...
4 4 1 4 0 0
.*..
.**.
.*..
*...
3 3 1 3 2 1
.*.
...
...
5 1 5 1 1 1
.
.
.
.
.
4 2 3 1 2 0
*.
..
..
.*
2 3 2 2 1 0
...
..*
5 2 1 1 2 0
..
..
**
*.
*.
3 3 3 3 2 1
...
.*.
...
1 3 1 1 1 2
...
4 3 2 1 2 1
...
..*
.*.
*..
5 5 4 5 1 2
*....
...*.
..*..
**.*.
.*...
4 2 2 1 2 0
.*
.*
*.
..
4 4 4 1 2 0
..*.
.**.
..*.
...*
4 4 4 1 0 1
*..*
*..*
.*.*
...*
5 2 1 2 0 0
*.
.*
..
**
..
4 3 3 2 1 2
...
...
..*
*..
4 2 4 1 1 2
..
*.
..
..
2 4 2 4 2 0
.**.
....
2 4 2 3 0 0
***.
...*
5 4 1 2 1 0
...*
**..
..**
**..
..*.
1 4 1 1 2 1
..*.
3 4 2 4 0 2
..**
.**.
..*.
3 3 1 2 2 1
*..
*..
...
4 1 2 1 1 0
.
.
.
.
4 2 2 2 1 2
..
..
.*
..
5 5 2 5 2 0
**..*
.*...
....*
.*...
...*.
2 1 1 1 2 2
.
.
5 4 1 2 1 2
*.*.
...*
....
..*.
***.
3 4 1 2 0 2
....
**..
.***
2 2 1 2 2 1
..
**
1 3 1 2 0 2
..*
3 4 1 1 2 1
....
...*
.**.
3 5 3 3 0 2
.*.**
*.**.
...*.
3 5 1 3 2 0
.....
*....
*....
`

type cell struct {
	r, c, lsteps, rsteps int
}

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

func computeExpected(n, m, sr, sc, x, y int, rows []string) string {
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = []byte(rows[i])
	}
	visited := make([][]bool, n)
	for i := 0; i < n; i++ {
		visited[i] = make([]bool, m)
	}
	sr--
	sc--
	q := make([]cell, 0, n*m)
	q = append(q, cell{sr, sc, 0, 0})
	visited[sr][sc] = true
	ans := 1
	for head := 0; head < len(q); head++ {
		cur := q[head]
		r, c := cur.r, cur.c
		ls, rs := cur.lsteps, cur.rsteps
		if c > 0 && ls < x && !visited[r][c-1] && grid[r][c-1] == '.' {
			visited[r][c-1] = true
			ans++
			q = append(q, cell{r, c - 1, ls + 1, rs})
		}
		if c+1 < m && rs < y && !visited[r][c+1] && grid[r][c+1] == '.' {
			visited[r][c+1] = true
			ans++
			q = append(q, cell{r, c + 1, ls, rs + 1})
		}
		for i := r - 1; i >= 0; i-- {
			if visited[i][c] || grid[i][c] != '.' {
				break
			}
			visited[i][c] = true
			ans++
			if c > 0 && ls < x && !visited[i][c-1] && grid[i][c-1] == '.' {
				visited[i][c-1] = true
				ans++
				q = append(q, cell{i, c - 1, ls + 1, rs})
			}
			if c+1 < m && rs < y && !visited[i][c+1] && grid[i][c+1] == '.' {
				visited[i][c+1] = true
				ans++
				q = append(q, cell{i, c + 1, ls, rs + 1})
			}
		}
		for i := r + 1; i < n; i++ {
			if visited[i][c] || grid[i][c] != '.' {
				break
			}
			visited[i][c] = true
			ans++
			if c > 0 && ls < x && !visited[i][c-1] && grid[i][c-1] == '.' {
				visited[i][c-1] = true
				ans++
				q = append(q, cell{i, c - 1, ls + 1, rs})
			}
			if c+1 < m && rs < y && !visited[i][c+1] && grid[i][c+1] == '.' {
				visited[i][c+1] = true
				ans++
				q = append(q, cell{i, c + 1, ls, rs + 1})
			}
		}
	}
	return strconv.Itoa(ans)
}

func main() {
	var _ = solutionSource
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	if !scanner.Scan() {
		fmt.Println("empty test data")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scanner.Scan() {
			fmt.Printf("unexpected EOF on case %d\n", caseNum)
			os.Exit(1)
		}
		parts := strings.Fields(scanner.Text())
		if len(parts) != 6 {
			fmt.Printf("bad header on case %d\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		sr, _ := strconv.Atoi(parts[2])
		sc, _ := strconv.Atoi(parts[3])
		x, _ := strconv.Atoi(parts[4])
		y, _ := strconv.Atoi(parts[5])
		rows := make([]string, n)
		for i := 0; i < n; i++ {
			if !scanner.Scan() {
				fmt.Printf("unexpected EOF on case %d\n", caseNum)
				os.Exit(1)
			}
			rows[i] = scanner.Text()
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		sb.WriteString(fmt.Sprintf("%d %d\n", sr, sc))
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		for _, r := range rows {
			sb.WriteString(r)
			sb.WriteByte('\n')
		}
		input := sb.String()
		exp := computeExpected(n, m, sr, sc, x, y, rows)
		out, err := run(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", caseNum, exp, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
