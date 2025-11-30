package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesE.txt.
const embeddedTestcasesE = `2 3 c baaa 1 2 2 1 1 2 2 2 1
3 2 c baa bcc 2 2 3 2 3 3
2 2 bb c 2 2 2 2 2 2
4 5 bba bb a b 3 4 3 2 4 2 1 4 2 4 4 2 3 4 3
5 2 aac bc aabab aabc b 3 3 3 3 4 2
4 5 b abbab cabc b 1 4 2 3 4 3 1 4 2 4 4 1 1 1 1
2 5 ca bcaba 1 2 2 1 2 1 1 2 2 2 2 1 2 2 1
4 2 aaab aca bb ba 3 3 4 1 4 4
1 4 ca 1 1 1 1 1 1 1 1 1 1 1 1 1 1
2 1 cabc a 2 2 2
4 2 a ac abb ca 2 2 3 4 4 1
1 1 abcbc 1 1 1
1 5 b 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
5 4 cbac cb baca aac bb 3 4 5 1 2 1 3 3 5 1 4 5
4 5 c bb aac c 1 3 4 3 4 4 4 4 3 3 4 2 4 4 2
5 2 bac cca b bc bab 1 5 1 5 5 3
1 2 cb 1 1 1 1 1 1
4 5 bc bb cb ccbbc 4 4 3 4 4 1 3 3 4 3 4 4 1 2 1
4 3 babba caaa b aa 4 4 1 2 2 1 2 2 3
1 1 acb 1 1 1
2 1 bcacc cb 1 1 2
4 5 b ab bccbb aba 4 4 3 1 3 2 1 2 1 1 2 2 1 4 2
4 4 acc ab bbcbc c 2 3 3 4 4 1 2 3 1 1 3 4
3 1 ca cccc aa 3 3 2
5 2 acc a cb acaab ab 2 4 3 4 5 1
2 4 accb cbacc 2 2 1 1 1 2 2 2 2 2 2 2
5 2 bca b cbc cb a 5 5 2 3 5 2
4 1 abc aa ca b 3 3 1
2 1 ca acb 2 2 1
2 2 a ca 2 2 1 1 1 2
4 5 bcbbb aaaa aba cac 3 4 3 1 2 2 3 3 2 1 1 1 2 3 3
5 4 b aabc bbab ba aabb 5 5 1 4 4 5 1 2 5 2 4 5
1 1 baca 1 1 1
5 3 bc ccb cb ccbbc cac 2 4 5 2 5 4 1 1 5
2 2 cc ca 2 2 1 1 1 2
5 1 bba bcca bb a ccb 1 2 4
1 4 b 1 1 1 1 1 1 1 1 1 1 1 1
5 3 ca ccaca a aacb bcbac 5 5 5 5 5 5 1 2 1
3 5 bcba bcbac abcc 3 3 3 3 3 2 3 3 2 1 2 1 2 3 2
2 5 cbbb ba 1 2 1 1 2 1 2 2 1 1 1 1 1 1 2
4 4 cb bcab a ab 1 4 2 4 4 2 3 4 3 2 2 3
3 1 cabb bccab ab 2 2 1
1 3 a 1 1 1 1 1 1 1 1 1
1 3 cacb 1 1 1 1 1 1 1 1 1
5 2 abacb bbcc b c bbcc 3 3 2 3 4 2
3 5 cbcab cbbc accbb 3 3 3 3 3 1 2 3 1 1 1 1 2 2 2
3 4 abaca b c 3 3 1 2 2 3 3 3 1 1 1 1
5 2 ac bbbc cacb c baa 3 5 1 4 4 2
2 1 cbba a 2 2 2
1 2 aa 1 1 1 1 1 1
3 3 abbcc cbc a 3 3 1 3 3 3 3 3 3
2 5 ca cabaa 2 2 1 2 2 1 1 1 2 2 2 2 1 1 2
3 5 cb cbcca aaab 3 3 1 1 3 3 2 3 1 3 3 2 3 3 3
4 3 cab bbaa cbc cbcba 3 4 1 3 3 2 4 4 4
4 2 b ca caab bbbb 4 4 1 1 2 1
5 4 baa acc bbab bbca caab 4 5 2 4 5 3 5 5 2 2 4 5
5 2 b accbc ab bbaa a 4 5 5 4 4 3
3 5 aabb bbb c 3 3 2 2 2 1 3 3 1 2 3 2 3 3 2
4 4 cbcac cabb cca cccc 2 3 4 3 4 3 1 3 4 3 4 2
2 2 bcabb bcbc 1 2 1 2 2 2
2 3 aaba caab 1 1 2 2 2 2 1 1 1
4 3 cb abc acbcb ccbaa 3 4 2 3 4 2 2 2 4
2 4 ca aca 2 2 1 2 2 1 1 2 2 1 2 1
3 2 ca ccbaa ac 1 1 3 2 2 1
3 3 bba c bcb 3 3 1 3 3 3 2 3 1
5 3 c cca cbcab bcbc bcac 2 3 4 5 5 5 3 4 1
3 2 acbb a bcac 3 3 1 1 1 3
1 4 bc 1 1 1 1 1 1 1 1 1 1 1 1
1 5 b 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
2 4 cc bb 2 2 2 2 2 2 1 2 1 1 1 1
5 3 cbcc c baaaa aaaa baaa 5 5 3 5 5 3 1 3 4
2 4 aba aaa 1 2 2 2 2 1 1 2 1 2 2 2
4 4 cac ccbc bbab cccc 1 2 1 4 4 3 2 4 4 2 3 3
4 5 ab a aca aab 4 4 2 1 1 1 1 1 4 1 2 2 3 4 1
2 3 cc bacbc 2 2 1 2 2 2 2 2 1
3 3 bcaaa c acc 3 3 2 3 3 2 2 3 3
3 1 abc acaa b 1 2 3
3 3 a bca b 2 3 3 3 3 1 3 3 2
2 3 ccbbc bcbbc 2 2 1 1 1 2 1 2 2
1 5 accb 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
5 3 cbbc aaca ab aabca bba 5 5 5 5 5 5 5 5 5
1 1 ba 1 1 1
2 5 bb abb 2 2 2 2 2 1 1 1 1 2 2 1 1 1 1
2 4 bca cabab 1 2 2 2 2 1 2 2 1 1 2 1
2 1 ca cabb 2 2 2
1 4 ccab 1 1 1 1 1 1 1 1 1 1 1 1
2 3 b b 1 1 1 2 2 2 1 2 2
3 3 bacca acaba cb 3 3 2 1 1 2 1 1 3
2 1 acaa ab 2 2 2
5 3 abaac acc accc bcacc bbaa 2 4 4 4 5 1 5 5 4
5 4 aabb bbaca ba aaa cc 3 4 3 2 5 3 5 5 3 5 5 4
5 2 acbba cc cb c c 3 5 3 5 5 5
2 3 cbc acac 1 1 1 2 2 1 2 2 1
3 2 aaca c a 2 2 2 1 1 3
2 3 a a 1 1 2 2 2 2 2 2 1
2 5 bca bbbab 2 2 1 1 2 1 1 2 2 1 2 2 1 2 2
3 3 ac abcba bc 2 3 2 3 3 3 2 3 2
3 3 b b a 2 3 1 3 3 2 2 3 3
2 1 abbcb a 1 1 1
3 5 aabbc ac abcaa 3 3 2 1 3 3 1 2 2 3 3 3 1 1 1`

