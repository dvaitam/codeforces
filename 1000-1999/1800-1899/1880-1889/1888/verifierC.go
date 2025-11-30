package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesRaw = `kqidto
apvnsq
ulmvierwao
kxb
iehypltjvlsutewjmx
ucatgwkfhhuomw
nbmwsnyvwbfociwfoqp
tyabpkjobzzngrucxe
m
kagawyavqtdgdt
jiwfdpm
aio
ieuq
deiabbgvirkl
bxwtupwuounlrfgmsja
eikkz
wckytbbifesj
mrejdpxhbjfq
jmk
nddrppkzzk
pdwp
bjkxvefusmzucc
gxh
ma
mrqj
pzswvgnclhisyfn
ldcwaqo
dpmigub
tedgoml
edtpesmuvnqpvkppuv
rthakwx
kbqeitzemsj
czcqbchebjayokfz
uolqm
qbscvzzqytcxnygjr
npzmtshzavaxfjqsikcp
jynmzmbfu
hjxkb
pn
ptwcv
qgzgruykannokk
mslxqvp
uodw
lsptm
mpqhkncggu
gpmqqqjowynicbj
kzgh
xhuj
awi
gxdyodmw
zuyrdhsewvlxt
liej
osyyjuezq
gabwz
jprklfevh
wz
opqekywsgzffmdbkihys
jgeymahutunc
oeifrqkybqtjpqvpon
owrdr
ityowmszduambni
qmnwwgytpi
ysapg
ezswjkkll
zq
ndsoiijctqqclkzy
vseja
xfz
vvkoqjnhjgfypmro
km
xnrrndfwtbkok
ouvsbrwifhldslczt
cm
iyvnmjskk
uoarsrqrsosqkp
jxkoyjpzkhmcg
gsqp
knjoxxpisrhi
wobdettpxsidprajepnd
aiypbdtsmwvwkjh
azbswus
qnxmt
gkxgjmsjwziyxftlyng
nsznkjsbfhvzjzolz
od
qrflwgmfszsyghplreqc
zldvepyldfddtmx
frkrjqlvfued
viezmzdabbd
oorjkmhbvdwdgcgu
roftbjrzcf
kccbfzwfypyqgmkp
gqlmfhjrxulyzrg
hsalcxrdezlkzcteo
eesuljboejqvwedshyhk
ejxd
jshjfsbuljutlobmdvd
nocjdplamlfzvxlxztb
brhsscdtobrcbtbqbkkr
srqpbgcwwtyqqenpasc
ptackudefidengicrqyn
rnlqknsgutbszibvlga
vjns
xxvdxcxwmzprebqapqi`

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func loadCases() []string {
	lines := strings.Fields(testcasesRaw)
	if len(lines) == 0 {
		fmt.Println("no embedded testcases")
		os.Exit(1)
	}
	return lines
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := loadCases()
	for i, s := range cases {
		input := fmt.Sprintf("1\n%s\n", s)
		expected := reverse(s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
