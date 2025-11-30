package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// referenceSolutionSource embeds the original 339E solution for traceability.
const referenceSolutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

// Pair represents a command segment [first, second]
type Pair struct{ first, second int }

var n int
var a []int
var path []Pair

func dfs(remain int) bool {
   ok := true
   for i := 1; i <= n; i++ {
       if a[i] != i {
           ok = false
           break
       }
   }
   if ok {
       return true
   }
   if remain == 0 {
       return false
   }
   // generate candidate segments
   vc := make([]Pair, 0)
   vc = append(vc, Pair{1, n})
   pos := make([]int, n+1)
   for i := 1; i <= n; i++ {
       pos[a[i]] = i
   }
   // consecutive runs
   for i := 1; i <= n; {
       j := i
       for j+1 <= n && a[j+1]-a[j] == a[i+1]-a[i] && abs(a[i+1]-a[i]) == 1 {
           j++
       }
       vc = append(vc, Pair{i, j})
       i = j + 1
   }
   // adjacent value positions
   for i := 1; i <= n; i++ {
       if a[i] != 1 {
           j := pos[a[i]-1]
           if j+1 < i {
               vc = append(vc, Pair{j + 1, i})
           }
           if i+1 < j {
               vc = append(vc, Pair{i, j - 1})
           }
       }
       if a[i] != n {
           j := pos[a[i]+1]
           if j+1 < i {
               vc = append(vc, Pair{j + 1, i})
           }
           if i+1 < j {
               vc = append(vc, Pair{i, j - 1})
           }
       }
   }
   // try each candidate
   for _, p := range vc {
       L, R := p.first, p.second
       path = append(path, p)
       reverse(L, R)
       if dfs(remain - 1) {
           return true
       }
       reverse(L, R)
       path = path[:len(path)-1]
   }
   return false
}

