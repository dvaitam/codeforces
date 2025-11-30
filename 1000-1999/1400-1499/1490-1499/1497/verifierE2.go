package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE2.txt so the verifier is self-contained.
const testcasesRaw = `2 3 98 34
1 0 19
8 2 41 99 3 35 63 26 94 53
2 1 73 71
5 0 55 43 12 47 53
5 3 90 13 97 26 90
5 0 6 76 26 84 47
8 1 66 74 83 90 65 4 82 47
4 3 39 46 76 16
2 1 15 78
5 2 93 26 49 62 29
3 1 90 67 2
4 1 3 83 43 72
5 2 49 68 50 38 17
8 0 24 55 77 96 51 13 57 32
2 3 58 49
2 3 61 39
7 0 25 96 88 35 58 63 93
3 0 3 69 16
5 2 25 33 65 58 43
5 3 54 79 63 35 78
8 3 63 19 93 49 64 40 82 60
6 2 85 21 80 49 89 78
5 2 83 51 63 94 21
5 0 80 58 8 24 4
2 2 47 64
1 1 20
5 0 55 68 63 10 61
4 0 48 47 19 88
4 2 19 5 81 86
2 0 6 83
8 3 93 9 81 4 90 71 17 79
2 1 80 81
5 3 64 89 3 82 18
4 3 62 100 57 58
5 3 92 75 9 38 47
5 2 5 11 65 36 36
5 3 56 57 48 5 92
2 1 78 82
4 0 95 23 44 48
1 0 95
4 3 44 10 21 94
2 3 71 18
8 3 82 10 17 10 45 67 4 24
3 1 37 95 20
7 1 29 90 73 92 17 54 89
3 2 31 72 74
1 3 25
6 3 68 57 2 51 68 71
3 1 83 70 96
6 2 86 24 20 22 89 81
3 3 10 8 2
1 3 98
1 0 84
2 2 77 5
1 1 53
4 2 12 61 10 93
1 1 85
8 0 15 90 23 7 88 4 97 87
7 2 33 57 30 26 84 100 69
4 2 76 2 33 68
6 0 75 77 83 5 47 47
2 1 98 4
6 2 56 88 7 45 18 4
1 2 9
8 0 60 34 12 87 76 79 35 30
3 1 54 70 24
4 3 77 55 73 67
8 2 33 18 39 5 25 98 58 14
3 0 1 66 21
2 1 76 51
8 3 68 99 89 70 61 28 83 50
3 3 2 44 54
5 2 18 37 48 76 87
6 1 1 100 40 23 11 53
4 0 25 49 53 17
3 0 67 68 100
8 2 43 92 76 50 35 84 85 9
5 1 6 15 76 78 79
3 0 89 98 100
2 0 86 97
2 1 45 44
3 1 23 24 29
3 3 95 59 80
2 1 59 26
6 1 43 81 34 9 36 97
1 0 84
1 3 81
8 2 8 2 36 93 13 50 80 44
2 2 19 97
3 2 86 98 61
7 0 69 20 14 58 48 2 17
8 2 10 62 16 52 91 26 70 22
8 2 33 70 100 45 38 70 89 35
8 0 65 96 32 85 80 1 68 64
5 1 5 42 2 27 74
5 1 20 36 8 35 94
5 2 28 91 94 22 77
1 0 46
3 2 24 62 82
`

type testCase struct {
	n   int
	k   int
	arr []int
}

func buildPrimes(limit int) []int {
	isComp := make([]bool, limit+1)
	primes := []int{}
	for i := 2; i <= limit; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			if i*i <= limit {
				for j := i * i; j <= limit; j += i {
					isComp[j] = true
				}
			}
		}
	}
	return primes
}

func canonical(x int, primes []int) int {
	res := 1
	for _, p := range primes {
		if p*p > x {
			break
		}
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt ^= 1
		}
		if cnt == 1 {
			res *= p
		}
	}
	if x > 1 {
		res *= x
	}
	return res
}

func solveCase(tc testCase, primes []int) int {
	a := make([]int, tc.n)
	for i, v := range tc.arr {
		a[i] = canonical(v, primes)
	}
	idMap := make(map[int]int)
	ids := make([]int, tc.n)
	idCnt := 0
	for i, v := range a {
		if idx, ok := idMap[v]; ok {
			ids[i] = idx
		} else {
			idMap[v] = idCnt
			ids[i] = idCnt
			idCnt++
		}
	}
	m := idCnt
	k := tc.k
	n := tc.n
	nxt := make([][]int, n)
	for i := range nxt {
		nxt[i] = make([]int, k+1)
	}
	freq := make([][]int, k+1)
	for i := 0; i <= k; i++ {
		freq[i] = make([]int, m)
	}
	r := make([]int, k+1)
	dup := make([]int, k+1)
	for i := 0; i < n; i++ {
		val := ids[i]
		for j := 0; j <= k; j++ {
			for r[j] < n {
				id := ids[r[j]]
				extra := 0
				if freq[j][id] > 0 {
					extra = 1
				}
				if dup[j]+extra > j {
					break
				}
				dup[j] += extra
				freq[j][id]++
				r[j]++
			}
			nxt[i][j] = r[j]
			if freq[j][val] > 1 {
				dup[j]--
			}
			freq[j][val]--
		}
	}
	const INF = int(1e9)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, k+1)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	dp[0][0] = 0
	for i := 0; i < n; i++ {
		for used := 0; used <= k; used++ {
			if dp[i][used] == INF {
				continue
			}
			for add := 0; add <= k-used; add++ {
				to := nxt[i][add]
				if dp[to][used+add] > dp[i][used]+1 {
					dp[to][used+add] = dp[i][used] + 1
				}
			}
		}
	}
	ans := INF
	for j := 0; j <= k; j++ {
		if dp[n][j] < ans {
			ans = dp[n][j]
		}
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var res []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		k, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse k: %v", idx+1, err)
		}
		if len(parts) != n+2 {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", idx+1, n+2, len(parts))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse value %d: %v", idx+1, i+1, err)
			}
			arr[i] = v
		}
		res = append(res, testCase{n: n, k: k, arr: arr})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return res, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	primes := buildPrimes(3200)

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d %d\n", tc.n, tc.k)
		for idx, v := range tc.arr {
			if idx > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		expected := strconv.Itoa(solveCase(tc, primes))
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
