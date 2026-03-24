package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func solveBuildSA(s string) ([]int, []int) {
	s += "$"
	n := len(s)
	p := make([]int, n)
	c := make([]int, n)
	cnt := make([]int, n)
	if 256 > n {
		cnt = make([]int, 256)
	}
	for i := 0; i < n; i++ {
		cnt[s[i]]++
	}
	for i := 1; i < len(cnt); i++ {
		cnt[i] += cnt[i-1]
	}
	for i := n - 1; i >= 0; i-- {
		cnt[s[i]]--
		p[cnt[s[i]]] = i
	}
	c[p[0]] = 0
	classes := 1
	for i := 1; i < n; i++ {
		if s[p[i]] != s[p[i-1]] {
			classes++
		}
		c[p[i]] = classes - 1
	}
	pn := make([]int, n)
	cn := make([]int, n)
	for k := 1; k < n; k *= 2 {
		for i := 0; i < n; i++ {
			pn[i] = p[i] - k
			if pn[i] < 0 {
				pn[i] += n
			}
		}
		for i := 0; i < classes; i++ {
			cnt[i] = 0
		}
		for i := 0; i < n; i++ {
			cnt[c[pn[i]]]++
		}
		for i := 1; i < classes; i++ {
			cnt[i] += cnt[i-1]
		}
		for i := n - 1; i >= 0; i-- {
			cnt[c[pn[i]]]--
			p[cnt[c[pn[i]]]] = pn[i]
		}
		cn[p[0]] = 0
		classes = 1
		for i := 1; i < n; i++ {
			mid1, mid2 := c[p[i]], c[p[i-1]]
			val1, val2 := c[(p[i]+k)%n], c[(p[i-1]+k)%n]
			if mid1 != mid2 || val1 != val2 {
				classes++
			}
			cn[p[i]] = classes - 1
		}
		copy(c, cn)
	}
	sa := make([]int, n-1)
	rank := make([]int, n-1)
	for i := 1; i < n; i++ {
		sa[i-1] = p[i]
		rank[p[i]] = i - 1
	}
	return sa, rank
}

func solveBuildLCP(s string, sa []int, rank []int) []int {
	n := len(s)
	lcp := make([]int, n)
	k := 0
	for i := 0; i < n; i++ {
		if rank[i] == n-1 {
			k = 0
			continue
		}
		j := sa[rank[i]+1]
		for i+k < n && j+k < n && s[i+k] == s[j+k] {
			k++
		}
		lcp[rank[i]+1] = k
		if k > 0 {
			k--
		}
	}
	return lcp
}

