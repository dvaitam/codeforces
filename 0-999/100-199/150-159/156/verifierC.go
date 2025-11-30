package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod = 1000000007
const maxN = 3000

var fact [maxN + 1]int64
var invfact [maxN + 1]int64

// Embedded testcases from testcasesC.txt.
const testcaseData = `
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
etbvmwqiqz
hgnsiopvuwzltk
psg
ghaxidwzgrt
hkqloq
puyvfkyhx
fmpzovf
ptvyoqjulmve
slkeroa
kxlwry
nhehyplzsjvlsu
hdhhu
dmwvsnbmwwnyvwbfoc
ioltyabpkjob
cnkagawszd
ivtdgdtugji
jdpmcai
nzzdieuquxdeiabb
zirkpsbxwt
vwounlrfqmsja
aeikxzlwcky
mnrqjdp
bjfqmmjm
kzzkdpdwqs
kxnbjkxvefamzucczcgxh
jdh
rqjopzsfncl
syfngobdcwaq
dpqombioubzg
hedgomlredttyesmuvnqp
uppuvgrmhakwxkkbqe
atzzemsylwwzpc
xqbqhe
ouiqsholpzwlgk
bkfzeuolqxqqbscvz
xtcxnygjrtvnpzmtshzav
qxjwqjs
zxkcpijpmzmbfue
jxkbep
fpudptwcvwzjzln
jmobdpyeaiat
eukdwrulgmb
gkzdwtotukevw
qzemzjxvzdizgbzmolyg
klzucbbpigvq
nsghcyuyqwqjldqjphmfensn
sdnzdn
xazvnvnapkxuisol
slcglwfrpm
rllfhlct
pgagjixa
lwathdfodd
hlwieaglkxjjrukfsc
qrsdfmeez
vqhhpxjlnr
reeycarmcwcenjrn
ssnjychoulu
jbmnznxkogjjpcfzd
rftrtwewmfynnfho
rqelouumy
jjawojxoagjdyujrt
annwyrrcvpyhrym
uqadivgkthaa
dmqswmmfdxiljyvg
xcbczijrkdqfyfcn
djqesqxfdnu
bmxyzijolsafdws
pmmmoervjlupxn
zppwqhppuboje
dbtgalpmapcvcvxvmal
rztaiuwjxheys
wgdnowmfknuv
geoweqkezfolz
mnzpmxhzfogswn
xmhb
mwfljx
zhtjbcwqyjylnob
vvurxnsopiogkibbbfl
vjgpqe
mdnztmrh
bgktdtczkrrokiaq
digcsgqlgg
gvxxjjwmiplwhbjr
yaozxobznpoo
nccqdyen
rotcnyymbfhphei
kkndrjtyzgwjyoqto
ruiihadtzwufx
hhgjxvaxmqnbdmuidx
mlhvwwrvjhxmxcqjv
`

func modpow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func initFact() {
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invfact[maxN] = modpow(fact[maxN], mod-2)
	for i := maxN; i > 0; i-- {
		invfact[i-1] = invfact[i] * int64(i) % mod
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n || n < 0 {
		return 0
	}
	return fact[n] * invfact[k] % mod * invfact[n-k] % mod
}

// solve mirrors 156C.go.
func solve(word string) string {
	n := len(word)
	sum := 0
	for j := 0; j < n; j++ {
		sum += int(word[j] - 'a')
	}
	var total int64
	for k := 0; k <= n; k++ {
		rem := sum - k*26
		if rem < 0 {
			break
		}
		ways := comb(n, k) * comb(rem+n-1, n-1) % mod
		if k&1 == 1 {
			total = (total - ways + mod) % mod
		} else {
			total = (total + ways) % mod
		}
	}
	ans := (total - 1 + mod) % mod
	return strconv.FormatInt(ans, 10)
}

func parseTestcases() ([]string, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]string, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		res = append(res, line)
		if len(line) > maxN {
			return nil, fmt.Errorf("case %d exceeds maxN", idx+1)
		}
	}
	return res, nil
}

func runCandidate(bin string, word string) (string, error) {
	input := "1\n" + word + "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	initFact()

	for idx, w := range tests {
		expect := solve(w)
		got, err := runCandidate(bin, w)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
