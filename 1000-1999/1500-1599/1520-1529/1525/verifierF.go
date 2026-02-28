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

func solveOptimal(n, m, l int, edges [][2]int, waves [][2]int64) (int64, error) {
	graph := make([][]bool, n)
	for i := 0; i < n; i++ {
		graph[i] = make([]bool, n)
	}
	for i := 0; i < m; i++ {
		a, b := edges[i][0], edges[i][1]
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

	cLen := 0
	for u := 0; u < n; u++ {
		used := make([]bool, n)
		if dfsMatch(u, used) {
			matchX[u] = true
			cLen++
		}
	}

	dp := make([][]int64, l+1)
	for i := 0; i <= l; i++ {
		dp[i] = make([]int64, n+1)
		for j := 0; j <= n; j++ {
			dp[i][j] = -1
		}
	}
	dp[0][cLen] = 0

	for i := 1; i <= l; i++ {
		x, y := waves[i-1][0], waves[i-1][1]
		for j := 0; j < n-i; j++ {
			for k := j; k <= n; k++ {
				prev := dp[i-1][k]
				if prev >= 0 {
					gain := x - int64(k-j)*y
					if gain < 0 {
						gain = 0
					}
					tot := prev + gain
					if tot > dp[i][j] {
						dp[i][j] = tot
					}
				}
			}
		}
	}

	var best int64 = -1
	for j := 0; j <= n; j++ {
		if dp[l][j] >= 0 && dp[l][j] > best {
			best = dp[l][j]
		}
	}
	if best == -1 {
		return -1, fmt.Errorf("no valid state")
	}

	return best, nil
}

func maxMatching(n int, edges [][2]int, blockedOut, blockedIn []bool) int {
	graph := make([][]bool, n)
	for i := 0; i < n; i++ {
		graph[i] = make([]bool, n)
	}
	for _, e := range edges {
		u, v := e[0]-1, e[1]-1
		if !blockedOut[u] && !blockedIn[v] {
			graph[u][v] = true
		}
	}

	matchY := make([]int, n)
	for i := range matchY {
		matchY[i] = -1
	}

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

	M := 0
	for u := 0; u < n; u++ {
		used := make([]bool, n)
		if dfsMatch(u, used) {
			M++
		}
	}
	return M
}

func evaluateGot(got string, n, m, l int, edges [][2]int, waves [][2]int64) (int64, error) {
	fields := strings.Fields(got)
	if len(fields) < 1 {
		return -1, fmt.Errorf("empty output")
	}
	var total int
	fmt.Sscanf(fields[0], "%d", &total)
	if len(fields) != total+1 {
		return -1, fmt.Errorf("expected %d elements, got %d", total, len(fields)-1)
	}

	seq := make([]int, total)
	for i := 0; i < total; i++ {
		fmt.Sscanf(fields[i+1], "%d", &seq[i])
	}

	blockedOut := make([]bool, n)
	blockedIn := make([]bool, n)

	var score int64 = 0
	waveIdx := 0
	opsBeforeWave := 0

	usedBlock := make(map[int]bool)

	for _, action := range seq {
		if action == 0 {
			waveIdx++
			if waveIdx > l {
				return -1, fmt.Errorf("called more than %d waves", l)
			}
			M := maxMatching(n, edges, blockedOut, blockedIn)
			if n-M <= waveIdx {
				return -1, fmt.Errorf("wave %d failed: M=%d, n=%d, n-M=%d <= %d", waveIdx, M, n, n-M, waveIdx)
			}
			x, y := waves[waveIdx-1][0], waves[waveIdx-1][1]
			gain := x - int64(opsBeforeWave)*y
			if gain > 0 {
				score += gain
			}
			opsBeforeWave = 0
		} else {
			if usedBlock[action] {
				return -1, fmt.Errorf("repeated block action %d", action)
			}
			usedBlock[action] = true
			if action > 0 {
				if action > n {
					return -1, fmt.Errorf("invalid action %d", action)
				}
				blockedOut[action-1] = true
			} else {
				if -action > n {
					return -1, fmt.Errorf("invalid action %d", action)
				}
				blockedIn[-action-1] = true
			}
			opsBeforeWave++
		}
	}
	if waveIdx != l {
		return -1, fmt.Errorf("called %d waves, expected %d", waveIdx, l)
	}
	return score, nil
}

func runCandidate(bin string, input string) (string, error) {
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	in := strings.NewReader(testcaseData)
	testIdx := 1
	for {
		var n, m, l int
		if _, err := fmt.Fscan(in, &n, &m, &l); err != nil {
			break
		}
		var inputBuf bytes.Buffer
		fmt.Fprintf(&inputBuf, "%d %d %d\n", n, m, l)

		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &edges[i][0], &edges[i][1])
			fmt.Fprintf(&inputBuf, "%d %d\n", edges[i][0], edges[i][1])
		}
		waves := make([][2]int64, l)
		for i := 0; i < l; i++ {
			fmt.Fscan(in, &waves[i][0], &waves[i][1])
			fmt.Fprintf(&inputBuf, "%d %d\n", waves[i][0], waves[i][1])
		}

		input := inputBuf.String()
		expectedScore, err := solveOptimal(n, m, l, edges, waves)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to solve embedded data for test %d: %v\n", testIdx, err)
			os.Exit(1)
		}

		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d candidate failed: %v\n", testIdx, err)
			os.Exit(1)
		}

		gotScore, err := evaluateGot(got, n, m, l, edges, waves)
		if err != nil {
			fmt.Printf("test %d output invalid: %v\n", testIdx, err)
			fmt.Printf("expected score: %d\ngot output:\n%s\n", expectedScore, got)
			os.Exit(1)
		}

		if gotScore != expectedScore {
			fmt.Printf("test %d output mismatch\n", testIdx)
			fmt.Printf("expected score: %d\n", expectedScore)
			fmt.Printf("got score: %d\n", gotScore)
			fmt.Printf("got output:\n%s\n", got)
			os.Exit(1)
		}

		testIdx++
	}

	fmt.Println("All tests passed")
}
