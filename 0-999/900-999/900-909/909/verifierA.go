package main

import (
	"bufio"
	"bytes"
	"fmt"
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
	oracle := filepath.Join(dir, "oracleA")
	cmd := exec.Command("go", "build", "-o", oracle, "909A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesRaw = `biqpmzj plsgqej
tzi rw
jdxcvkprdl nkt
poqi bzracxmwz
p khxkwcgshh
occ kqpdjrjwd
rgztrsjoc tzmksh
gfbtv ipc
eb cwr
iqzhgvs nsiopvuwz
ktdpsu kg
xidw h
knbdze whb
tvcadugtsd mcldbtagf
gx zbvarntd
hcujl nf
o btdwmgilx
fwvgybzv ffkqidtovf
v nsqjulmv
rwaox ckx
i ehypltjvl
ewjmxnucat gwkfhhuomw
bmwsnyvwbf ociwfoq
tyabpkjo bzzngrucx
mvn k
a wyav
dgdtugjiw fdpmucaioz
eu quuld
abb gvirk
bxwtup wuounlrfgm
aeeikkzlwc kytbb
esjlm rej
xh bjfqxcjm
nddrpp kzzkd
wpnbjkxv ef
zucczcgxhb madmrqj
zswvgncl hisyfngl
wa qo
pmig ub
edgo mlredtpesm
pvkppuv grthakwxk
qeitze m
wwzpczcqbc hebja
fzeuolqm qqbscv
cxnygjrtn pzmtshzava
qsi kcpij
zmbfueh jxkbbpn
twc vwezlnbt
obdpyeab tteukdw
gmzyypdbt wotuku
jz emzjxvzdqz
zmol y
lzuc bbpiaqvs
hcyuyqwqnq jden
cdncdnyexa zonvnap
clcdlw allfa
cteg agvvxd
l wathe
dpl wieaglkp
rukfs cdrsj
eez hkqhhyf
nvbet amcwce
rnxesnj ulcho
bmnanx kogljpcfz
dr twezw
fynnfhok qelouuc
jawotoag jdyu
tenwy pcvpyhrym
divba i
swmodxi ljyvgtcbc
rkdqh yfcnj
esqug rdnurmxyz
olsue fdwdm
oervjlu pxngppwqkp
j expbtgal
aqcvcvxv malbdta
xheys jgdno
fknuvn eoweqke
olzm nzp
zgogswb mbhu
lb xuv
tjt cwqy
nobuw qvurxn
piwpgkibbb flajuaec
tmrhogk tdtc
rokiaq bglcgq
givxxj jqmi
whbjrcao pxobzn
odcchdye ngotcnry
fhpheil k
rjtrzgw jy
toruiiha dtzwdfxnh
xvax rqnbd
dxslhvw wrvjh
qjvk hl
sfezarqk lsuaz
efq ceygzyp
hxehymltse updt
t lpojahruf
rkwcietm wgkzjmbg`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		input := "1\n" + parts[0] + " " + parts[1] + "\n"
		// run oracle
		cmdO := exec.Command(oracle)
		cmdO.Stdin = strings.NewReader(input)
		var outO bytes.Buffer
		cmdO.Stdout = &outO
		if err := cmdO.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "oracle run error: %v\n", err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(outO.String())

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed. Expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
