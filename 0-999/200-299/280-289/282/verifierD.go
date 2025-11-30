package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

var rawTestcases = []string{
	"2 215 20",
	"2 261 248",
	"2 155 244",
	"2 298 111",
	"3 71 144 71",
	"1 128",
	"3 75 158 50",
	"3 37 169 241",
	"3 51 181 222",
	"2 104 282",
	"2 226 266",
	"2 31 280",
	"1 47",
	"3 204 0 252",
	"2 124 166",
	"3 32 97 290",
	"1 122",
	"1 278",
	"2 46 41",
	"2 260 250",
	"1 154",
	"3 149 63 280",
	"2 276 104",
	"3 280 147 227",
	"1 197",
	"2 294 123",
	"2 94 96",
	"1 16",
	"3 133 243 35",
	"1 66",
	"1 19",
	"1 276",
	"3 200 268 141",
	"3 120 110 214",
	"3 140 230 252",
	"3 182 42 166",
	"3 59 249 171",
	"1 124",
	"1 138",
	"1 112",
	"2 87 170",
	"2 31 51",
	"1 112",
	"1 293",
	"3 273 37 13",
	"1 96",
	"3 294 61 200",
	"1 189",
	"1 18",
	"3 11 99 94",
	"3 63 245 107",
	"3 31 11 278",
	"2 51 133",
	"1 113",
	"1 154",
	"2 223 92",
	"1 257",
	"2 20 51",
	"3 200 102 133",
	"2 240 291",
	"1 104",
	"1 81",
	"1 175",
	"3 128 60 226",
	"3 89 6 241",
	"3 209 291 260",
	"2 182 198",
	"3 128 78 287",
	"3 6 234 40",
	"2 23 278",
	"2 69 122",
	"2 180 147",
	"3 183 67 158",
	"2 212 41",
	"1 98",
	"3 171 81 122",
	"1 229",
	"2 290 212",
	"1 205",
	"3 290 214 23",
	"1 228",
	"1 132",
	"3 80 228 270",
	"2 287 0",
	"1 253",
	"2 159 239",
	"1 212",
	"1 280",
	"3 42 66 7",
	"2 213 161",
	"1 109",
	"1 1",
	"3 270 50 97",
	"1 101",
	"2 143 93",
	"1 243",
	"2 41 11",
	"2 231 59",
	"2 68 266",
	"3 177 58 79",
}

func solveCase(nums []int) string {
	n := len(nums)
	switch n {
	case 1:
		if nums[0] > 0 {
			return "BitLGM"
		}
		return "BitAryo"
	case 2:
		sort.Ints(nums)
		phi := (1 + math.Sqrt(5)) / 2
		k := nums[1] - nums[0]
		t := int(math.Floor(float64(k) * phi))
		if t == nums[0] {
			return "BitAryo"
		}
		return "BitLGM"
	case 3:
		if nums[0]^nums[1]^nums[2] != 0 {
			return "BitLGM"
		}
		return "BitAryo"
	default:
		return "BitLGM"
	}
}

func parseCase(line string) ([]int, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}
	if n < 1 || n > 3 || len(fields) != 1+n {
		return nil, fmt.Errorf("invalid case")
	}
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		val, _ := strconv.Atoi(fields[1+i])
		nums[i] = val
	}
	return nums, nil
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, line := range rawTestcases {
		nums, err := parseCase(line)
		if err != nil {
			fmt.Printf("case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected := solveCase(nums)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(nums)))
		for i, v := range nums {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
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
