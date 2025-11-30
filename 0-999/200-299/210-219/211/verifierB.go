package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `a 2 bck nv
gon 2 fj kxz
manirrz 4 s n af etv
ckh 2 a x
wcn 1 ouw
ttb 3 xz amu p
lj 2 hq fx
kizpmaj 3 bpy irs bmo
ml 4 a bx sv jvy
yqqk 4 dg hkz rvx efl
xasbse 3 ju kp nt
aze 1 ek
x 4 ivy cgt fin cfv
sdq 4 ij an iqr gkw
zeayqez 4 bo ntu a fq
uluw 4 x his y cs
hxwoqzzw 1 f
cnumiiny 3 cjk p y
gymmn 4 bos es ikw m
qebcslla 1 w
vr 4 k k ey inv
tem 3 bef p uwy
x 1 uwx
flzsccr 2 gz kw
ziqoe 4 beu fqu y cg
thzoqfwk 2 ry r
qk 1 d
tlsokmq 3 dek f aek
bnub 3 bt cfl bo
ttyivj 4 af iz m l
da 3 f tw ako
pcbrmzzi 1 cdq
ldpbeq 3 a fk efw
yfu 2 akp bmv
hujk 2 l f
olsemsy 1 s
v 4 efz a aq b
st 2 v az
nsvkwhe 3 gmr e sv
dnnhpm 2 hu ms
ii 3 aty hi t
zmudrb 2 amy mx
dwotoff 3 fn jqy l
eulypu 1 i
xsk 3 bu nr v
imxglee 1 afl
smocuww 1 rx
feg 2 a elp
wzkvd 4 fk kqr emx hjr
lmpqjnn 1 aes
wy 2 qtu i
mcz 1 d
xpkdvo 3 ivw yz f
yq 2 u q
tvgowvm 3 s m o
m 1 kms
qobp 1 pq
pifhrlf 3 ceo p ns
ip 2 j l
ev 1 ers
abmr 4 dpw kly v h
pjih 1 lq
ccjsnh 3 ey j x
vzljmtey 1 lq
huwlulni 3 wy dp do
lhx 2 hp v
moqosth 4 hjp qv c km
nbsy 1 ci
xkfd 2 ahl a
raedttg 1 ago
crfhhnm 4 n m t as
wlkvou 2 ciq w
ia 2 evy gs
gnqqdr 1 dpq
pforkeni 4 s hko l ny
o 1 moy
nhip 4 h jl ept cgv
quvdb 2 bk fyz
mvgnpgfm 1 vy
xkfq 4 u ptv got h
tkrxgabb 2 o htv
yq 4 juv pv f lr
movkoceh 1 m
wrglwbam 2 n nsz
uave 4 mr ctx ay n
jcxayz 2 aqs eo
dsoih 4 qx hs bc k
cqpnkmxa 1 tv
yea 1 bhn
imdlj 3 guz bkr afr
auikeldt 4 i cdq bep hlz
rw 2 w ghm
lfixlrgg 4 ei avy v jw
bctajz 1 w
gqrlo 1 qr
dcsokfj 2 nvw ab
u 3 v kr bjz
wqei 1 in
ikzr 4 dw fl bhw am
cjdf 2 ay nq`

type testCase struct {
	s    string
	m    int
	sets []string
}

// solveLogic mirrors 211B.go.
func solveLogic(s string, sets []string) []int {
	m := len(sets)
	qMasks := make([]uint32, m)
	counts := make(map[uint32]int)
	for i, ci := range sets {
		var mask uint32
		for j := 0; j < len(ci); j++ {
			mask |= 1 << (ci[j] - 'a')
		}
		qMasks[i] = mask
		counts[mask] = 0
	}

	n := len(s)
	occ := make([][]int, 26)
	for i := 0; i < n; i++ {
		c := s[i] - 'a'
		occ[c] = append(occ[c], i)
	}

	ptr := make([]int, 26)
	for l := 0; l < n; l++ {
		if l > 0 {
			c0 := s[l-1] - 'a'
			if ptr[c0] < len(occ[c0]) && occ[c0][ptr[c0]] == l-1 {
				ptr[c0]++
			}
		}
		var mask uint32
		for {
			minPos := n
			minC := -1
			for c := 0; c < 26; c++ {
				if mask&(1<<c) != 0 {
					continue
				}
				pi := ptr[c]
				if pi < len(occ[c]) {
					p := occ[c][pi]
					if p < minPos {
						minPos = p
						minC = c
					}
				}
			}
			if minC < 0 {
				break
			}
			if mask != 0 {
				if l == 0 || (mask&(1<<uint(s[l-1]-'a'))) == 0 {
					if _, ok := counts[mask]; ok {
						counts[mask]++
					}
				}
			}
			mask |= 1 << uint(minC)
		}
		if mask != 0 {
			if l == 0 || (mask&(1<<uint(s[l-1]-'a'))) == 0 {
				if _, ok := counts[mask]; ok {
					counts[mask]++
				}
			}
		}
	}

	res := make([]int, m)
	for i, mask := range qMasks {
		res[i] = counts[mask]
	}
	return res
}

func loadTestcases() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	var cases []testCase
	for lineIdx := 1; scanner.Scan(); lineIdx++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: invalid testcase", lineIdx)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %w", lineIdx, err)
		}
		if len(fields) != 2+m {
			return nil, fmt.Errorf("line %d: wrong field count", lineIdx)
		}
		cases = append(cases, testCase{
			s:    fields[0],
			m:    m,
			sets: fields[2:],
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return cases, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range testcases {
		expected := solveLogic(tc.s, tc.sets)
		var input strings.Builder
		fmt.Fprintf(&input, "%s %d", tc.s, tc.m)
		for _, set := range tc.sets {
			input.WriteByte(' ')
			input.WriteString(set)
		}
		input.WriteByte('\n')
		gotStr, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotFields := strings.Fields(gotStr)
		if len(gotFields) != len(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\n", idx+1, len(expected), len(gotFields))
			os.Exit(1)
		}
		for i, exp := range expected {
			gotVal, err := strconv.Atoi(gotFields[i])
			if err != nil || gotVal != exp {
				fmt.Fprintf(os.Stderr, "case %d failed at index %d: expected %d got %s\n", idx+1, i, exp, gotFields[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
