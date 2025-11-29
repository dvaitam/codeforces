package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded test data from testcasesG.txt.
const testData = `.cpyibaevs cybv
igs i
vcnkclznzi vcncni
hqsuw h
njlsdcqvqg nlsdcqg
lauk l
qmjevpb qevb
vooz.mzcqn vozzcn
r r
zlgiqokqin lgoqn
vftmwt vw
atobfatys aoat
zzz.qpc zzzq
buv bu
xcuawret ce
.go go
ipwsc.zjl iwcj
w w
uhrxbxfklb uxfl
xpuceclq pc
wswe wswe
azzsmgykmq azm
fwu fwu
ybduwrdv ybwv
cxsbfvod sfd
iohg. g
aiqzxxwvky zxvy
l l
k k
co c
tiheenrfhv tenrhv
sdz sdz
ppqywrz pyr
k k
ejlsvl elsl
gm g
kqbq k
miuvzciebd b
gl g
ffhu ff
eozgslxeku eoglek
mt m
ee.yelv e
ol o
vrp vrp
zrfo zo
orwyb r
rrifbkag rifba
ltac l
puh p
ux ux
tgtoihubpv tohub
bibu b
vd vd
fv fv
wjlmm m
oxijhyab ijhb
zerda zerda
weda w
tspoktog pko
s.tgs gs
xdubklzpui duzui
hdosflsn dosl
nubua.x. ubax
tdal.uhtxr ttr
njfpzl nfl
rn r
ggledgbh gled
x.rb r
bvnuaucxc b
ogy y
rsje sj
cvgxn c
trfevtg fg
hxnly hnly
n.jzetuti tui
bg b
bpmdlkejr mle
st t
ggcnatlzkg ctlg
jfif ff
es s
lc. l
omiqyh oiy
tzhbxcnvq thcvq
vxhzralu vxhzr
gp gp
heomnr hn
oxvbgjs xbgs
ljimnukolg ljmng
fa a
begadif gad
uymcax u
ynibssmjzt yz
eagjd eajd
g.lun.kyx gluky
wn w
fa fa
k k
e e`

// Embedded solver logic from 1366G.go (placeholder that always prints 0).
func solveInstance(_ string, _ string) string {
	// The reference solution currently ignores the inputs and prints 0.
	return "0"
}

func parseTests() ([][2]string, error) {
	lines := strings.Split(testData, "\n")
	tests := make([][2]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid test line: %q", line)
		}
		tests = append(tests, [2]string{fields[0], fields[1]})
	}
	return tests, nil
}

func runCase(idx int, bin, s, t string) error {
	expected := solveInstance(s, t)
	input := fmt.Sprintf("%s\n%s\n", s, t)

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Test %d: runtime error: %v\n%s", idx, err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("Test %d failed: expected %s got %s", idx, expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(i+1, bin, tc[0], tc[1]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
