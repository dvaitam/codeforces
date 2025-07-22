package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const MOD = 1000000007

func solveD(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return ""
	}
	Spos := n - k - 1
	jumps := make([]int, 0, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		d := v - u
		if d == 1 {
			continue
		} else if d == k+1 {
			jumps = append(jumps, u)
		} else {
			return "0"
		}
	}
	sort.Ints(jumps)
	P := len(jumps)
	if Spos <= 0 {
		if P == 0 {
			return "1"
		} else {
			return "0"
		}
	}
	maxN := Spos
	if k+1 > maxN {
		maxN = k + 1
	}
	pow2 := make([]int, maxN+2)
	pow2[0] = 1
	for i := 1; i <= maxN+1; i++ {
		pow2[i] = pow2[i-1] * 2 % MOD
	}
	if P == 0 {
		if Spos <= k {
			return fmt.Sprint(pow2[Spos])
		}
		res := int64(Spos-k+1) * int64(pow2[k]) % MOD
		return fmt.Sprint(res)
	}
	imin := jumps[0]
	imax := jumps[P-1]
	if imax-imin > k {
		return "0"
	}
	low := imax - k
	if low < 1 {
		low = 1
	}
	high := imin
	res := 0
	for L0 := low; L0 <= high; L0++ {
		R0 := L0 + k
		if R0 > Spos {
			R0 = Spos
		}
		wlen := R0 - L0 + 1
		lo := sort.Search(len(jumps), func(i int) bool { return jumps[i] >= L0 })
		hi := sort.Search(len(jumps), func(i int) bool { return jumps[i] > R0 }) - 1
		inCnt := 0
		if lo < len(jumps) && hi >= lo {
			inCnt = hi - lo + 1
		}
		aSz := wlen - inCnt
		var add int
		if L0 == imin {
			add = pow2[aSz]
		} else {
			if aSz-1 >= 0 {
				add = pow2[aSz-1]
			}
		}
		res = (res + add) % MOD
	}
	return fmt.Sprint(res)
}

func genTestD() string {
	n := rand.Intn(8) + 2
	k := rand.Intn(n-1) + 1
	edges := make([][2]int, 0)
	for u := 1; u <= n; u++ {
		if u+1 <= n && rand.Intn(2) == 0 {
			edges = append(edges, [2]int{u, u + 1})
		}
		if u+k+1 <= n && rand.Intn(2) == 0 {
			edges = append(edges, [2]int{u, u + k + 1})
		}
	}
	m := len(edges)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d %d %d\n", n, m, k)
	for _, e := range edges {
		fmt.Fprintf(&buf, "%d %d\n", e[0], e[1])
	}
	return buf.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		in := genTestD()
		expected := solveD(in)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(in)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\noutput: %s\n", i, err, string(out))
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
