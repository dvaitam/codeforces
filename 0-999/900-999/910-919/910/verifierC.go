package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesRaw = `100
1
b
6
ee
djajc
gifi
ieaa
hfg
icic
4
ac
cci
fiich
ifjf
6
cghi
he
iifh
fjih
dfcj
hee
9
iijjg
dhi
jbf
adbaja
jdb
cedda
aaff
da
b
2
a
afeccc
9
agjadc
a
jbe
hae
ijae
jchd
f
a
cijg
8
fcfee
gaica
acc
bh
diaddh
e
j
jj
6
gei
c
g
cbib
bb
c
4
d
i
hheigd
dggiaj
10
g
jcbhf
i
j
efe
g
b
dah
g
hhdjjb
1
afe
2
hd
j
6
hcfg
e
b
j
gdb
j
8
h
fhc
ehi
gheg
ch
eigbj
jbbfci
gb
2
acegdf
cieb
3
gbfid
iecchd
fjch
8
ajgcgi
h
geg
hfifbd
jdgga
hih
cbagdj
gdbgi
4
jjd
jcaj
gheijc
dbfa
8
bjhfh
iha
j
cge
cachgh
ecaeih
f
i
7
hdehc
iebe
efe
gibidg
icibe
d
idie
1
b
7
dff
f
fchh
hch
defcbd
dfcf
cd
5
ggfej
jfgei
bfegh
ef
hbcf
7
ab
cfb
gaifdj
iehc
fdh
c
fe
3
febf
dd
jaffja
3
bg
ecfi
bfjgd
1
hhjf
9
jbjii
hghcgg
habhj
bi
bg
hae
f
ca
gb
6
hahdbh
ia
ii
a
ia
fidcf
8
c
bdbhd
j
gf
giici
c
dcgdef
cgcg
6
bib
eeih
dgc
ibajid
dg
acaeh
9
d
jf
d
c
icbhed
fe
jbgga
ebea
gf
5
gjidg
ch
fghj
jdj
hdhj
6
bcf
hdjjc
edieba
d
afi
fhb
7
aejj
dc
jg
bjhebh
dcje
dj
fjjgig
4
diaedc
ggbhg
hged
da
9
bjiaa
ggdi
bfi
ihj
h
deaaha
cd
dia
cebii
2
cgcaei
ehaiff
2
fbjff
eheijc
1
f
7
afiabi
jgggd
cj
a
fceaa
jd
bfbj
2
di
ba
7
i
jda
iiggc
gc
fajc
hgjhc
jcfc
1
ccg
10
ehhhdg
edfa
jage
i
jeed
hfij
diic
efgb
dgbgj
jihb
7
jfif
cd
cghhcg
fjg
efa
jadh
b
1
fjf
7
fb
jhgjd
bjaf
e
caaf
hjah
jddjfc
6
gjj
hhebid
fdffc
ij
dj
dgej
4
caghf
cdjchj
c
da
10
h
dcegja
f
e
j
hgb
h
ch
bechib
eeaii
4
g
cj
hda
fhacdj
6
ebdcda
beb
bij
fc
deia
f
3
biachg
fheebc
babcfc
5
icbdj
dfbcgd
afgh
ifciba
gc
7
dejibb
jjcdj
bifcid
hg
bfda
bga
befid
6
dacc
cg
be
habg
b
d
7
d
fcce
ejd
eejgh
gegc
gejhig
fibja
5
jhbi
hiecd
bdd
cidi
cbeh
3
bfcdi
g
agied
3
gddjdi
iej
caj
7
aahdi
eah
cha
fe
heieaf
g
dcghcg
5
c
g
i
caafd
jihgjj
6
a
d
i
acaee
a
jafebg
3
addc
giaje
d
10
ad
dcjacb
ddefb
dcjdbf
cddg
abd
efeie
aebhic
jfejg
gbicj
4
cbhbj
ha
gd
ficg
4
f
egc
je
faaa
2
ea
iddd
2
a
g
2
jjgi
a
3
d
be
iebcj
5
jb
aa
bhhaja
fb
eb
3
bg
hbe
jcfgd
10
dh
bhgf
ahaghd
giaig
bhe
gdc
dbaeic
gjac
aadaei
ch
6
bgijcb
edad
add
hihhaa
gij
idgbhh
9
cicejc
ibicaf
ajjegf
jjcbj
cj
fad
adfaa
hgjcaj
gbfcii
4
e
fhe
dfia
ajjbg
7
jhje
aeaa
jf
db
fefab
ecc
ehi
9
eejfi
ihggb
bicb
ghj
g
hjeh
i
ifgd
bacea
5
g
gdg
ejgge
e
ea
4
d
fe
dfgjc
a
2
gb
eced
7
g
fbdh
ih
hejd
hdhe
hg
e
9
efd
gabbc
fe
ddbbi
ihf
cddbe
hhhif
chh
g
9
cdjf
gehce
gac
ijd
jgfid
cej
djbai
haidab
fcgcj
10
gfc
aijha
ccdb
dfdb
abe
bbdg
gj
dhdbci
f
ddfhfe
3
ieij
a
icfid
4
dg
geceag
hcefi
bcigj
6
g
gbbh
hacd
fgfjfj
jih
ggcij
3
dd
a
g
2
ihcf
caag
7
cciec
b
dd
dbfadf
i
icfaie
ci
4
dhaiig
hia
jd
f
1
jagb
4
fjhb
hb
iidca
hfjh
8
ecjj
jhe
b
g
bfaead
cffcb
jbb
jd
1
jca
7
beh
hcjadg
efcjjh
bee
bb
ihi
c`

type testCaseC struct {
	words []string
}

func parseCases() ([]testCaseC, error) {
	in := strings.NewReader(testcasesRaw)
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return nil, err
	}
	cases := make([]testCaseC, T)
	for i := 0; i < T; i++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return nil, err
		}
		words := make([]string, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(in, &words[j]); err != nil {
				return nil, err
			}
		}
		cases[i] = testCaseC{words: words}
	}
	return cases, nil
}

func solve(tc testCaseC) int64 {
	weights := make([]int64, 10)
	leading := make([]bool, 10)
	pow10 := [7]int64{1}
	for i := 1; i < 7; i++ {
		pow10[i] = pow10[i-1] * 10
	}
	for _, s := range tc.words {
		leading[s[0]-'a'] = true
		L := len(s)
		for j := L - 1; j >= 0; j-- {
			letter := s[j] - 'a'
			pos := L - 1 - j
			weights[letter] += pow10[pos]
		}
	}
	best := int64(^uint64(0) >> 1)
	for zero := 0; zero < 10; zero++ {
		if leading[zero] {
			continue
		}
		type pair struct {
			w   int64
			idx int
		}
		pairs := make([]pair, 0, 9)
		for i := 0; i < 10; i++ {
			if i == zero {
				continue
			}
			pairs = append(pairs, pair{weights[i], i})
		}
		sort.Slice(pairs, func(i, j int) bool {
			if pairs[i].w == pairs[j].w {
				return pairs[i].idx < pairs[j].idx
			}
			return pairs[i].w > pairs[j].w
		})
		digits := make([]int, 10)
		digits[zero] = 0
		for i, p := range pairs {
			digits[p.idx] = i + 1
		}
		var sum int64
		for i := 0; i < 10; i++ {
			sum += weights[i] * int64(digits[i])
		}
		if sum < best {
			best = sum
		}
	}
	return best
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.words)))
		for _, w := range tc.words {
			sb.WriteString(w)
			sb.WriteByte('\n')
		}
		expected := solve(tc)
		gotStr, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(gotStr), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
