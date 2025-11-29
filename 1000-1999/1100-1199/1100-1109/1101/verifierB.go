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

const testcasesRaw = `100
vwbfociwfoq:p|rtya:bpkj[o
z[
zngruc[xeamvnkagawya[vqtdgdtu
]jiwfdp
|mucai|ozzd]ieuq[uuld]ei]abb
virkls|
b]xwtupwu:ounl]rfgmsjaeeikkz
wckytbbifesj
mrejdpxh|bjf
qxcjm[kjnddr|ppk[zzkdpdwpnbj
xv:e|fusmzu
czc
xhbmadm
qjo|pzswvgnclhisyf
gldc[]]wa:qoyv
dpmigub
gtedgomlr[edtpesmuvn:qpv|:
[ppuv]grt|h
k
x:k[kbqe]itze[msjwwzpcz
q]|
ch
bjay]
k]fze]uolqm:qqb
cvzzqytcxn:ygjr:tn[
]zmtsh]]zav:axfj
sikcp]i[jynmzmbfu
ehjx[kbbpnep:]twcvwezlnbtomobd
yeabtteukdwrulgm
yypd:btwotu|kudvwtjz]e|mzj
x]vzdq]zgbzmolygolzucb:|bpi:aq
ss:ghcy[::uyqwqnqjden:
n|c|dncdnyexazonvna
|]kxiclcdlwallfa
h[lct:egagvvx:dxajlwa|th]ef
dplwieaglkpjj|:
ukfscdrsjfm:eezhkq
hyfjlnvb
etamcwcenjrnx|esnjulchouluqb
nanx:kogljpcf
didrtwezwo|mfynnfhok|qelou
cpygja[wotoagjdyujrte
wypcvpyhrymiua
i:vb
i
q:swmodxilj]y
|gtcbczijrkdq]h:|yfcn]
jqesqugrdn
rmxy:zijolsuefdwdmmso
rvjlu
xngppwqkpubojexp
:t
alpm:a]
cv]cvx|vmalbdtai]
wjx:heysjgdnowkmfknu:
neo|weqkeg|folzmnzpmxh
gogswb:mbhuc]flbxuvfhtjtcw
qyjy::lnobuwqvur|x|nsopiwpgk
bbbflajua
cznvh
mr|hogktdt:czkkro:ki
q
gl
g]q
g[givxxjjq]m
pl]whbj|r
aop
ob|znpoodcchd[yen:|gotc[
ry:[mbfhphe[]i
lkndr:jtrzgwjyoqtoruiihadtzwdf
nhgj:xvaxrqn]bdmuidxslhv
wrvjhx[:hcqjvkhlupjsfez
r
klsuaze:mefqceygz
ypsywghxe[hymltseup:|dt[aqtlpo
ahrufv:pzx
rkw]cietmwg[kzjm
gb
xxhk]ovxvvh
lvfjalsrb
uelapubahbahukc[blvne:|g
neljfukxzxnman:
rrzxvwoyb
|dnmfaqet]vq[xweckh
[[fhaz:fxz[vrfwcn|]tdtuowet
bikzxxmau|:bpcljveoh
lfxymkizpmajq:jrp
yr
ri|vbomxdmlpbaixbiv
swy|jvygyqqkm[igdskz|h
vxrv:lf|]:ek]x[as[b
elljujkpzmtnfazesbo
k|axp
v[vyixtgcr|ni[fqfcvufsdquzrtmy
ijjany[wiirqrk
gwznze]ayqe
vw|zsmlobr|nutzzyh:alqfvgu
uwpaxxhs:hif
|ncsohxwo]q:zzwdgfo|c|num
in:yyltkc
`

// solve mirrors 1101B.go logic for a single string input.
func solve(s string) int {
	n := len(s)
	st1 := -1
	for i := 0; i < n; i++ {
		if s[i] == '[' {
			st1 = i
			break
		}
	}
	if st1 < 0 {
		return -1
	}
	st2 := -1
	for i := st1 + 1; i < n; i++ {
		if s[i] == ':' {
			st2 = i
			break
		}
	}
	if st2 < 0 {
		return -1
	}
	en1 := -1
	for i := n - 1; i >= 0; i-- {
		if s[i] == ']' {
			en1 = i
			break
		}
	}
	if en1 < 0 || en1 <= st2 {
		return -1
	}
	en2 := -1
	for i := en1 - 1; i >= 0; i-- {
		if s[i] == ':' {
			en2 = i
			break
		}
	}
	if en2 < 0 || en2 <= st2 {
		return -1
	}
	cnt := 4
	for i := st2 + 1; i < en2; i++ {
		if s[i] == '|' {
			cnt++
		}
	}
	return cnt
}

type testCase struct {
	s string
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		return nil, fmt.Errorf("invalid test file")
	}
	tLine := strings.TrimSpace(scan.Text())
	t, err := strconv.Atoi(tLine)
	if err != nil {
		return nil, err
	}
	cases := make([]testCase, 0, t)
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if line == "" {
			continue
		}
		cases = append(cases, testCase{s: line})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	input := tc.s + "\n"
	expected := solve(strings.TrimSpace(tc.s))
	gotStr, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	got, err := strconv.Atoi(strings.TrimSpace(gotStr))
	if err != nil {
		return fmt.Errorf("invalid output %q", gotStr)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) > 2 {
		bin = os.Args[2]
	}
	cases, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
