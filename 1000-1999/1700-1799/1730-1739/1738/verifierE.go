package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded test cases from testcasesE.txt. Each line: n followed by n integers.
const testcasesRaw = `4 4 1 6 7
3 1 1 0
7 8 4 0 3 8 8 5
5 2 1 4 3 0
5 4 3 2 4 4
6 1 9 5 10 6 8
4 2 3 7 4
2 8 4
1 4
5 8 3 6 6 9
5 6 7 2 3 4
5 0 1 0 7 10
5 8 8 10 7 5
3 10 3 1
7 3 10 10 7 4 2 5
7 9 5 10 8 3 5 1
1 3
5 9 9 3 1 5
3 4 7 0
1 5
2 4 10
6 0 5 4 5 2 10
7 9 10 1 4 9 3 7
5 2 4 6 9 2
6 9 0 5 0 7 2
6 5 4 9 1 7 3
7 3 1 0 0 0 2 9
3 9 0 8
8 9 3 5 0 1 8 4 6
4 7 3 3 7
7 7 0 3 6 7 3 10
7 3 7 3 0 0 4 4
4 8 3 3 6
5 2 5 0 5 9
2 9 6
1 7
7 1 6 3 9 2 5 4
8 10 5 6 8 3 10 10 4
6 6 7 1 4 10 10
4 0 6 9 2
5 10 0 2 10 7
8 6 6 3 0 3 2 0 9
5 1 6 6 3 8
1 3
3 10 9 5
8 8 7 0 1 0 9 1 7
5 9 2 0 5 1
1 4
6 1 1 8 7 6 3
5 6 3 7 6 1
2 1 9
6 8 6 6 7 1 10
4 10 4 7 6
2 8 2
6 2 2 2 5 7 5
5 8 0 2 0 10
5 1 8 1 7 9
8 8 1 8 3 6 4 5 3
3 10 0 10
1 9
6 8 7 9 4 8 7
8 6 2 4 9 5 10 5 2
7 1 9 2 10 9 2 4
6 3 9 5 10 9 1
2 6 10
3 5 10 5
6 2 4 0 9 0 8
2 5 1
3 2 9 7
2 1 2
8 10 3 9 10 4 6 9 3
8 3 4 5 3 5 8 10 8
8 6 8 6 5 4 7 6 9
1 4
3 8 7 8
6 6 6 9 0 2 8
2 7 10
6 9 4 1 4 7 3
8 9 1 2 3 1 4 2 0
3 6 9 10
5 0 8 3 2 1
4 0 5 1 1
5 0 10 4 9 10
3 2 10 6
3 1 8 5
1 8
3 8 6 2
4 4 7 8 1
7 2 2 4 8 6 8 10
5 6 5 2 6 0
7 0 4 10 0 4 2 1
3 1 9 0
4 3 8 0 7
3 7 10 6
6 2 8 9 3 1 7
6 4 10 6 9 0 10
4 4 9 10 9
5 7 10 5 1 3
6 1 10 0 9 5 8
8 5 1 2 0 7 8 8 9`

type testCase struct {
	n   int
	arr []int64
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d parse n: %v", idx+1, err)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d expected %d numbers got %d", idx+1, n, len(fields)-1)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[1+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d parse arr[%d]: %v", idx+1, i, err)
			}
			arr[i] = v
		}
		cases = append(cases, testCase{n: n, arr: arr})
	}
	return cases, nil
}

const MOD int64 = 998244353
const MAXN = 100000

var fact [MAXN + 1]int64
var invFact [MAXN + 1]int64

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= MOD
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func init() {
	fact[0] = 1
	for i := 1; i <= MAXN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[MAXN] = modPow(fact[MAXN], MOD-2)
	for i := MAXN; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
}

func comb(n, k int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
}

// Embedded solver logic from 1738E.go.
func solve(tc testCase) int64 {
	prefix := make([]int64, tc.n+1)
	for i := 0; i < tc.n; i++ {
		prefix[i+1] = prefix[i] + tc.arr[i]
	}
	total := prefix[tc.n]
	if total == 0 {
		return modPow(2, int64(tc.n-1))
	}

	cnt := make(map[int64]int)
	for i := 1; i < tc.n; i++ {
		cnt[prefix[i]]++
	}

	ans := int64(1)
	for s, l := range cnt {
		if l == 0 || 2*s >= total {
			continue
		}
		r := cnt[total-s]
		ans = ans * comb(int64(l+r), int64(l)) % MOD
		cnt[s] = 0
		cnt[total-s] = 0
	}
	if total%2 == 0 {
		c := cnt[total/2]
		if c > 0 {
			ans = ans * modPow(2, int64(c)) % MOD
		}
	}
	return ans % MOD
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		want := solve(tc)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(strconv.Itoa(tc.n))
		input.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatInt(v, 10))
		}
		input.WriteByte('\n')

		gotStr, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(gotStr), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d output parse error: %v (got %q)\n", idx+1, err, gotStr)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
