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

const testcasesARaw = `865
395
777
912
431
42
266
989
524
498
415
941
803
850
311
992
489
367
598
914
930
224
517
143
289
144
774
98
634
819
257
932
546
723
830
617
924
151
318
102
748
76
921
871
701
339
484
574
104
363
445
324
626
656
935
210
990
566
489
454
887
534
267
64
825
941
562
938
15
96
737
861
409
728
845
804
685
641
2
627
506
848
889
342
250
748
334
721
892
65
196
940
582
228
245
823
991
146
823
557
`

func isPrime(x int64) bool {
	if x < 2 {
		return false
	}
	for i := int64(2); i*i <= x; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func solveA(n int64) int64 {
	for m := int64(1); m <= 1000; m++ {
		if !isPrime(n*m + 1) {
			return m
		}
	}
	return -1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]


	scanner := bufio.NewScanner(strings.NewReader(testcasesARaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, _ := strconv.ParseInt(line, 10, 64)
		input := fmt.Sprintf("%d\n", n)

		cmd := exec.Command(exe)
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
		got := strings.TrimSpace(out.String())
		m, err := strconv.ParseInt(got, 10, 64)
		if err != nil || m < 1 || m > 1000 || isPrime(n*m+1) {
			fmt.Printf("test %d failed: invalid output %q\n", idx, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
