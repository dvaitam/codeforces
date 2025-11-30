package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `kfzeuolqmqqbscv
ytcxnygjrtnpzmtsh
v
x
jqsikc
ijynmzmbfuehjxkb
pn
ptwcv
zlnbt
mobdpyeabtteukd
ulgmzyypdbtwotukud
jzemzjxvzdqzgbzmolyg
lzucbbpiaqvssgh
yuy
wqnqjdensncdncdny
xazon
apkxiclcdlwall
ahlcte
agvvxdx
j
wathefodplwi
aglkp
jrukfscdrs
fmeezhkqhh
jlnvbe
amcwcenjrnxesnjulcho
uqbmnanxkogl
pcfzdidrtw
zwomf
nfhokqelouucpy
jawotoa
jdyujrt
nwypc
yhrymiuadivbaimq
wmodxiljyvgtcbczijr
dqhyfcnjjqe
qugrdnurmxyzijolsue
dwdmms
ervjlupxngppwqk
ubojexpbtgalpmaq
vcv
albdtaiuwjxhe
jgdnowkmfknuvneoweq
egfolzmnzpm
zgogswbm
hu
flb
htjtcw
yjylnobuwqvurxnso
iwpgkibbbflajuae
znv
tmrhogkt
tczk
rokiaqbglcg
lggivxxjjqmiplwhb
rcaopxobzn
oodcchdyengotcnr
bfhpheilkndrj
rzgwjyoqtoruiihadtzw
fxnh
jxvaxrq
bdmuidxslhvwwr
hxhcqjvkhl
jsfezarqklsuazem
fqcey
zypsywg
xehymlts
updta
tlpojahrufvpzxprk
iet
wgkzjmbgbkxxh
ovxvvhilvfj
l
rbxuelapubahbahukcb
vnegoneljfuk
manirrzxvwoybs
nmfa
etvqxweckhfhazfxz
fwcntdtuowettbikzx
aubpcljveohql
xymkiz
majqjrpbyrsrivbo
xdmlpbaixbivv
wyjvygyqqkmigdskzhs
vlfekxasbselljujkp
tnfazesboekax
vvyixtgcrnifqfcv
sdquzr
mynijjanywiirqrkkgwz
zeayqezvwzsmlo
rn
zzyhalqfvguluwpaxxhs
ifyncsoh
qzzwdgfocnumiin
tkcjapayigym
nyuuvmwbsolse
wikampqebcsllacgwdv
pbkakmeyuinvetemjq
fe
pwuwy`

type testCase struct {
	s string
}

func solveCase(tc testCase) string {
	r := []rune(tc.s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		if r[i] != r[j] {
			return "0"
		}
	}
	return "1"
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		got, err := runCandidate(bin, tc.s+"\n")
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
