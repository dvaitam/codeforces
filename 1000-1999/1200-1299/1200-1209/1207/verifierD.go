package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const solution1207DSource = `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	pairs := make([][2]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pairs[i][0], &pairs[i][1])
	}

	fac := make([]int64, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}

	tmp := make([][2]int, n)
	copy(tmp, pairs)
	sort.Slice(tmp, func(i, j int) bool { return tmp[i][0] < tmp[j][0] })
	cntA := int64(1)
	for i := 0; i < n; {
		j := i
		for j < n && tmp[j][0] == tmp[i][0] {
			j++
		}
		cntA = cntA * fac[j-i] % mod
		i = j
	}

	copy(tmp, pairs)
	sort.Slice(tmp, func(i, j int) bool { return tmp[i][1] < tmp[j][1] })
	cntB := int64(1)
	for i := 0; i < n; {
		j := i
		for j < n && tmp[j][1] == tmp[i][1] {
			j++
		}
		cntB = cntB * fac[j-i] % mod
		i = j
	}

	copy(tmp, pairs)
	sort.Slice(tmp, func(i, j int) bool {
		if tmp[i][0] == tmp[j][0] {
			return tmp[i][1] < tmp[j][1]
		}
		return tmp[i][0] < tmp[j][0]
	})
	valid := true
	for i := 1; i < n; i++ {
		if tmp[i][1] < tmp[i-1][1] {
			valid = false
			break
		}
	}
	cntAB := int64(0)
	if valid {
		cntAB = 1
		for i := 0; i < n; {
			j := i
			for j < n && tmp[j] == tmp[i] {
				j++
			}
			cntAB = cntAB * fac[j-i] % mod
			i = j
		}
	}

	ans := (fac[n] - cntA - cntB + cntAB) % mod
	if ans < 0 {
		ans += mod
	}
	fmt.Fprintln(out, ans)
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1207DSource

const (
	mod = 998244353
	inf = int64(4e18)
)

type pair struct {
	a int
	b int
}

type testCase struct {
	n     int
	pairs []pair
}

var testcases = []testCase{
	{n: 2, pairs: []pair{{a: 2, b: 1}, {a: 3, b: 2}}},
	{n: 4, pairs: []pair{{a: 1, b: 1}, {a: 1, b: 1}, {a: 2, b: 3}, {a: 2, b: 1}}},
	{n: 2, pairs: []pair{{a: 3, b: 3}, {a: 2, b: 2}}},
	{n: 2, pairs: []pair{{a: 1, b: 2}, {a: 1, b: 1}}},
	{n: 6, pairs: []pair{{a: 2, b: 2}, {a: 1, b: 1}, {a: 2, b: 2}, {a: 3, b: 3}, {a: 2, b: 1}, {a: 3, b: 2}}},
	{n: 6, pairs: []pair{{a: 2, b: 3}, {a: 1, b: 1}, {a: 1, b: 2}, {a: 2, b: 1}, {a: 3, b: 2}, {a: 1, b: 2}}},
	{n: 5, pairs: []pair{{a: 3, b: 2}, {a: 3, b: 1}, {a: 2, b: 2}, {a: 3, b: 2}, {a: 2, b: 2}}},
	{n: 2, pairs: []pair{{a: 1, b: 2}, {a: 2, b: 1}}},
	{n: 1, pairs: []pair{{a: 1, b: 2}}},
	{n: 6, pairs: []pair{{a: 2, b: 3}, {a: 3, b: 3}, {a: 2, b: 3}, {a: 2, b: 1}, {a: 3, b: 1}, {a: 1, b: 2}}},
	{n: 2, pairs: []pair{{a: 3, b: 3}, {a: 2, b: 2}}},
	{n: 2, pairs: []pair{{a: 2, b: 2}, {a: 3, b: 3}}},
	{n: 3, pairs: []pair{{a: 3, b: 3}, {a: 1, b: 2}, {a: 1, b: 1}}},
	{n: 6, pairs: []pair{{a: 1, b: 2}, {a: 3, b: 3}, {a: 1, b: 1}, {a: 2, b: 1}, {a: 2, b: 2}, {a: 1, b: 1}}},
	{n: 3, pairs: []pair{{a: 3, b: 1}, {a: 2, b: 3}, {a: 3, b: 2}}},
	{n: 1, pairs: []pair{{a: 2, b: 2}}},
	{n: 3, pairs: []pair{{a: 1, b: 3}, {a: 2, b: 3}, {a: 3, b: 1}}},
	{n: 3, pairs: []pair{{a: 3, b: 1}, {a: 2, b: 2}, {a: 1, b: 2}}},
	{n: 4, pairs: []pair{{a: 3, b: 1}, {a: 2, b: 3}, {a: 1, b: 2}, {a: 1, b: 2}}},
	{n: 2, pairs: []pair{{a: 2, b: 2}, {a: 2, b: 3}}},
	{n: 1, pairs: []pair{{a: 2, b: 1}}},
	{n: 4, pairs: []pair{{a: 1, b: 1}, {a: 1, b: 1}, {a: 1, b: 3}, {a: 1, b: 3}}},
	{n: 6, pairs: []pair{{a: 1, b: 3}, {a: 1, b: 3}, {a: 2, b: 3}, {a: 1, b: 2}, {a: 1, b: 1}, {a: 3, b: 2}}},
	{n: 4, pairs: []pair{{a: 3, b: 1}, {a: 2, b: 1}, {a: 1, b: 2}, {a: 2, b: 2}}},
	{n: 1, pairs: []pair{{a: 1, b: 2}}},
	{n: 4, pairs: []pair{{a: 1, b: 3}, {a: 2, b: 1}, {a: 2, b: 1}, {a: 1, b: 1}}},
	{n: 3, pairs: []pair{{a: 2, b: 1}, {a: 3, b: 1}, {a: 1, b: 2}}},
	{n: 3, pairs: []pair{{a: 1, b: 2}, {a: 1, b: 2}, {a: 3, b: 1}}},
	{n: 5, pairs: []pair{{a: 2, b: 3}, {a: 3, b: 3}, {a: 3, b: 1}, {a: 2, b: 2}, {a: 1, b: 2}}},
	{n: 2, pairs: []pair{{a: 3, b: 1}, {a: 2, b: 2}}},
	{n: 6, pairs: []pair{{a: 2, b: 3}, {a: 2, b: 2}, {a: 3, b: 1}, {a: 3, b: 3}, {a: 2, b: 2}, {a: 2, b: 2}}},
	{n: 1, pairs: []pair{{a: 2, b: 3}}},
	{n: 6, pairs: []pair{{a: 1, b: 1}, {a: 2, b: 3}, {a: 1, b: 2}, {a: 3, b: 1}, {a: 1, b: 3}, {a: 3, b: 2}}},
	{n: 5, pairs: []pair{{a: 2, b: 3}, {a: 2, b: 2}, {a: 1, b: 1}, {a: 1, b: 1}, {a: 1, b: 3}}},
	{n: 3, pairs: []pair{{a: 1, b: 2}, {a: 2, b: 1}, {a: 3, b: 1}}},
	{n: 2, pairs: []pair{{a: 1, b: 3}, {a: 3, b: 2}}},
	{n: 5, pairs: []pair{{a: 2, b: 3}, {a: 2, b: 1}, {a: 1, b: 1}, {a: 3, b: 3}, {a: 1, b: 2}}},
	{n: 5, pairs: []pair{{a: 2, b: 3}, {a: 1, b: 1}, {a: 2, b: 1}, {a: 3, b: 1}, {a: 2, b: 2}}},
	{n: 1, pairs: []pair{{a: 1, b: 3}}},
	{n: 4, pairs: []pair{{a: 2, b: 1}, {a: 2, b: 2}, {a: 1, b: 2}, {a: 2, b: 1}}},
	{n: 1, pairs: []pair{{a: 1, b: 3}}},
	{n: 3, pairs: []pair{{a: 3, b: 2}, {a: 2, b: 3}, {a: 2, b: 1}}},
	{n: 6, pairs: []pair{{a: 1, b: 3}, {a: 2, b: 2}, {a: 2, b: 1}, {a: 3, b: 1}, {a: 2, b: 1}, {a: 1, b: 3}}},
	{n: 2, pairs: []pair{{a: 2, b: 2}, {a: 2, b: 2}}},
	{n: 5, pairs: []pair{{a: 1, b: 3}, {a: 1, b: 1}, {a: 3, b: 2}, {a: 1, b: 3}, {a: 1, b: 2}}},
	{n: 6, pairs: []pair{{a: 3, b: 2}, {a: 3, b: 1}, {a: 3, b: 1}, {a: 2, b: 2}, {a: 2, b: 1}, {a: 1, b: 3}}},
	{n: 1, pairs: []pair{{a: 3, b: 1}}},
	{n: 5, pairs: []pair{{a: 2, b: 3}, {a: 2, b: 3}, {a: 2, b: 3}, {a: 2, b: 3}, {a: 3, b: 2}}},
	{n: 4, pairs: []pair{{a: 1, b: 2}, {a: 3, b: 2}, {a: 3, b: 2}, {a: 1, b: 2}}},
	{n: 1, pairs: []pair{{a: 3, b: 1}}},
	{n: 6, pairs: []pair{{a: 3, b: 1}, {a: 2, b: 2}, {a: 1, b: 3}, {a: 2, b: 3}, {a: 3, b: 1}, {a: 1, b: 2}}},
	{n: 6, pairs: []pair{{a: 1, b: 2}, {a: 3, b: 2}, {a: 2, b: 1}, {a: 2, b: 1}, {a: 3, b: 1}, {a: 3, b: 1}}},
	{n: 3, pairs: []pair{{a: 1, b: 1}, {a: 1, b: 3}, {a: 2, b: 3}}},
	{n: 5, pairs: []pair{{a: 1, b: 1}, {a: 1, b: 3}, {a: 2, b: 3}, {a: 3, b: 1}, {a: 3, b: 3}}},
	{n: 3, pairs: []pair{{a: 3, b: 2}, {a: 3, b: 1}, {a: 2, b: 3}}},
	{n: 2, pairs: []pair{{a: 2, b: 2}, {a: 1, b: 3}}},
	{n: 3, pairs: []pair{{a: 3, b: 3}, {a: 3, b: 2}, {a: 2, b: 3}}},
	{n: 4, pairs: []pair{{a: 2, b: 2}, {a: 2, b: 2}, {a: 3, b: 1}, {a: 2, b: 1}}},
	{n: 5, pairs: []pair{{a: 2, b: 3}, {a: 3, b: 3}, {a: 2, b: 2}, {a: 2, b: 3}, {a: 1, b: 1}}},
	{n: 5, pairs: []pair{{a: 1, b: 2}, {a: 3, b: 2}, {a: 3, b: 2}, {a: 1, b: 2}, {a: 2, b: 1}}},
	{n: 6, pairs: []pair{{a: 2, b: 3}, {a: 1, b: 1}, {a: 1, b: 1}, {a: 2, b: 1}, {a: 1, b: 1}, {a: 2, b: 3}}},
	{n: 6, pairs: []pair{{a: 3, b: 3}, {a: 3, b: 3}, {a: 2, b: 1}, {a: 3, b: 1}, {a: 1, b: 1}, {a: 1, b: 1}}},
	{n: 3, pairs: []pair{{a: 1, b: 1}, {a: 2, b: 1}, {a: 3, b: 2}}},
	{n: 5, pairs: []pair{{a: 3, b: 3}, {a: 1, b: 3}, {a: 1, b: 3}, {a: 2, b: 1}, {a: 3, b: 3}}},
	{n: 1, pairs: []pair{{a: 3, b: 2}}},
	{n: 6, pairs: []pair{{a: 1, b: 3}, {a: 3, b: 1}, {a: 3, b: 2}, {a: 1, b: 1}, {a: 2, b: 2}, {a: 3, b: 1}}},
	{n: 4, pairs: []pair{{a: 1, b: 1}, {a: 2, b: 3}, {a: 2, b: 3}, {a: 3, b: 2}}},
	{n: 4, pairs: []pair{{a: 2, b: 1}, {a: 2, b: 1}, {a: 2, b: 1}, {a: 2, b: 3}}},
	{n: 1, pairs: []pair{{a: 2, b: 1}}},
	{n: 1, pairs: []pair{{a: 1, b: 1}}},
	{n: 5, pairs: []pair{{a: 1, b: 1}, {a: 1, b: 3}, {a: 1, b: 2}, {a: 3, b: 3}, {a: 1, b: 2}}},
	{n: 6, pairs: []pair{{a: 3, b: 3}, {a: 3, b: 2}, {a: 2, b: 1}, {a: 3, b: 3}, {a: 1, b: 1}, {a: 2, b: 3}}},
	{n: 3, pairs: []pair{{a: 2, b: 3}, {a: 2, b: 3}, {a: 3, b: 1}}},
	{n: 6, pairs: []pair{{a: 1, b: 2}, {a: 3, b: 3}, {a: 3, b: 2}, {a: 2, b: 3}, {a: 2, b: 1}, {a: 1, b: 2}}},
	{n: 1, pairs: []pair{{a: 3, b: 3}}},
	{n: 1, pairs: []pair{{a: 3, b: 2}}},
	{n: 5, pairs: []pair{{a: 2, b: 2}, {a: 1, b: 1}, {a: 3, b: 1}, {a: 2, b: 3}, {a: 3, b: 3}}},
	{n: 6, pairs: []pair{{a: 1, b: 1}, {a: 1, b: 3}, {a: 1, b: 2}, {a: 3, b: 1}, {a: 3, b: 2}, {a: 1, b: 3}}},
	{n: 5, pairs: []pair{{a: 1, b: 1}, {a: 1, b: 3}, {a: 2, b: 3}, {a: 1, b: 3}, {a: 1, b: 1}}},
	{n: 4, pairs: []pair{{a: 3, b: 1}, {a: 1, b: 2}, {a: 1, b: 3}, {a: 3, b: 2}}},
	{n: 3, pairs: []pair{{a: 2, b: 3}, {a: 3, b: 2}, {a: 1, b: 1}}},
	{n: 4, pairs: []pair{{a: 2, b: 1}, {a: 3, b: 3}, {a: 2, b: 3}, {a: 3, b: 2}}},
	{n: 2, pairs: []pair{{a: 3, b: 3}, {a: 2, b: 2}}},
	{n: 2, pairs: []pair{{a: 1, b: 1}, {a: 2, b: 2}}},
	{n: 3, pairs: []pair{{a: 2, b: 2}, {a: 3, b: 2}, {a: 2, b: 3}}},
	{n: 2, pairs: []pair{{a: 2, b: 2}, {a: 3, b: 1}}},
	{n: 1, pairs: []pair{{a: 2, b: 2}}},
	{n: 5, pairs: []pair{{a: 2, b: 3}, {a: 3, b: 1}, {a: 3, b: 2}, {a: 1, b: 1}, {a: 3, b: 1}}},
	{n: 6, pairs: []pair{{a: 1, b: 2}, {a: 3, b: 1}, {a: 3, b: 2}, {a: 1, b: 3}, {a: 2, b: 1}, {a: 2, b: 3}}},
	{n: 2, pairs: []pair{{a: 1, b: 2}, {a: 1, b: 2}}},
	{n: 1, pairs: []pair{{a: 1, b: 2}}},
	{n: 3, pairs: []pair{{a: 3, b: 3}, {a: 3, b: 1}, {a: 3, b: 1}}},
	{n: 6, pairs: []pair{{a: 1, b: 3}, {a: 3, b: 1}, {a: 1, b: 3}, {a: 1, b: 2}, {a: 1, b: 3}, {a: 2, b: 3}}},
	{n: 3, pairs: []pair{{a: 3, b: 3}, {a: 2, b: 1}, {a: 1, b: 1}}},
	{n: 1, pairs: []pair{{a: 1, b: 1}}},
	{n: 6, pairs: []pair{{a: 1, b: 3}, {a: 2, b: 3}, {a: 3, b: 3}, {a: 2, b: 1}, {a: 1, b: 1}, {a: 2, b: 2}}},
	{n: 5, pairs: []pair{{a: 3, b: 2}, {a: 3, b: 2}, {a: 2, b: 2}, {a: 3, b: 2}, {a: 1, b: 1}}},
	{n: 2, pairs: []pair{{a: 1, b: 2}, {a: 1, b: 2}}},
	{n: 2, pairs: []pair{{a: 2, b: 3}, {a: 1, b: 3}}},
	{n: 6, pairs: []pair{{a: 2, b: 2}, {a: 2, b: 3}, {a: 3, b: 1}, {a: 1, b: 1}, {a: 2, b: 1}, {a: 3, b: 1}}},
}

func solveExpected(tc testCase) int64 {
	n := tc.n
	pairs := make([][2]int, n)
	for i, p := range tc.pairs {
		pairs[i] = [2]int{p.a, p.b}
	}
	fac := make([]int64, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	tmp := make([][2]int, n)
	copy(tmp, pairs)
	sort.Slice(tmp, func(i, j int) bool { return tmp[i][0] < tmp[j][0] })
	cntA := int64(1)
	for i := 0; i < n; {
		j := i
		for j < n && tmp[j][0] == tmp[i][0] {
			j++
		}
		cntA = cntA * fac[j-i] % mod
		i = j
	}
	copy(tmp, pairs)
	sort.Slice(tmp, func(i, j int) bool { return tmp[i][1] < tmp[j][1] })
	cntB := int64(1)
	for i := 0; i < n; {
		j := i
		for j < n && tmp[j][1] == tmp[i][1] {
			j++
		}
		cntB = cntB * fac[j-i] % mod
		i = j
	}
	copy(tmp, pairs)
	sort.Slice(tmp, func(i, j int) bool {
		if tmp[i][0] == tmp[j][0] {
			return tmp[i][1] < tmp[j][1]
		}
		return tmp[i][0] < tmp[j][0]
	})
	valid := true
	for i := 1; i < n; i++ {
		if tmp[i][1] < tmp[i-1][1] {
			valid = false
			break
		}
	}
	cntAB := int64(0)
	if valid {
		cntAB = 1
		for i := 0; i < n; {
			j := i
			for j < n && tmp[j] == tmp[i] {
				j++
			}
			cntAB = cntAB * fac[j-i] % mod
			i = j
		}
	}
	ans := (fac[n] - cntA - cntB + cntAB) % mod
	if ans < 0 {
		ans += mod
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, tc := range testcases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for _, p := range tc.pairs {
			fmt.Fprintf(&sb, "%d %d\n", p.a, p.b)
		}
		input := sb.String()
		want := solveExpected(tc)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d failed: %v\nstderr: %s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != fmt.Sprintf("%d", want) {
			fmt.Printf("case %d failed: expected %d got %s\ninput:\n%s", idx+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
