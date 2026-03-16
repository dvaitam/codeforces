package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleG")
	cmd := exec.Command("go", "build", "-o", oracle, "808G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return oracle, nil
}

func readCases(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	cases := []string{}
	var sb strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			if sb.Len() > 0 {
				cases = append(cases, sb.String())
				sb.Reset()
			}
			continue
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	if sb.Len() > 0 {
		cases = append(cases, sb.String())
	}
	return cases, scanner.Err()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesRaw = `kfuaqywakvoeh
dwet

nkvnej
fdokbeho

fahuqoz
qapw

mpqcpemaeenjw
uaedrqcmj

nlqaowlsktt
tmfl

qyszknpbuj
nzhnlbb

qowyzwjcrbx
jebl

yvhxanpeumbxqvwexy
pjvt

cwjqsogixabkfgwmtyfp
qxmsw

vubjlkmllomwu
mzrk

fjglfwjpknqxvysnga?k
eafdxvo

gkoepuhgyzfvj
fxkvx

djfchksrtwnijgmmz
n

q
hyecm

?y
ze

gqzknoxzxvo
y

hdnjyw?xhjmzxj
jk

l?qo
tgt

clwgqwtfmam
hzndj

y
sjvjrila

kalkjtprrkzzjakl
olsq

wvvgibcjuobmreedqj
nyl

hviykh
puuhapii

csgcfkvltlpl
mdqcu

fuhop?svuxogpbdapvv
nevm

cig?q?zz
ufxv

lxp?tu
kagnppq

xvtmoxhaaqn
fiymkapwgx

rxwipjbks
etnfowt

myifttjrbeqh?mbjmfu
guhbjli

?
oayt

qkffhdda
rtk

haynebmqgvv
ycn

ztaluzjrfpaorqfaixc
paixozodzh

ddjrfkwpvkrlrxxcnc
vnnmxewduz

wnwcqhol?x
uribqhrc

dqk?skpeagzdayq
lqdeuufhs

ljoojutcz
yron

ieosdpkity?wktoc
xs

h
ooply

frjmrrijysanzd
ucglonsrv

znonkorqiiw
drnswmatn

cyp?ot?wceyvluz
lbrxqhlb

thafqtq
hljaiq

eahtxbq?bft
izsi

uwzlvmiewimb
nmmvddwjqv

hcatwhvkn?sfrjlfepe
o

purmmwvklazbwgtrpw
btbxusd

hoyulq
zm

uvnqxemyhzrk
idvnryizl

ky?nmx
hx

vbxhajpswhewxjnn?elp
ste

ifnbzdwceqxlw
zbk

yovjwqm
namyspabg

zypwoxbtcv
jcyynu

amly
lfwoub

prgfktupslkfkds
rnmdoj

jrcajjidacbtdh
f

tktthbzlvpmyoagz
vx

fdty
wyghngkwsm

tnednrdjkjxnkef
ff

eiqqv
adxvsud

?szcufukhotwswupv?
rbhnajodr

rdwqobkgshohiacnnqev
dvlrf

uwnvsygo?vjzsvqk
vyvrvgb

fzmaekxkubp
ke

tb?fpfm
nakufmzxzk

uscmx
wcve

oqlaq??zkeyvwhk
kcdhsaf

iflubwzrnfooxvz
enfpwmvvpl

lcznezrbnkyggmca
nhokujpv

jdbspfln?whmzrjji?
weskjseonj

cd?
jt

ofidwmwhfcbaxzwf
twgju

sdapbn
vi

vfmmhzzvnfqco
w

ppku
zi

kgnifgifftkem
dwjsxyjy

cl??yphrbr?k?s
ipjkhzmyow

motjgxwyjkund?nefim
hpt

eb
fdolaajmux

hfvgjabpo
jone

nyk
flw

ysiuaihubyti
bs

ninzuzemuqquqqm
gvvcuxje

jvltbxwsw
sgiyawok

rfuylazlpv
yoeng

tdpi
umc

ka
kpykdcxgj

muouaiysqvurbvm
osa

llewdr
mmcxhjkqnk

ncixmrksxgqhr
hnlqslxc

wqgmefkzpi
kufeua

zymeru
yqbhva

?tzcoosbvvldr
vcr

gihotespvnbvtaqfkd
rxvqnaqq

xpjalkzicpvynqojf
hnsxvpjzur

nsxrxrcthnal
owhx

zuquhywqwvspgojxg
girmjw

xubdxo?oyyvyqbqdk
hdijvwue`

	cases, err := readCases(strings.NewReader(testcasesRaw))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read cases: %v\n", err)
		os.Exit(1)
	}

	for i, c := range cases {
		idx := i + 1
		cmdO := exec.Command(oracle)
		cmdO.Stdin = strings.NewReader(c)
		var outO bytes.Buffer
		cmdO.Stdout = &outO
		if err := cmdO.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "oracle run error on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(outO.String())

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(c)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("case %d failed\nexpected: %s\n got: %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
