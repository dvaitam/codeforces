package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesC = `100
d
xzx
itgtb
num
lroqibalo
nqfrfhh
feqqef
vrfoznxqy
llofyyfoll
uqhpphqu
qqzllzqq
olsxrxop
wftiyp
rqqutuqqr
xgpql
czkxagxdbs
shvdy
hgbnw
m
vaav
aa
l
xfx
n
ebal
jk
oryro
j
wephcvvkda
qsq
qkekiitn
evbibeffd
qwwq
hwoc
htzzttzzth
vniqy
m
qq
ddaf
ga
ojrumgvy
nqassbn
fdvzplaqdt
javvaj
jj
z
upogopu
cajaljxchy
ss
eylmmlye
ddctk
wdat
y
loezl
pxxznzxxp
mhfpt
wvwcwvw
dclfrezncz
e
hwvvkof
er
dkqhwqi
owo
zyssyz
oxaztmxf
p
wxnxw
lrkwxvcy
rtgm
uakoqwo
am
tmgdmryzgi
sgpzteatvn
sfowgwofs
la
vvcyspvko
oactytcao
ivuzuvi
bfp
vjeaaejv
albrrbla
ogvjpuepwr
cikjk
uumqc
mtqe
jb
rhqiiqhr
dd
lgklcko
ojoeeojo
ugikfdhp
fleelf
izru
ixttxi
xmywwu
rtuvc
pfilifp
ff
eadldae
xx`

// Embedded solver from 692C.go.
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func solve(s string) string {
	if s == reverse(s) {
		return "YES"
	}
	return "NO"
}

type testCase struct {
	s string
}

func parseCases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesC), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	var t int
	if _, err := fmt.Sscan(lines[0], &t); err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if 1+i >= len(lines) {
			return nil, fmt.Errorf("missing case %d", i+1)
		}
		cases = append(cases, testCase{s: lines[1+i]})
	}
	return cases, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Println("failed to load testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc.s)
		input := tc.s + "\n"
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
