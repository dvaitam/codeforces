package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var rawTestcases = []string{
	"2 5 8272 33433 15456 64938 99741 58916 61899 85406 49757 27520",
	"1 4 3716 51094 56724 79619",
	"1 4 34909 94574 29985 77484",
	"1 3 4010 2926 3336",
	"5 1 49966 89979 28391 55328 95139",
	"1 5 29058 57395 64988 72465 30551",
	"3 2 88716 28677 99739 60242 37983 2817",
	"4 5 84187 13108 24368 82491 94849 38849 15846 97406 43608 94567 93218 65641 55327 66548 87859 24884 39764 37246 77016 65453",
	"5 4 77202 4526 62945 31817 97483 52991 54305 87130 22677 48120 71933 92149 88407 96760 49114 11334 57536 87001 66641 14147",
	"2 5 51545 48566 64186 96046 3877 61515 5700 40440 92194 80585",
	"5 5 51590 84825 22329 22098 65830 29746 1613 26152 70729 71872 30432 53013 67342 45066 75733 46305 60180 35295 86405 71827 79816 95604 749 50291 97060",
	"5 2 67985 73579 26934 55849 7357 63059 47807 74711 72667 26194",
	"5 4 63561 46766 54320 45362 208 70580 70794 81723 80276 43403 60051 78625 3667 30095 83280 23228 72189 76607 23696 12007",
	"5 3 4255 88227 9235 10910 2188 59376 1909 98848 99037 36858 32711 35212 14351 81895 24198",
	"3 3 9112 21951 20923 33452 69125 22040 86070 35772 84962",
	"3 4 92095 42206 65077 62099 14968 3098 40896 50667 45003 55171 24647 33872",
	"1 3 95703 66862 27406",
	"5 4 2729 29541 2342 52077 19198 4631 94220 21002 58415 92355 66363 88890 55924 71396 28915 82677 91102 67712 59094 29255",
	"5 1 51761 88461 75478 42107 86485",
	"4 1 96660 39139 16474 27805",
	"1 3 9271 10020 40680",
	"3 2 54549 74048 33078 17091 1112 73495",
	"1 5 28521 74748 60405 22482 92278",
	"5 5 4906 49542 26268 45473 12980 26970 75155 88363 56748 77518 25444 64534 13688 87289 51127 38807 66075 65510 2255 42644 80233 52734 36878 2372 20574",
	"2 3 73839 17714 44446 56262 27923 34936",
	"1 4 71779 45070 90061 70036",
	"4 5 30755 8562 95089 5296 11100 17435 22243 21831 70545 27915 35129 99499 43547 78671 66308 33462 48249 44414 44602 14931",
	"3 2 79166 93731 64068 17741 76017 72244",
	"1 3 5130 53294 9594",
	"4 2 16387 44683 15033 80634 76993 49551 10047 74814",
	"5 2 74183 10715 34961 47828 38739 73984 70032 14984 60001 36331",
	"1 1 38763",
	"1 5 87873 1907 12018 54203 15087",
	"1 2 31410 76913",
	"4 2 15147 59102 21940 89246 31644 20834 97519 13479",
	"4 4 71163 38539 72118 33215 93273 62523 41217 13125 27213 85466 41605 5194 3574 1378 38739 95222",
	"5 3 58963 51285 41063 52240 8253 8414 41596 78833 59751 14597 32777 28206 80978 71161 90203",
	"4 3 33959 24016 70989 27242 40282 26112 32294 47247 10666 36804 11720 98735",
	"4 1 85461 75283 84341 44419",
	"2 4 40211 5381 42893 24486 41516 75892 39690 32224",
	"3 1 71333 80137 75889",
	"5 1 32126 28857 2671 31951 52662",
	"1 3 72248 9296 95574",
	"1 1 83281",
	"1 3 98400 47080 64653",
	"4 2 13230 65724 43004 10107 66752 87196 22708 23537",
	"2 2 41915 40059 14009 92973",
	"5 5 38469 16555 27098 18571 71499 94717 4163 41428 81728 88107 72477 97804 90387 26927 23352 39181 56707 70451 20696 6365 93694 87528 32414 33108 8443",
	"4 4 71994 32797 70960 57593 70525 59417 1425 51867 44391 22482 33813 63673 3200 84731 54616 74791",
	"1 1 90663",
}

func solveCase(n, m int, matrix [][]int) int {
	maxA := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if matrix[i][j] > maxA {
				maxA = matrix[i][j]
			}
		}
	}
	buf := 500
	limit := maxA + buf
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	nextPrime := make([]int, limit+2)
	next := -1
	for i := limit; i >= 0; i-- {
		if isPrime[i] {
			next = i
		}
		nextPrime[i] = next
	}
	rowSum := make([]int, n)
	colSum := make([]int, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			v := matrix[i][j]
			np := nextPrime[v]
			if np < 0 {
				np = v
			}
			delta := np - v
			rowSum[i] += delta
			colSum[j] += delta
		}
	}
	ans := rowSum[0]
	for i := 0; i < n; i++ {
		if rowSum[i] < ans {
			ans = rowSum[i]
		}
	}
	for j := 0; j < m; j++ {
		if colSum[j] < ans {
			ans = colSum[j]
		}
	}
	return ans
}

func parseCase(line string) (int, int, [][]int, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) < 2 {
		return 0, 0, nil, fmt.Errorf("invalid line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, nil, err
	}
	m, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, nil, err
	}
	if len(fields) != 2+n*m {
		return 0, 0, nil, fmt.Errorf("expected %d values, got %d", 2+n*m, len(fields))
	}
	matrix := make([][]int, n)
	idx := 2
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, m)
		for j := 0; j < m; j++ {
			v, _ := strconv.Atoi(fields[idx])
			matrix[i][j] = v
			idx++
		}
	}
	return n, m, matrix, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, line := range rawTestcases {
		n, m, matrix, err := parseCase(line)
		if err != nil {
			fmt.Printf("case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected := strconv.Itoa(solveCase(n, m, matrix))
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(matrix[i][j]))
			}
			sb.WriteByte('\n')
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}
