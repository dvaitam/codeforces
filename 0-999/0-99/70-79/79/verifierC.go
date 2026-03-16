package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

const testcasesCRaw = `100
ff
2
VZ
qMn
cLRkBOzZU3G8xI7CGr5c
1
D7u
B54HkJlpoblul
1
GxGRJ
5CYAVH
2
wx_29
k9Wz
PHpFr7FGG1YwQ4D
2
KU6
UDFQo
0S1k46NrX6E
2
9ZT
JHGPN
AtUnFGx7RN4eY0vUa60
1
d
PdrLoR46gWHi2rp0n84
1
5TWc
xx
1
Rb
h9e
1
U
x
2
07
Vl
SayLcYpj_caw8NOVV
1
vFb
CJXMVc5qWz
3
E9
fQ
1gb
Y38iHLXzFGuj39v
2
M_A
S
9iQdqcikkgDOoG6_T7
1
oT
eqfLoNYZNTxqRBr
3
j
y
khGU
pgg
1
Wo
nbHQ
2
tIPy
R6
UZBBGbL
3
4
7HLl
Q
xbH96hNxsS_7xtb3
3
ggtm
Z
dAOF
n4LNeasbxt7UeoW
2
hK
zTD
Wwz4h
2
h
N
Pz9nSgbNQEX
1
sw9D
ZxrEH
2
7F1R
zok
MqJBSRS6fLU0Kgew
1
jZA5e
65R
3
i
yoT
ClH
hjI8_WBgvH
1
qk5kD
z3wYWKUj
2
UbZM
4Vlz
dErzqTU7ATPEx9JvT
3
W
IN
z0
_4O6auDHT64DP
1
b
8nUKMy5n4_gy0
3
rV
LmFZN
aN
EqGKlDTn_9Wewa
2
1QQeW
F7R7v
r4GDbfNWwlWW8Yz
2
dk
yDRs
asJD9
1
cI2
KCn3RtFPiESIT
2
q
tvP
PPzH16fGOn
2
H22jZ
Oftco
JoHrd9hhR0Yy2xn
2
evD
kFC
D5i6TCOn7r
2
g4
Em
lwiYiorZJOyz
3
r4U
GLSU8
VzW
INOQextzEl
2
CEf
u9
i9bgwkwe644UW
3
aIup
yIsEO
xu
7Fg_jYn
2
jAx
fvm
TpUNcv8x
3
d3jl2
B
Xriu
K2hvPXTN9zodzXEFN
2
1NMfL
IQFz1
kAyHCc4gCLih76R
3
ez
DZT
q
Qwol
1
BQ
v_0
3
d25E
eE
Jb
SGIdd
1
6a00H
RH4pixFaiIh
1
D
ZdNnOyv
3
57TH
X6RkG
1
O9nly
1
vBj
izuZtZgJgErsHW
2
oAS
SJ
bMJW
1
mz
cPiObVqSUEIdVX30o1j
3
cSm
i
6lW9f
Os9nk0uS0r25HKe
2
QUcD
Q3h
86bnAvqIU
2
HV75m
XiSk
3D_wyENqNmLECmW
2
vtekx
O3EoX
PQK4iR7t3nI10tgaY3bm
2
u
qYUQv
eAE1TbsLKinjkMX
2
O
C6rPf
EZY6pjKt2927omNT
2
NSzHA
Pn
dqQpi8NUz0Bh_DzzE9
2
npo
I
04fM8IRadyTBzoGrg
2
xH8YF
eTDVS
rbbEciPj
1
pId
jPs8WgPJIfRQiBTictGQ
2
d8Jw
R7g
xgMYwxY88OrZE5sGM9jb
1
B1O
w
3
Td3Qe
GNXBB
pYlk
cbLW9U0wRlsb53c9pKY7
1
e6xh
6epoJmgaSzfG2rLPodHH
3
B6Yi
Bi
Vxd_
lH0C8BM80Z3PCkFMi2w
1
q
_jOAKO
2
EDmB
rXow
0z
3
B
71b
EKqRr
DTDxH2ND
3
J2
kDs87
AhG
V_PQ9yhB
3
NHDf
CN8X
6ZJ
kjoQ5PlACFTU
1
q6uK
tPTY_5qQuQazL
1
Dg
a81x
2
uUzlV
YfH
EzN6pC4gMbwbtFiTTdav
2
P5N_
4
O6KnoVNukvtX1zKM
3
DWrf
nKxpx
_Zl
IZPUNonL
2
z3
LnG
azXEVx
3
Rm
7l2FN
Y
6nTna
3
D
Vl
zN1
a
2
t
K
Eyhe9DljZ44X
2
r
3E
e6Z5HTZss6bJ4JneA6
1
Lt
7m2cSvDYdj_o2Nw
3
ZSP
2
j0
buhrgxXe
3
iwi9B
tJ
W
v
3
Ue
Wdi_V
_ARw
Qste4Oj9Vha0f5jw
3
8t
R47Ii
7
KUXVoZvf
1
nAPQ
1_
2
7FPJ
7lH
Q3ci
3
jAXW
Or
IS5hf
1KUNlmQ_NRTeI8uiG
1
RD
DguodwgzbKhVZ
2
Jpv
ocik
7U18k
2
gP
WBD
h
2
S
1
mB
1
Y2
0u2jlr2rZqKo9
3
tKA
E31y
zlP
sYKDHBKSuIeKaq
2
DUf51
IEHOr
9nx0gV
1
Bi
4mGy50k3et0DjWGfV
2
o5
cBJc5
`

