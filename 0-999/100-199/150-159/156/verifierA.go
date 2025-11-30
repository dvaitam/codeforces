package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesA.txt.
const testcaseData = `
ynbiqpm plsgq
jey tz
rwzte dxcvk
rdlnktug poqibzrac
wzvuatp hxkwcg
hhzezrocck pdjrjwdrk
gztrsjoct kshjfgf
t pccvy
ebc vmwqiqzhg
nsiopvuwzl kt
ps ghaxid
lzfk bdzewhb
urtvcadugt dmcldbtagf
pg v
r tdichcu
lnfbq btdwmgil
sfwvgybz fkq
dtovf p
sqjulmv erwao
kx r
ehypl jvlsutewjm
ucatgwk hhu
mwvsnbmw nyvwbfociw
oqp tyabpkjob
grucxea vnkagaw
v tdgdtugji
dpm ai
zzdieuqu deiabb
virk sbxwtu
wuounlrf msja
eik zlwcky
bbifesjlmr jdp
bjfq jm
jnddrp kzzkdpdw
nbjkxvef mzucczcgxh
m d
rqjopzs nclh
syfng dcwaqo
dpmi ubzg
edgomlredt esmuvnqp
ppuvgr hakwxkkbqe
tzems wwzpc
qb he
j y
kfzeuolq qqbscvz
ytcxnygjr npzmtshzav
x jqs
kcpij mzmbfue
jxkb p
eptwcvw zln
t mobdpyea
t eukdwrulgm
dbtwotuk vw
jzemzjxvzd zgbzmolyg
lzucbbpi q
sghcyuyqwq qjdensn
dn dn
xaz nvnapkxi
lc lw
l fahlct
gag xa
lwath fod
lwieaglk jjrukfsc
rs fmeez
kqhh jln
e amcwcenjrn
snj choulu
bmnanxkog jpcfzd
drtwe mfynnfho
qelouu py
jawo oagjdyujrt
nwy cvpyhrym
uadiv a
mqswm dxiljyvg
cbczijrkdq yfcn
jqesq rdnu
mxyzijols fdw
mm oervjlupxn
ppwq puboje
btgalpma cvcvxvmal
d aiuwjxheys
gdnow mfknuv
eoweqke folz
nzpmxhz ogsw
m h
fl x
htj cwqyjylnob
vurxnsopi gkibbbfl
j e
zn tmrh
gktdtczk rokiaq
g cgqlgg
vxxjj miplwhbjr
ao xobznpoo
cc dyen
otcn ymbfhphei
kndrjt zgwjyoqto
uiihadtzw fx
hgjxvax qnbdmuidx
lhvwwrvjhx cqjv
`

// solve mirrors 156A.go.
func solve(s, u string) int {
	n := len(s)
	m := len(u)
	best := 0
	for d := -m + 1; d < n; d++ {
		i := d
		j := 0
		if i < 0 {
			j = -i
			i = 0
		}
		count := 0
		for i < n && j < m {
			if s[i] == u[j] {
				count++
			}
			i++
			j++
		}
		if count > best {
			best = count
		}
	}
	return m - best
}

func parseTestcases() ([][2]string, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([][2]string, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("case %d: expected 2 strings", idx+1)
		}
		res = append(res, [2]string{parts[0], parts[1]})
	}
	return res, nil
}

func runCandidate(bin string, s, u string) (string, error) {
	input := fmt.Sprintf("%s\n%s\n", s, u)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := strconv.Itoa(solve(tc[0], tc[1]))
		got, err := runCandidate(bin, tc[0], tc[1])
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
