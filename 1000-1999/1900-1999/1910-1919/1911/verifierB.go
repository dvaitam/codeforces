package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcases = `7 omlredt
16 esmuvnqpvkppuvgr
20 hakwxkkbqeitzemsjwwz
16 czcqbchebjayokfz
5 uolqm
17 qbscvzzqytcxnygjr
20 npzmtshzavaxfjqsikcp
9 jynmzmbfu
5 hjxkb
2 pn
5 ptwcv
5 zlnbt
15 mobdpyeabtteukd
18 ulgmzyypdbtwotukud
20 jzemzjxvzdqzgbzmolyg
15 lzucbbpiaqvssgh
3 yuy
17 wqnqjdensncdncdny
5 xazon
14 apkxiclcdlwall
6 ahlcte
7 agvvxdx
1 j
12 wathefodplwi
5 aglkp
10 jrukfscdrs
10 fmeezhkqhh
6 jlnvbe
20 amcwcenjrnxesnjulcho
12 uqbmnanxkogl
10 pcfzdidrtw
5 zwomf
14 nfhokqelouucpy
7 jawotoa
7 jdyujrt
5 nwypc
16 yhrymiuadivbaimq
19 wmodxiljyvgtcbczijr
11 dqhyfcnjjqe
19 qugrdnurmxyzijolsue
6 dwdmms
15 ervjlupxngppwqk
16 ubojexpbtgalpmaq
3 vcv
13 albdtaiuwjxhe
19 jgdnowkmfknuvneoweq
11 egfolzmnzpm
8 zgogswbm
2 hu
3 flb
6 htjtcw
17 yjylnobuwqvurxnso
16 iwpgkibbbflajuae
3 znv
8 tmrhogkt
4 tczk
11 rokiaqbglcg
17 lggivxxjjqmiplwhb
10 rcaopxobzn
16 oodcchdyengotcnr
13 bfhpheilkndrj
20 rzgwjyoqtoruiihadtzw
4 fxnh
7 jxvaxrq
14 bdmuidxslhvwwr
10 hxhcqjvkhl
16 jsfezarqklsuazem
5 fqcey
7 zypsywg
8 xehymlts
5 updta
17 tlpojahrufvpzxprk
3 iet
13 wgkzjmbgbkxxh
11 ovxvvhilvfj
1 l
19 rbxuelapubahbahukcb
12 vnegoneljfuk
14 manirrzxvwoybs
4 nmfa
17 etvqxweckhfhazfxz
18 fwcntdtuowettbikzx
13 aubpcljveohql
6 xymkiz
16 majqjrpbyrsrivbo
13 xdmlpbaixbivv
19 wyjvygyqqkmigdskzhs
18 vlfekxasbselljujkp
13 tnfazesboekax
16 vvyixtgcrnifqfcv
6 sdquzr
20 mynijjanywiirqrkkgwz
14 zeayqezvwzsmlo
2 rn
20 zzyhalqfvguluwpaxxhs
8 ifyncsoh
15 qzzwdgfocnumiin
12 tkcjapayigym
13 nyuuvmwbsolse
19 wikampqebcsllacgwdv`

func decode(n int, s string) string {
	var res []byte
	step := 1
	for idx := 0; idx < n; idx += step {
		res = append(res, s[idx])
		step++
	}
	return string(res)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(testcases, "\n")
	count := 0
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		count++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "invalid line %d\n", idx+1)
			os.Exit(1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid n on line %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		s := parts[1]
		expected := decode(n, s)

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n%s\n", n, s))

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx+1, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		if outStr != expected {
			fmt.Printf("Test %d failed: expected %s got %s\n", idx+1, expected, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", count)
}