// reverse reverses segment [l, r] in a
func reverse(l, r int) {
   for l < r {
       a[l], a[r] = a[r], a[l]
       l++
       r--
   }
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n)
   a = make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   dfs(3)
   k := len(path)
   fmt.Fprintln(writer, k)
   for i := k - 1; i >= 0; i-- {
       fmt.Fprintln(writer, path[i].first, path[i].second)
   }
}
`

var _ = referenceSolutionSource

const rawTestcases = `6 3 6 4 1 5 2
6 5 2 4 1 3 6
5 3 5 1 2 4
1 1
3 2 3 1
1 1
2 1 2
1 1
4 1 3 4 2
5 1 2 4 3 5
8 4 3 8 5 6 7 2 1
2 2 1
2 1 2
8 2 6 7 3 5 1 4 8
3 1 2 3
5 3 5 2 1 4
7 4 6 2 1 7 5 3
5 3 2 1 5 4
6 3 5 2 4 1 6
2 1 2
2 2 1
1 1
1 1
1 1
5 2 5 3 1 4
8 8 2 1 6 7 5 3 4
8 4 2 5 6 1 7 3 8
2 2 1
7 2 1 7 5 6 3 4
3 1 2 3
8 7 8 3 2 6 1 4 5
5 4 1 2 5 3
7 4 1 5 7 2 3 6
5 1 2 4 5 3
2 1 2
1 1
8 4 1 2 6 7 5 8 3
8 5 4 6 1 7 8 2 3
7 3 6 2 1 5 7 4
1 1
8 3 8 1 4 7 6 2 5
1 1
5 4 3 5 2 1
5 5 2 1 3 4
8 5 4 8 1 2 7 3 6
8 2 7 5 8 3 4 6 1
2 2 1
1 1
6 1 5 3 4 6 2
7 5 6 1 7 3 4 2
5 4 2 1 5 3
1 1
1 1
2 2 1
6 6 3 4 1 2 5
3 2 3 1
3 2 1 3
7 7 6 2 5 3 1 4
4 4 3 2 1
6 4 6 2 3 5 1
7 7 4 1 6 5 3 2
3 2 1 3
2 1 2
2 2 1
6 5 1 2 6 4 3
7 7 2 5 6 1 3 4
3 1 3 2
2 1 2
6 5 1 4 2 6 3
4 1 4 3 2
7 3 4 7 1 2 5 6
1 1
8 3 8 2 4 6 1 5 7
7 7 3 4 2 5 6 1
8 2 6 5 3 1 8 4 7
6 3 2 6 4 5 1
4 3 1 4 2
8 8 1 7 2 5 3 4 6
8 4 5 6 1 7 2 8 3
4 3 2 4 1
3 3 1 2
3 1 3 2
6 3 5 1 2 4 6
7 1 4 5 6 2 7 3
7 7 3 1 4 2 5 6
8 7 5 3 2 4 6 1 8
7 4 5 7 2 6 3 1
5 3 1 4 5 2
3 2 3 1
3 1 3 2
5 4 3 5 1 2
2 1 2
5 2 4 1 5 3
1 1
6 3 1 2 5 4 6
2 1 2
7 7 2 4 3 5 1 6
3 2 3 1
3 1 3 2
1 1`

var testcases = parseTestcases(rawTestcases)

type Pair struct{ first, second int }

func parseTestcases(data string) []string {
	lines := strings.Split(data, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}

func parseCase(line string) ([]int, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty case")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}
	if len(fields) != 1+n {
		return nil, fmt.Errorf("expected %d numbers got %d", 1+n, len(fields))
	}
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[1+i])
		if err != nil {
			return nil, err
		}
		perm[i] = v
	}
	return perm, nil
}

func solveCase(perm []int) (string, error) {
	n := len(perm)
	a := make([]int, n+1)
	for i, v := range perm {
		a[i+1] = v
	}
	path := make([]Pair, 0, 3)

	reverse := func(l, r int) {
		for l < r {
			a[l], a[r] = a[r], a[l]
			l++
			r--
		}
	}
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	var dfs func(remain int) bool
	dfs = func(remain int) bool {
		sorted := true
		for i := 1; i <= n; i++ {
			if a[i] != i {
				sorted = false
				break
			}
		}
		if sorted {
			return true
		}
		if remain == 0 {
			return false
		}

		vc := []Pair{{1, n}}
		pos := make([]int, n+1)
		for i := 1; i <= n; i++ {
			pos[a[i]] = i
		}
		for i := 1; i <= n; {
			j := i
			for j+1 <= n && a[j+1]-a[j] == a[i+1]-a[i] && abs(a[i+1]-a[i]) == 1 {
				j++
			}
			vc = append(vc, Pair{i, j})
			i = j + 1
		}
		for i := 1; i <= n; i++ {
			if a[i] != 1 {
				j := pos[a[i]-1]
				if j+1 < i {
					vc = append(vc, Pair{j + 1, i})
				}
				if i+1 < j {
					vc = append(vc, Pair{i, j - 1})
				}
			}
			if a[i] != n {
				j := pos[a[i]+1]
				if j+1 < i {
					vc = append(vc, Pair{j + 1, i})
				}
				if i+1 < j {
					vc = append(vc, Pair{i, j - 1})
				}
			}
		}

		for _, p := range vc {
			l, r := p.first, p.second
			path = append(path, p)
			reverse(l, r)
			if dfs(remain - 1) {
				return true
			}
			reverse(l, r)
			path = path[:len(path)-1]
		}
		return false
	}

	dfs(3)
	k := len(path)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(k))
	sb.WriteByte('\n')
	for i := k - 1; i >= 0; i-- {
		sb.WriteString(fmt.Sprintf("%d %d\n", path[i].first, path[i].second))
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, line := range testcases {
		perm, err := parseCase(line)
		if err != nil {
			fmt.Printf("case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}

		expect, err := solveCase(perm)
		if err != nil {
			fmt.Printf("case %d solve error: %v\n", idx+1, err)
			os.Exit(1)
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(perm)))
		for i, v := range perm {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(testcases))
}
