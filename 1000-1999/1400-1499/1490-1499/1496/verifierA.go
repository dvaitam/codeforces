package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solve mirrors 1496A.go.
func solve(n, k int, s string) string {
	if n < 2*k+1 {
		return "NO"
	}
	for i := 0; i < k; i++ {
		if s[i] != s[n-1-i] {
			return "NO"
		}
	}
	return "YES"
}

type testCase struct {
	n int
	k int
	s string
}

// Embedded testcases from testcasesA.txt.
const testcaseData = `
100
28 6
ynbiqpmzjplsgqejeydtzirwztej
4 2
cvkp
18 1
lnktugrpoqibzracxm
23 10
uatpkhxkwcgshhzezrocckq
30 15
djrjwdrkrgztrsjoctzmkshjfgfbtv
9 3
ccvyeebcw
30 12
wqiqzhgvsnsiopvuwzlcktdpsukgha
24 4
dwhlzfknbdzewhbsurtvcadu
7 0
mcldbta
7 1
wdpgxzb
30 0
rntdichcujlnfbqobtdwmgilxpsfwv
7 0
zvffkqi
4 2
ovfa
16 6
sqjulmvierwaoxck
24 0
riehypltjvlsutewjmxnucat
7 2
fhhuomw
29 10
snbmwsnyvwbfociwfoqprtyabpkjo
2 1
gr
21 1
xeamvnkagawyavqtdgdtu
7 2
iwfdpmu
3 0
ioz
26 1
ieuquuldeiabbgvirklsbxwtup
23 10
ounlrfgmsjaeeikkzlwckyt
2 0
if
5 2
jlmre
10 0
pxhbjfqxcj
13 6
kjnddrppkzzkd
16 1
wpnbjkxvefusmzuc
3 0
gxh
2 1
ad
13 4
qjopzswvgnclh
9 4
yfngldcwa
29 8
oyvgdpmigubzgtedgomlredtpesmu
22 6
qpvkppuvgrthakwxkkbqei
20 2
msjwwzpczcqbchebjayo
11 1
zeuolqmqqbs
3 0
xny
7 2
rtnpzmt
19 3
zavaxfjqsikcpijynmz
13 0
fuehjxkbbpnep
29 13
twcvwezlnbtomobdpyeabtteukdwr
21 5
gmzyypdbtwotukudvwtjz
28 2
mzjxvzdqzgbzmolygolzucbbpiaq
22 9
sghcyuyqwqnqjdensncdnc
4 1
yexa
26 7
nvnapkxiclcdlwallfahlctega
7 0
xajlwat
8 1
fodplwie
1 0
l
11 3
jjrukfscdrs
10 1
meezhkqhhy
6 2
lnvbet
1 0
c
23 1
enjrnxesnjulchouluqbmna
14 5
ogljpcfzdidrtw
5 2
omfyn
14 2
hokqelouucpygj
1 0
t
15 0
gjdyujrtenwypcv
16 3
rymiuadivbaimqsw
13 3
dxiljyvgtcbcz
9 2
rkdqhyfcn
28 4
jqesqugrdnurmxyzijolsuefdwdm
13 4
oervjlupxngpp
23 8
kpubojexpbtgalpmaqcvcvx
30 12
albdtaiuwjxheysjgdnowkmfknuvne
15 2
qkegfolzmnzpmxh
26 3
ogswbmbhucflbxuvfhtjtcwqyj
25 5
nobuwqvurxnsopiwpgkibbbfl
1 0
u
1 0
c
26 6
vhtmrhogktdtczkkrokiaqbglc
7 2
ggivxxj
10 4
miplwhbjrc
1 0
p
24 7
bznpoodcchdyengotcnrymbf
8 3
heilkndr
29 4
trzgwjyoqtoruiihadtzwdfxnhgjx
22 0
xrqnbdmuidxslhvwwrvjhx
27 3
cqjvkhlupjsfezarqklsuazemef
17 1
eygzypsywghxehyml
20 9
eupdtaqtlpojahrufvpz
24 7
rkwcietmwgkzjmbgbkxxhkov
24 10
vhilvfjalsrbxuelapubahba
8 2
cblvnego
14 2
ljfukxzxnmanir
18 7
ybsdnmfaqetvqxweck
8 1
hazfxzvr
6 0
ntdtuo
23 2
ttbikzxxmaubpcljveohqlf
24 12
mkizpmajqjrpbyrsrivbomxd
13 2
pbaixbivvswyj
22 3
yqqkmigdskzhsvxrvlfekx
27 0
sbselljujkpzmtnfazesboekaxp
30 8
xtgcrnifqfcvufsdquzrtmynijjany
27 11
iirqrkkgwznzeayqezvwzsmlobr
30 13
utzzyhalqfvguluwpaxxhshifyncso
8 3
qzzwdgfo
30 2
numiinyyltkcjapayigymmnyuuvmwb
19 7
lseswikampqebcsllac
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	pos++
	res := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos+2 >= len(fields) {
			return nil, fmt.Errorf("case %d incomplete", i+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		k, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return nil, fmt.Errorf("case %d bad k: %v", i+1, err)
		}
		s := fields[pos+2]
		pos += 3
		if len(s) != n {
			return nil, fmt.Errorf("case %d: string length %d != n %d", i+1, len(s), n)
		}
		res = append(res, testCase{n: n, k: k, s: s})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra tokens at end")
	}
	return res, nil
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := fmt.Sprintf("1\n%d %d\n%s\n", tc.n, tc.k, tc.s)
		expected := solve(tc.n, tc.k, tc.s)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
