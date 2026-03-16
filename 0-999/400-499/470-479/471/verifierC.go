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
const testcasesCRaw = `
641178
154870
816134
160592
706312
825072
398328
825967
27727
439078
457456
595638
958079
719776
349096
745753
260004
967350
144527
381710
535039
227188
559257
419837
76597
136103
427566
593609
690678
367624
101248
452504
939215
457217
257578
494020
400304
235965
960850
414018
251716
676610
504557
416347
965129
609617
72371
876172
264628
289277
928032
552166
389858
568685
22783
633375
641314
817465
497603
919341
248531
289703
42262
640645
337191
909465
836839
409758
656516
111842
560435
920272
893977
51085
151213
748134
411493
28098
805039
441503
766512
912504
412204
453988
969476
109911
747840
483065
633534
484430
169219
177848
357036
497402
431066
166734
618875
963026
906590
297492
798838
528024
117982
978792
386173
362200
149625
659140
374994
803348
495640
660886
543915
790030
56524
205199
919747
270148
185491
758600
`


func expectedC(n int64) int64 {
	var count int64
	for h := int64(1); ; h++ {
		req := h * (3*h + 1) / 2
		if req > n {
			break
		}
		if (n+h)%3 == 0 {
			count++
		}
	}
	return count
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesCRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, _ := strconv.ParseInt(line, 10, 64)
		expect := expectedC(n)
		input := line + "\n"
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
		gotStr := strings.TrimSpace(out.String())
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil {
			fmt.Printf("test %d: invalid output %q\n", idx, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
