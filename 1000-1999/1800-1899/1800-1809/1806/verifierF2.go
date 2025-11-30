package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesF2 = `2 665 1 45 94
3 405 2 307 90 376
3 669 1 31 42 404
5 840 3 176 390 27 687 327
5 536 2 527 365 252 34 426
5 530 1 205 347 469 189 202
3 39 1 8 19 21
4 754 2 148 185 591 738
5 287 1 196 57 255 84 195
3 422 2 66 287 324
5 909 3 48 659 632 538 234
3 528 2 498 230 81
2 158 1 78 8
4 120 1 80 82 68 80
3 562 2 263 14 376
3 48 1 46 29 43
4 802 2 793 62 770 256
4 201 2 171 27 178 69
5 385 3 30 136 357 56 184
4 203 2 60 62 52 49
5 129 3 44 89 101 7 12
5 706 1 261 206 524 308 88
2 514 1 362 184
3 819 1 295 750 818
4 182 2 76 64 138 149
3 82 2 25 46 55
4 957 2 951 911 106 811
3 477 2 260 67 33
5 415 4 164 313 160 248 73
4 50 2 45 5 30 38
3 72 1 51 30 38
2 437 1 278 16
2 415 1 269 189
3 911 1 167 670 668
2 974 1 549 620
3 306 2 249 262 127
5 801 2 613 614 192 350 186
2 344 1 53 263
4 624 1 105 576 107 489
5 298 3 47 107 22 207 133
4 229 3 88 67 186 201
2 824 1 819 131
2 689 1 450 288
4 264 2 10 176 165 125
4 752 2 368 736 385 207
2 405 1 206 309
2 560 1 322 25
5 151 1 33 13 143 75 51
4 252 3 22 196 58 171
5 330 3 11 34 236 141 324
5 223 2 186 116 3 41 222
2 172 1 24 95
2 685 1 232 509
4 323 3 192 275 161 255
2 251 1 190 55
3 103 1 14 6 77
2 606 1 465 461
5 19 3 17 7 5 10 18
4 692 2 610 535 52 470
3 950 2 271 571 679
4 759 1 322 588 722 262
4 467 2 53 301 418 160
4 708 1 447 350 7 112
4 856 3 108 571 202 585
2 961 1 19 318
5 117 1 64 95 31 36 110
4 807 2 206 763 412 618
5 934 4 315 339 16 777 880
4 113 2 89 12 20 50
5 273 3 143 114 94 233 239
4 491 1 43 476 352 89
2 796 1 369 341
4 62 2 43 22 12 19
5 50 3 21 18 49 40 41
4 80 2 3 58 80 23
4 924 3 302 712 271 503
2 458 1 125 178
3 119 1 49 57 102
5 667 1 333 27 13 471 257
3 716 1 477 106 713
5 580 1 271 108 141 60 11
4 437 3 339 397 424 298
4 825 1 342 494 599 347
4 346 1 25 267 122 227
2 696 1 327 44
3 735 1 575 250 31
4 157 3 62 147 127 139
5 448 2 229 135 89 20 185
4 873 3 350 642 281 235
5 10 1 2 3 10 6 10
3 22 1 18 6 19
3 83 1 20 71 24
5 403 1 190 148 177 399 289
2 244 1 95 199
4 354 1 295 347 351 190
4 70 2 25 47 33 20
3 859 1 615 139 843
2 328 1 316 85
4 265 2 135 229 95 203
5 374 1 163 77 133 92 173`

type int64Heap []int64

func (h int64Heap) Len() int           { return len(h) }
func (h int64Heap) Less(i, j int) bool { return h[i] < h[j] }
func (h int64Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *int64Heap) Push(x any)        { *h = append(*h, x.(int64)) }
func (h *int64Heap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type testCase struct {
	input    string
	expected string
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func solve(n int, k int, arr []int64) int64 {
	h := &int64Heap{}
	for _, v := range arr {
		*h = append(*h, v)
	}
	heap.Init(h)
	var sum int64
	for _, v := range arr {
		sum += v
	}
	for i := 0; i < k; i++ {
		x := heap.Pop(h).(int64)
		y := heap.Pop(h).(int64)
		g := gcd(x, y)
		sum -= x + y - g
		heap.Push(h, g)
	}
	return sum
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(testcasesF2, "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		p := 0
		n, err := strconv.Atoi(parts[p])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %w", idx+1, err)
		}
		p++
		// m exists but unused in solve
		p++ // skip m
		k, err := strconv.Atoi(parts[p])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad k: %w", idx+1, err)
		}
		p++
		if len(parts) != 3+n {
			return nil, fmt.Errorf("line %d: expected %d values got %d", idx+1, 3+n, len(parts))
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(parts[p+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: bad value: %w", idx+1, err)
			}
			arr[i] = v
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		// include m even if unused in solve to match input format
		sb.WriteString(parts[0])
		sb.WriteByte(' ')
		sb.WriteString(parts[1])
		sb.WriteByte(' ')
		sb.WriteString(parts[2])
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: strconv.FormatInt(solve(n, k, arr), 10),
		})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: verifierF2 /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
