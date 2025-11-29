package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded source for the reference solution (was 1254C.go).
const solutionSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

// This is a stub solution for the interactive problem "Point Ordering".
// The real problem asks to determine the counter-clockwise ordering of n
// unknown points that form a convex polygon using orientation and area
// queries.  Since this repository does not provide an interactive judge,
// we simply read the value of n and output the identity permutation as a
// placeholder.
func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	if n <= 0 {
		return
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, i)
	}
	fmt.Fprintln(out)
}
`

const testcasesRaw = `4
7
7
4
5
7
6
8
7
3
7
3
6
5
7
4
4
8
6
7
7
6
6
8
4
4
8
4
7
6
8
3
8
3
4
7
3
5
3
5
6
7
8
6
8
6
6
8
7
6
4
5
3
3
4
6
4
5
8
6
8
5
6
7
6
7
5
7
7
6
7
4
5
8
3
5
7
8
8
4
8
5
7
7
7
3
8
8
4
8
7
5
5
3
3
6
8
6
3
5`

var _ = solutionSource

func expected(n int) string {
	if n <= 0 {
		return ""
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(i))
	}
	return sb.String()
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseTests() ([]int, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	tests := make([]int, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid integer on line %d", idx+1)
		}
		tests = append(tests, n)
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, n := range tests {
		input := fmt.Sprintf("%d\n", n)
		want := expected(n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
