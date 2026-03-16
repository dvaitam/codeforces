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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesRaw = `38439 48383
90370 3231
79279 30522
18598 24510
59513 14732
62482 45148
92731 33846
17061 3662
27279 47480
43911 62046
38356 38843
72510 83351
42859 24124
77729 10592
13446 69892
76134 40340
20501 49361
19258 16415
29206 41405
66629 31830
31020 98941
24110 38150
48820 55023
86927 6063
17333 78802
2694 51617
10214 92035
9596 17305
55081 39242
72192 54625
97109 18650
77474 55338
39064 83489
46474 11083
32514 58303
82940 48400
83499 69365
7584 49340
53560 1105
54687 95445
42031 57842
26738 48706
38445 61721
11935 24309
14237 36304
14701 73162
79371 90204
20170 92091
58480 52269
24299 55274
56586 22904
32501 59435
44616 68597
18687 46579
60613 82775
83533 11338
63346 98964
26692 38631
242 91619
58863 81099
60557 1020
28670 39133
15003 82434
39474 71462
79844 20477
55602 92576
98482 61792
12132 88891
65232 99627
30457 71270
99877 53125
36712 82847
2834 15838
35400 87621
5315 34
33630 52223
68958 76234
93057 51936
58273 13374
97852 33110
46377 37146
98989 88153
25695 78042
11176 4647
9235 34405
40057 69946
44551 15491
69493 32644
99877 21438
8926 54373
37971 37061
68136 17619
75233 68559
82201 27562
69649 13803
53845 83178
71248 52877
97181 36519
38307 57982
48734 74480`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		input := fmt.Sprintf("%d %d\n", a, b)
		expected := strconv.Itoa(gcd(a, b))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
