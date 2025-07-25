package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type CaseD struct {
	input    string
	expected float64
}

func solveD(n int, mat [][]int) float64 {
	size := 1 << n
	prob := make([][]float64, size)
	for i := range prob {
		prob[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			prob[i][j] = float64(mat[i][j]) / 100.0
		}
	}
	v := make([][][2]int, n+1)
	dp := make([][][2]float64, n+1)
	used := make([][]bool, n+1)
	for lvl := 0; lvl <= n; lvl++ {
		v[lvl] = make([][2]int, size)
		dp[lvl] = make([][2]float64, size)
		used[lvl] = make([]bool, size)
	}

	var buildSegments func(level, l, r int)
	buildSegments = func(level, l, r int) {
		if level < 0 {
			return
		}
		for i := l; i <= r; i++ {
			v[level][i][0] = l
			v[level][i][1] = r
		}
		if level == 0 {
			return
		}
		mid := (l + r) / 2
		buildSegments(level-1, l, mid)
		buildSegments(level-1, mid+1, r)
	}

	var getDp func(x, b int) (float64, float64)
	getDp = func(x, b int) (float64, float64) {
		if used[x][b] {
			return dp[x][b][0], dp[x][b][1]
		}
		if x == 0 {
			used[x][b] = true
			dp[x][b][0] = 0
			dp[x][b][1] = 1
			return 0, 1
		}
		pF, pS := getDp(x-1, b)
		var bestF, sumS float64
		l, r := v[x][b][0], v[x][b][1]
		for i := l; i <= r; i++ {
			if v[x-1][b][0] <= i && i <= v[x-1][b][1] {
				continue
			}
			tmpF, tmpS := getDp(x-1, i)
			sumS += tmpS * prob[b][i]
			if tmpF+pF > bestF {
				bestF = tmpF + pF
			}
		}
		sumS *= pS
		score := float64(int64(1) << uint(x-1))
		bestF += sumS * score
		used[x][b] = true
		dp[x][b][0] = bestF
		dp[x][b][1] = sumS
		return bestF, sumS
	}

	buildSegments(n, 0, size-1)
	ans := 0.0
	for i := 0; i < size; i++ {
		f, _ := getDp(n, i)
		if f > ans {
			ans = f
		}
	}
	return ans
}

func generateCaseD(rng *rand.Rand) CaseD {
	n := rng.Intn(5) + 2 // 2..6
	size := 1 << n
	mat := make([][]int, size)
	for i := range mat {
		mat[i] = make([]int, size)
		for j := 0; j < size; j++ {
			mat[i][j] = rng.Intn(101)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", mat[i][j]))
		}
		sb.WriteByte('\n')
	}
	return CaseD{sb.String(), solveD(n, mat)}
}

func runCase(exe, input string, expected float64) error {
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
	var got float64
	if _, err := fmt.Fscan(strings.NewReader(strings.TrimSpace(out.String())), &got); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	if math.Abs(got-expected) > 1e-6 {
		return fmt.Errorf("expected %.6f got %.6f", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	// deterministic simple case
	mat := [][]int{{0, 50, 50, 50}, {50, 0, 50, 50}, {50, 50, 0, 50}, {50, 50, 50, 0}}
	n := 2
	var sb strings.Builder
	sb.WriteString("2\n")
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", mat[i][j]))
		}
		sb.WriteByte('\n')
	}
	cases := []CaseD{{sb.String(), solveD(n, mat)}}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseD(rng))
	}
	for i, tc := range cases {
		if err := runCase(exe, tc.input, tc.expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			fmt.Fprint(os.Stderr, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