func solveBuildST(lcp []int) ([][]int, []int) {
	n := len(lcp)
	log := make([]int, n+1)
	for i := 2; i <= n; i++ {
		log[i] = log[i/2] + 1
	}
	K := log[n] + 1
	st := make([][]int, K)
	for i := 0; i < K; i++ {
		st[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		st[0][i] = lcp[i]
	}
	for j := 1; j < K; j++ {
		for i := 0; i+(1<<j) <= n; i++ {
			a := st[j-1][i]
			b := st[j-1][i+(1<<(j-1))]
			if a < b {
				st[j][i] = a
			} else {
				st[j][i] = b
			}
		}
	}
	return st, log
}

func solveQueryST(st [][]int, log2 []int, L, R int) int {
	j := log2[R-L+1]
	a := st[j][L]
	b := st[j][R-(1<<j)+1]
	if a < b {
		return a
	}
	return b
}

func solve(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)
	scanner.Split(bufio.ScanWords)

	readInt := func() int {
		scanner.Scan()
		res := 0
		for _, b := range scanner.Bytes() {
			res = res*10 + int(b-'0')
		}
		return res
	}

	readString := func() string {
		scanner.Scan()
		return scanner.Text()
	}

	n := readInt()
	q := readInt()
	s := readString()

	sa, rank := solveBuildSA(s)
	lcp := solveBuildLCP(s, sa, rank)
	st, log2 := solveBuildST(lcp)

	var out strings.Builder

	type Element struct {
		r   int
		isA bool
	}

	type Node struct {
		r  int
		cA int
		cB int
	}

	type StackItem struct {
		val   int
		count int
	}

	for qi := 0; qi < q; qi++ {
		k := readInt()
		l := readInt()
		elements := make([]Element, 0, k+l)
		for i := 0; i < k; i++ {
			idx := readInt() - 1
			elements = append(elements, Element{r: rank[idx], isA: true})
		}
		for i := 0; i < l; i++ {
			idx := readInt() - 1
			elements = append(elements, Element{r: rank[idx], isA: false})
		}
		sort.Slice(elements, func(i, j int) bool {
			return elements[i].r < elements[j].r
		})

		nodes := make([]Node, 0)
		for i := 0; i < len(elements); i++ {
			r := elements[i].r
			cA, cB := 0, 0
			if elements[i].isA {
				cA++
			} else {
				cB++
			}
			for i+1 < len(elements) && elements[i+1].r == r {
				i++
				if elements[i].isA {
					cA++
				} else {
					cB++
				}
			}
			nodes = append(nodes, Node{r: r, cA: cA, cB: cB})
		}

		var ans int64 = 0
		for _, nd := range nodes {
			ans += int64(nd.cA) * int64(nd.cB) * int64(n-sa[nd.r])
		}

		var stack []StackItem
		var sum int64 = 0

		for i := 0; i < len(nodes); i++ {
			nd := nodes[i]
			if i > 0 {
				w := solveQueryST(st, log2, nodes[i-1].r+1, nd.r)
				countAccum := 0
				for len(stack) > 0 && stack[len(stack)-1].val > w {
					top := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					sum -= int64(top.val) * int64(top.count)
					countAccum += top.count
				}
				if countAccum > 0 {
					stack = append(stack, StackItem{val: w, count: countAccum})
					sum += int64(w) * int64(countAccum)
				}
			}
			if nd.cB > 0 {
				ans += sum * int64(nd.cB)
			}
			if nd.cA > 0 {
				lenA := n - sa[nd.r]
				stack = append(stack, StackItem{val: lenA, count: nd.cA})
				sum += int64(lenA) * int64(nd.cA)
			}
		}

		stack = stack[:0]
		sum = 0
		for i := 0; i < len(nodes); i++ {
			nd := nodes[i]
			if i > 0 {
				w := solveQueryST(st, log2, nodes[i-1].r+1, nd.r)
				countAccum := 0
				for len(stack) > 0 && stack[len(stack)-1].val > w {
					top := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					sum -= int64(top.val) * int64(top.count)
					countAccum += top.count
				}
				if countAccum > 0 {
					stack = append(stack, StackItem{val: w, count: countAccum})
					sum += int64(w) * int64(countAccum)
				}
			}
			if nd.cA > 0 {
				ans += sum * int64(nd.cA)
			}
			if nd.cB > 0 {
				lenB := n - sa[nd.r]
				stack = append(stack, StackItem{val: lenB, count: nd.cB})
				sum += int64(lenB) * int64(nd.cB)
			}
		}

		fmt.Fprintln(&out, ans)
	}

	return strings.TrimSpace(out.String())
}

func runProg(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return out.String() + stderr.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	q := rng.Intn(5) + 1
	letters := []byte{'a', 'b', 'c'}
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(len(letters))]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	sb.WriteString(string(b))
	sb.WriteByte('\n')
	for qi := 0; qi < q; qi++ {
		k := rng.Intn(n) + 1
		l := rng.Intn(n) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", k, l))
		pa := rng.Perm(n)[:k]
		for i, v := range pa {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v + 1))
		}
		sb.WriteByte('\n')
		pb := rng.Perm(n)[:l]
		for i, v := range pb {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v + 1))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp := solve(input)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected: %s\n got: %s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
