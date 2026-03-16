package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesBRaw = `ZYCIDPYOP
DPA?MNTYYAWOIXZHSDKAAAURA
NXAQHYOPRHLHVHYOJAN?RUDFU
KXWQNQ?VGJJSPQMSBPH
T?CB?LDZOEWNJ?I?VSRGAHQY??
I
ZTFLJCFFIQFVIUWJ
PPDAJMKNZGIDIXQGTN?AHAMEBXFOWQ
HUZWQOHQUAMVSZKVUNBXJEGBJCCJ
NSIEARBS?GSOF?YWTQBM
GSVNSGPDVMJQP
MJ
KZ
KNGIVD?MRL?VRPYRHCXBCEFFRGIYKTQ?ILKKD
TYWPESRYDKBNCMZE?EK
ZMCSRHSC
JSRDOIDZB?JATVACND
ZSN
FVHFXDNMZR?
WPKDGUKBAAZJXTKOMKM
KTODI
YRWPVLIFRGJGHL
CYOCUS
MJBKFKZSJHKDRTSZTCHHAZ
IRCXCAUAJYZLPPED
KCQVFFYEE?KJDWQ?TJEGERXBYK?TZVR?X
NRFBWVHIYCVOZN
ORBSNCRMAVEJI?GLDX?YTHW?U?KQSMFTCHPA
CEYT?VV
HMZNMF
TPGDNTRNDVJIHMXRAGQOS
TH
JERGIJSYI?VOZZFRL
DYGSMGUFWCXEN?GBH?PTZYMRJLVODAKQ
MZIFYC?YTA
ZWNVRJEO?IPFOQBIQDXSNCL
AFQWF
WITJGQ
KICCW?QVLOQRXB
XW?RILTXHMR
PC?V?LANKBDGMWZ?RYFQOITEXS
SDGCBAZAPKMSJGMF?
AAMEVRBSMI
U?JABRBQ?
IYDNCGAPU
V?GVOMKUI
HHBSZSFLNTWRUQBLR
WRNVCWIXTXYCIFDEBGNBBUC?QPQ
KBERBOVEMYWOAXQICIZKCJBM
KXEI
VJDN?HQRGKKQZMSPDEU?OQRX?
RA?JXFGLMQKDNLESCBJ?ZURKNJKLIKXXQQAQDE
KZKSCOIPOLXM?CSZBEBQP
IZHWSXKLZULMJOTKRQFAEIVHSEDFYNXTBZDRV
GICUSQUCCZGU?FQNAS
PWZJHGTPHNOVLRGZPXC??ING
Y
PCMTQZSSNBLOAGJWWUARD?JQX
RUSRJQNR?QNTUSJOJEQOS
FIUANXVSB
JVYVACCAMIOIZZLUXPYKMOZD
OFJTRW?QN?BVC?KMUPXI?AH?GEPMWNC
PUZRZ?FWP?CJDSMOLH??Y?QXBTU?NGWHNEE
YAQ?ENZ?TM?OGVPAFYI?K?DWC?BS
MWDXLCURLRRZXQVSA
VEECSEV?GP?FYCZEKGSX?QB??JMULIVTNAODRUEM
QEBKVDQFRUUPKYWD
GMU?FMSZTKRWLY?HXPBIVEJQMOF?UADNPI?EEA
QN?TBADH?GU?M?R?KX?FVSYE?IPMWMZ
PAOCJYU??WGMZRLHNSD??KXBTI
DCVXBN?OGFTQGQMQLGHLVS
BOB?T
JPBSQC?SMCMZ
UJMILPBRPANJSXKZETSRICTZZYLNMQZAS
SQADKDZAQSFX?ROG?PWV?JMY?CBLTI?SDMKQN?
SQDXIOGZ
OKTX
UAAPBFIRBAHYCQ?FBQGGOJH
KMUCGTFGVTJSNTPLAPADVUSVTN?WSKKC
QNV?EUPKG?FMQH?DLYO?ZAJXIRT
QOZCGNHTBTHUHHWMMGTEXJXXL
V
FVEALNRKZQPKTDSUJZRV
?JYCUPDQHTXUXINLZH
QQQF
BCGAVBNXW
AB
ZA?R?WTZIYUNXPLAVSOCE?BDJF?GASFQUCY
ZNIRJKY?LNOLLKMP
EJFJSGXNS?HUJMIZTW?LQBR??CKVEDOER
PFQLGNZCIGHYEEYGAFPLFBZLCTHVW
UUGTKFSW??VWAG
RBBLPRLEPCQKVXSVJTKZS
NCICVU
FKHKI??IJPNAJFUJBDNNT
YUXSPSJTIVFKEL
QXSWGMO
WHBXUHCXCB
SPWKQZFSWPMAMRXR?XOFSSLB?XL?LOHWUV
COYLGDAGYL?WMKUS??P?VONRHJQTXFEKIEF
`

func possible(s string) bool {
	n := len(s)
	if n < 26 {
		return false
	}
	arr := []rune(s)
	for i := 0; i+26 <= n; i++ {
		var freq [26]int
		q := 0
		ok := true
		for j := i; j < i+26; j++ {
			ch := arr[j]
			if ch == '?' {
				q++
			} else {
				idx := ch - 'A'
				if idx < 0 || idx >= 26 {
					ok = false
					break
				}
				if freq[idx] > 0 {
					ok = false
					break
				}
				freq[idx]++
			}
		}
		if !ok {
			continue
		}
		used := 0
		for k := 0; k < 26; k++ {
			if freq[k] > 0 {
				used++
			}
		}
		if used+q == 26 {
			return true
		}
	}
	return false
}

func validOutput(orig, cand string) bool {
	if cand == "-1" {
		return false
	}
	if len(orig) != len(cand) {
		return false
	}
	for i := 0; i < len(orig); i++ {
		if orig[i] != '?' && orig[i] != cand[i] {
			return false
		}
		if cand[i] < 'A' || cand[i] > 'Z' {
			return false
		}
	}
	n := len(cand)
	if n < 26 {
		return false
	}
	for i := 0; i+26 <= n; i++ {
		var freq [26]int
		ok := true
		for j := i; j < i+26; j++ {
			idx := cand[j] - 'A'
			if idx < 0 || idx >= 26 {
				ok = false
				break
			}
			if freq[idx] > 0 {
				ok = false
				break
			}
			freq[idx]++
		}
		if !ok {
			continue
		}
		valid := true
		for k := 0; k < 26; k++ {
			if freq[k] != 1 {
				valid = false
				break
			}
		}
		if valid {
			return true
		}
	}
	return false
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	idx := 0
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" {
			continue
		}
		idx++
		possibleFlag := possible(s)
		out, err := run(bin, s+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if !possibleFlag {
			if out != "-1" {
				fmt.Printf("test %d failed: expected -1 got %s\n", idx, out)
				os.Exit(1)
			}
			continue
		}
		if !validOutput(s, out) {
			fmt.Printf("test %d failed: invalid output %q\n", idx, out)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
