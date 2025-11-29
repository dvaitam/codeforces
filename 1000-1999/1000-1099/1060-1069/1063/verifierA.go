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

const testcases = `100
13
ynbiqpmzjplsg
17
ejeydtzirwztejdxc
11
prdlnktugrp
15
qibzracxmwzvuat
16
khxkwcgshhzezroc
3
kqp
4
jrjw
4
rkrg
20
rsjoctzmkshjfgfbtvip
3
cvy
5
ebcwr
13
wqiqzhgvsnsio
16
vuwzlcktdpsukgha
9
dwhlzfknb
4
zewh
2
su
18
tvcadugtsdmcldbtag
6
wdpgxz
2
va
18
ntdichcujlnfbqobtd
13
gilxpsfwvgybz
6
fkqidt
15
vfapvnsqjulmvie
18
waoxckxbriehypltjv
12
sutewjmxnuca
20
gwkfhhuomwvsnbmwsnyv
2
fo
3
iwf
15
qprtyabpkjobzzn
7
rucxeam
14
kagawyavqtdgdt
7
jiwfdpm
3
aio
4
ieuq
12
deiabbgvirkl
19
bxwtupwuounlrfgmsja
5
eikkz
12
wckytbbifesj
12
mrejdpxhbjfq
3
jmk
10
nddrppkzzk
4
pdwp
14
bjkxvefusmzucc
3
gxh
2
ma
4
mrqj
15
pzswvgnclhisyfn
7
ldcwaqo
7
dpmigub
7
tedgoml
18
edtpesmuvnqpvkppuv
7
rthakwx
11
kbqeitzemsj
16
czcqbchebjayokfz
5
uolqm
17
qbscvzzqytcxnygjr
20
npzmtshzavaxfjqsikcp
9
jynmzmbfu
5
hjxkb
2
pn
5
ptwcv
5
zlnbt
15
mobdpyeabtteukd
18
ulgmzyypdbtwotukud
20
jzemzjxvzdqzgbzmolyg
15
lzucbbpiaqvssgh
3
yuy
17
wqnqjdensncdncdny
5
xazon
14
apkxiclcdlwall
6
ahlcte
7
agvvxdx
1
j
12
wathefodplwi
5
aglkp
10
jrukfscdrs
10
fmeezhkqhh
6
jlnvbe
20
amcwcenjrnxesnjulcho
12
uqbmnanxkogl
10
pcfzdidrtw
5
zwomf
14
nfhokqelouucpy
7
jawotoa
7
jdyujrt
5
nwypc
16
yhrymiuadivbaimq
19
wmodxiljyvgtcbczijr
11
dqhyfcnjjqe
19
qugrdnurmxyzijolsue
6
dwdmms
15
ervjlupxngppwqk
16
ubojexpbtgalpmaq
3
vcv
13
albdtaiuwjxhe
19
jgdnowkmfknuvneoweq
11
egfolzmnzpm
8
zgogswbm
2
hu
3
flb`

type testCase struct {
	n int
	s string
}

func parseTestcases() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcases))
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty testcases")
	}
	t, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return nil, fmt.Errorf("failed to parse t: %w", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected EOF on case %d (n)", i+1)
		}
		n, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			return nil, fmt.Errorf("failed to parse n on case %d: %w", i+1, err)
		}
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected EOF on case %d (s)", i+1)
		}
		s := strings.TrimSpace(scanner.Text())
		cases = append(cases, testCase{n: n, s: s})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return cases, nil
}

func referenceSolve(n int, s string) (string, error) {
	if len(s) != n {
		return "", fmt.Errorf("length mismatch: declared %d, got %d", n, len(s))
	}
	freq := make([]int, 26)
	for i := 0; i < len(s); i++ {
		c := s[i] - 'a'
		if c < 26 {
			freq[c]++
		}
	}
	var out bytes.Buffer
	out.Grow(len(s))
	for i := 0; i < 26; i++ {
		for cnt := 0; cnt < freq[i]; cnt++ {
			out.WriteByte(byte('a' + i))
		}
	}
	return out.String(), nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		exp, err := referenceSolve(tc.n, tc.s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
		out, err := run(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, exp, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
