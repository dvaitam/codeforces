package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const modF = 998244353

type Fenwick struct {
	n    int
	tree []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int, n+1)}
}

func (f *Fenwick) Add(i, v int) {
	for ; i <= f.n; i += i & -i {
		f.tree[i] += v
	}
}

func (f *Fenwick) Sum(i int) int {
	s := 0
	for ; i > 0; i -= i & -i {
		s += f.tree[i]
	}
	return s
}

func fpowF(a, b int) int {
	res := 1
	for b > 1 {
		if b&1 == 1 {
			res = res * a % modF
		}
		b >>= 1
		a = a * a % modF
	}
	return a * res % modF
}

func solveF(a []int) int {
	n := len(a) - 1
	vis := make([]bool, n+1)
	tot := 0
	for i := 1; i <= n; i++ {
		if a[i] == -1 {
			tot++
		} else if a[i] >= 1 && a[i] <= n {
			vis[a[i]] = true
		}
	}
	ans := tot * (tot - 1) % modF * fpowF(4, modF-2) % modF
	cnt := tot
	fenw := NewFenwick(n)
	for i := n; i >= 1; i-- {
		if a[i] != -1 {
			ans = (ans + fenw.Sum(a[i]-1)) % modF
			fenw.Add(a[i], 1)
		}
	}
	sum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		sum[i] = sum[i-1]
		if !vis[i] {
			sum[i]++
		}
	}
	invTot := 0
	if tot > 0 {
		invTot = fpowF(tot, modF-2)
	}
	for i := 1; i <= n; i++ {
		if a[i] == -1 {
			cnt--
		} else {
			leftMiss := sum[a[i]]
			rightMiss := tot - leftMiss
			add := (leftMiss*cnt%modF + rightMiss*(tot-cnt)%modF) % modF
			ans = (ans + add*invTot%modF) % modF
		}
	}
	if ans < 0 {
		ans += modF
	}
	return ans
}

type testCaseF struct {
	input    string
	expected int
}

func generateCaseF(rng *rand.Rand) testCaseF {
	n := rng.Intn(6) + 1
	arr := make([]int, n+1)
	used := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		if rng.Intn(2) == 0 {
			arr[i] = -1
		} else {
			v := rng.Intn(n) + 1
			for used[v] {
				v = rng.Intn(n) + 1
			}
			used[v] = true
			arr[i] = v
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	return testCaseF{input: sb.String(), expected: solveF(arr)}
}

func runCaseF(bin string, tc testCaseF) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := []testCaseF{generateCaseF(rng)}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseF(rng))
	}
	for i, tc := range cases {
		if err := runCaseF(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
