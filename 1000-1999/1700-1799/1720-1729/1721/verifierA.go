package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	s1 string
	s2 string
}

// solve embeds the logic from 1721A.go.
func solve(tc testCase) string {
	seen := make(map[byte]struct{})
	seen[tc.s1[0]] = struct{}{}
	seen[tc.s1[1]] = struct{}{}
	seen[tc.s2[0]] = struct{}{}
	seen[tc.s2[1]] = struct{}{}
	return strconv.Itoa(len(seen) - 1)
}

// Embedded copy of testcasesA.txt.
const testcaseData = `
my
nb
iq
pm
zj
pl
sg
qe
je
yd
tz
ir
wz
te
jd
xc
vk
pr
dl
nk
tu
gr
po
qi
bz
ra
cx
mw
zv
ua
tp
kh
xk
wc
gs
hh
ze
zr
oc
ck
qp
dj
rj
wd
rk
rg
zt
rs
jo
ct
zm
ks
hj
fg
fb
tv
ip
cc
vy
ee
bc
wr
vm
wq
iq
zh
gv
sn
si
op
vu
wz
lc
kt
dp
su
kg
ha
xi
dw
hl
zf
kn
bd
ze
wh
bs
ur
tv
ca
du
gt
sd
mc
ld
bt
ag
fw
dp
gx
zb
va
rn
td
ic
hc
uj
ln
fb
qo
bt
dw
mg
il
xp
sf
wv
gy
bz
vf
fk
qi
dt
ov
fa
pv
ns
qj
ul
mv
ie
rw
ao
xc
kx
br
ie
hy
pl
tj
vl
su
te
wj
mx
nu
ca
tg
wk
fh
hu
om
wv
sn
bm
ws
ny
vw
bf
oc
iw
fo
qp
rt
ya
bp
kj
ob
zz
ng
ru
cx
ea
mv
nk
ag
aw
ya
vq
td
gd
tu
gj
iw
fd
pm
uc
ai
oz
zd
ie
uq
uu
ld
ei
ab
bg
vi
rk
ls
`

var expectedOutputs = []string{
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"2",
	"2",
	"2",
	"3",
	"3",
	"2",
	"3",
	"3",
	"3",
	"3",
	"3",
	"2",
	"2",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"2",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"3",
	"2",
	"3",
	"3",
	"3",
	"2",
	"3",
	"3",
	"3",
	"3",
	"3",
	"2",
	"3",
	"2",
	"3",
	"3",
	"3",
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines)/2)
	for i := 0; i < len(lines); i += 2 {
		s1 := strings.TrimSpace(lines[i])
		if i+1 >= len(lines) {
			return nil, fmt.Errorf("incomplete pair at line %d", i+1)
		}
		s2 := strings.TrimSpace(lines[i+1])
		if len(s1) != 2 || len(s2) != 2 {
			return nil, fmt.Errorf("bad string length at pair starting line %d", i+1)
		}
		tests = append(tests, testCase{s1: s1, s2: s2})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(tc.s1)
	input.WriteByte('\n')
	input.WriteString(tc.s2)
	input.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}
	if len(tests) != len(expectedOutputs) {
		fmt.Fprintf(os.Stderr, "testcase/expected mismatch: %d vs %d\n", len(tests), len(expectedOutputs))
		os.Exit(1)
	}

	for i, tc := range tests {
		if err := runCase(bin, tc, expectedOutputs[i]); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