// BIT implements a Fenwick tree for int64 values.
type BIT struct {
	n    int
	tree []int64
}

func newBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+1)}
}

func (b *BIT) add(i int, v int64) {
	for x := i; x <= b.n; x += x & -x {
		b.tree[x] += v
	}
}

func (b *BIT) sum(i int) int64 {
	var s int64
	for x := i; x > 0; x -= x & -x {
		s += b.tree[x]
	}
	return s
}

func solve547E(n, q int, strs []string, queries [][3]int) []int64 {
	next := make([][]int, 1)
	next[0] = make([]int, 26)
	fail := make([]int, 1)
	endpoints := make([]int, n+1)

	for i := 1; i <= n; i++ {
		u := 0
		for _, ch := range strs[i-1] {
			c := int(ch - 'a')
			if next[u][c] == 0 {
				next = append(next, make([]int, 26))
				fail = append(fail, 0)
				next[u][c] = len(next) - 1
			}
			u = next[u][c]
		}
		endpoints[i] = u
	}

	queue := make([]int, 0, len(next))
	for c := 0; c < 26; c++ {
		v := next[0][c]
		if v != 0 {
			queue = append(queue, v)
		}
	}
	for idx := 0; idx < len(queue); idx++ {
		u := queue[idx]
		for c := 0; c < 26; c++ {
			v := next[u][c]
			if v != 0 {
				f := fail[u]
				next[u][c] = v
				fail[v] = next[f][c]
				queue = append(queue, v)
			} else {
				next[u][c] = next[fail[u]][c]
			}
		}
	}

	sz := len(next)
	children := make([][]int, sz)
	for u := 1; u < sz; u++ {
		p := fail[u]
		children[p] = append(children[p], u)
	}
	tin := make([]int, sz)
	tout := make([]int, sz)
	time := 0
	type frame struct {
		u, idx int
	}
	stack := []frame{{0, 0}}
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		u := top.u
		if top.idx == 0 {
			time++
			tin[u] = time
		}
		if top.idx < len(children[u]) {
			v := children[u][top.idx]
			stack[len(stack)-1].idx++
			stack = append(stack, frame{v, 0})
		} else {
			tout[u] = time
			stack = stack[:len(stack)-1]
		}
	}

	visits := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		u := 0
		for _, ch := range strs[i-1] {
			u = next[u][int(ch-'a')]
			visits[i] = append(visits[i], u)
		}
	}

	type event struct {
		p, id, sign int
	}
	events := make([][]event, n+1)
	for id, qu := range queries {
		l, r, k := qu[0], qu[1], qu[2]
		p := endpoints[k]
		events[r] = append(events[r], event{p: p, id: id, sign: 1})
		if l > 1 {
			events[l-1] = append(events[l-1], event{p: p, id: id, sign: -1})
		}
	}

	ans := make([]int64, q)
	bit := newBIT(sz + 2)
	for t := 1; t <= n; t++ {
		for _, u := range visits[t] {
			bit.add(tin[u], 1)
		}
		for _, ev := range events[t] {
			lo := tin[ev.p]
			hi := tout[ev.p]
			cnt := bit.sum(hi) - bit.sum(lo-1)
			ans[ev.id] += int64(ev.sign) * cnt
		}
	}
	return ans
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesE), "\n")
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "case %d: malformed line\n", idx+1)
			os.Exit(1)
		}
		n, err1 := strconv.Atoi(fields[0])
		q, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid n/q\n", idx+1)
			os.Exit(1)
		}
		expectedCount := 2 + n + 3*q
		if len(fields) != expectedCount {
			fmt.Fprintf(os.Stderr, "case %d: expected %d tokens, got %d\n", idx+1, expectedCount, len(fields))
			os.Exit(1)
		}
		strs := make([]string, n)
		pos := 2
		for i := 0; i < n; i++ {
			strs[i] = fields[pos]
			pos++
		}
		queries := make([][3]int, q)
		for i := 0; i < q; i++ {
			l := mustAtoi(fields[pos])
			r := mustAtoi(fields[pos+1])
			k := mustAtoi(fields[pos+2])
			queries[i] = [3]int{l, r, k}
			pos += 3
		}

		wantVals := solve547E(n, q, strs, queries)
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", n, q)
		for _, s := range strs {
			fmt.Fprintln(&input, s)
		}
		for _, qu := range queries {
			fmt.Fprintf(&input, "%d %d %d\n", qu[0], qu[1], qu[2])
		}

		gotRaw, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotFields := strings.Fields(gotRaw)
		if len(gotFields) != len(wantVals) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d outputs, got %d\n", idx+1, len(wantVals), len(gotFields))
			os.Exit(1)
		}
		for i, g := range gotFields {
			v, err := strconv.ParseInt(g, 10, 64)
			if err != nil || v != wantVals[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at query %d: expected %d got %s\n", idx+1, i+1, wantVals[i], g)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}

func mustAtoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	return v
}
