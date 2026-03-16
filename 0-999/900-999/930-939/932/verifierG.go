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
	oracle := filepath.Join(dir, "oracleG")
	cmd := exec.Command("go", "build", "-o", oracle, "932G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func buildInput(fields []string) (string, error) {
	if len(fields) != 1 {
		return "", fmt.Errorf("expected 1 field")
	}
	return fields[0] + "\n", nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	const testcasesRaw = `pyibae
yxlkyaipzgxnrrvdgs
wzxivztvcnkclznziowd
wuzjdbsgul
gqsuwzqaulhtnjlsdc
vqgdtvijxgmphetgwqag
aukrvttj
mqmjevpbfntxmd
hzctvoozmzcqnpjwnc
xviopxzfaa
diszlgiqokqinntpitpv
pexmpjuzoklvftmwti
umpxzfjratobfa
vwllps
geit
nzzz
pcphdllevuhtkebuvdwd
upox
uawret
ftuiuo
wazuegopyooripwscz
ljlbcqiiipno
bwchtuhrxbxfkl
bxyy
wnkcfxvdmr
oxpucecl
afqwefjxeqsmzhhwswen
ykhrsazz
gykmqoamqrvewhur
ivfefwuemcbabp
duwr
vltbbr
znhyswjcpc
fvod
wfbvay
uwjiohguyrzugisa
qzxxwvkyqcst
llwd
yalwknvble
akcp
oicv
heenrfhvmtns
upiiejbgyodzfsdzaqfd
smppqywr
gumeyoaknikejlsvlh
yjfc
ttyhucgmneyeqbqq
lkwsmiuvzciebdsttx
qqdzwyyc
vydglk
xhffhuep
tdxteozgslxekuiciy
buap
kbaixdmtkdxleyelvy
xmzaredolaevrpicpd
wgzrfokizryljrwi
rwybqxhvtaqprrifbk
gsjg
ibixjwlg
taclsflf
uhpmrtlduxcnbstgto
hubpvuqqiqiu
jznxldyrnufibu
kvqbcvdbonrcfvilaciw
lmmswwcxdjzp
xijhyabyqrmehzjuux
dize
damipgcxiyfeedazgtkn
tspoktogzfrdhvemhs
osyuistgshxanh
sxdubklzpuinmjlq
sbnusfjjsvphdo
lsnvozkc
npadtcpnubuaxqqllf
prcmstda
uhtxritqzizmlw
tlnj
pzlbarzg
ybap
rnpkpg
ledgbhinhl
lelwyuxigxrbqggmqb
uaucxcdvtvzgqydq
kvfogyyahkgrsjetaa
dmqgvgxnrtru
cntrfevtguwjnjwmvb
ggjhxnlybuklxlmnlq
jzetutiugtoqfura
chriotag
jbrbpmdlkejrasizyeiu
cxuqtldstahjtg
cnatlzkgsy
mzsxbbfurthqhjfi
tlljsyhf
eszrob
crbllo
iqyhnssnkuomzrrt
bxcnvqmdtq
kstvfdwcdggpvxhzra
ufbbmajowywytd`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		input, err := buildInput(fields)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad testcase %d: %v\n", idx, err)
			os.Exit(1)
		}
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
