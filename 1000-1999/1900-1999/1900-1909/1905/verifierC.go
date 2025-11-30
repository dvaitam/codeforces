package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcases = `c
lf
itgtb
vfnumzxqlr
qibalokm
qfrfhha
kfe
qlqvrfozn
ylzsllofy
wxouqhp
pqqzl
olsxrxop
kwft
ypjjz
rqqutsnjx
pqlv
czkxagxdbs
i
hvdyqeihgb
wybbllf
vacd
ab
l
efxfq
m
bzhebaltux
jk
ajorytxb
ymtwe
hcvvkdao
qsy
pqkekii
nuawrevbib
ffd
uhqwbhhw
cicshtzz
wlivniqyae
m
fdqxchd
af
dgaq
ojrumgvy
xznn
assbnqsfd
laqdtljw
javndd
gyvaz
bnupogst
aj
l
xchyp
dslm
eylmdidd
tk
gwdatvp
x
jloezlip
pxxznpvjm
fpti
nwvwcsxsd
lf
eznczcvzu
e
mhwvv
ofqjde
yndkqhwqi
fow
mlzy
xeooxaztmx
mqb
imiwxnwu
lrkwxvcy
rtgm
uakoqwo
dam
xstm
dmry
ixss
pzte
t
piqsfow
ycla
rvvcyspv
oiqoac
ylfyyzmivu
bfp
ovjeajr
albrmsog
puepw
wjcikjkuz
uumqc
ugmtqezqu
jb
orhq
bddvz
lgklcko
fpojoe
ugikfdhp
yvlf
ezehiz`

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

func solve(s string) int {
	n := len(s)
	if sort.SliceIsSorted([]byte(s), func(i, j int) bool { return s[i] < s[j] }) {
		return 0
	}

	// Collect lexicographically largest non-increasing subsequence from right to left.
	var subIdx []int
	suffixMax := byte(0)
	for i := n - 1; i >= 0; i-- {
		if s[i] >= suffixMax {
			subIdx = append(subIdx, i)
			suffixMax = s[i]
		}
	}
	// Reverse to maintain increasing index order.
	for i, j := 0, len(subIdx)-1; i < j; i, j = i+1, j-1 {
		subIdx[i], subIdx[j] = subIdx[j], subIdx[i]
	}

	m := len(subIdx)
	orig := []byte(s)
	subVals := make([]byte, m)
	for i, idx := range subIdx {
		subVals[i] = orig[idx]
	}

	// Apply cyclic right shift on selected indices.
	for i, idx := range subIdx {
		prev := (i - 1 + m) % m
		orig[idx] = subVals[prev]
	}

	if !sort.SliceIsSorted(orig, func(i, j int) bool { return orig[i] < orig[j] }) {
		return -1
	}

	// Cost is length minus count of maximal char in the chosen subsequence.
	maxVal := subVals[0]
	countMax := 0
	for _, v := range subVals {
		if v == maxVal {
			countMax++
		}
	}
	return m - countMax
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := strings.Split(testcases, "\n")
	count := 0
	for idx, line := range cases {
		s := strings.TrimSpace(line)
		if s == "" {
			continue
		}
		count++
		want := solve(s)
		input := fmt.Sprintf("1\n%d\n%s\n", len(s), s)
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(gotStr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d output parse error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\nInput: %s\n", idx+1, want, got, s)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", count)
}
