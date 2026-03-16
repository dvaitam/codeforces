package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)
const testcasesARaw = `
100
885452
403970
794784
933500
441013
42462
271505
536122
509544
424616
962850
821884
870175
318058
499760
375453
611732
934985
952237
229065
529214
146051
295540
146546
792530
99449
648418
838246
262686
953950
558445
739438
849586
631152
946001
154112
325225
103572
765296
77336
942512
891798
717221
346248
495089
587019
105604
370989
455274
331568
640573
671544
957373
214422
579375
500193
464209
907355
546690
273157
65316
844144
963092
575364
960501
14735
97814
754677
880911
418208
744766
864924
823194
700621
655650
1210
641632
517565
868299
909759
349329
255771
765764
341013
737834
912767
66055
200360
961576
595090
232485
250218
842380
149428
842206
569378
469742
95658
84365
335613
`


func isComposite(x int) bool {
	if x < 4 {
		return false
	}
	for i := 2; i*i <= x; i++ {
		if x%i == 0 {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data := []byte(testcasesARaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	inputs := make([]int, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		inputs[i] = n
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewReader(bytes.NewReader(out))
	for i := 0; i < t; i++ {
		var x, y int
		if _, err := fmt.Fscan(outScan, &x, &y); err != nil {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		n := inputs[i]
		if x+y != n || !isComposite(x) || !isComposite(y) {
			fmt.Printf("test %d failed: n=%d x=%d y=%d\n", i+1, n, x, y)
			os.Exit(1)
		}
	}
	if _, err := fmt.Fscan(outScan, new(int)); err == nil {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
