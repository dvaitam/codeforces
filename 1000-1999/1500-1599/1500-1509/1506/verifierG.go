package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded copy of testcasesG.txt so the verifier is self-contained.
const testcasesRaw = `cpyibaevspyxlkyaipz
xnrrvdg
rwzxivztvcnkclznzio
ygwu
dbsgulpgqs
aulhtnjlsdcqvqgdt
jxgmphetg
agyfaukrvttjlmqmj
vpbfn
xmdohzctvoozmzcqnpjw
cgxviopxzfaard
szlgiqokq
nntpitpvp
exmpjuzoklvftmwt
kumpxzfjr
t
bfatysdvwllpsbg
itanz
pcphdllevuhtkebuv
wdbu
oxcuawretcftuiuo
wazuegopyooripws
zjl
lbcqiiipno
bwchtuhrxbxf
lbbxyygwnkc
xvdmre
xpuceclqafqwefj
qsmzh
wswenwfy
hrsazzsmgyk
qoamqrvewhurz
ivfefwuemcba
py
du
dvltbbrhznhyswjcpc
bfvoddwfbvayvxnuwji
hguyrzugisaiqzx
yqcstubllwd
yalwknvb
eaakcpaoicvs
iheenrfhvmtnsqupiiej
gy
dzfsdzaqfdesmpp
ywrzzpgumeyoaknik
jlsvl
ayjfcntt
ucgmneye
bqqyplkwsmiuvzcie
ds
txeqqdzwyyccvydglkfx
ffhuepxo
dxteozgslxekuiciyabu
p
kbaixdmtkdxleye
vypxmzaredol
e
picpdmwgzrfokizryl
rwiorwybqx
vtaqprri
bkagsj
eibixjw
gwxftaclsflf
uhpmrtlduxcnbstg
oihubpvuqqiqiuljznxl
yrnu
ibuqkv
bcvdbonrcfvilaciw
lmmswwcxdj
oxijhyabyqrmehzj
d
zerdamipg
xiy
eedazg
knwptspoktogzfrdhvem
slosyuis
gshxanhmsxdubklzpuin
jlqwslsbnusfj
svphdosfls
vozkconpadtcpn
ua
qllffprcmstdaluht
itqzizmlwatlnjfpzl
ar
aybapud
npkpggledgbhinhlpl
lwyux
gxrbqggmq
vn
u
xcd
vzgqydqsspkvfogyyahk
rsjetaa
dmqgvgxnrt
uytpcntrfevtguwjnj
vboggjhxnlybu
lxlmnlqnjze
utiugtoqfuraechriota
rjbrbpm
lkej`

type testCase struct {
	s string
}

func solveCase(tc testCase) string {
	last := make([]int, 26)
	for i := range last {
		last[i] = -1
	}
	for i := 0; i < len(tc.s); i++ {
		last[int(tc.s[i]-'a')] = i
	}
	used := make([]bool, 26)
	stack := make([]byte, 0, 26)
	for i := 0; i < len(tc.s); i++ {
		c := tc.s[i]
		idx := int(c - 'a')
		if used[idx] {
			continue
		}
		for len(stack) > 0 && stack[len(stack)-1] < c && last[int(stack[len(stack)-1]-'a')] > i {
			used[int(stack[len(stack)-1]-'a')] = false
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, c)
		used[idx] = true
	}
	return string(stack)
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cases = append(cases, testCase{s: line})
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return cases, nil
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
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
		expected := solveCase(tc)
		input := fmt.Sprintf("1\n%s\n", tc.s)
		got, err := runCandidate(bin, input)
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
