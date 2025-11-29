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

type testcase struct {
	n int32
	s string
}

const testcasesRaw = `100
8
jdxmpecc
1
m
18
jzybhqrliyfdigauzi
9
gfjjuxlct
11
vmqhfhpicrj
1
j
19
wjyqgnntjnofhjizbcb
15
uiqrupwkevgcngu
15
iflnxskurgkdbwh
9
ysthdkfjo
1
b
12
wcjxvkakjkey
14
tvcjtgojeimtfk
19
alboflzzljsdogngdbb
2
xf
20
vetbrpshkbdqjynugpgh
15
npbhnohungpgbbi
9
hqgyhniek
2
ks
4
smuu
2
pm
3
ngs
6
kjvpzu
11
nqguzvzikmp
3
iuv
7
bmteyiv
2
fw
15
spxmmgzagfatidm
13
hrbgfvtkryzpq
15
acbwtdprityeblc
17
ajlccromgzjmhypmd
3
dtz
12
qnnwyocugujp
14
dzrflffwekpkir
1
w
6
aujdrd
16
zwtpqcqhnjlhyfua
2
tk
18
ozsjqzottomeiytlvk
5
nctev
20
fjlgszlvtccmufkulkfj
1
t
1
q
3
zlz
4
ffsp
19
cydfupvxhtzvjwmthpw
8
jlhxkruq
15
zmqmkjonsaifrow
18
tlmmtaeqcoultzjcyi
16
huptycehcjebfmzs
20
rwriaqhfydgbkcdibujt
6
xeunew
3
rlw
1
x
18
eqnegjpqcmffiqmrvj
13
kfmbnaivajzec
6
dtahhr
1
p
18
foxwvwmkfqsgdptlju
14
txbugiszvsjzpu
11
zchkdvwbtkq
16
lcfwbpqrtwhbgvck
20
ewjdqqachtitaxbapweg
12
hxxljmutmbfn
16
cyqqjrwjzgwwnkfg
1
p
14
ikynunkytgiiqc
1
m
9
tjvuetykc
6
tcwgiv
7
upyduka
15
sgfchqddzkjssve
3
vfw
1
y
6
vgibsk
11
xytlffygbbf
2
vi
16
dgzemptrprklnxmb
2
fe
12
hjhlsgtvlmlr
4
hglf
18
gaadgjwxkrpbjzxzmf
5
mtoll
1
o
5
oaxzx
14
gdkjtrbcpfvqbf
20
htktaxvamnhzsetbeftj
8
nxkrulfx
12
vrvslitkrbre
12
sgppdcutgfer
5
tbrlc
2
zt
5
rfitb
10
pigupyjqwy
2
ej
14
balqtlcgnvcznz
2
jg
7
fgwxcca
1
j
17
qqnccqukvfugqwrbw`

var testcases = mustParseTestcases(testcasesRaw)

func mustParseTestcases(raw string) []testcase {
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
	cases := make([]testcase, 0, t)
	for i := 0; i < t; i++ {
		n := nextInt()
		if !scanner.Scan() {
			panic(fmt.Sprintf("missing string for case %d", i+1))
		}
		s := scanner.Text()
		if len(s) != n {
			panic(fmt.Sprintf("length mismatch on case %d: got %d want %d", i+1, len(s), n))
		}
		cases = append(cases, testcase{n: int32(n), s: s})
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("scanner error: %v", err))
	}
	if len(cases) == 0 {
		panic("no testcases parsed")
	}
	return cases
}

// solve replicates 1063F.go logic for a single testcase.
func solve(tc testcase) int32 {
	n := tc.n
	s := []byte(tc.s)
	for i, j := 0, int(n)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	maxStates := int(2*n) + 5
	trans := make([][26]int32, maxStates)
	fail := make([]int32, maxStates)
	lengthArr := make([]int32, maxStates)
	pos := make([]int32, n+1)
	lstpos := make([]int32, n+1)
	g := make([]int32, maxStates)
	fArr := make([]int32, n+1)

	const root int32 = 1
	var tot int32 = 1
	pos[0] = root
	fail[root] = 0
	lengthArr[root] = 0

	insertSAM := func(cnt, c int32) {
		cur := tot + 1
		tot = cur
		lst := pos[cnt-1]
		pos[cnt] = cur
		lengthArr[cur] = cnt

		for lst != 0 && trans[lst][c] == 0 {
			trans[lst][c] = cur
			lst = fail[lst]
		}
		if lst == 0 {
			fail[cur] = root
			return
		}
		p := trans[lst][c]
		if lengthArr[p] != lengthArr[lst]+1 {
			q := tot + 1
			tot = q
			lengthArr[q] = lengthArr[lst] + 1
			fail[q] = fail[p]
			trans[q] = trans[p]
			for lst != 0 && trans[lst][c] == p {
				trans[lst][c] = q
				lst = fail[lst]
			}
			fail[p] = q
			fail[cur] = q
		} else {
			fail[cur] = p
		}
	}

	moveSAM := func(u *int32, length int32) {
		for *u != root && lengthArr[fail[*u]] >= length {
			*u = fail[*u]
		}
	}

	for i := int32(1); i <= n; i++ {
		insertSAM(i, int32(s[i-1]-'a'))
	}

	lst := root
	for i := int32(1); i <= n; i++ {
		c := int32(s[i-1] - 'a')
		cur := trans[lst][c]
		curfa := cur
		fArr[i] = fArr[i-1] + 1
		moveSAM(&curfa, fArr[i]-1)
		for fArr[i] != 1 {
			if g[lst] >= fArr[i]-1 || g[curfa] >= fArr[i]-1 {
				break
			}
			fArr[i]--
			moveSAM(&cur, fArr[i])
			moveSAM(&lst, fArr[i]-1)
			moveSAM(&curfa, fArr[i]-1)
			p := i - fArr[i]
			u := lstpos[p]
			if g[u] < fArr[p] {
				g[u] = fArr[p]
			}
			u = fail[u]
			for u != 0 && g[u] < lengthArr[u] {
				g[u] = lengthArr[u]
				u = fail[u]
			}
		}
		lst = cur
		lstpos[i] = cur
	}

	var ans int32
	for i := int32(1); i <= n; i++ {
		if fArr[i] > ans {
			ans = fArr[i]
		}
	}
	return ans
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
	return out.String(), nil
}

func parseCandidateOutput(out string) (int32, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return 0, fmt.Errorf("no output")
	}
	v, err := strconv.ParseInt(scanner.Text(), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse output: %v", err)
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scanner error: %v", err)
	}
	return int32(v), nil
}

func checkCase(bin string, idx int, tc testcase) error {
	input := fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
	expected := solve(tc)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	got, err := parseCandidateOutput(out)
	if err != nil {
		return err
	}
	if got != expected {
		return fmt.Errorf("case %d: expected %d got %d", idx+1, expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range testcases {
		if err := checkCase(bin, i, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
