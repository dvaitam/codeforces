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

func expected(n, a int) int {
	if a%2 == 1 {
		return (a + 1) / 2
	}
	return (n-a)/2 + 1
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesARaw = `138 17
262 61
508 390
462 242
668 389
808 215
98 63
30 29
856 400
444 312
782 3
714 457
274 118
606 105
924 326
32 2
28 21
556 10
962 903
392 352
222 109
744 30
542 228
784 449
962 508
568 239
354 119
694 225
780 471
976 297
950 23
428 285
946 658
104 24
646 304
124 96
342 257
960 433
520 195
312 146
602 512
868 518
404 302
874 36
492 125
762 414
426 341
178 94
562 384
90 57
680 521
112 100
168 134
862 403
380 251
752 31
482 23
316 315
608 593
404 332
176 44
516 233
14 13
206 139
944 881
562 238
416 264
354 296
362 236
932 276
676 562
624 6
394 380
526 133
532 211
438 29
494 446
374 292
568 205
964 517
424 249
834 366
426 178
2 2
470 308
30 26
236 163
182 141
600 186
882 94
818 565
818 262
34 5
86 3
464 8
774 774
288 128
276 57
818 640
190 89`

	scanner := bufio.NewScanner(strings.NewReader(testcasesARaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		n, _ := strconv.Atoi(parts[0])
		a, _ := strconv.Atoi(parts[1])
		exp := expected(n, a)
		input := fmt.Sprintf("%d %d\n", n, a)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		got, err2 := strconv.Atoi(gotStr)
		if err2 != nil || got != exp {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
