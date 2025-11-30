package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `100
szy
id
yopumzgd
amntyyaw
ixzhsdka
a
m
nxaq
yopr
lhvh
janrudfu
dxkxw
qvgjjsp
sbphxzm
vflrwyv
covqdy
qml
xapbjwts
uffqhay
rrhm
sloivr
m
qyr
nbpl
qnpl
larrtzt
otazhu
rsf
zr
bvcca
ayyihidz
ljc
fiq
viu
owkpp
aj
knzgidi
tnah
m
bxf
wqvnrhuz
hquamvsz
vunbxj
gbj
cj
xfnsi
arb
sofy
m
ldgs
sgpdvmj
aktmjafg
zszekn
ivdm
vrpyrh
xb
ef
rgi
tqilkk
jh
esrydkbn
mz
ekd
csrhsci
jsrdoi
zb
atvac
dzbghzs
fdofvhf
nm
riwpk
gu
baazjx
omkmcc
todigz
vlifrgjg
lcic
cusukhmj
k
kzs
hkdrt
hh
z
mcir
xc
u
j
ppedqy
cqvffy
ekj
wq
egerx
y
tzvrxw
fjnr
bwv
iycv
znriroro
m
fipazu
`

func solveCase(s string) string {
	n := len(s)
	present := make([]bool, 26)
	for i := 0; i < n; i++ {
		idx := int(s[i] - 'a')
		if idx < 0 || idx >= 26 || present[idx] {
			return "NO"
		}
		present[idx] = true
	}
	for i := 0; i < n; i++ {
		if !present[i] {
			return "NO"
		}
	}
	l, r := 0, n-1
	ch := byte('a' + n - 1)
	for ch >= 'a' {
		if l <= r && s[l] == ch {
			l++
		} else if l <= r && s[r] == ch {
			r--
		} else {
			return "NO"
		}
		ch--
	}
	return "YES"
}

func parseTestcases() ([]string, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	if len(fields) != t+1 {
		return nil, fmt.Errorf("expected %d cases, got %d", t, len(fields)-1)
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
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(tc)
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		vals := strings.Fields(got)
		if len(vals) != 1 {
			fmt.Printf("case %d: expected single token output, got %q\n", i+1, got)
			os.Exit(1)
		}
		if vals[0] != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expected, vals[0])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
