package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded test cases from testcasesB.txt. Each line is a single string input.
const testcasesRaw = `lygolzuc.biaqvssgh.yqjdensn.dcdnyexa.nv
xiclcd.walfahlc.egavvxd.ajlt.e
lwieaglk.jj
ukfscdrsjfmeezhkqh
lnvbe.amcen.rn
snjul
uqbmna.xkgljpcfzd.drzwo.fynfhokqe.ou
ygja.wotoagjdyuj
nwypc
rymiuadi
ai
dxiljyvg.cbcjrkdq.ycnj.qerdnu.mxy
jolsuefdw
ervjlupx.gpwqkpuboj.xbtgalpma.cvcalbdtai.wjx
ysjgd
fknuvne.weegfolz.nzmxhzgogs.bmb
lbx.vfh
cwqyjylnobuwqvurxnso
kibb.fajuaec.vhrhogktd.czkrokiaq.g
ggivxx.jqiplwhbj.cao
z.podcchdyen.onr.bfphei.kn
wjyo.torihadt.dfxhgjxvax.qnb
xs.hvhxhcq.vklupj.fez
suazem.fey.zsywghxeh.lt
dt.aqtlpojahrufv
rkwcietmwg.kzjmb
kovx.vhivfjals.bxulap.bah
cblvne.oeljfukx.nma
ybsdnmfa.etvckh.hz.x
fwcntdtuowettbikzx
u
pc
hqlfxymk.zpajqjrpb.sri
om
lpbaixbivvswy
yqqkm.i
svxr.lfexasbse.ljkpzmt.fa
boekaxpv.vyixtgcrni
qfcvuf
uzrtmynijjany.wii
kkgwznzeayqezvwzsm
utzzyha.qf
uwpaxxhshify
soh
qzzwdgfocnumiin
tkcjapayigym
s.lsswi.amqebcslla.gvr.bk
k
nvete.jqf.e
x.
wxwumflzs.ccrf.igz
qoeyo.ebuuqb.cgthzoqfwk.epy
c.qka.ycdntlsok.mq
kaf.eak
bnubwjz
lcn.o
szonf.ogmcld.aafmtwu.k
cbrmzzia.qdcldpbeq.zjbm.frffy.u
kapuvmbhhujkfhlhfnol
msyaf
v
aak.abbst.y
n.svwhelqg.mcesvldnnh.mhhupmsci.ql
typhibtkzmudrbewmayn
ndwotoffkpnfs
yqdlleulyp
gi
mubjr.bvimxglee.tfansmo.u
ynr
fegfh
ltxjwzkv.zyifkuqk.wem
jhmlmpqjnn.dveeast
dwyuwgtutqdiwxtfm
d.xp
swivpz.fyq.aqnt.gowiasemfo.bmz
hqobpt.irumpifhr.fjocc.msrciphdj.l
vb.ve
abmrwyt
rlykdvah.p
ihaplqkccj
xlmyehjg
vzljmte.ydmlqphu
nilmyywjdpjd
xfkp.dmoqosth.mqjhkqvacpk.hn
ybncigxkfdfzwlahbam
d.ttg
aqtn.rhhn.pzn.mt.as
vouetq.iwd.awtyv.m
nqqd.dwpqo.foenimcs.kho
aoxarmy.hnippe.olupte.vcgquvdb.k
mvgnpgfm.zyvqgxk.q
sneuwtpvsgtobhpt.krx
e.o
hvvdyqnuujvkpvgflr.l
ovkocehydyemo
lwbamgc
sznrguavequpjmrrcct
yagntyljcxayzhta
leoidsoi.hmmxqkh.
bek`

type testCase struct {
	line string
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseTestcases() []testCase {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cases = append(cases, testCase{line: line})
	}
	return cases
}

// Embedded solver logic from 174B.go.
func solveCase(s string) (string, error) {
	n := len(s)
	dots := make([]int, 0)
	for i := 0; i < n; i++ {
		if s[i] == '.' {
			dots = append(dots, i)
		}
	}
	m := len(dots)
	if m == 0 {
		return "NO", nil
	}
	if dots[0] < 1 || dots[0] > 8 {
		return "NO", nil
	}
	segLen := make([]int, m)
	for k := 0; k < m-1; k++ {
		delta := dots[k+1] - dots[k]
		low := delta - 9
		if low < 1 {
			low = 1
		}
		high := delta - 2
		if high > 3 {
			high = 3
		}
		if low > high {
			return "NO", nil
		}
		found := false
		for l := low; l <= high; l++ {
			ok := true
			for j := 1; j <= l; j++ {
				if dots[k]+j >= n || s[dots[k]+j] == '.' {
					ok = false
					break
				}
			}
			if ok {
				segLen[k] = l
				found = true
				break
			}
		}
		if !found {
			return "NO", nil
		}
	}
	lastL := n - 1 - dots[m-1]
	if lastL < 1 || lastL > 3 {
		return "NO", nil
	}
	for j := 1; j <= lastL; j++ {
		if dots[m-1]+j >= n || s[dots[m-1]+j] == '.' {
			return "NO", nil
		}
	}
	segLen[m-1] = lastL
	res := []string{"YES"}
	start := 0
	for k := 0; k < m; k++ {
		end := dots[k] + segLen[k]
		res = append(res, s[start:end+1])
		start = end + 1
	}
	return strings.Join(res, "\n"), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	cases := parseTestcases()

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		expected, err := solveCase(tc.line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(bin, tc.line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
