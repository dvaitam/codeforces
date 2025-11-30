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

const testcases = `
864
394
776
911
430
41
265
988
523
497
414
940
802
849
310
991
488
366
597
913
929
223
516
142
288
143
773
97
633
818
256
931
545
722
829
616
923
150
317
101
747
75
920
870
700
338
483
573
103
362
444
323
625
655
934
209
989
565
488
453
886
533
266
63
824
940
561
937
14
95
736
860
408
727
844
803
684
640
1
626
505
847
888
341
249
747
333
720
891
64
195
939
581
227
244
822
990
145
822
556

`

func referenceSolve(n int64) string {
	return strconv.FormatInt(n*(n+1)/2, 10)
}

func parseCases() ([]int64, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcases))
	scanner.Split(bufio.ScanWords)
	var res []int64
	for scanner.Scan() {
		val, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return nil, err
		}
		res = append(res, val)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), stderr.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	values, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, n := range values {
		input := fmt.Sprintf("%d\n", n)
		expected := referenceSolve(n)
		out, stderr, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(values))
}
