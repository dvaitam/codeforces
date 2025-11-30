package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
144272510
611178003
909925048
861425549
820096754
67760437
273878288
126614243
531969375
817077202
482637353
507069465
699642631
407608742
846885254
225437260
100780964
523832097
30437867
959191866
670312030
947340757
927158472
313642286
939433170
429422320
52562369
281407542
266676470
736893142
441740010
965475928
672517463
738713522
215763396
930634681
922238406
26196150
413425094
235078188
570545509
109168469
777409675
98237502
224417573
765765101
554800190
630365666
189340049
22187556
34580242
613711000
381168288
656840451
217980929
166520159
453771295
762828919
787454982
238202994
756712051
288097985
202069090
284637636
646799679
186274903
385415970
402455512
306495785
917098836
793593294
805568658
280388245
324584131
817323617
189176439
846448240
861144390
396834261
77924061
756992530
550534635
988115071
594439257
545913961
382709469
86388227
548946910
149999767
563455189
21308057
692409896
868526645
823172039
74109215
192057478
176410614
644804696
295494373
674373179
`

func solve(x int) ([]int, error) {
	if x < 0 {
		return nil, fmt.Errorf("negative x not supported in embedded solver")
	}
	vc := make([]int, 0, 32)
	for x > 0 {
		if x%2 == 0 {
			vc = append(vc, 0)
			x /= 2
		} else {
			if x%4 == 1 {
				vc = append(vc, 1)
				x--
			} else {
				vc = append(vc, -1)
				x++
			}
			x /= 2
		}
	}
	if len(vc) == 0 {
		vc = append(vc, 0)
	}
	return vc, nil
}

func parseTests(raw string) ([]int, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	res := make([]int, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("bad testcase line: %q", line)
		}
		res = append(res, v)
	}
	return res, nil
}

func buildInput(x int) string {
	return fmt.Sprintf("1\n%d\n", x)
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

func check(x int, out string) error {
	scan := bufio.NewScanner(strings.NewReader(out))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return fmt.Errorf("missing n")
	}
	n, err := strconv.Atoi(scan.Text())
	if err != nil || n < 1 || n > 32 {
		return fmt.Errorf("invalid n")
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if !scan.Scan() {
			return fmt.Errorf("missing array element")
		}
		v, err := strconv.Atoi(scan.Text())
		if err != nil || (v != -1 && v != 0 && v != 1) {
			return fmt.Errorf("invalid element")
		}
		arr[i] = v
		if i > 0 && arr[i] != 0 && arr[i-1] != 0 {
			return fmt.Errorf("consecutive non-zero elements")
		}
	}
	if scan.Scan() {
		return fmt.Errorf("extra output")
	}
	sum := int64(0)
	pow := int64(1)
	for i := 0; i < n; i++ {
		sum += int64(arr[i]) * pow
		pow <<= 1
	}
	if sum != int64(x) {
		return fmt.Errorf("representation does not match")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, x := range tests {
		input := buildInput(x)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := check(x, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected, _ := solve(x)
		// optional: quick sanity check using embedded solver output when lengths differ
		_ = expected
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
