package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"
)

// Embedded source for the reference solution (was 1145F.go).
const solutionSource = `package main

import (
   "fmt"
)

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   inSet := map[rune]bool{
       'A': true, 'E': true, 'F': true, 'H': true,
       'I': true, 'K': true, 'L': true, 'M': true,
       'N': true, 'T': true, 'V': true, 'W': true,
       'X': true, 'Y': true, 'Z': true,
   }
   allIn, noneIn := true, true
   for _, c := range s {
       if inSet[c] {
           noneIn = false
       } else {
           allIn = false
       }
   }
   if allIn || noneIn {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
`

const testcasesRaw = `EMSYAFSAVM NO
EAA YES
QABBYD NO
TEYEVZMZAN YES
SVKWHEL NO
GRMCENSVL NO
NN YES
PMHM NO
UPMS NO
II YES
LRATTYPHI NO
T YES
ZMUDRB NO
WMA YES
XMNDWOT NO
FFKPNFSJ NO
DLLEULYPU NO
YBGIFXSKJ NO
UBJRNBV NO
IMXGLEE NO
TL YES
ANS NO
OCUWWVC NO
RXREFEG NO
HAQ NO
PLT NO
WZKVD NO
YIFKUQK NO
WEMXXRJHM NO
MPQJNN NO
VE YES
AST NO
DWYUWGTUT NO
DIWXTFMCZ NO
A YES
LX YES
KDVOLSWI NO
ZYHFRRCY NO
FAUFQNTVG NO
WVMIASEM NO
OSB YES
ZCUSMKH NO
OBPTDIQPR NO
PIFHRLF NO
TEOCC NO
MSNRCIPH NO
JE NO
ZDEVBV NO
SRG YES
B YES
RWYTPWW NO
PR YES
YKDVAH NO
PJIH NO
P YES
QKCCJS NO
HXLMYEH YES
GXYPV NO
JMTEYD NO
LQPHUWL NO
NILMYY YES
DPJDO YES
LHX YES
KPH NO
VM YES
OQOSTHV NO
QJPHKQV NO
C YES
KMHNBSYB NO
CIGXKFD NO
ZWL YES
H YES
A YES
RAEDTTG NO
XO NO
AQTN NO
RF NO
HNMP NO
N YES
MYBT NO
ASLWL NO
VOUETQ NO
ID NO
IA YES
TYV YES
MGS NO
GNQQDR NO
WP NO
QO YES
FORKENIM NO
SQ YES
HOHLPY NO
AOXARMY NO
HNHIPPEH NO
JLUPTEQV NO
GJ YES
UVDBEKBKT NO
YOM NO
NPGF NO
BZZYKYV NO`

func main() {
	var _ = solutionSource
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	total := 0
	passed := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Println("invalid test case:", line)
			os.Exit(1)
		}
		if utf8.RuneCountInString(parts[0]) != len(parts[0]) {
			fmt.Println("non-ascii input not supported in verifier:", parts[0])
			os.Exit(1)
		}
		input := parts[0] + "\n"
		expected := parts[1]
		var cmd *exec.Cmd
		if strings.HasSuffix(bin, ".go") {
			cmd = exec.Command("go", "run", bin)
		} else {
			cmd = exec.Command(bin)
		}
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		result := strings.TrimSpace(string(out))
		total++
		if err != nil {
			fmt.Printf("Case %d: runtime error: %v\n", total, err)
			fmt.Printf("Output: %s\n", result)
			continue
		}
		if result == expected {
			passed++
		} else {
			fmt.Printf("Case %d failed: expected %s got %s\n", total, expected, result)
		}
	}
	fmt.Printf("%d/%d cases passed\n", passed, total)
}
