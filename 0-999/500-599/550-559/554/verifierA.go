package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(s string) int {
	set := make(map[string]struct{})
	for i := 0; i <= len(s); i++ {
		prefix := s[:i]
		suffix := s[i:]
		for c := 'a'; c <= 'z'; c++ {
			str := prefix + string(byte(c)) + suffix
			set[str] = struct{}{}
		}
	}
	return len(set)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

const testcasesARaw = `ynbiqpmzjplsg
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
ntdichcujlnfbqobtd
gilxpsfwvgybz
fkqidt
vfapvnsqjulmvie
waoxckxbriehypltjv
sutewjmxnuca
gwkfhhuomwvsnbmwsnyv
fo
iwf
qprtyabpkjobzzn
rucxeam
kagawyavqtdgdt
jiwfdpm
aio
ieuq
deiabbgvirkl
bxwtupwuounlrfgmsja
eikkz
wckytbbifesj
mrejdpxhbjfq
jmk
nddrppkzzk
pdwp
bjkxvefusmzucc
gxh
ma
mrqj
pzswvgnclhisyfn
ldcwaqo
dpmigub
tedgoml
edtpesmuvnqpvkppuv
rthakwx
kbqeitzemsj
czcqbchebjayokfz
uolqm
qbscvzzqytcxnygjr
npzmtshzavaxfjqsikcp
jynmzmbfu
hjxkb
pn
ptwcv
zlnbt
mobdpyeabtteukd
ulgmzyypdbtwotukud
jzemzjxvzdqzgbzmolyg
lzucbbpiaqvssgh
yuy
wqnqjdensncdncdny
xazon
apkxiclcdlwall
ahlcte
agvvxdx
j
wathefodplwi
aglkp
jrukfscdrs
fmeezhkqhh
jlnvbe
amcwcenjrnxesnjulcho
uqbmnanxkogl
pcfzdidrtw
zwomf
nfhokqelouucpy
jawotoa
jdyujrt
nwypc
yhrymiuadivbaimq
wmodxiljyvgtcbczijr
dqhyfcnjjqe
qugrdnurmxyzijolsue
dwdmms
ervjlupxngppwqk
ubojexpbtgalpmaq
vcv
albdtaiuwjxhe
jgdnowkmfknuvneoweq
egfolzmnzpm
zgogswbm
hu
flb
`

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesARaw))
	idx := 0
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" {
			continue
		}
		idx++
		want := expected(s)
		gotStr, err := run(bin, s+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Sscan(gotStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\n", idx)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %d\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
