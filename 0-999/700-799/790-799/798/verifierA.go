package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func run(bin, input string) (string, error) {
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

// Embedded testcases from testcasesA.txt.
const testcasesRaw = `mynbiqpmzjplsg
ejeydtzir
ztejdxcvkprd
nktugr
oqibzrac
mwzvuatpkhxk
cgshhzezrocc
qpdjrj
drkrgztrsjoc
zmkshjfgfb
vipccvyeeb
cwrvmwqiqzhgvs
nsiopvuwzlcktd
sukghaxi
wh
zfknbd
ewhbsurtvcadu
tsdm
ld
t
g
wdp
xzbv
r
tdichcu
lnfbq
btdwmgil
xpsfwvgybzvffkq
dtovf
p
nsqjulmvier
aoxckxbriehy
ltjvlsut
wjm
nucatgwkfhhu
mwvsnbmw
nyvwbfociw
oqp
rtyabpkjobzzngr
cxeamvnkaga
yavqtdgdtugj
wfdpm
caiozzdieuq
uuldeiabbgvirk
sbxwtu
wuounlrf
msja
eik
zlwcky
bbifesjlmr
jdp
hbjfqxcjmkjn
dr
ppkzzkdpdwpnbjk
vefusmzucczc
xhbm
d
rqjopzs
vgnclhisyfng
dcwaqo
vgdpmigubzgte
go
lredtpe
muvnqpvkpp
vgrthakwxkk
q
itz
msj
wzpczcqbcheb
ayokf
euolqmqqbscvz
qytcxnygjrtnp
zmtshzavaxfjqs
kcpij
nmzmbfuehjxkb
p
eptwcvw
zln
t
mobdpyea
t
eukdwrulgm
yypdbtwotukud
wtjzemzjxvz
qz
bzmo
ygolzu
bb
iaqvssgh
yu
qwqnqjdensncd
cdnyexa
onvnapkxiclcd
wallfa
hlctegagvvxdxa
lwath
efodplwieaglkp
jrukf
cdrsjfmeez
kqhh`

type testcase struct {
	s string
}

func parseTestcases() ([]testcase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var res []testcase
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		res = append(res, testcase{s: ln})
	}
	return res, nil
}

// Embedded solver logic from 798A.go.
func solve(s string) string {
	n := len(s)
	diff := 0
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			diff++
		}
	}
	if diff == 1 || (diff == 0 && n%2 == 1) {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}


	for idx, tc := range cases {
		input := tc.s + "\n"
		expected := solve(tc.s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
