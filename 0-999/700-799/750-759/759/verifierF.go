package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcases = `
szy
id
yopumzgd
amntyyaw
ixzhsdka
a
amvgnxaqh
prhlhvhy
janrudfu
dxkxw
nqvgjjspq
sbphxzm
vflrwyv
covqdy
qml
xapbjwts
muffqhaygr
hmqlsloiv
txamzxqze
yrgnbplsr
qnpl
larrtzt
otazhu
rsf
zr
bvcca
ayyihidz
fljcffiqfv
uwjow
ppdajm
nzgidi
gtnahameb
owq
rhuzwqo
quam
zkvunbxjeg
j
cj
xfnsi
arb
gsofywtqbm
ldgs
sgpdvmj
paktmjafg
zszekn
ivdm
lvrpyrhcx
c
ffr
iykt
ilkkdjhty
esrydkbn
mz
ekd
szmcsrhsci
jsrdoi
zb
atvac
dzbghzs
fdofvhf
nm
jriwpkdgu
baazjx
komkmcckto
ig
yrwpvlifrg
ghlci
yo
us
hmjbkf
zsjhkd
tsztchhaz
mcir
xc
u
j
ppedqy
cqvffy
ekj
wq
jegerxbykt
xwgfjnrfb
iycv
znriroro
m
fipazu
sabwlse
eeiimsmftc
pafq
quovux
hkpv
hwnkrtxu
uhbcy
ulfqyzgjj
rlfww
tcdtqsmf
ing
xyzbpvmwul
qfrxbqc
udixc

`

func referenceSolve(s string) string {
	count := 0
	for _, c := range s {
		switch c {
		case 'a', 'e', 'i', 'o', 'u':
			count++
		}
	}
	return fmt.Sprintf("%d", count)
}

func parseCases() ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcases))
	var cases []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		cases = append(cases, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), stderr.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, s := range cases {
		input := s + "\n"
		expected := referenceSolve(s)
		out, stderr, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
