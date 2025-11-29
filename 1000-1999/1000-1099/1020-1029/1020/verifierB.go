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

// Embedded source for the reference solution (was 1020B.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &p[i])
   }

   for i := 1; i <= n; i++ {
       visited := make([]bool, n+1)
       x := p[i]
       visited[i] = true
       for !visited[x] {
           visited[x] = true
           x = p[x]
       }
       fmt.Fprint(writer, x, " ")
   }
   fmt.Fprintln(writer)
}
`

const testcasesRaw = `100
1 1
2 2 1
5 3 5 2 5 1
10 3 7 7 9 6 9 8 9 5 1
1 1
8 6 7 7 3 3 4 4 1
3 2 1 1
9 9 6 9 9 3 8 7 9 6
10 6 6 8 3 7 8 9 4 8 5
8 6 8 8 6 8 8 4 6
3 3 2 2
5 3 5 5 5 5
10 10 7 5 4 8 9 6 10 2 6
1 1
2 1 1
5 5 2 1 5 2
5 2 2 1 4 1
1 1
6 2 2 6 1 1 1
2 1 1
1 1
5 2 2 2 5 1
7 5 1 7 2 2 1 1
6 5 6 6 6 1 3
6 4 1 3 4 5 5
1 1
7 7 5 6 2 4 2 1
6 1 1 4 2 5 5
7 4 5 3 2 7 3 3
5 5 4 1 5 2
1 1
1 1
3 1 1 2
4 1 2 2 4
2 2 1
10 4 10 10 6 5 7 5 9 1 3
1 1
7 2 1 5 6 1 2 1
2 1 1
4 1 2 1 4
8 5 7 4 4 7 7 1 1
7 5 5 2 1 6 7 4
6 1 5 1 5 3 3
6 3 1 6 4 1 1
5 2 1 4 1 4
8 8 4 2 1 5 1 6 5
2 1 2
4 1 3 4 4
3 2 2 1
5 1 1 1 5 3
7 2 6 1 1 5 6 4
1 1
5 3 4 2 3 3
8 8 7 8 5 7 4 3 8
10 5 9 7 2 10 10 2 2 6 3
9 3 7 2 2 1 3 5 7 4
6 4 2 5 3 1 2
9 7 2 6 9 4 9 5 3 3
8 4 7 6 3 8 8 1 7
3 2 3 1
8 5 7 5 7 8 6 6 2
4 2 4 4 1
6 4 5 6 4 6 2
2 1 2
4 4 2 1 4
9 4 5 4 8 3 1 7 8 5
9 3 8 4 2 6 1 8 9 2
10 8 6 8 5 9 8 1 2 10 6
3 2 2 3
3 1 1 2
7 4 6 3 2 1 3 5
8 1 6 1 7 8 4 5 8
3 2 3 3
5 1 3 3 3 3
5 4 5 1 5 2
7 5 5 7 7 2 7 5
2 2 1
4 4 2 3 1
2 1 2
6 2 3 3 1 3 4
6 2 4 4 3 4 2
8 4 5 6 3 2 4 8 4
6 2 3 2 2 2 3
9 7 7 6 5 9 6 7 5 9
10 2 6 5 7 8 3 5 6 8 8
2 1 2
7 2 1 1 3 2 3 1
7 1 5 3 2 7 7 5
7 5 3 4 6 2 3 3
4 4 1 2 2
6 3 2 4 3 3 1
6 2 2 6 2 6 5
1 1
6 6 5 1 2 2 1
7 4 7 3 2 3 5 5
2 2 2
4 1 4 4 4
10 6 9 10 10 2 10 9 9 8 7
8 3 7 7 8 1 2 8 3
2 1 1`

func expected(p []int) string {
	n := len(p) - 1
	res := make([]string, n)
	for i := 1; i <= n; i++ {
		visited := make([]bool, n+1)
		x := p[i]
		visited[i] = true
		for !visited[x] {
			visited[x] = true
			x = p[x]
		}
		res[i-1] = strconv.Itoa(x)
	}
	return strings.Join(res, " ")
}

func runCase(exe string, p []int) error {
	n := len(p) - 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		sb.WriteString(strconv.Itoa(p[i]))
		if i < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	input := sb.String()
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := expected(p)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	var _ = solutionSource
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		p := make([]int, n+1)
		for i := 1; i <= n; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			p[i], _ = strconv.Atoi(scan.Text())
		}
		if err := runCase(exe, p); err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
