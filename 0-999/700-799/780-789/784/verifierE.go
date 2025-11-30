package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesE.txt.
const embeddedTestcasesE = `241 310 105 738 844
405 490 158 92 189
68 20 411 562 1017
939 296 819 783 703
60 227 532 549 238
368 283 798 176 965
846 108 268 219 757
965 949 26 848 826
656 826 266 819 915
278 198 168 317 69
296 642 888 749 575
983 875 868 901 93
381 88 865 620 40
345 687 397 518 125
254 182 253 484 337
286 91 967 957 319
837 886 964 560 455
860 307 7 930 458
298 586 721 903 566
319 868 783 520 860
199 423 433 613 692
295 441 462 165 501
238 312 265 832 927
816 44 83 47 864
473 641 287 531 84
547 663 482 717 923
351 148 689 200 946
68 422 935 207 650
650 647 451 282 212
188 364 446 764 658
602 328 649 571 928
203 926 331 103 633
859 63 725 234 347
284 783 596 630 561
882 243 125 339 687
957 181 297 470 1015
26 43 365 714 918
84 917 987 977 971
292 752 691 982 689
334 18 330 295 305
329 989 156 793 273
667 420 882 965 904
890 635 696 837 252
79 300 632 196 991
915 454 299 139 1013
256 390 613 984 315
162 339 586 9 946
372 45 465 173 37
373 802 823 371 19
297 585 99 449 706
212 434 938 212 536
116 60 63 56 79
754 172 609 693 650
993 153 621 41 316
559 502 596 255 370
329 36 125 854 582
541 299 792 419 397
667 205 489 206 881
247 449 420 503 357
37 224 431 454 172
254 662 910 438 80
850 220 510 192 688
32 37 260 259 2
248 538 213 790 289
237 427 890 267 823
145 332 52 912 633
950 322 578 119 193
583 412 985 922 920
668 669 891 809 83
734 762 41 506 503
396 95 440 215 188
952 894 586 921 277
966 946 169 344 389
303 672 482 818 351
922 657 322 834 779
430 540 220 670 496
820 702 820 274 940
346 944 401 979 168
508 76 879 994 317
286 984 642 684 744
195 45 404 920 738
634 130 785 991 566
275 683 876 61 233
890 171 704 650 923
475 582 483 764 130
413 949 399 223 888
817 3 216 944 90
160 13 624 897 348
263 118 405 829 985
788 801 390 898 561
227 563 54 899 357
973 206 165 687 265
622 338 838 913 1003
575 792 803 482 998
937 539 450 27 107
80 35 712 608 219
115 500 574 885 204
263 622 792 142 255
42 371 81 786 538
900 535 926 10 519`

func solve784E(a, b, c, d int) int {
	return a ^ b ^ c ^ d
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesE), "\n")
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 5 {
			fmt.Fprintf(os.Stderr, "case %d malformed\n", idx+1)
			os.Exit(1)
		}
		a, err1 := strconv.Atoi(fields[0])
		b, err2 := strconv.Atoi(fields[1])
		c, err3 := strconv.Atoi(fields[2])
		d, err4 := strconv.Atoi(fields[3])
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error\n", idx+1)
			os.Exit(1)
		}
		want := strconv.Itoa(solve784E(a, b, c, d))
		input := fmt.Sprintf("%d %d %d %d\n", a, b, c, d)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
