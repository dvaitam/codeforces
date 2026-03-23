package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = `package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Node struct {
	next [27]int32
	link int32
	len  int32
}

type SegNode struct {
	l, r   int32
	maxVal int32
	maxIdx int32
}

var (
	st   []Node
	seg  []SegNode
	sz   int32
	last int32

	up   [300005][20]int32
	root [300005]int32
)

func initSAM() {
	st = make([]Node, 1, 300005)
	st[0].link = -1
	sz = 1
	last = 0
}

func addChar(c int32) {
	cur := sz
	sz++
	st = append(st, Node{})
	st[cur].len = st[last].len + 1

	p := last
	for p != -1 && st[p].next[c] == 0 {
		st[p].next[c] = cur
		p = st[p].link
	}
	if p == -1 {
		st[cur].link = 0
	} else {
		q := st[p].next[c]
		if st[p].len+1 == st[q].len {
			st[cur].link = q
		} else {
			clone := sz
			sz++
			st = append(st, Node{})
			st[clone].len = st[p].len + 1
			st[clone].next = st[q].next
			st[clone].link = st[q].link
			for p != -1 && st[p].next[c] == q {
				st[p].next[c] = clone
				p = st[p].link
			}
			st[q].link = clone
			st[cur].link = clone
		}
	}
	last = cur
}

func initSeg() {
	seg = make([]SegNode, 1, 6000000)
}

func allocNode() int32 {
	seg = append(seg, SegNode{})
	return int32(len(seg) - 1)
}

func pushUp(u int32) {
	lc := seg[u].l
	rc := seg[u].r
	if lc == 0 && rc == 0 {
		return
	}
	if lc == 0 {
		seg[u].maxVal = seg[rc].maxVal
		seg[u].maxIdx = seg[rc].maxIdx
		return
	}
	if rc == 0 {
		seg[u].maxVal = seg[lc].maxVal
		seg[u].maxIdx = seg[lc].maxIdx
		return
	}
	if seg[lc].maxVal >= seg[rc].maxVal {
		seg[u].maxVal = seg[lc].maxVal
		seg[u].maxIdx = seg[lc].maxIdx
	} else {
		seg[u].maxVal = seg[rc].maxVal
		seg[u].maxIdx = seg[rc].maxIdx
	}
}

func insert(l, r, pos, val int32) int32 {
	newNode := allocNode()
	if l == r {
		seg[newNode].maxVal += val
		seg[newNode].maxIdx = l
		return newNode
	}
	mid := (l + r) / 2
	if pos <= mid {
		seg[newNode].l = insert(l, mid, pos, val)
	} else {
		seg[newNode].r = insert(mid+1, r, pos, val)
	}
	pushUp(newNode)
	return newNode
}

func merge(u, v, l, r int32) int32 {
	if u == 0 || v == 0 {
		return u ^ v
	}
	newNode := allocNode()
	if l == r {
		seg[newNode].maxVal = seg[u].maxVal + seg[v].maxVal
		seg[newNode].maxIdx = l
		return newNode
	}
	mid := (l + r) / 2
	seg[newNode].l = merge(seg[u].l, seg[v].l, l, mid)
	seg[newNode].r = merge(seg[u].r, seg[v].r, mid+1, r)
	pushUp(newNode)
	return newNode
}

func query(u, left, right, ql, qr int32) (int32, int32) {
	if ql <= left && right <= qr {
		if u == 0 || seg[u].maxVal == 0 {
			return 0, left
		}
		return seg[u].maxVal, seg[u].maxIdx
	}
	mid := (left + right) / 2
	if qr <= mid {
		return query(seg[u].l, left, mid, ql, qr)
	} else if ql > mid {
		return query(seg[u].r, mid+1, right, ql, qr)
	} else {
		v1, id1 := query(seg[u].l, left, mid, ql, qr)
		v2, id2 := query(seg[u].r, mid+1, right, ql, qr)
		if v1 >= v2 {
			return v1, id1
		}
		return v2, id2
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		return
	}
	s := scanner.Text()

	if !scanner.Scan() {
		return
	}
	m_int, _ := strconv.Atoi(scanner.Text())
	m := int32(m_int)

	initSAM()
	initSeg()

	for i := int32(1); i <= m; i++ {
		scanner.Scan()
		t_i := scanner.Text()
		for j := 0; j < len(t_i); j++ {
			addChar(int32(t_i[j] - 'a'))
			root[last] = insert(1, m, i, 1)
		}
		if i < m {
			addChar(26)
		}
	}

	maxLen := int32(0)
	for i := int32(0); i < sz; i++ {
		if st[i].len > maxLen {
			maxLen = st[i].len
		}
	}
	head := make([]int32, maxLen+1)
	for i := int32(0); i <= maxLen; i++ {
		head[i] = -1
	}
	nextNode := make([]int32, sz)
	for i := int32(0); i < sz; i++ {
		l := st[i].len
		nextNode[i] = head[l]
		head[l] = i
	}

	order := make([]int32, 0, sz)
	for l := maxLen; l >= 0; l-- {
		for i := head[l]; i != -1; i = nextNode[i] {
			order = append(order, i)
		}
	}

	for i := int32(sz - 1); i >= 0; i-- {
		u := order[i]
		up[u][0] = st[u].link
		if up[u][0] == -1 {
			up[u][0] = 0
		}
		for j := 1; j < 20; j++ {
			up[u][j] = up[up[u][j-1]][j-1]
		}
	}

	for _, u := range order {
		p := st[u].link
		if p != -1 {
			root[p] = merge(root[p], root[u], 1, m)
		}
	}

	stateArr := make([]int32, len(s)+1)
	lengthArr := make([]int32, len(s)+1)

	curr := int32(0)
	currLen := int32(0)
	for i := 0; i < len(s); i++ {
		c := int32(s[i] - 'a')
		for curr != -1 && st[curr].next[c] == 0 {
			curr = st[curr].link
			if curr != -1 {
				currLen = st[curr].len
			}
		}
		if curr == -1 {
			curr = 0
			currLen = 0
		} else {
			curr = st[curr].next[c]
			currLen++
		}
		stateArr[i+1] = curr
		lengthArr[i+1] = currLen
	}

	scanner.Scan()
	q_int, _ := strconv.Atoi(scanner.Text())
	q := int(q_int)

	out := bufio.NewWriterSize(os.Stdout, 1024*1024)
	defer out.Flush()

	for i := 0; i < q; i++ {
		scanner.Scan()
		lInt, _ := strconv.Atoi(scanner.Text())
		l := int32(lInt)

		scanner.Scan()
		rInt, _ := strconv.Atoi(scanner.Text())
		r := int32(rInt)

		scanner.Scan()
		plInt, _ := strconv.Atoi(scanner.Text())
		pl := int32(plInt)

		scanner.Scan()
		prInt, _ := strconv.Atoi(scanner.Text())
		pr := int32(prInt)

		pl--
		pr--
		L := pr - pl + 1

		if lengthArr[pr+1] < L {
			fmt.Fprintf(out, "%d 0\n", l)
			continue
		}

		cState := stateArr[pr+1]
		for j := 19; j >= 0; j-- {
			if st[up[cState][j]].len >= L {
				cState = up[cState][j]
			}
		}

		maxVal, maxIdx := query(root[cState], 1, m, l, r)
		fmt.Fprintf(out, "%d %d\n", maxIdx, maxVal)
	}
}
`

