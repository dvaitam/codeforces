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

const testcaseData = `2 665 1 45 94
3 405 2 307 90 376
3 669 1 31 42 404
5 840 3 176 390 27 687 327
5 536 2 527 365 252 34 426
3 656 2 67 300 499
4 22 3 48 259 256 437
2 564 1 453 495
5 503 2 173 240 91 462 417
2 208 1 326 31
3 994 2 875 312 987
5 655 1 422 583 149 472 783
4 564 1 492 275 854 777
2 135 2 220 113
2 894 2 621 681
5 34 1 20 53 553 246 865
4 123 2 646 580 884 557
2 391 2 341 13
4 30 2 620 11 211 109
2 904 2 351 1000
3 2 1 264 695 229
4 850 2 45 898 717 811
5 653 3 45 561 544 123 103
5 37 1 971 103 620 18 948
5 590 1 798 603 373 241 872
4 527 2 197 801 125 677
2 611 1 1000 502
2 118 3 62 27
5 682 2 684 576 639 387 843
5 136 2 111 878 251 937 64
2 17 2 862 42
5 647 1 932 753 287 808 315
4 196 2 585 638 204 452
2 331 3 105 541
2 570 2 540 757
2 252 2 788 474
4 629 2 387 472 794 322
4 557 2 884 232 119 701
2 985 1 958 112
2 776 1 388 323
3 35 2 896 242 116
3 809 2 788 634 685
2 496 2 752 203
5 913 2 422 142 304 78 784
5 604 2 876 841 485 292 808
2 282 1 491 60
5 970 2 575 304 213 315 650
2 650 1 728 711
5 336 2 247 484 687 501 170
3 785 3 777 348 749
5 733 2 297 780 34 625 797
5 94 3 854 978 828 573 15
4 991 2 123 154 955 702
3 418 2 37 111 469
4 182 3 178 358 615 364
2 610 2 580 140
5 854 1 336 7 533 572 310
4 10 2 134 308 332 675
2 621 1 552 92
4 102 1 327 489 394 148
3 48 2 19 366 304
2 127 1 408 867
4 992 3 706 990 895 618
2 815 1 841 987
5 726 3 812 651 770 618 570
5 357 2 209 821 52 977 753
5 284 1 503 846 938 160 579
2 56 1 802 233
4 363 1 335 273 130 468
2 786 1 375 180
5 375 2 813 845 49 215 665
5 928 1 992 852 30 122 572
5 945 1 695 86 789 562 823
4 878 1 934 414 521 184
4 235 1 903 928 792 282
5 337 2 393 946 79 63 492
4 636 1 808 423 569 708
2 69 1 661 376
3 468 1 763 595 154
5 703 2 558 884 194 789 369
2 826 1 450 592
3 121 1 467 842 73
4 744 1 173 929 766 203
5 347 2 829 159 185 471 809
2 822 1 963 148
4 372 1 330 978 179 278
5 77 1 749 617 472 598 901
3 402 1 877 219 805
4 256 2 135 113 201 550
3 934 1 863 543 204
3 198 2 113 526 92
5 53 2 199 955 185 396 128
3 698 2 918 942 54
2 383 1 609 699
5 124 2 967 452 51 286 4
2 125 1 57 461
5 345 1 841 174 826 570 543
5 40 1 350 659 833 365 822
4 552 2 525 829 348 229
4 202 1 172 70 571 322
4 365 1 174 496 985 394
2 837 2 237 62`

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
		if h.Len() < 2 {
			break
		}
		x := heap.Pop(h).(int64)
		y := heap.Pop(h).(int64)
		g := gcd(x, y)
		sum -= x + y - g
		heap.Push(h, g)
	}
	return sum
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(testcaseData, "\n")
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
