package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesRaw = `gbo
ennc
lvujl
zlipukb
hozcwjfwcw
qd
vdeshi
biw
othlretw
igc
flfqog
fnvalpmh
nyyj
vuywnfjn
knas
cdx
ehzio
uempxtc
hkqyerhlb
redkcln
xst
w
tgwryt
s
oglenhd
pwylcn
mzbrucfrz
lkpai
whlniuncj
t
tcedctne
bwesjqs
bji
ylxg
onnusjjerf
xslcioto
rujzgryhah
lbtior
sfwpvry
iuxa
pentc
bmbvgtu
vvli
kdaotatyoo
iyxntui
lfdop
e
jjecu
rkyoxwxxn
htig
dqbbtau
txnopibuk
gzwxi
nszl
ur
rvjzoo
zjydvnwnj
zfdewqd
bjsgsoqxj
kkvy
xfvxz
syrzmygaj
s
hnjjddef
mclcmsvkdg
kmnjrtyaz
dm
gu
qrskt
hy
jwzlhch
bzwac
uhrmxjtsa
cqkrid
mwhjqun
dlkjgpc
foe
zjldz
ycnttyr
j
dxcue
kjudi
wnt
jsz
yebcdjsoxm
wupwtgi
yqcchdh
leaiujghof
hjtspggmjk
kd
iiksewd
umk
xg
taxdqruoy
zbgfenpi
xmgjczzohn
plkxqaceqy
kvvlxb
kdmygs
foduy`

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func isPal(s string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" {
			continue
		}
		idx++
		expected := "NO"
		if isPal(s) {
			expected = "YES"
		}
		input := s + "\n"
		out, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		if out != expected {
			fmt.Printf("Test %d failed: expected %s got %s\n", idx, expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
