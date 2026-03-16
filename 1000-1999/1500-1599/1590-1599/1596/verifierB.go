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

const testcasesRaw = `199 470
67 719
356 764
397 540
415 817
261 6
983 286
15 53
87 688
957 405
666 433
928 45
923 171
553 749
106 257
435 400
275 897
109 354
111 403
118 617
249 21
755 54
585 541
233 355
507 67
492 125
696 25
271 221
440 152
303 554
786 82
35 412
356 406
195 265
703 891
34 647
256 119
184 785
469 338
694 485
934 698
616 947
860 85
448 469
504 95
289 356
333 506
589 592
877 432
960 118
869 367
544 237
281 489
779 904
644 652
807 917
965 746
574 324
87 371
494 784
825 166
173 319
298 690
565 436
885 415
49 989
531 493
584 560
429 834
478 405
731 599
499 674
168 342
240 328
564 992
264 979
298 688
424 593
992 625
131 804
922 987
181 972
379 421
629 428
161 521
516 605
403 90
662 301
212 975
435 371
528 754
3 624
825 0
413 819
871 361
834 929
767 159
938 91
992 578
185 476`

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

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
			panic("invalid test case")
		}
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		expected := gcd(a, b)
		input := fmt.Sprintf("%d %d\n", a, b)
		out, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		if out != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
