package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (same format as original file).
const embeddedTestcases = `2
ax

1
hexd

4
nbacghqtar

3
nhosizayz

1
kiegykd

6
dlltizb

6
rdmcrjut

2
gwcbvhyjch

2
ioulfll

2
wvuct

4
xhfomiuwr

3
yybhbz

6
icgswkg

5
muoeiehx

1
nsmlheqpcy

1
eu

3
tcmmtoq

6
avxdvryiy

1
jn

2
axxiqyfq

3
uqtge

1
ryq

10
kpadlzjhbh

2
cx

1
yr

1
vpr

2
qtngr

5
mvulo

4
dhhckasr

2
hacwubhcbk

8
hivpgrexs

2
zpzn

2
dv

6
nnoxbv

2
bm

1
hg

9
oenfiohco

1
bu

1
yh

1
ppgmbfm

6
izzojnw

9
vwpegjgbs

1
x

1
b

4
qfbqcfctc

1
dshstbt

6
vssqkig

3
himevu

2
kycaotsd

3
gqielchlj

5
rwjtzuqa

1
ei

1
xr

2
jtgwk

4
qpibc

3
bakye

3
xor

1
radcwer

3
sreneb

6
zblgvh

2
ly

6
xehfzzfnaf

7
zvxzhif

2
mb

2
golj

3
avgm

5
cyilu

5
vrkadif

3
bdtnl

1
tqdmsgi

6
aqzrvxx

3
lncv

7
kvdxjqjvnk

4
regnv

2
tsj

6
ajjgnzstu

3
oovgqpzz

2
jq

2
yh

1
zgea

4
ptyc

4
usgwwmp

1
euwa

5
hfzwqob

4
doez

9
rtkyotxqn

5
fxpoiyhu

5
puhiocwjh

4
krceeh

2
wgc

4
kronbgn

7
waysmpaljy

4
rxxrzth

4
inpa

4
vvzmxf

2
tra

5
svacuneofb

3
kgokkym

3
icpaxrb

1
ucyu`

func feasible(m int, s string, maxc byte) bool {
	n := len(s)
	last := -1
	for i := 0; i < n; i++ {
		if s[i] <= maxc {
			last = i
		}
		if i >= m-1 {
			if last < i-m+1 {
				return false
			}
		}
	}
	return true
}

// Embedded reference solution (adapted from 724D.go).
func solve(m int, s string) string {
	n := len(s)
	lo, hi := byte('a'), byte('z')
	ans := byte('z')
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if feasible(m, s, mid) {
			ans = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	cnt := make([]int, 26)
	for i := 0; i < n; i++ {
		if s[i] < ans {
			cnt[int(s[i]-'a')]++
		}
	}
	var b strings.Builder
	for c := 0; c < int(ans-'a'); c++ {
		if cnt[c] > 0 {
			b.WriteString(strings.Repeat(string('a'+byte(c)), cnt[c]))
		}
	}
	lastSmall, lastEq := -1, -1
	lastSmallPos := make([]int, n)
	lastEqPos := make([]int, n)
	for i := 0; i < n; i++ {
		if s[i] < ans {
			lastSmall = i
		}
		if s[i] == ans {
			lastEq = i
		}
		lastSmallPos[i] = lastSmall
		lastEqPos[i] = lastEq
	}
	picksAns := 0
	pos := -1
	for pos < n-m {
		reach := pos + m
		ns := lastSmallPos[reach]
		if ns > pos {
			pos = ns
		} else {
			picksAns++
			pos = lastEqPos[reach]
		}
	}
	if picksAns > 0 {
		b.WriteString(strings.Repeat(string(ans), picksAns))
	}
	return b.String()
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	scan := bufio.NewScanner(strings.NewReader(embeddedTestcases))
	scan.Split(bufio.ScanWords)
	cases := make([]struct {
		m int
		s string
	}, 0)
	for {
		if !scan.Scan() {
			break
		}
		m, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		s := scan.Text()
		cases = append(cases, struct {
			m int
			s string
		}{m: m, s: s})
	}

	for i, tc := range cases {
		input := fmt.Sprintf("%d\n%s\n", tc.m, tc.s)
		expect := solve(tc.m, tc.s)
		got, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
