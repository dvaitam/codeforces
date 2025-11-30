package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `cl fxz
tgtbsvfnum zxqlroqib
o kmnqfrfhhafk
qqlqvr foznx
zsllofymwxouqhpip qqzlvoolsxrx
vhkwftiypjjzwqr qqutsnjxgpqlvtcz
gxdbsubishv d
ihgbnwybbllfhvacd cabxa
efxfqwamsbzh ebaltuxxd
pajorytxbi ymtwephcvvk
ozeq s
qkekiitnuawre vbibeffdouhqwbhh
icshtzztwlivniq yae
nf dqxchddafyhdg
v oojrumgvygxznnqas
nqsfdvzplaqdtljwlja vn
jgyv azob
ogstcajaljxchy pgdslmwoeylmdidd
kum gwdatvpybxwpjloezlip
xxznpvjmhfptirnwv wcsxsdclfreznczc
jm hwvvk
qjderyndkqhwqif fowhml
ooxaztmxfmqbpimiwxn wuplr
yxhrtgmvmua koq
damgxstmgdmryzg ixssgp
atvnpiqsfowgyclaprvv cyspv
iqoactylfyy zmivuzxebfpmovj
jroal b
sogvjpuepwrwjcikjk uzjuumqcqugmt
zqucjbhorhqibddvz mlgkl
olf pojoewougik
hpgyvl flez
izrum mzxkixtq
xmywwuywjrtuvcljmpf ilopcfkmead
lcxyunarkhtm rjpuel
pdezgkienli ckghwhx
kluytbefcnoyiekqsdku yw
hbmypptkrttcsqrvpmwo fnmqobdosedvq
mjozwa idv
faenvckuobph cperaewq
bgraqkvqhelpaerdhd og
gu mktumwqqyvfqdeugfmgj
enemkzjzdrd pijqypihnwewrv
tryg g
sbueuax iwxprbxyhetkb
euwrfyc vouj
kwiqscn nvxboj
agnk irxmsqxgn
zoolm ptitgs
gypskjcfltuphytu vsevjgrjdazagkb
izxvkocnpwa jssegeftymxcusoiuc
zzhesjhgtwkstwmq nhugrbivhetxmndo
pmjghhbrqctrv abmwnmhqidlql
scwoxwhiaapbeuegk hrbteujydurrcvve
bjqvipbrlykvdt ldtzl
zpjqteabknua lvrwbvcwr
ynnnhzfftbasyxlvf jabhszhmcldtchhrgdaw
qisuhbqqqmnze ene
bsfqontzuofptel eaiwfeunsuio
gnniyhlyubmtanja rpsivihowolqtov
rfojylnd qvhxuvmdntotqocxmo
zrlfehvufnopwxzfniks mjuwzivkvams
od dalktkx
xkzcqtpmthodt alajpe
km p
pusghxtkfkjymstxpoyi c
slhllzfhrzuxthgso hmisgqf
y pxlufvgxwtfpt
g wgatc
xfimtaaldjbslpm dcofezy
ieprczqwzjjarrg cne
jogbwk ozbehtlujzwudgehakd
lycqzwtel enxx
ryayaku excrybexpn
vjjcuexdacel wejysvrechsxyxhz
fvngnuvbknp url
dvbevm enyyxhuitrwdcqsxt
vtvwcr keqhhvo
dkhbldmasdxzi lrhknhbefexfmed
oadnawdbgn agzmkefiizishq
snppmyjmfu znjzsoqnsw
csainsoxcqr pquifgldxhgneqgqmf
oey qcxlvfhzrb
bmryjhfxmzhgsh qiritifautmqwbbzpg
jxubokfobglixoivqi bxkanuhfmwv
nwjcvfyzbncrseaa lgvus
nvytwskdacgaxqyqae ajidausalidmvvuv
ahgftw nrbsidhshbgvhv
afdsgu hjkwcugueszhdlpwfggm
chsjuljryj v
cp qvetsljtnt
rfuthtuecoctep aemg
quengdkkjnehsipz kaxbadgjbzoq
hcbadmdn ssnuqay
d ghcwiq
dyzfytjgxtcfvaavdu zooausahw
hwizydfgcni pdi
ekmhshgpodpmukwbobno xvwwhzrnruybzvzqmido
xnhyvyfwhd yairwemnsx
wu bbgbw`

type testCase struct {
	a string
	b string
}

func longestCommonSubstr(a, b string) int {
	n := len(a)
	m := len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	maxLen := 0
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				if dp[i][j] > maxLen {
					maxLen = dp[i][j]
				}
			}
		}
	}
	return maxLen
}

func solveCase(tc testCase) int {
	l := longestCommonSubstr(tc.a, tc.b)
	return len(tc.a) + len(tc.b) - 2*l
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 || len(fields)%2 != 0 {
		return nil, fmt.Errorf("malformed embedded testcases")
	}
	res := make([]testCase, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		res = append(res, testCase{a: fields[i], b: fields[i+1]})
	}
	return res, nil
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

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
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

	for i, tc := range cases {
		expected := strconv.Itoa(solveCase(tc))
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(tc.a)
		sb.WriteByte('\n')
		sb.WriteString(tc.b)
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
