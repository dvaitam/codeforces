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

func runCandidate(bin, input string) (string, error) {
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

func expected(n, b, p int) string {
	bottles := (n - 1) * (2*b + 1)
	towels := n * p
	return fmt.Sprintf("%d %d", bottles, towels)
}

const testcasesARaw = `100
433 198 389
456 216 21
133 495 262
249 208 471
402 425 156
496 245 184
299 457 465
112 259 72
145 72 387
49 317 410
129 466 273
362 415 309
462 76 159
51 374 38
461 436 351
170 242 287
52 182 223
162 313 328
468 105 495
283 245 227
444 267 134
32 413 471
281 469 8
48 369 431
205 364 423
402 343 321
1 314 253
424 445 171
125 374 167
361 446 33
98 470 291
114 123 412
496 73 412
279 230 47
42 164 449
261 478 251
56 155 283
150 362 64
281 171 418
473 277 105
494 410 309
281 301 148
228 47 306
409 198 163
295 124 149
95 97 421
96 17 314
337 134 244
36 46 348
388 67 449
77 473 20
432 42 460
359 473 425
277 350 201
429 362 269
142 268 416
121 435 111
459 348 302
423 487 215
297 141 231
253 339 329
359 470 407
183 43 167
314 60 250
301 323 172
433 98 125
9 375 139
60 362 113
191 407 88
171 219 418
32 52 401
75 438 358
113 24 419
294 325 466
480 274 309
349 38 14
64 326 97
311 426 295
62 201 47
190 427 60
19 311 12
100 493 497
95 368 64
246 108 373
410 32 480
348 12 279
218 318 52
428 134 36
114 37 332
155 180 224
93 32 258
240 21 306
52 359 201
103 134 184
464 375 241
430 462 471
292 87 358
345 105 495
393 30 404
347 82 433`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data := []byte(testcasesARaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 1; i <= t; i++ {
		if !scan.Scan() {
			fmt.Printf("missing n for case %d\n", i)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		b, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		pVal, _ := strconv.Atoi(scan.Text())
		input := fmt.Sprintf("%d %d %d\n", n, b, pVal)
		expect := expected(n, b, pVal)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i, err, input)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", i, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
