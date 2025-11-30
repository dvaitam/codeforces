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

const embeddedSolutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, p int
   if _, err := fmt.Fscan(in, &n, &p); err != nil {
       return
   }
   xs := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &xs[i])
   }
   const maxT = 32
   f := make([][][]bool, maxT)
   g := make([][][]bool, maxT)
   for i := 0; i < maxT; i++ {
       f[i] = make([][]bool, maxT)
       g[i] = make([][]bool, maxT)
       for j := 0; j < maxT; j++ {
           f[i][j] = make([]bool, p)
           g[i][j] = make([]bool, p)
       }
   }
   add := func(cur, y int) int {
       v := cur * 10
       if y > 9 {
           v *= 10
       }
       return (v + y) % p
   }
   findPrev := func(y, curMod, curXor, idx int) int {
       for t := 0; t < p; t++ {
           if f[idx][curXor][t] && add(t, y) == curMod {
               return t
           }
       }
       return 0
   }
   f[0][0][0] = true
   num := make([]int, maxT)
   pos := make([]int, maxT)
   tot := 0
   found := false
   for i, x := range xs {
       if x < maxT {
           tot++
           num[tot] = x
           pos[tot] = i + 1
           for j := 0; j < maxT && !found; j++ {
               for k := 0; k < p; k++ {
                   if f[tot-1][j][k] {
                       f[tot][j][k] = true
                       t := add(k, x)
                       nx := j ^ x
                       f[tot][nx][t] = true
                       g[tot][nx][t] = true
                       if j == x && t == 0 {
                           found = true
                           break
                       }
                   }
               }
           }
       }
       if found {
           break
       }
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   if !found {
       fmt.Fprintln(out, "No")
       return
   }
   fmt.Fprintln(out, "Yes")
   curXor, curMod := 0, 0
   var ansPos []int
   for i := tot; i > 0; i-- {
       if g[i][curXor][curMod] {
           ansPos = append(ansPos, pos[i])
           y := num[i]
           curXor ^= y
           curMod = findPrev(y, curMod, curXor, i-1)
       }
   }
   fmt.Fprintln(out, len(ansPos))
   for i := len(ansPos) - 1; i >= 0; i-- {
       if i < len(ansPos)-1 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, ansPos[i])
   }
   fmt.Fprintln(out)
}`

const testcasesRaw = `
5 2 4 3 1 2 5
4 2 3 1 4 2
3 2 1 2 3
5 2 3 1 5 4 2
3 2 3 1 2
4 2 3 4 1 2
3 2 3 2 1
1 2 1
3 2 1 3 2
4 2 4 1 3 2
3 2 1 3 2
4 2 1 2 3 4
5 2 2 4 5 3 1
4 2 4 3 1 2
5 2 3 2 1 5 4
3 2 1 3 2
5 2 1 2 5 4 3
2 2 1 2
2 2 2 1
4 2 2 1 4 3
4 2 4 1 2 3
1 2 1
1 2 1
1 2 1
2 2 2 1
3 2 2 3 1
3 2 3 1 2
5 2 2 5 1 3 4
4 2 2 4 1 3
1 2 1
3 2 3 2 1
3 2 3 2 1
1 2 1
2 2 2 1
4 2 3 2 1 4
3 2 3 2 1
5 2 4 5 3 2 1
5 2 2 5 4 1 3
5 2 2 5 3 4 1
3 2 2 3 1
5 2 1 2 4 3 5
1 2 1
5 2 1 3 5 4 2
5 2 2 3 4 5 1
1 2 1
5 2 2 3 1 5 4
3 2 2 1 3
1 2 1
2 2 1 2
3 2 3 1 2
2 2 2 1
2 2 2 1
5 2 5 4 3 1 2
1 2 1
5 2 1 2 3 5 4
3 2 2 1 3
3 2 1 2 3
4 2 4 3 2 1
1 2 1
2 2 1 2
5 2 4 5 2 1 3
2 2 1 2
2 2 2 1
1 2 1
3 2 2 1 3
5 2 5 1 4 3 2
2 2 1 2
4 2 4 1 2 3
1 2 1
5 2 2 1 4 5 3
1 2 1
2 2 2 1
4 2 1 3 2 4
5 2 2 5 3 1 4
5 2 4 3 2 5 1
4 2 1 3 4 2
1 2 1
2 2 1 2
1 2 1
1 2 1
2 2 2 1
1 2 1
2 2 2 1
3 2 3 2 1
5 2 5 1 4 2 3
4 2 4 3 2 1
5 2 1 2 5 4 3
2 2 1 2
2 2 1 2
5 2 1 5 3 2 4
3 2 3 2 1
5 2 3 4 5 2 1
5 2 1 4 2 3 5
2 2 1 2
1 2 1
4 2 2 3 1 4
1 2 1
1 2 1
2 2 2 1
4 2 4 2 1 3
`

var (
	_            = embeddedSolutionSource
	rawTestcases = func() []string {
		scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
		scanner.Buffer(make([]byte, 0, 1024), 1024*1024)
		var cases []string
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				cases = append(cases, line)
			}
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		return cases
	}()
)

const maxT = 32

func addMod(cur, y, p int) int {
	v := cur * 10
	if y > 9 {
		v *= 10
	}
	return (v + y) % p
}

func solveCase(n, p int, xs []int) (string, error) {
	f := make([][][]bool, maxT)
	g := make([][][]bool, maxT)
	for i := 0; i < maxT; i++ {
		f[i] = make([][]bool, maxT)
		g[i] = make([][]bool, maxT)
		for j := 0; j < maxT; j++ {
			f[i][j] = make([]bool, p)
			g[i][j] = make([]bool, p)
		}
	}
	findPrev := func(y, curMod, curXor, idx int) int {
		for t := 0; t < p; t++ {
			if f[idx][curXor][t] && addMod(t, y, p) == curMod {
				return t
			}
		}
		return 0
	}

	f[0][0][0] = true
	num := make([]int, maxT)
	pos := make([]int, maxT)
	tot := 0
	found := false
	for i, x := range xs {
		if x < maxT {
			tot++
			num[tot] = x
			pos[tot] = i + 1
			for j := 0; j < maxT && !found; j++ {
				for k := 0; k < p; k++ {
					if f[tot-1][j][k] {
						f[tot][j][k] = true
						t := addMod(k, x, p)
						nx := j ^ x
						f[tot][nx][t] = true
						g[tot][nx][t] = true
						if j == x && t == 0 {
							found = true
							break
						}
					}
				}
			}
		}
		if found {
			break
		}
	}

	if !found {
		return "No", nil
	}
	var sb strings.Builder
	sb.WriteString("Yes\n")
	curXor, curMod := 0, 0
	var ansPos []int
	for i := tot; i > 0; i-- {
		if g[i][curXor][curMod] {
			ansPos = append(ansPos, pos[i])
			y := num[i]
			curXor ^= y
			curMod = findPrev(y, curMod, curXor, i-1)
		}
	}
	sb.WriteString(strconv.Itoa(len(ansPos)))
	sb.WriteByte('\n')
	for i := len(ansPos) - 1; i >= 0; i-- {
		if i != len(ansPos)-1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(ansPos[i]))
	}
	if len(ansPos) > 0 {
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}

func parseCase(line string) (int, int, []int, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) < 2 {
		return 0, 0, nil, fmt.Errorf("invalid line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, nil, err
	}
	p, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, nil, err
	}
	if len(fields) != 2+n {
		return 0, 0, nil, fmt.Errorf("expected %d numbers got %d", 2+n, len(fields))
	}
	xs := make([]int, n)
	for i := 0; i < n; i++ {
		val, _ := strconv.Atoi(fields[2+i])
		xs[i] = val
	}
	return n, p, xs, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, line := range rawTestcases {
		n, p, xs, err := parseCase(line)
		if err != nil {
			fmt.Printf("case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected, err := solveCase(n, p, xs)
		if err != nil {
			fmt.Printf("case %d solve error: %v\n", idx+1, err)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d %d\n%s\n", n, p, strings.Join(strings.Fields(line)[2:], " "))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}
