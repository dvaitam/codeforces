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

type caseF struct {
	n, m   int
	pos    []int
	trucks [][4]int
}

const testcasesRaw = `100
3 2
10 12 18
2 3 3 0
2 3 4 1
2 3
8 11
1 2 1 2
1 2 2 1
1 2 2 2
2 1
1 6
1 2 2 1
3 2
2 10 18
1 3 2 0
1 3 3 2
3 2
2 9 13
2 3 1 0
2 3 4 2
3 3
3 7 13
2 3 1 0
2 3 3 1
1 2 1 2
2 1
6 20
1 2 5 0
2 3
5 18
1 2 1 2
1 2 1 1
1 2 1 1
4 2
4 6 17 18
1 3 4 2
1 2 5 0
4 3
4 8 12 13
3 4 5 2
3 4 3 0
3 4 2 2
3 2
10 12 17
2 3 1 0
1 3 2 2
2 1
4 16
1 2 3 1
2 1
6 12
1 2 1 2
3 1
7 9 11
1 2 3 1
2 2
13 18
1 2 3 0
1 2 4 2
3 1
2 9 20
1 3 4 0
4 2
6 8 9 12
3 4 2 1
1 4 5 2
2 2
2 10
1 2 1 2
1 2 2 1
2 3
7 18
1 2 4 2
1 2 3 1
1 2 4 1
2 2
8 13
1 2 5 2
1 2 4 0
2 3
13 14
1 2 3 0
1 2 3 2
1 2 2 1
2 2
9 17
1 2 5 2
1 2 4 0
3 3
9 11 19
2 3 5 2
1 3 5 2
1 3 3 2
4 2
1 2 3 14
2 3 1 1
3 4 5 1
4 2
15 16 17 18
3 4 5 0
3 4 2 1
3 1
4 6 8
2 3 4 1
3 3
5 9 16
2 3 5 2
1 2 3 2
2 3 5 2
4 1
5 10 17 19
1 2 3 0
2 1
16 20
1 2 2 2
3 3
5 9 19
1 3 3 0
2 3 1 1
1 2 5 2
2 2
10 19
1 2 1 2
1 2 4 2
2 1
9 14
1 2 4 1
2 2
7 14
1 2 1 0
1 2 3 1
3 2
6 8 15
1 3 2 1
1 2 5 2
4 2
3 6 11 20
1 3 3 2
3 4 1 1
3 2
1 10 11
1 3 2 0
1 3 2 0
4 3
2 4 9 13
3 4 4 1
2 3 2 0
1 2 1 1
4 1
3 4 7 10
3 4 4 0
4 2
6 9 11 13
2 3 2 2
1 2 5 1
4 3
3 11 15 17
2 4 4 0
1 3 2 1
1 2 1 0
3 2
4 8 12
1 3 2 0
1 3 4 1
4 3
9 10 16 17
1 4 5 0
2 3 5 2
3 4 3 1
2 1
2 15
1 2 3 2
3 1
5 11 16
1 3 2 0
2 3
1 5
1 2 3 1
1 2 1 0
1 2 4 2
4 3
7 12 13 14
3 4 5 2
1 4 1 1
3 4 4 0
2 1
5 19
1 2 3 1
3 3
11 14 16
2 3 4 2
2 3 3 1
2 3 2 2
3 2
5 13 14
1 3 2 0
1 2 3 2
4 1
4 6 16 19
2 3 4 0
2 1
2 17
1 2 5 0
2 3
7 12
1 2 1 1
1 2 3 2
1 2 5 2
3 3
4 5 10
2 3 2 2
1 2 3 1
2 3 4 2
3 3
3 4 13
1 2 2 1
2 3 3 1
1 3 1 1
4 2
2 8 10 13
2 4 2 0
3 4 4 0
4 3
5 6 11 18
3 4 1 1
2 4 3 0
2 3 4 0
2 2
10 11
1 2 4 0
1 2 4 0
3 1
3 6 19
2 3 5 2
2 2
5 11
1 2 1 0
1 2 2 1
4 3
1 7 11 13
2 4 5 1
2 3 3 1
2 3 4 0
4 3
7 11 12 13
2 3 2 2
3 4 2 0
2 4 2 1
3 3
6 13 20
1 2 2 1
2 3 1 2
1 2 2 1
4 2
4 6 18 20
3 4 4 1
1 4 1 1
2 1
2 7
1 2 2 2
3 3
3 5 19
2 3 5 1
2 3 5 1
1 3 3 0
3 2
4 10 12
1 3 1 2
1 3 1 1
2 2
15 19
1 2 2 2
1 2 1 2
2 1
5 13
1 2 2 2
4 1
1 2 13 18
2 3 4 0
4 2
7 14 17 19
2 3 1 1
2 3 2 1
3 3
6 7 19
1 2 2 0
2 3 3 2
2 3 2 2
4 2
2 4 9 13
3 4 2 1
1 3 4 0
4 3
2 12 15 17
1 2 5 1
2 3 5 2
1 4 1 0
2 2
1 11
1 2 5 2
1 2 3 0
3 1
2 7 17
2 3 4 0
2 3
7 9
1 2 3 1
1 2 1 1
1 2 3 0
3 1
7 8 19
2 3 5 1
4 1
1 2 12 13
3 4 1 1
4 1
3 4 10 17
3 4 5 1
2 1
4 10
1 2 4 1
2 3
8 19
1 2 2 1
1 2 3 0
1 2 3 2
2 1
13 15
1 2 5 1
3 1
2 10 14
1 3 5 0
2 1
8 13
1 2 4 2
4 2
1 2 13 18
3 4 5 1
1 2 2 2
4 1
4 5 10 18
2 4 5 0
3 1
6 11 20
1 2 3 0
2 3
8 11
1 2 4 2
1 2 3 2
1 2 1 2
3 2
8 11 17
2 3 1 1
1 2 1 1
4 2
1 9 11 15
2 4 2 1
3 4 3 2
3 1
1 13 18
2 3 5 0
4 1
1 5 11 14
1 2 1 2
4 2
3 7 8 12
1 4 4 1
3 4 1 0
3 2
6 7 13
2 3 2 0
1 2 1 1
2 3
8 20
1 2 3 2
1 2 4 2
1 2 2 2
2 1
4 5
1 2 5 2
2 3
14 20
1 2 4 0
1 2 2 2
1 2 4 2
4 1
8 13 15 19
2 4 1 1
4 3
9 15 17 20
2 3 5 1
2 3 4 1
3 4 1 2
3 3
4 7 18
1 2 1 1
2 3 5 2
1 3 2 1`

