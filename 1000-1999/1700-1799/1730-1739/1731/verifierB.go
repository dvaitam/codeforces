package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod int64 = 1000000007
const inv6 int64 = 166666668

type testCase struct {
	n int64
}

// solve embeds the logic from 1731B.go.
func solve(tc testCase) string {
	x := tc.n % mod
	t1 := (4 * x % mod) * x % mod
	t2 := (3 * x) % mod
	term := (t1 + t2 - 1) % mod
	if term < 0 {
		term += mod
	}
	ans := x * term % mod
	ans = ans * inv6 % mod
	ans = ans * 2022 % mod
	return strconv.FormatInt(ans, 10)
}

// Embedded copy of testcasesB.txt.
const testcaseData = `
695470141
588687857
811193832
271559602
464111031
565813033
908827385
992905248
202332127
360080634
696713113
963163988
542118850
928410378
519381548
998843703
863912864
3118816
822190380
164961036
415102559
797865291
791642876
268161558
756423301
140153964
860793501
760967030
139561573
870773673
826746213
348153325
500895392
806744962
22397959
34309796
978583180
87878182
5445349
26541720
514729605
136945736
690676201
843291007
828116152
930803107
982994678
692063543
525515595
886648342
751740872
410477703
178120310
479958295
121329000
693191618
333092995
803479174
12904270
903987313
686798995
804812900
190570275
724207784
670019792
634519131
265462208
8540028
868942313
951905717
91791632
469252172
308893844
51206247
984572442
107673405
885582183
11444945
532617784
102954885
152202601
724467275
218033108
47495980
725735507
387224775
660856888
708114410
911617769
446034646
18666177
957315721
408070459
161778776
264908105
94187102
60818775
694281599
589241097
929783836
`

var expectedOutputs = []string{
	"150400642",
	"388016161",
	"128365114",
	"709207592",
	"552330588",
	"580192334",
	"258009936",
	"400151430",
	"739823690",
	"893674733",
	"184833083",
	"816923806",
	"89608014",
	"734566032",
	"853816906",
	"932466968",
	"114854504",
	"463478747",
	"488703630",
	"456154297",
	"217538492",
	"252713112",
	"355638492",
	"544758211",
	"41058866",
	"920050290",
	"678108311",
	"743107223",
	"330131668",
	"438706328",
	"287463800",
	"990544446",
	"95162144",
	"993318146",
	"562744331",
	"800058788",
	"665802558",
	"670593828",
	"12638787",
	"529831256",
	"816977667",
	"36322681",
	"676289188",
	"504914660",
	"477282927",
	"516319036",
	"463924148",
	"766688277",
	"885686202",
	"821954906",
	"164378228",
	"154293018",
	"389925653",
	"856153225",
	"704068447",
	"629913876",
	"890354453",
	"519451233",
	"38364410",
	"551668161",
	"720023831",
	"959513978",
	"343708517",
	"354570804",
	"938211806",
	"236395002",
	"515423832",
	"48382302",
	"416318613",
	"83031131",
	"757464232",
	"49625026",
	"579316784",
	"33077576",
	"214025004",
	"542504735",
	"97684104",
	"50583559",
	"322262029",
	"714190963",
	"998299096",
	"69107449",
	"939443106",
	"618315606",
	"608012382",
	"908448660",
	"496532206",
	"541497062",
	"142983771",
	"666601157",
	"626950873",
	"90638115",
	"151377578",
	"470363170",
	"89652402",
	"631428977",
	"715170665",
	"471327119",
	"967940799",
	"466323265",
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		tests = append(tests, testCase{n: n})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	input := fmt.Sprintf("1\n%d\n", tc.n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}
	if len(tests) != len(expectedOutputs) {
		fmt.Fprintf(os.Stderr, "testcase/expected mismatch: %d vs %d\n", len(tests), len(expectedOutputs))
		os.Exit(1)
	}

	for i, tc := range tests {
		exp := solve(tc)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