type testCase struct {
	s       string
	texts   []string
	queries [][4]int
}

func genCase(rng *rand.Rand) testCase {
	slen := rng.Intn(5) + 1
	sb := make([]byte, slen)
	for i := 0; i < slen; i++ {
		sb[i] = byte('a' + rng.Intn(3))
	}
	s := string(sb)
	m := rng.Intn(3) + 1
	texts := make([]string, m)
	for i := 0; i < m; i++ {
		l := rng.Intn(5) + 1
		tb := make([]byte, l)
		for j := 0; j < l; j++ {
			tb[j] = byte('a' + rng.Intn(3))
		}
		texts[i] = string(tb)
	}
	q := rng.Intn(3) + 1
	queries := make([][4]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(m) + 1
		r := l + rng.Intn(m-l+1)
		pl := rng.Intn(slen) + 1
		pr := pl + rng.Intn(slen-pl+1)
		queries[i] = [4]int{l, r, pl, pr}
	}
	return testCase{s: s, texts: texts, queries: queries}
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(tc.s + "\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.texts)))
	for _, t := range tc.texts {
		sb.WriteString(t + "\n")
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.queries)))
	for _, q := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", q[0], q[1], q[2], q[3]))
	}
	return sb.String()
}

func buildRef() (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "refbuild")
	if err != nil {
		return "", nil, err
	}
	srcPath := filepath.Join(tmpDir, "ref.go")
	if err := os.WriteFile(srcPath, []byte(refSource), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref")
	cmd := exec.Command("go", "build", "-o", binPath, srcPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return binPath, func() { os.RemoveAll(tmpDir) }, nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, cleanup, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		input := buildInput(tc)
		exp, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexp:\n%s\n---\ngot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
