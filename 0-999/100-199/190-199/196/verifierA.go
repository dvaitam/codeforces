package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
100
ynbiqpmzjplsg
ejeydtzirwztejdxc
prdlnktugrp
qibzracxmwzvuat
khxkwcgshhzezroc
kqp
jrjw
rkrg
rsjoctzmkshjfgfbtvip
cvy
ebcwr
wqiqzhgvsnsio
vuwzlcktdpsukgha
dwhlzfknb
zewh
su
tvcadugtsdmcldbtag
wdpgxz
va
cquynkvjwewehhq
jwxgkzpwcqnzqttbzdai
hjtelnftwfazvitslbq
wfchruwgoxyjnodgq
howxfagvezofvb
hbjdirplwwejnjmttp
yhnhafexjntmq
tgoq
xffq
mwthsqqheutrq
ckiqnjtmz
ikw
wctrrtpoauc
scbkndervpfupiot
xntjlr
gfvozfc
dslybbswfh
urrqdiujpzddat
uifonfmmyjxcv
gqpnpmgwa
airfgmdduxqijm
uwzsufvpzncaqemqahe
qjhlscqcidmlntfqsxq
iqbwsmmjjxqy
mcanaos
lxkydnpeyhymlsjhz
lxsewyzsoncjmyetpvbe
nuwlisilohzfcpmuf
birxnqtnngselhwz
chpbxtuqkfq
edtjzifernompcxxbzuc
ntfefpijorkeomkonc
yzzpfrpyfcbgqq
rtchgqyljpbrqdtc
dfebpphuze
vedopmaiuiughdiqw
aikpmspxekxddsacq
ztaaokmcwm
yolheivczhatfjnsldc
rgpnvop
av
cvbdyhlnynmdorawjzm
xolkhqgokyvmvjptrzq
gvs
zmvaehbbtksgwooluxo
pevlhrkmqvfudhr
ncxqoeimsip
aqqjzlyemnugltfxeurd
pocwnlpjtkzuzrszid
nwgqijmze
xiltmdkwdsaalnddll
vznfxeg
syrdxbfkjq
xsmlbrjvuegnyexajz
oeexnxoh
yjcfgehufsikcrkcto
depnocivngwednjhl
upkkigjhypo
qspohwylokumvxrbc
zssuczgynnvwsglipkx
ytplrkgfzm
dimxwqgkpzjsrrmkeqw
wskheneowkxm
hssbnwosubwscthmk
gphxkmqtpmahkrhgvik
dlhsuqxqmoba
zwatjkqrgcsjggklc
blszhvlvccpvnrileylb
qknoxlelicsaeijfrr
ikvxr
gvpiidgashjh
yayvvvcdvhyfxm
inxalqalhxwaqtud
cqkrlntbrnyzekeabz
dyjfezvl
iqujwurysakjygrhdu
tolqdqmwzknvlapcjy
jqtzipqoczrgogqfh
dvwcpdlwxnyzvjyt
dansqcdsbwfoikck
jeslbg
zlphiikcoxkjs
mfjpcjpkm
igwkn
`

func maxSubseq(s string) string {
	n := len(s)
	if n == 0 {
		return ""
	}
	maxChar := byte(0)
	res := make([]byte, 0, n)
	for i := n - 1; i >= 0; i-- {
		if s[i] >= maxChar {
			maxChar = s[i]
			res = append(res, s[i])
		}
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return string(res)
}

func parseTests(raw string) ([]string, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty testcases")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	if len(lines)-1 < t {
		return nil, fmt.Errorf("not enough cases")
	}
	cases := make([]string, 0, t)
	for i := 1; i <= t; i++ {
		cases = append(cases, strings.TrimSpace(lines[i]))
	}
	return cases, nil
}

func buildInput(s string) string {
	return s + "\n"
}

func runCandidate(bin, input string) (string, error) {
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
	bin := os.Args[1]

	cases, err := parseTests(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, s := range cases {
		input := buildInput(s)
		want := maxSubseq(s)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed: expected %s got %s\ninput:\n%s", idx+1, want, strings.TrimSpace(got), input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
