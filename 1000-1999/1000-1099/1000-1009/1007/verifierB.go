package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	maxN      = 100000
	testcases = `100
5 19 3
9 4 16
15 16 13
7 4 16
1 13 14
20 1 15
9 8 19
4 11 1
1 1 18
1 13 7
14 1 17
8 15 16
18 8 12
8 8 15
10 1 14
18 4 6
10 4 11
17 14 17
7 10 10
19 16 17
13 19 2
16 8 13
14 6 12
18 12 3
15 17 4
6 17 13
12 16 1
16 2 10
20 19 19
13 6 6
17 8 1
7 18 18
8 13 17
12 19 12
15 9 18
20 1 13
17 5 17
18 7 14
2 16 12
19 18 7
17 14 16
12 14 12
1 18 18
20 20 11
15 20 1
8 6 18
19 6 3
18 9 2
3 3 1
15 1 9
8 9 4
20 6 12
10 3 6
6 9 17
6 9 10
15 11 16
16 4 1
10 13 11
14 7 9
4 9 17
7 20 14
1 8 1
13 5 2
6 15 17
14 18 8
17 15 8
17 1 13
19 11 14
2 10 5
7 2 10
3 3 10
10 6 14
19 9 5
1 18 2
19 7 19
15 6 20
17 2 13
7 12 4
7 19 14
19 7 16
4 13 10
17 16 1
11 20 13
10 1 6
7 11 19
5 11 14
7 9 4
13 18 12
18 16 18
8 3 2
3 5 6
6 18 7
9 11 20
17 9 12
11 11 4
10 8 20
16 5 19
18 4 11
2 14 3
13 5 5`
)

var divisorCounts [maxN + 1]int

func init() {
	for i := 1; i <= maxN; i++ {
		for j := i; j <= maxN; j += i {
			divisorCounts[j]++
		}
	}
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func referenceSolve(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return "", fmt.Errorf("failed to read t: %w", err)
	}
	for i := 0; i < t; i++ {
		var a, b, c int
		if _, err := fmt.Fscan(reader, &a, &b, &c); err != nil {
			return "", fmt.Errorf("failed to read triple %d: %w", i+1, err)
		}
		x := gcd(gcd(a, b), c)
		m := gcd(a, b)
		n := gcd(b, c)
		p := gcd(c, a)
		var ans int64
		ans += int64(divisorCounts[a]) * int64(divisorCounts[b]) * int64(divisorCounts[c])
		tmp := int64(divisorCounts[a] - divisorCounts[m] - divisorCounts[p] + divisorCounts[x])
		ans -= tmp * int64(divisorCounts[n]) * int64(divisorCounts[n]-1) / 2
		tmp = int64(divisorCounts[b] - divisorCounts[m] - divisorCounts[n] + divisorCounts[x])
		ans -= tmp * int64(divisorCounts[p]) * int64(divisorCounts[p]-1) / 2
		tmp = int64(divisorCounts[c] - divisorCounts[n] - divisorCounts[p] + divisorCounts[x])
		ans -= tmp * int64(divisorCounts[m]) * int64(divisorCounts[m]-1) / 2
		ans -= 5 * int64(divisorCounts[x]) * int64(divisorCounts[x]-1) * int64(divisorCounts[x]-2) / 6
		ans -= 2 * int64(divisorCounts[x]) * int64(divisorCounts[x]-1)
		tmp = int64(divisorCounts[m] - divisorCounts[x])
		ans -= 3 * tmp * int64(divisorCounts[x]) * int64(divisorCounts[x]-1) / 2
		tmp = int64(divisorCounts[n] - divisorCounts[x])
		ans -= 3 * tmp * int64(divisorCounts[x]) * int64(divisorCounts[x]-1) / 2
		tmp = int64(divisorCounts[p] - divisorCounts[x])
		ans -= 3 * tmp * int64(divisorCounts[x]) * int64(divisorCounts[x]-1) / 2
		tmp = int64(divisorCounts[m] - divisorCounts[x])
		ans -= int64(divisorCounts[x]) * tmp
		tmp = int64(divisorCounts[n] - divisorCounts[x])
		ans -= int64(divisorCounts[x]) * tmp
		tmp = int64(divisorCounts[p] - divisorCounts[x])
		ans -= int64(divisorCounts[x]) * tmp
		tmp = int64(divisorCounts[m] - divisorCounts[x])
		ans -= int64(divisorCounts[x]) * tmp * (tmp - 1) / 2
		tmp = int64(divisorCounts[n] - divisorCounts[x])
		ans -= int64(divisorCounts[x]) * tmp * (tmp - 1) / 2
		tmp = int64(divisorCounts[p] - divisorCounts[x])
		ans -= int64(divisorCounts[x]) * tmp * (tmp - 1) / 2
		ans -= 2 * int64(divisorCounts[x]) * int64(divisorCounts[m]-divisorCounts[x]) * int64(divisorCounts[n]-divisorCounts[x])
		ans -= 2 * int64(divisorCounts[x]) * int64(divisorCounts[m]-divisorCounts[x]) * int64(divisorCounts[p]-divisorCounts[x])
		ans -= 2 * int64(divisorCounts[x]) * int64(divisorCounts[n]-divisorCounts[x]) * int64(divisorCounts[p]-divisorCounts[x])
		tmp = int64(divisorCounts[n] - divisorCounts[x] + divisorCounts[p] - divisorCounts[x])
		ans -= tmp * int64(divisorCounts[m]-divisorCounts[x]) * int64(divisorCounts[m]-divisorCounts[x]-1) / 2
		tmp = int64(divisorCounts[m] - divisorCounts[x] + divisorCounts[n] - divisorCounts[x])
		ans -= tmp * int64(divisorCounts[p]-divisorCounts[x]) * int64(divisorCounts[p]-divisorCounts[x]-1) / 2
		tmp = int64(divisorCounts[m] - divisorCounts[x] + divisorCounts[p] - divisorCounts[x])
		ans -= tmp * int64(divisorCounts[n]-divisorCounts[x]) * int64(divisorCounts[n]-divisorCounts[x]-1) / 2
		ans -= int64(divisorCounts[m]-divisorCounts[x]) * int64(divisorCounts[n]-divisorCounts[x]) * int64(divisorCounts[p]-divisorCounts[x])
		fmt.Fprintln(writer, ans)
	}
	if err := writer.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush output: %w", err)
	}
	return strings.TrimSpace(buf.String()), nil
}

func run(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	want, err := referenceSolve(testcases)
	if err != nil {
		fmt.Println("reference error:", err)
		os.Exit(1)
	}
	got, err := run(bin, []byte(testcases))
	if err != nil {
		fmt.Println("candidate runtime error:", err)
		os.Exit(1)
	}
	if want != got {
		fmt.Println("outputs differ")
		fmt.Println("expected:\n", want)
		fmt.Println("got:\n", got)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
