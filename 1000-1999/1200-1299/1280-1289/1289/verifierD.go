package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesData = `
srel
pusctapirh
wprr
muehueqm
avycfysbjyai
txmwznmx
oeldbepgiv
yujnqms
rsnshk
aitvwfwkrss
wu
usij
cp
pclzcneajny
dbttybm
skriqhbjacdt
bgnjtiewb
klemmo
mutvrdtzq
nuxwh
niqjr
aznska
tsuebuu
olvltw
xpasb
aliuojstkfl
kyl
ijzmdyasvx
jqh
zihkfvnuwdd
kkvhozfckx
gsoihzdbqgk
fikzucztls
njq
olunj
snbnega
tqnrwhbx
yvxqjrkh
sj
zh
b
qgnsbapxdfqj
vaqr
btdkeir
zzblhgdr
fh
zeapu
mbyihitqqn
p
yabyeb
bc
bwcqqpkf
clmums
ligkn
er
w
mzcsfblotuzr
uzbtnbl
pywknwnoahg
iwscznhne
k
rzidow
xv
zmvdxksrd
wapehymbqc
dvmfakdadv
wjsjzcby
qqwhdrxdrb
ksfchfuho
wymiltmlrn
mq
nxfn
sysvqvpeumef
px
wqosxfei
esqk
wryj
wntssigjaip
gfslhkp
nwp
tgosurapxcmz
bohhuwyvcgih
yief
wvbifbkfnc
zcdcijblosxv
aakknm
cgusxpme
kdicvndoq
dqwlv
yojvvv
zidykvsrqdv
qlbwjvxs
fuuxuefluodd
ekuxutnrj
fopjzfwcdwf
rsxmldiim
e
p
ihwyqlkmo
zyclpdeis
`

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func parseTestcases() []string {
	lines := strings.Split(strings.TrimSpace(testcasesData), "\n")
	var cases []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cases = append(cases, line)
	}
	return cases
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases := parseTestcases()
	for idx, s := range testcases {
		input := s + "\n"
		want := reverse(s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
