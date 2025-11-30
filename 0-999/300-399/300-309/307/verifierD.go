package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded testcases from testcasesD.txt.
const embeddedTestcasesD = `c
lf
itgtb
vfnumzxqlr
qibalokm
qfrfhha
kfe
qlqvrfozn
ylzsllofy
wxouqhp
pqqzl
olsxrxop
kwft
ypjjz
rqqutsnjx
pqlv
czkxagxdbs
i
hvdyqeihgb
wybbllf
vacd
ab
l
efxfq
m
bzhebaltux
jk
ajorytxb
ymtwe
hcvvkdao
qsy
pqkekii
nuawrevbib
ffd
uhqwbhhw
cicshtzz
wlivniqyae
m
fdqxchd
af
dgaq
ojrumgvy
xznn
assbnqsfd
laqdtljw
javndd
gyvaz
bnupogst
aj
l
xchyp
dslm
eylmdidd
tk
gwdatvp
x
jloezlip
pxxznpvjm
fpti
nwvwcsxsd
lf
eznczcvzu
e
mhwvv
ofqjde
yndkqhwqi
fow
mlzy
xeooxaztmx
mqb
imiwxnwu
lrkwxvcy
rtgm
uakoqwo
dam
xstm
dmry
ixss
pzte
t
piqsfow
ycla
rvvcyspv
oiqoac
ylfyyzmivu
bfp
ovjeajr
albrmsog
puepw
wjcikjkuz
uumqc
ugmtqezqu
jb
orhq
bddvz
lgklcko
fpojoe
ugikfdhp
yvlf
ezehiz`

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func runCandidate(bin, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesD), "\n")
	for idx, line := range lines {
		input := line + "\n"
		want := reverse(line)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
