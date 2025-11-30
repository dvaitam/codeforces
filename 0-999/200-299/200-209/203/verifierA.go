package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input    string
	expected string
}

var rawTestcases = []string{
	"523 28 194 222 7 7",
	"597 16 227 166 7 5",
	"143 29 147 270 4 9",
	"545 25 252 256 2 10",
	"75 23 259 91 10 3",
	"444 29 193 269 6 8",
	"453 11 241 284 10 4",
	"14 28 246 275 9 5",
	"249 3 266 172 7 1",
	"581 24 186 280 6 2",
	"82 8 257 67 4 3",
	"564 11 117 157 9 8",
	"553 10 188 289 2 9",
	"455 7 210 127 10 9",
	"297 3 182 137 10 7",
	"487 6 36 148 4 3",
	"82 3 80 25 2 3",
	"282 29 297 263 9 7",
	"429 17 293 237 4 4",
	"84 19 216 235 5 8",
	"194 11 224 191 10 2",
	"380 8 66 147 1 5",
	"103 26 184 165 3 6",
	"547 26 86 246 3 4",
	"589 20 101 115 2 1",
	"37 4 210 65 7 2",
	"126 20 113 259 1 4",
	"557 16 299 26 4 1",
	"226 14 196 61 10 2",
	"62 3 233 104 5 6",
	"400 17 280 67 8 1",
	"208 7 270 122 5 6",
	"542 25 106 159 1 3",
	"179 9 242 250 2 10",
	"318 1 291 260 8 7",
	"157 21 288 204 6 7",
	"46 18 57 221 1 8",
	"493 18 146 246 5 3",
	"135 12 282 146 10 5",
	"1 23 216 174 5 7",
	"228 20 116 175 4 6",
	"411 21 266 148 8 7",
	"47 28 294 273 10 7",
	"457 6 172 90 8 2",
	"39 17 282 144 8 9",
	"425 16 209 87 6 5",
	"15 7 268 78 9 2",
	"14 13 84 126 7 6",
	"121 23 72 222 1 9",
	"186 20 147 271 4 5",
	"281 4 65 32 8 7",
	"533 30 264 92 8 2",
	"19 27 195 123 6 2",
	"322 2 133 289 1 4",
	"469 12 289 264 10 1",
	"212 21 277 165 7 6",
	"154 13 122 95 10 5",
	"95 9 250 142 6 6",
	"167 11 110 148 10 1",
	"561 5 224 222 10 5",
	"49 5 264 130 5 2",
	"412 10 63 158 3 9",
	"101 27 262 157 6 5",
	"127 18 222 223 8 8",
	"309 16 248 129 2 8",
	"67 11 222 74 3 3",
	"62 26 241 156 2 4",
	"296 13 213 289 1 2",
	"433 15 294 195 8 10",
	"441 3 145 93 6 4",
	"538 7 68 26 6 2",
	"262 15 183 231 4 2",
	"202 7 80 77 1 4",
	"154 15 237 295 7 6",
	"433 4 105 231 10 8",
	"510 29 295 265 9 8",
	"9 16 216 191 4 9",
	"151 11 69 194 6 6",
	"388 28 174 296 5 10",
	"528 19 106 165 5 8",
	"41 28 142 87 1 2",
	"164 10 178 292 1 8",
	"391 26 170 264 3 8",
	"92 29 254 288 9 9",
	"211 22 198 264 9 10",
	"397 10 187 213 9 10",
	"20 20 298 297 10 4",
	"583 22 175 192 1 3",
	"310 9 174 150 6 2",
	"167 25 266 183 7 7",
	"342 21 133 266 3 4",
	"503 2 214 80 1 8",
	"361 29 284 263 10 2",
	"469 14 251 228 1 10",
	"33 2 79 18 2 8",
	"107 20 270 139 10 3",
	"502 23 222 230 9 6",
	"127 4 242 202 1 10",
	"300 22 226 204 10 5",
	"39 30 278 273 2 9",
}

func solveCase(input string) (string, error) {
	var x, t, a, b, da, db int
	if _, err := fmt.Fscan(strings.NewReader(input), &x, &t, &a, &b, &da, &db); err != nil {
		return "", err
	}

	set1 := make(map[int]struct{})
	set2 := make(map[int]struct{})
	for i := 0; i < t; i++ {
		s1 := a - i*da
		s2 := b - i*db
		set1[s1] = struct{}{}
		set2[s2] = struct{}{}
	}

	possible := x == 0
	if !possible {
		if _, ok := set1[x]; ok {
			possible = true
		} else if _, ok := set2[x]; ok {
			possible = true
		}
	}
	if !possible {
		for s1 := range set1 {
			if _, ok := set2[x-s1]; ok {
				possible = true
				break
			}
		}
	}

	if possible {
		return "YES", nil
	}
	return "NO", nil
}

func buildTestCases() ([]testCase, error) {
	testcases := make([]testCase, 0, len(rawTestcases))
	for _, in := range rawTestcases {
		expected, err := solveCase(in)
		if err != nil {
			return nil, fmt.Errorf("compute expected for %q: %w", in, err)
		}
		testcases = append(testcases, testCase{
			input:    in,
			expected: expected,
		})
	}
	return testcases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := buildTestCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := tc.input + "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != tc.expected {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
