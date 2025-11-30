package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (one n per line).
const testcasesRaw = `729
-212
552
823
-139
-918
-470
977
47
-5
-171
880
605
699
-379
982
-24
-267
194
826
859
-553
33
-715
-423
-714
547
-806
266
637
-487
863
90
444
659
232
847
-700
-365
-798
494
-849
840
741
400
-324
-34
146
-794
-276
-111
-353
251
311
869
-582
979
131
-24
-94
772
67
-467
-873
648
881
123
875
-972
-809
473
720
-184
454
689
607
368
280
-998
253
10
695
776
-318
-501
495
-334
441
782
-872
-609
878
162
-546
-512
645
981
-709
644
112`

func expected(n int64) string {
	return fmt.Sprintf("%d", n*n)
}

func run(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func loadCases() ([]int64, error) {
	lines := strings.Fields(testcasesRaw)
	cases := make([]int64, 0, len(lines))
	for _, line := range lines {
		n, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, err
		}
		cases = append(cases, n)
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load embedded testcases: %v\n", err)
		os.Exit(1)
	}
	for i, n := range cases {
		input := fmt.Sprintf("1\n%d\n", n)
		want := expected(n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
