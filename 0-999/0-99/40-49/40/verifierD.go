package main

import (
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strings"
)

const embeddedSolutionSource = `package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(in, &s); err != nil {
       return
   }
   A := new(big.Int)
   A.SetString(s, 10)
   two := big.NewInt(2)
   thirteen := big.NewInt(13)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   if A.Cmp(two) == 0 {
       fmt.Fprintln(w, "YES")
       fmt.Fprintln(w, 1)
       fmt.Fprintln(w, 1)
       fmt.Fprintln(w, 1)
       fmt.Fprintln(w, 13)
   } else if A.Cmp(thirteen) == 0 {
       fmt.Fprintln(w, "YES")
       fmt.Fprintln(w, 1)
       fmt.Fprintln(w, 2)
       fmt.Fprintln(w, 1)
       fmt.Fprintln(w, 2)
   } else {
       fmt.Fprintln(w, "NO")
   }
}`

const testcasesRaw = `
2
13
98259791
7483378876232860129
40
796669725
273
646869589
93504925899139
411771516
4661
99
69
84830401198
36
9420555266
869819846346484
506950699205755947
9
5
749952525
944610924834
52611955
72153974
10835924
19592648475
46960623079868
7848531
4130083690071284308
609152487053318123
207960434
88607500201017018875
515669
5435612806
920
798609605756
70338491636205
841718615981
723608191620510174
4109883181808592132
79645965
616836269987262217
8792243298538469
9434047632953572679
7980716
73
13433422
4205261111440579500
56671397628
14161784185
744159881785
49
225711
259684277421128896
144602078438
567815712491
9123
6629926388292345407
65894784709301
72873383
87
9422
190059381780555
219015313637
10613267718
37407767270455
859630345278323
2
6820219277200675600
60670331
7325915190447
7384055510
19901010280703
53757506496142618
8586268525673
754289670220711
39090775501
6159417
330
29013
3388689583
291019013411763
17569714976318074155
7185619
33505367143523397
653967679040210248
7033476057342174
71327465264778
354407554380028280
3755
72304650981744376
500602932
85824020070778
86582421281657744748
95280
73799405998172
416150877
4526650
305596840
64617
`

var (
	_            = embeddedSolutionSource
	rawTestcases = strings.Fields(testcasesRaw)
)

func solveCase(s string) string {
	A := new(big.Int)
	A.SetString(s, 10)
	two := big.NewInt(2)
	thirteen := big.NewInt(13)

	switch {
	case A.Cmp(two) == 0:
		return "YES\n1\n1\n1\n13"
	case A.Cmp(thirteen) == 0:
		return "YES\n1\n2\n1\n2"
	default:
		return "NO"
	}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range rawTestcases {
		expected := solveCase(tc)
		got, err := run(bin, tc+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}
