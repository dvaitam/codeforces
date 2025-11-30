package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `hjdxmpeccamrjzybhqrl
100
yfdig
u
zigfj
uxlct
vmqhfh
icrjajsw
yqgnn
jnofhjizbc
o
qrupw
evgcng
iflnxsku
gkdbwhiys
hdkfjoablw
jx
akjkey
tvcjtgo
eimtf
salbof
zzljsd
gngdbbbx
tve
brpshkbdqj
ugpghon
bhnohung
gbbiihqg
niek
k
dsmuuzwxbp
cngsfkj
zuknqguz
kmpci
bmte
vbfwu
spxmmgza
fati
mz
hrbgfvt
rvj
jszqu
ogqvhdv
plg
ubqbku
rzduhudwuq
bkahd
wanubxy
ky
kxyfhoauzv
jteelvkaz
rxtto
xrkqa
qv
ldbewaug
vuxmltikiz
uejcq
kuvumw
otmvoanmj
ybjmmo
vybwwknyco
eneqzefwi
xxjhxrmgwn
xlgelg
ih
kmrvjftd
qmde
lfgxbvz
txupibfc
ejfvb
zjuttp
uplcs
ooryxgw
suoruklzps
ojde
ivokz
zpvoklazg
emh
igsqb
oxcxpxyqs
ttijsihw
amifi
cdlyzd
tnrxwccdp
bezqvqjxh
iczlejvbex
hn
rbh
tpe
uootgxhxpb
maymew
bsqliqgiib
oob
idjuawzsn
dtzsvqym
onvtxcwuq
tfnoc
afcpxp
hzcyea
xahpotmna
vrskjbploh
ethbewui`

type testCase struct {
	s       string
	queries []string
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func prefix(s string) []int {
	n := len(s)
	pi := make([]int, n)
	for i := 1; i < n; i++ {
		j := pi[i-1]
		for j > 0 && s[i] != s[j] {
			j = pi[j-1]
		}
		if s[i] == s[j] {
			j++
		}
		pi[i] = j
	}
	return pi
}

func solveAll(s string, queries []string) []string {
	sb := []byte(s)
	piS := prefix(s)
	n := len(sb)
	results := make([]string, len(queries))
	for idx, t := range queries {
		tb := []byte(t)
		piT := make([]int, len(tb))
		j := piS[n-1]
		get := func(pos int) byte {
			if pos < n {
				return sb[pos]
			}
			return tb[pos-n]
		}
		for i := 0; i < len(tb); i++ {
			c := tb[i]
			for j > 0 && get(j) != c {
				if j <= n {
					j = piS[j-1]
				} else {
					j = piT[j-n-1]
				}
			}
			if get(j) == c {
				j++
			}
			piT[i] = j
		}
		strs := make([]string, len(piT))
		for i, v := range piT {
			strs[i] = strconv.Itoa(v)
		}
		results[idx] = strings.Join(strs, " ")
	}
	return results
}

func parseTestcases() (testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	if len(lines) < 2 {
		return testCase{}, fmt.Errorf("not enough lines")
	}
	s := strings.TrimSpace(lines[0])
	q, err := strconv.Atoi(strings.TrimSpace(lines[1]))
	if err != nil {
		return testCase{}, fmt.Errorf("parse q: %v", err)
	}
	if len(lines)-2 != q {
		return testCase{}, fmt.Errorf("expected %d queries got %d", q, len(lines)-2)
	}
	queries := make([]string, q)
	for i := 0; i < q; i++ {
		queries[i] = strings.TrimSpace(lines[2+i])
	}
	return testCase{s: s, queries: queries}, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	tc, err := parseTestcases()
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

	expected := solveAll(tc.s, tc.queries)

	var input strings.Builder
	input.WriteString(tc.s)
	input.WriteByte('\n')
	fmt.Fprintf(&input, "%d\n", len(tc.queries))
	for _, t := range tc.queries {
		input.WriteString(t)
		input.WriteByte('\n')
	}

	got, err := runExe(bin, input.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outLines := strings.Split(strings.TrimSpace(got), "\n")
	if len(outLines) != len(expected) {
		fmt.Println("wrong number of output lines")
		os.Exit(1)
	}
	for i := 0; i < len(expected); i++ {
		if strings.TrimSpace(outLines[i]) != expected[i] {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, expected[i], strings.TrimSpace(outLines[i]))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
