package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesD = `100
2
ltpu
2
tapirhgwpr
8
muehueqmx
1
f
1
jyaiptxmwz
4
xzsoeld
1
p
3
vnyu
4
qmslr
7
shkvaitvwf
5
ssdwug
5
jdcpupclzc
2
ajnyndb
10
ybmwskriqh
1
a
1
tr
1
n
5
iewbk
3
emmoqm
9
dtzqinuxwh
4
iqjrk
1
s
1
mtsueb
4
lvltwi
1
sbvaliuo
5
tkflf
3
tijzmd
1
j
4
uzihkfvnu
1
tk
6
hozfck
4
ihzd
1
k
3
ikzucztlse
3
qziolun
4
snbne
1
ptqn
4
bxoyvxqjr
2
csjdzh
1
z
4
nsbapxdfq
2
vaqrn
1
k
2
rpz
1
h
1
rdfh
2
apusm
1
h
5
qqnbp
1
byebdb
2
bw
2
kf
1
lmumsj
3
gknder
1
z
1
bl
7
uzbtnblu
6
nwnoahgr
5
cznhn
1
klr
1
owdxv
7
vdxksrdzs
1
e
4
bqcs
2
fa
1
advpwj
5
zcbysqqwhd
2
rbrksfchf
4
twym
3
tmlrn
2
qh
7
xfnwsys
8
eumefdpxp
8
sxfeiyges
6
hwryjvwnt
10
igjaipzmgf
6
hkpyenwpwt
4
sura
2
mzxbohhu
1
ih
3
eftw
1
f
1
f
1
xzcdcij
1
o
6
aakknmpcgu
8
merkdicvnd
5
dqwlvyly
5
vvvuzidy
6
srqdvp
6
bwjvxsxfu
1
luo
1
re
6
xutnrj
2
pjz
3
cdw
3
rsx
3
diimbeb
5
hwyqlkmo
2
lpdeisdvd`

func nextString(d int, s string) string {
	n := len(s)
	if d == 1 {
		return "Impossible"
	}
	t := []byte(s)
	ok := func(pos int) bool {
		if pos-d+1 >= 0 && t[pos] == t[pos-d+1] {
			return false
		}
		if pos-d >= 0 && t[pos] == t[pos-d] {
			return false
		}
		return true
	}
	orig := []byte(s)
	for i := n - 1; i >= 0; i-- {
		for c := orig[i] + 1; c <= 'z'; c++ {
			t[i] = c
			if !ok(i) {
				continue
			}
			viable := true
			for j := i + 1; j < n && viable; j++ {
				placed := false
				for c2 := byte('a'); c2 <= 'z'; c2++ {
					t[j] = c2
					if ok(j) {
						placed = true
						break
					}
				}
				if !placed {
					viable = false
				}
			}
			if viable {
				return string(t)
			}
		}
		t[i] = orig[i]
	}
	return "Impossible"
}

type testCase struct {
	d int
	s string
}

func parseTests() ([]testCase, error) {
	reader := bufio.NewReader(strings.NewReader(testcasesD))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		var d int
		var s string
		if _, err := fmt.Fscan(reader, &d); err != nil {
			return nil, err
		}
		if _, err := fmt.Fscan(reader, &s); err != nil {
			return nil, err
		}
		tests[i] = testCase{d: d, s: s}
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d\n%s\n", tc.d, tc.s)
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse tests error:", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		want := nextString(tc.d, tc.s)
		input := buildInput(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\n", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
