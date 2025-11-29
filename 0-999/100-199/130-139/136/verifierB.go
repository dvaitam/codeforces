package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded test data from testcasesB.txt.
const testData = `144272509 611178002
909925047 861425548
820096753 67760436
273878287 126614242
531969374 817077201
482637352 507069464
699642630 407608741
846885253 225437259
100780963 523832096
30437866 959191865
897395948 418554019
464680097 652231581
818492001 823729238
2261353 747144854
478230859 285970256
774747711 860954509
245631564 634746160
109765575 967900366
340837476 32845751
23968184 27322286
697444855 581337223
9883727 946217654
409314931 737106430
232571831 453244221
779378296 31182305
566537775 238039615
820017699 470178216
532374341 593628450
250272526 371192992
247891063 726760591
234914346 817061414
493495461 311150634
994828918 23074397
446869806 899342503
983837244 597488273
990192429 689658324
107374479 199615329
675762534 777001467
923360559 318246764
129804604 797947650
357228733 961616757
774687978 763636349
537729581 453233942
545157245 891244035
977303767 719735122
203849597 325739463
305113796 630909864
947554609 536185925
908597559 542544369
422360239 632436358
916210962 37071829
515639791 260640056
798574707 856206294
434101039 444866269
713762923 185765286
394196212 589268179
947826293 754884265
833049334 724223642
792652820 402334307
92843870 471331461
712704513 545918789
115890309 835846392
175769705 559353361
901891103 422254446
397845686 525804415
786801296 31755873
503928666 46694123
331280949 755250767
910856872 660147977
636926179 620811651
422624440 694878646
182911059 181026748
539274549 243672113
13208723 827342927
214229068 579409818
987935283 923729113
588773950 249297217
434280104 551658122
369180232 909954310
620402451 379325246
492988938 976842008
289136634 707826512
588406554 653849522
783187487 6130128
411983601 841443398
920142113 880990038
951528077 795109485
550292614 868807354
138780515 556926566
834723866 602753419
220638116 457511382
60261934 516579138
934166291 391632347
612032126 595283749
214575938 541939476
443884919 520684370
873329535 383100304
444984939 371598338
1701610 578187200
579938223 669466698`

// Embedded solver logic from 136B.go.
func solve(a, c int64) int64 {
	var b int64
	base := int64(1)
	for a > 0 || c > 0 {
		da := a % 3
		dc := c % 3
		// Choose db so that (da + db) % 3 == dc.
		db := (dc - da + 3) % 3
		b += db * base
		base *= 3
		a /= 3
		c /= 3
	}
	return b
}

func parseTests() ([][2]int64, error) {
	lines := strings.Split(testData, "\n")
	tests := make([][2]int64, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid test line: %q", line)
		}
		a, err1 := strconv.ParseInt(fields[0], 10, 64)
		c, err2 := strconv.ParseInt(fields[1], 10, 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("failed to parse test line: %q", line)
		}
		tests = append(tests, [2]int64{a, c})
	}
	return tests, nil
}

func runCase(idx int, bin string, a, c int64) error {
	expect := solve(a, c)
	input := fmt.Sprintf("%d %d\n", a, c)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("test %d: runtime error: %v\nstderr: %s", idx, err, stderr.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
		return fmt.Errorf("test %d: failed to parse output %q", idx, outStr)
	}
	if got != expect {
		return fmt.Errorf("test %d failed: expected %d got %d", idx, expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(i+1, bin, tc[0], tc[1]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
