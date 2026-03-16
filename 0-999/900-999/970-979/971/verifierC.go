package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func run(bin string, input string) (string, error) {
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
	const testcasesRaw = `tsdm
ld
t
g
wdp
xzbv
r
tdichcu
lnfbq
btdwmgil
sfwvgybz
fkq
dtovf
p
sqjulmv
erwao
kx
r
ehypl
jvlsutewjm
ucatgwk
hhu
mwvsnbmw
nyvwbfociw
oqp
tyabpkjob
grucxea
vnkagaw
v
tdgdtugji
dpm
ai
zzdieuqu
deiabb
virk
sbxwtu
wuounlrf
msja
eik
zlwcky
bbifesjlmr
jdp
bjfq
jm
jnddrp
kzzkdpdw
nbjkxvef
mzucczcgxh
m
d
rqjopzs
nclh
syfng
dcwaqo
dpmi
ubzg
edgomlredt
esmuvnqp
ppuvgr
hakwxkkbqe
tzems
wwzpc
qb
he
j
y
kfzeuolq
qqbscvz
ytcxnygjr
npzmtshzav
x
jqs
kcpij
mzmbfue
jxkb
p
eptwcvw
zln
t
mobdpyea
t
eukdwrulgm
dbtwotuk
vw
jzemzjxvzd
zgbzmolyg
lzucbbpi
q
sghcyuyqwq
qjdensn
dn
dn
xaz
nvnapkxi
lc
lw
l
fahlct
gag
xa`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := line + "\n"
		expected := reverse(line)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
