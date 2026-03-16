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

func runRef(input string) (string, error) {
	cmd := exec.Command("go", "run", "919B.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runBin(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) == 3 {
		bin = os.Args[2]
	}
	const testcasesRaw = `100
116
489
353
725
983
265
134
29
990
214
371
344
485
984
300
304
960
900
982
567
652
335
189
608
83
106
547
595
316
161
386
920
151
969
129
824
229
324
521
249
243
773
189
299
382
430
680
48
882
136
616
22
404
80
720
75
136
431
307
564
427
759
949
146
606
433
306
653
364
87
255
456
648
379
653
972
542
60
386
419
9
428
986
746
923
329
452
209
381
301
976
483
94
974
190
816
112
284
115
572`

	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Fprintln(os.Stderr, "bad file")
			os.Exit(1)
		}
		k := scan.Text()
		input := k + "\n"
		exp, err := runRef(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on case %d: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:%sexpected:%s\ngot:%s\n", caseIdx+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
