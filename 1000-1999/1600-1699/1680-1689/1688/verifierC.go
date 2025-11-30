package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
1
c
fxz
itg

3
s
nu
zxql
oqiba
okm
qfrf
ha

1
feq
lqvrf
znxq

2
llofy
wxou
hpipq
zlvoo
sxr

3
pvhk
ti
jjzw
rqqut
njxgp
lvtcz
xag

3
b
ubish
y
eihgb
wybb
lfh
c

1
a
x
l

2
fx
qw
m
bzheb
l

3
j
paj
rytx
i
twep
cv
dao

1
sympq
eki
tnu

1
evbib
ff
o

3
qw
h
wo
i
s
tz
wlivn

2
yaebm
fdqx
h
d
f

1
g
q
ojru

2
vy
xz
nqas
bnqsf
v

2
aqd
ljwlj
v
ddjg
z

2
n
ogst
a
alj
h

2
ds
mwo
yl
didd
t

2
gwda
vpybx
jloe
ipq
xxzn

2
mhf
tirn
s
dclfr
zn

1
v
e
mhw

3
ofq
der
dkqh
iffow
ml
xeoox
z

3
xfmq
p
miw
wupl
kwxvc
rt
mv

2
k
qwou
da
gxst
gdmr

1
xss
pz
eatvn

2
qsf
wgyc
apr
y
pvkoi

3
acty
fyy
ivuz
bf
movj
aj
oalbr

2
ogvjp
pw
wjcik
kuz
uum

3
q
mt
ezquc
bho
hqibd
v
lgkl

1
olf
ojoe
ugik

1
h
gyvl
le

1
iz
ummzx
ixt

3
wxkxm
rtu
l
mpf
lop
f
mea

1
flc
arkh
mrjpu

1
kgp
e
ki

1
lick
hw
xt

1
luy
befcn
yiek

3
dkuyw
mhbmy
ptkr
tcsqr
mwof
mqob
o

3
dv
fcmjo
i
v
hfa
nv
k

3
bphc
erae
rbbgr
q
vqh
lp
e

3
h
o
zb
gumkt
wqqy
qd
ug

1
gjkn
ne
kzjz

1
dpijq
ihnw
wr

3
a
ryggg
sbue
x
wxp
bxyhe
kbwgd

1
fycvo
gfk
qsc

2
vxbo
vdu
agn
irx
sqxg

2
wf
olmp
itgsp
gyps
jcf

2
uphyt
evjgr
daz
g
bkr

2
ocn
wajs
egeft
xcus
iucp

2
es
hgt
stw
qnhu
rb

2
et
ndom
pmjg
hb
qctrv

1
m
mhqi
l

3
qzp
cwoxw
ia
p
e
gk
rb

3
uj
u
rcvve
webj
vipbr
ykv
t

2
t
lzu
zpj
teabk
ualv

3
v
w
qtynn
hzff
basyx
vfj
b

1
zhmcl
t
h

1
gdawm
q
suh

1
qqmnz
en
ox

2
s
qo
tzuo
pt
le

1
wfe
suio
ognn

2
ly
m
anjar
sivi
ow

2
qto
rr
oj
ndq
xu

3
dnto
qocxm
tyww
zrl
eh
no
wxzf

2
ksm
uwz
vkv
m
bgodd

1
ktk
fxkz
q

3
mtho
t
l
j
ewwb
k
puta

2
ghxtk
kj
stxp
yicq
sl

1
lzf
rz
hgsoh

2
sgq
am
xluf
xw
fptaz

3
gw
at
o
xf
mta
l
j

1
lpmdc
fezy
ucie

2
czqwz
jar
gcnef
jogbw
ozb

1
tl
zwu
g

1
ak
i
l

1
zwtel
nx
jr

1
k
xc
ybexp

2
pvj
cue
a
e
wej

3
echsx
zk
f
gnuv
k
purl
qd

3
e
enyy
ui
rwdcq
xtfgv
vwcrk
qh

1
vmod
hbl
m

1
dxzil
hknhb
fe

3
me
u
yno
d
awdb
na
zm

2
fi
zis
qj
snp
myjm

1
jzso
nswkr
s

1
nso
q
pquif

1
dxh
ne
gqmfc

2
eyqc
vfh
bnrbm
yjhfx
zhgs

1
iriti
au
mqwbb

2
rk
xub
kfob
li
ivqi

1
anu
fm
enwj

1
yz
n
r

3
aa
gvu
yrpnv
wskda
g
x
yqaea

2
dau
alidm
na
gf
wnrbs

2
h
hbgvh
ta
ds
uh

2
wcu
ue
zhdlp
gg
jach

3
ulj
yjvby
cpq
ts
jtn
nudrf
htuec

2
t
pa
mg
lque
gdkk

2
ehsi
zkax
a
g
bzo

3
gh
b
d
dnss
uqay
f
g
`

func expected(list []string, final string) string {
	cnt := make([]int, 26)
	for _, s := range list {
		for i := 0; i < len(s); i++ {
			cnt[s[i]-'a'] ^= 1
		}
	}
	for i := 0; i < len(final); i++ {
		cnt[final[i]-'a'] ^= 1
	}
	for i := 0; i < 26; i++ {
		if cnt[i] == 1 {
			return string('a' + byte(i))
		}
	}
	return ""
}

type testCase struct {
	n     int
	list  []string
	final string
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(strings.TrimSpace(testcasesRaw))
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	idx := 0
	t, err := strconv.Atoi(fields[idx])
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	idx++
	cases := make([]testCase, 0, t)
	for caseNum := 1; caseNum <= t; caseNum++ {
		if idx >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseNum)
		}
		n, err := strconv.Atoi(fields[idx])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", caseNum, err)
		}
		idx++
		need := 2*n + 1 // 2n strings plus final string
		if idx+need > len(fields) {
			return nil, fmt.Errorf("case %d: not enough strings", caseNum)
		}
		list := make([]string, 2*n)
		copy(list, fields[idx:idx+2*n])
		idx += 2 * n
		final := fields[idx]
		idx++
		cases = append(cases, testCase{n: n, list: list, final: final})
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra data after parsing: %d fields remain", len(fields)-idx)
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

func runCase(exe string, n int, list []string, final, exp string) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, s := range list {
		sb.WriteString(s)
		sb.WriteByte('\n')
	}
	sb.WriteString(final)
	sb.WriteByte('\n')
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()

	cases, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i, tc := range cases {
		exp := expected(tc.list, tc.final)
		if err := runCase(exe, tc.n, tc.list, tc.final, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
