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

// Embedded source for the reference solution (was 1092D2.go).
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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n+2)
	// Sentinel large value to terminate stack walk.
	a[0] = 1 << 30
	maxv := 0
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] > maxv {
			maxv = a[i]
		}
	}
	n++
	a[n] = maxv
	fa := make([]int, n+1)
	for i := 1; i <= n; i++ {
		j := i - 1
		for a[i] > a[j] {
			if (i-fa[j])%2 == 0 {
				fmt.Fprintln(writer, "NO")
				return
			}
			j = fa[j]
		}
		if a[i] == a[j] {
			fa[i] = fa[j]
		} else {
			fa[i] = j
		}
	}
	fmt.Fprintln(writer, "YES")
}
`

const testcasesRaw = `4 5 2 7 8
3 2 2 1
7 9 5 1 4 9 9 6
5 3 2 5 4 1
5 5 4 3 5 5
6 2 10 6 7 9 4
3 4 8 5
2 9 5
1 5
10 5 9 4 7 7 10 5 7 8 3
4 5 5 1 2
1 8
5 9 9 8 6 3
4 2 7 4 8
5 3 6 7 10 6
9 4 6 2 1 4 5 10 10 4
2 6 3
5 8 1 1 6 2
5 6 1 6 5 6
3 7 10 2
5 10 4 8 5 3
5 7 10 3 6 10
1 6
1 8
3 6 6 5
10 2 8 4 7 4 2 1 1 1 3
10 3 10 1 9 8 10 4 6 1 2
9 5 7 4 8 4 4 8 7 8
1 4
7 8 4 7 4 8 4 1
1 5
5 4 9 4 4 7
5 3 6 1 6 10
2 10 7
1 8
7 2 7 4 10 3 6 5
8 6 7 9 4 5 6 7 8
2 5 4
1 7
10 3 5 1 3 8 10 8 7 7 4
1 4
3 1 10 5
2 7 7
4 9 1 4 3
10 6 9 8 9 8 1 2 1 10 2
8 9 5 10 3 1 6 2 9
1 5
6 2 2 9 8 7 4
5 7 4 8 7 2
2 2 10
6 9 7 7 8 2 4
5 8 7 2 9 3
6 3 3 3 6 8 6
5 9 1 3 1 5
2 9 2
8 10 8 9 2 9 4 7 5
6 4 3 1 1 10 6
9 8 10 5 9 8 10 10 8 7
3 5 10 6
6 3 7 2 10 3 10
3 5 6 4
10 6 10 2 2 7 3 6 6 6 3
5 1 10 1 9 2
6 2 3 3 10 8 10
2 2 3
8 4 10 5 7 10 4 8 4
5 6 4 6 9 9
8 7 9 7 6 5 8 7 10
1 5
3 9 8 9
10 6 7 7 10 1 3 9 2 8 6
10 5 2 5 8 4 8 10 2 3 4
2 5 3
1 3
7 10 10 9 9 5 1 9
4 3 2 4 1
6 2 2 5 1 5 10
3 3 7 3
2 9 6
1 9
3 9 7 3
4 5 8 9 2
7 3 3 5 9 7 9 5
7 6 3 7 1 7 1 5
1 5
3 2 3 2
10 1 4 4 9 1 8 9 3 8 7
6 3 9 10 4 2 8
10 6 5 7 10 1 4 5 10 10 5
8 6 2 4 6 2 1 10 6
9 8 6 2 3 1 8 9 9 10
4 1 4 2 6
10 3 5 2 9 9 1 2 4 10 5
10 1 1 1 8 3 4 6 4 6 5
7 10 7 1 3 7 8 2
9 9 5 9 5 4 7 6 3 4
1 8
7 5 6 7 7 6 10 4
5 5 9 2 1 7
5 10 5 3 10 6`

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func expected(arr []int) string {
	n := len(arr)
	maxv := 0
	for _, v := range arr {
		if v > maxv {
			maxv = v
		}
	}
	a := make([]int, n+2)
	a[0] = 1 << 30
	copy(a[1:], arr)
	n++
	a[n] = maxv
	fa := make([]int, n+1)
	for i := 1; i <= n; i++ {
		j := i - 1
		for a[i] > a[j] {
			if (i-fa[j])%2 == 0 {
				return "NO"
			}
			j = fa[j]
		}
		if a[i] == a[j] {
			fa[i] = fa[j]
		} else {
			fa[i] = j
		}
	}
	return "YES"
}

func main() {
	var _ = solutionSource
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil || len(fields) != n+1 {
			fmt.Fprintf(os.Stderr, "test %d: invalid length\n", idx)
			os.Exit(1)
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i], _ = strconv.Atoi(fields[i+1])
		}
		input := line + "\n"
		want := expected(a)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx, input, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