func ok(s string, b []string, l []int, k int) int {
	mx := -1
	for t := 0; t < len(b); t++ {
		start := k - l[t] + 1
		if start < 0 {
			continue
		}
		match := true
		for i := 0; i < l[t]; i++ {
			if s[start+i] != b[t][i] {
				match = false
				break
			}
		}
		if match && start > mx {
			mx = start
		}
	}
	return mx
}

func solveCaseC(s string, b []string) (int, int) {
	l := make([]int, len(b))
	for i := range b {
		l[i] = len(b[i])
	}
	p, mx, st := 0, -1, 0
	for i := 0; i < len(s); i++ {
		k := ok(s, b, l, i)
		if k != -1 {
			if i-p > mx {
				mx = i - p
				st = p
			}
			if k+1 > p {
				p = k + 1
			}
		}
	}
	if len(s)-p > mx {
		mx = len(s) - p
		st = p
	}
	return mx, st
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	data := []byte(testcasesCRaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	T, _ := strconv.Atoi(scan.Text())
	cases := make([]struct {
		s string
		b []string
	}, T)
	expected := make([][2]int, T)
	for tc := 0; tc < T; tc++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		s := scan.Text()
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		b := make([]string, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			b[i] = scan.Text()
		}
		cases[tc] = struct {
			s string
			b []string
		}{s: s, b: b}
		mx, st := solveCaseC(s, b)
		expected[tc] = [2]int{mx, st}
	}
	for i, c := range cases {
		var buf bytes.Buffer
		fmt.Fprintln(&buf, c.s)
		fmt.Fprintln(&buf, len(c.b))
		for _, x := range c.b {
			fmt.Fprintln(&buf, x)
		}
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("execution failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		outScan := bufio.NewScanner(bytes.NewReader(out))
		outScan.Split(bufio.ScanWords)
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		mx, _ := strconv.Atoi(outScan.Text())
		if !outScan.Scan() {
			fmt.Printf("missing second value for test %d\n", i+1)
			os.Exit(1)
		}
		st, _ := strconv.Atoi(outScan.Text())
		if mx != expected[i][0] || st != expected[i][1] {
			fmt.Printf("test %d failed: expected %d %d got %d %d\n", i+1, expected[i][0], expected[i][1], mx, st)
			os.Exit(1)
		}
		if outScan.Scan() {
			fmt.Printf("extra output detected on case %d\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
