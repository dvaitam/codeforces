package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func compileRef() (string, error) {
	exe := filepath.Join(os.TempDir(), "ref628D")
	cmd := exec.Command("go", "build", "-o", exe, "628D.go")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return exe, nil
}

func run(bin, input string) (string, error) {
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
	ref, err := compileRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	const testcasesDRaw = `100
1 1 3 8
13 7 289185 721672
12 6 743 782
6 9 94 95
2 0 08786 99634
15 3 2203 7116
17 4 9 9
12 5 849 925
17 4 02 63
7 9 5560 9439
10 7 60 78
11 5 632386 831985
13 9 8216 8826
19 3 55 68
9 8 65 73
20 8 66363 97548
11 4 4942 8679
17 0 203763 473238
15 1 3 7
4 4 3 4
20 2 377885 658561
1 0 9 9
10 3 8188 9700
12 5 9 9
4 0 01 77
18 4 756 809
13 3 188 371
16 2 77517 94085
14 9 986265 989979
1 0 9903 9907
5 8 6 8
11 9 91 91
11 2 672 812
16 7 43330 60996
3 7 13 78
9 8 424212 893755
14 2 54939 63222
4 9 2 9
10 1 801474 913086
9 4 5 9
18 1 0 9
1 2 6 9
16 0 647 718
3 6 607341 725365
2 3 0452 1224
20 6 71 94
10 6 8951 9063
15 8 71 99
5 3 22066 33768
19 7 274 532
7 6 009430 916357
2 8 0 4
5 2 17640 32649
11 0 860394 945150
5 5 98443 98571
3 5 32 89
11 7 8776 9854
3 1 449 841
13 8 9032 9468
17 1 18 69
17 5 09 27
5 0 9260 9378
3 8 578 732
12 7 9690 9853
20 7 45248 60471
17 0 402523 786286
11 8 3 7
3 7 5 8
20 3 03183 71167
17 1 9 9
2 5 43 60
16 7 0021 3890
16 5 841795 883925
2 7 86310 89011
13 4 921380 997436
16 3 7 9
20 2 57891 78551
19 5 9304 9947
5 8 6413 7384
12 6 246 462
9 3 76959 89150
15 3 44266 85159
5 5 18 50
18 2 92360 98180
20 1 7861 8699
6 9 619 746
3 8 3445 5851
13 4 26481 72873
14 5 021744 067175
7 3 21437 94472
15 1 513 660
16 6 868 966
2 1 5 6
9 9 585924 815626
10 8 7892 9427
8 1 084072 912843
3 8 9 9
20 1 30792 33439
2 1 327249 423171
19 6 134629 944305`

	scan := bufio.NewScanner(strings.NewReader(testcasesDRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 1; i <= t; i++ {
		if !scan.Scan() {
			fmt.Printf("missing m for case %d\n", i)
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		d, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		a := scan.Text()
		scan.Scan()
		b := scan.Text()
		input := fmt.Sprintf("%d %d\n%s\n%s\n", m, d, a, b)
		expect, err := run(ref, input)
		if err != nil {
			fmt.Printf("reference failed on case %d: %v\n", i, err)
			os.Exit(1)
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", i, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
