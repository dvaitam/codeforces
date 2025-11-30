package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded test cases from testcasesA.txt so the verifier is self-contained.
const testcasesA = `biqpm z_jpls,
ejeydtzi rwzt,ejdlala.
miao.prdlnktu gr
.qibz r acx_mw_
miao.p_.khxkw.cg shhzezro
miao.k,q
jrjw
miao.k_ rgztrsjoctzmksh
miao.g_fbtvlala.
miao.ccvye,e b_c,w _rlala.
iqzh.g,vs_nsiopvu
miao.ktdlala.
uk.ghaxidwhlzfkn_bd
b_su  rtlala.
miao.d
_sdmcl_dbtagfwdpgxzblala.
miao.ntd_ichcujlnfbqobt
miao.gil,xp_, sfwv
zvlala.
kqidtlala.
apvn,s
miao.m_vierwaoxcklala.
miao.iehypltjvls,u.tewj
_ucatgwkfhhuom
m.
fo
miao.wfoq,p rt
miao.pk
miao.bz_,zngruc_xeamlala.
agawya_vqtd
tug.lala.
miao.wfdp. muclala.
miao.ozzd.ieu
miao..ei.lala.
miao.bg
miao.kls .b.xwtupwu,ounlala.
miao.fgmsjaeeikkzlwckytlala.
miao.if
jlmrejdpxh bjf.qxcjlala.
jnddr ppk_z
miao.pdwplala.
jk
usmzuclala.
miao.gxh
miao.admrqjo pzswv
clhisyfngldc_.lala.
miao.oyvgdpmigubzgtedg
lr_edtpesmuvn
_ppuv.grt h
miao.wx,k_kbqe.i
sjwwzpczcq. blala.
miao.ebjay.ok
.uolqlala.
qbscvzzqytcxn,ygj
.zmtsh..zav,axfj
miao.cp.i_jynmzmlala.
miao.u ehjx
miao.bplala.
p,.tw
miao.zlnbt
obdpyeabtteuk
miao.ulgmzyypd,btwotu k
miao.jz.e mzj x.vzdq.zgbz
lygolzucb, bpi,
miao.vss,ghcy_,,uyqwqn
miao.en,slala.
dn
miao.nyex
miao.nvnap .kxiclcdl
miao.lfa_h_lct,eg
miao.vvx,dxa
miao.wa th.efodpllala.
miao.aglkplala.
miao.,rukfscdrlala.
miao.m,eezhlala.
miao.hhyfjlnvb.etamcwclala.
jrnx esnjulcholala.
miao.bmnanx,kogljpcfzdlala.
miao.rtwelala.
fynnfhok qelo
miao.ygja_wotoagjdyuj
wypcvpyhrymiualala.
miao.,vbaimq,s
dxilj.yv gtcbcz
miao.rkdq.h, yflala.
miao..jjqesqugrdnur
jolsuefdw
miao.msoervjlupxng
wqkpubojexpb,tga
miao.m,a.qcv.cvx vmallala.
miao.tai.
miao.eysjgdnolala.
miao.fknu,vneo weqlala.
miao.g follala.
zpmxhzgogswb,m
miao.uc.flbxu
tjtcw.qylala.
miao.nobuwqvur xlala.
opiwpgkibbbflajuaec
tmr hogk`

func runCandidate(bin string, input string) (string, error) {
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

// Embedded solution from 312A.go.
func classify(sentence string) string {
	hasMiao := strings.HasPrefix(sentence, "miao.")
	hasLala := strings.HasSuffix(sentence, "lala.")
	switch {
	case hasMiao && hasLala:
		return "OMG>.< I don't know!"
	case hasLala:
		return "Freda's"
	case hasMiao:
		return "Rainbow's"
	default:
		return "OMG>.< I don't know!"
	}
}

func loadCases() []string {
	// Preserve blank lines inside the file; drop only a trailing empty from a final newline.
	lines := strings.Split(testcasesA, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases := loadCases()
	for idx, sentence := range cases {
		input := fmt.Sprintf("1\n%s\n", sentence)
		want := classify(sentence)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
