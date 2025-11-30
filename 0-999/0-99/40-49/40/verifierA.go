package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesA.txt.
const embeddedTestcasesA = `729 -212
552 823
-139 -918
-470 977
47 -5
-171 880
605 699
-379 982
-24 -267
194 826
859 -553
33 -715
-423 -714
547 -806
266 637
-487 863
90 444
659 232
847 -700
-365 -798
494 -849
840 741
400 -324
-34 146
-794 -276
-111 -353
251 311
869 -582
979 131
-24 -94
772 67
-467 -873
648 881
123 875
-972 -809
473 720
-184 454
689 607
368 280
-998 253
10 695
776 -318
-501 495
-334 441
782 -872
-609 878
162 -546
-512 645
981 -709
644 112
-83 -814
-836 -345
792 40
910 2
-777 -383
128 -404
447 -745
121 -319
668 888
106 -584
973 637
235 120
203 -411
-89 -813
221 634
-212 -351
178 -505
-406 -624
-613 682
-618 -933
254 344
-468 -25
-859 -817
390 551
-734 795
-694 891
-921 725
-836 839
432 890
698 107
399 -199
715 444
74 -436
68 662
-518 739
-560 833
391 207
690 945
-142 187
-437 -78
8 352
313 434
877 624
-269 -832
-336 254
-764 -4
202 290
-314 730
-611 -503
-967 498`

func solve40A(x, y int) string {
	d2 := x*x + y*y
	fs := math.Sqrt(float64(d2))
	if math.Abs(fs-math.Round(fs)) < 1e-9 {
		return "black"
	}
	k := int(math.Floor(fs))
	if k%2 == 0 {
		return "white"
	}
	return "black"
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewBufferString(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesA), "\n")
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Fprintf(os.Stderr, "case %d: expected 2 integers, got %d\n", idx+1, len(fields))
			os.Exit(1)
		}
		x, err1 := strconv.Atoi(fields[0])
		y, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			fmt.Fprintf(os.Stderr, "case %d: parse error\n", idx+1)
			os.Exit(1)
		}
		want := solve40A(x, y)
		input := fmt.Sprintf("%d %d\n", x, y)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
