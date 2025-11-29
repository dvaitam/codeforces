package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
4019036145
1785165348
4721075254
3323526175
6213609078
7243110601
4424392847
8427767308
7398264892
4632287772
302062763
4421142744
1494123483
4184687620
7812574303
8892805415
3837274504
893258887
9469504383
9119023543
4325824772
2597803408
3690131154
5098090550
4777683105
8414577615
559023048
4151168623
5850775681
5551812114
7026204473
9380404208
4264695486
9030514218
6789696414
4966713327
3855486119
4061657773
3454830043
9946677907
1042984255
3242092208
5545058411
9157874477
4383227254
8924603261
314439520
6099848714
6660537266
9201037469
6108254151
7030715614
363163949
4543461724
1755036262
5672238776
5171098853
2022446623
4081371410
3421014592
1189578158
2955797805
1712733788
7591959182
1854208235
6242507469
4907290008
2737197801
5169594032
7297113850
6952416109
33400473
1282305980
6996159765
4965931852
2024797583
7207722819
9587924891
7567710967
9792905026
4226312847
4813935130
174148580
5396952056
5996783591
9028157507
5379910779
3903156018
2557266134
152260744
7697597044
9902493486
1459818207
702462341
5539184744
2232670495
2693544449
2282241725
7680163570
5550189389
`

func countDivisors(b int64) int64 {
	var cnt int64
	for i := int64(1); i*i <= b; i++ {
		if b%i == 0 {
			cnt++
			if b/i != i {
				cnt++
			}
		}
	}
	return cnt
}

func parseTestcases() ([]int64, error) {
	lines := strings.Split(testcasesData, "\n")
	var cases []int64
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, err
		}
		cases = append(cases, v)
	}
	return cases, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, b := range testcases {
		input := fmt.Sprintf("%d\n", b)
		expected := strconv.FormatInt(countDivisors(b), 10)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
