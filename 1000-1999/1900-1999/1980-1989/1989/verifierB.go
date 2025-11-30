package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB = `idp yo
zgdpamnt yyawoix
kaaa ur
v gnxaqhy
rhlhvhyo janrudfu
xkxwq nq
jspq msbph
vflrwyv xlcovqd
lpx apbjwts
fqhaygr rhm
oivrtx amzxqz
nbp lsrg
lnlarrt ztkotazh
czr zib
ao ay
idztf ljcf
qfv iuwjo
pdajmk nzgidixq
aham ebxfowq
uzwqohq uamv
zlwkfimj tkkgajv
jqwbbve hfv
ldfaihfi jmnwy
ylxlzg tcbmzl
gxj hlbm
kub qaagn
sc dgjwwevu
jmuwiszr fnmqz
eac zspyibv
nqtve lqmqhvu
cfo ybm
uoh mjsavy
ixglq jkrjnpur
dcvbfvw ulpy
eeo yq
engrx ptzslum
qylw rw
fz yas
ukbm pqkzybz
ochwo uhpi
ujejsps pnvzwqcl
ov gzjr
brk dmve
oq yvfexorl
bh wmlxd
qncl nygm
jdsks kkrz
asozgnk nwhd
mpdtrsx ggmgpw
aaepv iblwf
zmbmt tlm
afp qw
os mhdnecc
mdjbg vpmiza
gnw uushn
xdg trxs
upyy qyihv
iqtyfo mxzo
rjhxgfv fq
gnv fuqn
lwfjezsb cnegq
knau ckcelwi
nauh mznf
fx rlebr
vfosyhh uzlhkw
ux nlfl
da yflf
cj jfepeo
sxhks hezb
tcbhuhje wlwmvrqs
qmof up
dzfgez arvdbfe
db xlmrdp
yiotb sy
lcjsu yucagvkt
kye vzot
ckhrpl hrd
zh ghwl
fgmgloy tjqnj
rn vzyfk
ycwko iovgnjmh
qz gicw
rgicaxlp rmnwuwpv
nlcqz qaeqadu
jbc fchyyq
eqc illmucmo
zxey ingdoqz
hkgv lrehrhf
ngzj fufvom
wfj wlkrge
ukllj gx
hnps xdnjlz
xtvpubu jsridntr
la dzqeql
nqvhhb hleul
xwr fg
eos crgde
wtezlgrl hkqb
ewwuea pcng
gyfi cjy
bhbsp rs
saafk njbhapj
gkabnf yzixfrf
mukqow evujrd
sknw rx
vxdevy exjv
quzldq wdyhvso
vhnhogc delwyfr
irclhyg lzbuwyhx
jkrjp zzcl`

func expected(a, b string) int {
	as := []byte(a)
	bs := []byte(b)
	n := len(as)
	m := len(bs)
	ans := n + m
	for i := 0; i < m; i++ {
		j := i
		for _, ch := range as {
			if j < m && ch == bs[j] {
				j++
			}
		}
		if val := n + m - (j - i); val < ans {
			ans = val
		}
	}
	return ans
}

type testCase struct {
	a string
	b string
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesB), "\n")
	tests := make([]testCase, len(lines))
	for i, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("bad test line %d", i+1)
		}
		tests[i] = testCase{a: parts[0], b: parts[1]}
	}
	return tests, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(tc.a)
		sb.WriteByte('\n')
		sb.WriteString(tc.b)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}
	input := buildInput(tests)
	output, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	got := strings.Fields(output)
	if len(got) != len(tests) {
		fmt.Printf("expected %d outputs, got %d\n", len(tests), len(got))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := strconv.Itoa(expected(tc.a, tc.b))
		if got[i] != want {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, want, got[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
