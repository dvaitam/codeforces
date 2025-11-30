package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded testcases from testcasesF.txt.
const testcaseData = `
4 5 2 1 4 2 3 1 3 2 1 4 2 32 49 70 14
3 0 1 24 99
5 2 1 1 5 5 3 41 26
3 1 1 1 2 47 54
3 1 2 3 2 78 76 1 77
4 0 2 24 62 61 91
3 0 2 52 3 71 54
4 3 3 1 3 2 4 1 4 60 45 66 46 68 33
5 1 3 5 1 90 70 12 40 88 41
4 1 1 2 3 11 77
5 0 2 64 43 27 17
3 3 1 3 2 1 3 1 2 54 38
3 3 1 2 1 3 1 3 2 36 38
5 6 2 3 2 5 3 1 2 5 4 3 4 5 1 47 72 91 86
4 6 2 1 3 2 1 3 2 4 2 4 3 4 1 37 65 83 87
4 2 2 4 2 2 1 43 67 19 68
3 1 2 1 2 93 86 94 54
3 3 2 1 2 1 3 2 3 76 57 80 84
3 1 1 3 1 29 22
2 0 1 41 24
5 3 1 1 3 5 4 1 4 45 52
4 6 2 2 3 3 1 4 3 4 1 4 2 1 2 54 53 73 83
4 3 1 1 2 3 4 1 4 45 85
2 1 1 2 1 73 79
3 2 2 1 2 3 2 38 62 91 49
5 2 3 2 4 5 2 4 16 85 80 57 32
4 0 1 32 61
2 0 1 35 20
4 2 2 3 1 4 1 69 91 44 79
4 5 3 2 4 3 4 3 1 4 1 2 3 88 70 25 2 55 100
5 9 4 5 1 1 2 3 2 4 1 3 1 5 3 2 4 4 5 4 3 89 71 64 43 69 20 55 75
2 0 1 35 11
2 0 1 93 55
2 1 1 1 2 16 16
3 0 1 20 24
3 3 1 1 3 2 3 2 1 96 10
4 0 3 2 56 12 41 88 77
5 7 3 3 2 3 5 1 5 5 2 3 1 3 4 4 5 48 59 95 50 23 17
2 1 1 1 2 12 16
5 3 4 3 1 1 4 3 5 11 86 63 28 18 90 80 49
4 6 1 1 2 3 2 3 1 4 2 4 3 1 4 76 74
3 1 2 1 3 5 24 26 26
2 1 1 2 1 1 55
5 4 4 4 1 5 1 2 5 3 2 15 36 21 97 37 87 31 5
5 0 3 86 42 38 59 87 31
3 0 1 74 30
2 1 1 2 1 24 60
4 2 1 2 4 1 4 83 78
5 4 1 3 2 1 4 3 5 1 3 63 12
5 0 2 40 82 7 5
3 1 2 1 3 41 14 12 70
3 0 1 57 70
4 3 2 4 3 2 1 1 4 71 15 18 40
5 1 3 5 4 74 86 32 36 8 64
3 0 1 13 10
5 5 4 1 4 1 5 5 2 5 4 3 4 49 72 17 42 9 40 22 18
4 3 3 1 3 2 4 2 3 17 46 56 34 79 3
3 3 1 3 2 2 1 1 3 21 75
3 0 2 81 91 19 58
5 6 1 1 4 2 5 2 1 5 3 2 3 1 5 17 54
2 0 1 5 74
3 3 2 1 2 2 3 3 1 26 39 72 67
4 6 1 3 2 2 1 4 2 1 3 3 4 1 4 63 49
2 0 1 51 13
4 1 2 1 4 14 66 78 76
3 3 1 3 2 2 1 3 1 70 91
5 7 1 2 4 3 4 1 5 3 1 2 1 3 2 4 5 36 92
3 0 1 45 89
3 3 1 3 1 3 2 2 1 1 86
5 6 2 1 3 5 2 3 2 4 3 1 2 1 4 10 81 53 99
3 1 1 2 1 30 55
5 8 2 2 4 3 5 5 1 1 3 4 1 3 4 2 1 4 5 83 89 43 94
2 0 1 18 67
2 1 1 2 1 83 50
3 3 1 1 3 2 3 1 2 43 20
4 4 1 3 1 3 4 1 2 1 4 27 42
4 1 2 4 2 1 16 77 48
3 0 2 96 10 63 98
4 0 1 19 38
5 10 1 5 3 2 4 4 3 2 5 1 2 4 1 1 5 3 2 3 1 4 5 70 50
3 0 2 11 1 99 19
4 6 1 2 4 3 4 1 3 3 2 1 4 2 1 67 11
4 4 3 3 1 4 1 3 4 1 2 4 54 99 93 53 56
5 0 2 51 54 25 35
3 3 2 2 3 1 2 3 1 10 10 88 49
2 1 1 1 2 83 3
5 5 4 3 1 2 4 5 3 2 5 3 2 43 100 28 20 38 67 90 72
2 1 1 1 2 83 52
4 2 2 3 1 1 4 70 2 94 50
3 1 2 1 3 2 99 28 54
3 1 2 1 3 89 50 45 100
2 1 1 2 1 14 84
2 1 1 2 1 49 39
3 0 2 68 80 50 13
4 6 3 1 2 4 3 2 4 1 3 3 2 1 4 8 95 12 15 44 19
4 2 2 4 2 2 1 4 81 33 50
3 1 2 1 2 77 18 73 7
4 4 1 2 4 1 3 4 1 3 4 53 93
5 9 2 4 3 4 1 1 5 2 1 2 3 5 2 4 2 5 4 3 1 39 39 1 98
4 4 3 3 1 4 1 1 2 4 2 12 20 46 11 82 55`

