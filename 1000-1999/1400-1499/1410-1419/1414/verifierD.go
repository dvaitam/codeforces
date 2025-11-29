package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesRaw = `100
kkzlw
ky
bbifesjlmr
jdp
bjfq
jm
jnddrp
kzzkdpdw
nbjkxvef
mzucczcgxh
m
d
rqjopzs
nclh
syfng
dcwaqo
dpmi
ubzg
edgomlredt
esmuvnqp
ppuvgr
hakwxkkbqe
tzems
wwzpc
qb
he
j
y
kfzeuolq
qqbscvz
ytcxnygjr
npzmtshzav
x
jqs
kcpij
mzmbfue
jxkb
p
eptwcvw
zln
t
mobdpyea
t
eukdwrulgm
dbtwotuk
vw
jzemzjxvzd
zgbzmolyg
lzucbbpi
q
sghcyuyqwq
qjdensn
dn
dn
xaz
nvnapkxi
lc
lw
l
fahlct
gag
xa
lwath
fod
lwieaglk
jjrukfsc
rs
fmeez
kqhh
jln
e
amcwcenjrn
snj
choulu
bmnanxkog
jpcfzd
drtwe
mfynnfho
qelouu
py
jawo
oagjdyujrt
nwy
cvpyhrym
uadiv
a
mqswm
dxiljyvg
cbczijrkdq
yfcn
jqesq
rdnu
mxyzijols
fdw
mm
oervjlupxn
ppwq
puboje
btgalpma
cvcvxvmal`

func solveCase(s string) string {
	r := []rune(strings.TrimSpace(s))
	for l, rpos := 0, len(r)-1; l < rpos; l, rpos = l+1, rpos-1 {
		if r[l] != r[rpos] {
			return "No"
		}
	}
	return "Yes"
}

func parseTestcases() ([]string, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %v", err)
	}
	if len(fields) != t+1 {
		return nil, fmt.Errorf("malformed testcases: expected %d cases, got %d", t, len(fields)-1)
	}
	cases := make([]string, t)
	copy(cases, fields[1:])
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(cases))
	for _, s := range cases {
		fmt.Fprintf(&input, "%s\n", s)
	}

	got, err := runCandidate(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var expected strings.Builder
	for i, s := range cases {
		if i > 0 {
			expected.WriteByte('\n')
		}
		expected.WriteString(solveCase(s))
	}

	if strings.TrimSpace(got) != strings.TrimSpace(expected.String()) {
		fmt.Printf("output mismatch\nexpected:\n%s\n\ngot:\n%s\n", expected.String(), got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