var testcases = mustParseTestcases(testcasesRaw)

func mustParseTestcases(raw string) []caseF {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(raw)))
	scanner.Split(bufio.ScanWords)

	nextInt := func() int {
		if !scanner.Scan() {
			panic("unexpected EOF while reading testcases")
		}
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(fmt.Sprintf("invalid integer %q: %v", scanner.Text(), err))
		}
		return v
	}

	t := nextInt()
	cases := make([]caseF, 0, t)
	for i := 0; i < t; i++ {
		n := nextInt()
		m := nextInt()
		pos := make([]int, n)
		for j := 0; j < n; j++ {
			pos[j] = nextInt()
		}
		trucks := make([][4]int, m)
		for j := 0; j < m; j++ {
			s := nextInt()
			f := nextInt()
			c := nextInt()
			r := nextInt()
			trucks[j] = [4]int{s, f, c, r}
		}
		cases = append(cases, caseF{n: n, m: m, pos: pos, trucks: trucks})
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("scanner error: %v", err))
	}
	if len(cases) == 0 {
		panic("no testcases parsed")
	}
	return cases
}

// solveCase embeds the reference logic from 1101F.go.
func solveCase(cf caseF) int64 {
	n := cf.n
	A := cf.pos
	N := n
	NN := N * N
	size := NN * N
	dp := make([]int32, size)
	const inf32 = int32(1000000005)
	for i := 0; i < N; i++ {
		iNN := i * NN
		for j := i; j < N; j++ {
			ij := iNN + j*N
			dp[ij] = int32(A[j] - A[i])
			s := i
			for k := 1; k <= j-i; k++ {
				for s < j-1 && dp[iNN+s*N+(k-1)] < int32(A[j]-A[s]) {
					s++
				}
				v := inf32
				if s != i {
					v1 := dp[iNN+(s-1)*N+(k-1)]
					d1 := int32(A[j] - A[s-1])
					if v1 > d1 {
						v = v1
					} else {
						v = d1
					}
				}
				v2 := dp[iNN+s*N+(k-1)]
				d2 := int32(A[j] - A[s])
				t32 := v2
				if d2 > t32 {
					t32 = d2
				}
				if v > t32 {
					v = t32
				}
				dp[ij+k] = v
			}
		}
	}
	var ans int64
	for _, tr := range cf.trucks {
		s := tr[0] - 1
		t := tr[1] - 1
		c := tr[2]
		r := tr[3]
		maxRef := t - s
		if r > maxRef {
			r = maxRef
		}
		val := int64(dp[s*NN+t*N+r]) * int64(c)
		if val > ans {
			ans = val
		}
	}
	return ans
}

func runCandidate(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseCandidateOutput(out string) (int64, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return 0, fmt.Errorf("no output")
	}
	v, err := strconv.ParseInt(scanner.Text(), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse output: %v", err)
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scanner error: %v", err)
	}
	return v, nil
}

func checkCase(bin string, idx int, cs caseF) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", cs.n, cs.m))
	for i, v := range cs.pos {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, tr := range cs.trucks {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tr[0], tr[1], tr[2], tr[3]))
	}
	input := sb.String()

	expected := solveCase(cs)
	out, err := runCandidate(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\ninput:\n%s", err, input)
	}
	got, err := parseCandidateOutput(out)
	if err != nil {
		return fmt.Errorf("output parse error: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d\ninput:\n%s", expected, got, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, cs := range testcases {
		if err := checkCase(bin, i, cs); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
