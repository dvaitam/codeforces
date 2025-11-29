package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesRaw = `tjn tjn
r r
nsglgtoc nsglgtc
qhl qhl
h h
tj tj
jnzmfm jnzmfm
svhs svhs
yktgringj yktgrngj
gidmo gdm
t t
hmc hmc
q q
zn zn
huno hn
ekwtmquxg kwtmqxg
akr kr
bejhpj bjhpj
er r
nbl nbl
epxddzn pxddzn
nsggay nsggy
rcchhpnhj rcchhpnhj
unvsmxp nvsmxp
mgh mgh
wdedhp wddhp
txnpgjrs txnpgjrs
ybhlgds ybhlgds
ngeyz ngyz
gxhtfspw gxhtfspw
epfvyhv pfvyhv
qqaupa qqp
vfmnit vfmnt
exihpyk xhpyk
psfwiykyj psfwykyj
geueepx gpx
oazovxtwue zvxtw
jxsubwma jxsbwm
ozblnw zblnw
izdzwtkik zdzwtkk
mseplycw msplycw
dsatoga dstg
ryx ryx
w w
kmqfpfpxgn kmqfpfpxgn
qbmrdw qbmrdw
ew w
ixfibcmb xfbcmb
qwtrgeh qwtrgh
audbscjp dbscjp
pkmmlx pkmmlx
kdof kdf
rcoij rcj
fgs fgs
grm grm
cts cts
cm cm
cdz cdz
ovtjjrnlqd vtjjrnlqd
xsdbvcouk xsdbvck
kdbowaglpc kdbwglpc
ymz ymz
kqglgh kqglgh
oiggcu ggc
cfqccan cfqccn
svofcjyob svfcjyb
f f
repts rpts
swoejp swjp
ibhexpsaj bhxpsj
litv ltv
frpzul frpzl
juktb jktb
dwendqfcg dwndqfcg
sqkh sqkh
veqzve vqzv
btex btx
ddasdcvw ddsdcvw
rjoiaugn rjgn
srciqrocv srcqrcv
ya y
vl vl
kykcvznfsf kykcvznfsf
erhwzmsw rhwzmsw
pez pz
cqnsnnw cqnsnnw
b b
lzmc lzmc
alicz lcz
gspsmfhz gspsmfhz
gzc gzc
o -
lwv lwv
zvke zvk
vsuejaug vsjg
stspjaaz stspjz
cebuhunenn cbhnnn
h h
xmirs xmrs
vwo vw`

// referenceSolve mirrors 1089L.go.
func referenceSolve(s string) string {
	vowels := map[rune]bool{
		'a': true,
		'e': true,
		'i': true,
		'o': true,
		'u': true,
	}
	var res []rune
	for _, ch := range s {
		if !vowels[ch] {
			res = append(res, ch)
		}
	}
	if len(res) == 0 {
		return "-"
	}
	return string(res)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input + "\n")
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierL.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	caseNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Printf("Case %d: invalid test line\n", caseNum+1)
			os.Exit(1)
		}
		in := parts[0]
		expected := parts[1]
		caseNum++
		if err := runCase(bin, in, expected); err != nil {
			fmt.Printf("Case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d cases passed\n", caseNum)
}