type pair struct {
	val   int64
	prevJ int
}

func solveInput(data string) (string, error) {
	in := strings.NewReader(data)
	out := &bytes.Buffer{}
	var n, m, l int
	if _, err := fmt.Fscan(in, &n, &m, &l); err != nil {
		return "", err
	}

	graph := make([][]bool, n)
	for i := 0; i < n; i++ {
		graph[i] = make([]bool, n)
	}
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		graph[a][b] = true
	}

	matchY := make([]int, n)
	for i := range matchY {
		matchY[i] = -1
	}
	matchX := make([]bool, n)

	var dfsMatch func(int, []bool) bool
	dfsMatch = func(u int, used []bool) bool {
		if used[u] {
			return false
		}
		used[u] = true
		for v := 0; v < n; v++ {
			if graph[u][v] {
				if matchY[v] < 0 || dfsMatch(matchY[v], used) {
					matchY[v] = u
					return true
				}
			}
		}
		return false
	}

	for u := 0; u < n; u++ {
		used := make([]bool, n)
		if dfsMatch(u, used) {
			matchX[u] = true
		}
	}

	visL := make([]bool, n)
	visR := make([]bool, n)
	var dfsCover func(int)
	dfsCover = func(u int) {
		if visL[u] {
			return
		}
		visL[u] = true
		for v := 0; v < n; v++ {
			if graph[u][v] && !visR[v] {
				visR[v] = true
				if matchY[v] >= 0 {
					dfsCover(matchY[v])
				}
			}
		}
	}
	for u := 0; u < n; u++ {
		if !matchX[u] {
			dfsCover(u)
		}
	}

	cuts := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if !visL[i] {
			cuts = append(cuts, i+1)
		}
		if visR[i] {
			cuts = append(cuts, -(i + 1))
		}
	}
	cLen := len(cuts)

	dp := make([][]pair, l+1)
	for i := 0; i <= l; i++ {
		dp[i] = make([]pair, n+1)
		for j := 0; j <= n; j++ {
			dp[i][j].val = -1
			dp[i][j].prevJ = -1
		}
	}
	dp[0][cLen].val = 0

	for i := 1; i <= l; i++ {
		var x, y int64
		fmt.Fscan(in, &x, &y)
		for j := 0; j < n-i; j++ {
			for k := j; k <= n; k++ {
				prev := dp[i-1][k].val
				if prev >= 0 {
					gain := x - int64(k-j)*y
					if gain < 0 {
						gain = 0
					}
					tot := prev + gain
					if tot > dp[i][j].val {
						dp[i][j].val = tot
						dp[i][j].prevJ = k
					}
				}
			}
		}
	}

	lastJ := -1
	best := pair{val: -1}
	for j := 0; j <= n; j++ {
		if dp[l][j].val >= 0 && dp[l][j].val > best.val {
			best.val = dp[l][j].val
			lastJ = j
		}
	}
	if lastJ == -1 {
		return "", fmt.Errorf("no valid state")
	}

	seq := make([]int, 0)
	for i := l; i >= 1; i-- {
		prevJ := dp[i][lastJ].prevJ
		seq = append(seq, 0)
		for lastJ < prevJ {
			seq = append(seq, cuts[lastJ])
			lastJ++
		}
		lastJ = prevJ
	}
	total := len(seq)
	fmt.Fprintln(out, total)
	for i := total - 1; i >= 0; i-- {
		if i > 0 {
			fmt.Fprint(out, seq[i], " ")
		} else {
			fmt.Fprintln(out, seq[i])
		}
	}
	return strings.TrimSpace(out.String()), nil
}

func runCandidate(bin string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(testcaseData)
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	expect, err := solveInput(testcaseData)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to solve embedded data:", err)
		os.Exit(1)
	}

	got, err := runCandidate(bin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if got != expect {
		fmt.Println("output mismatch")
		fmt.Println("expected:")
		fmt.Println(expect)
		fmt.Println("got:")
		fmt.Println(got)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
