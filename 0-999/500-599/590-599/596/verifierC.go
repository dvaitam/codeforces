package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesC.txt.
const embeddedTestcasesC = `1 0 0 0
3 0 2 0 1 0 0 0 1 2
6 2 1 1 0 1 1 0 1 0 0 2 0 0 1 -1 0 -2 -1
12 0 1 2 3 1 3 2 2 1 2 0 2 0 3 1 1 2 1 0 0 1 0 2 0 0 1 -1 2 0 -2 3 1 -1 2 0 1
4 1 1 1 0 0 0 0 1 0 1 -1 0
6 2 1 0 0 1 1 1 0 2 0 0 1 0 1 -1 0 -2 -1
16 0 0 3 1 2 0 0 3 0 1 1 0 3 3 2 2 1 2 0 2 1 3 3 2 1 1 2 1 3 0 2 3 0 1 -1 2 0 -2 3 1 -1 -3 2 0 -2 1 -1 0
12 1 0 0 0 1 2 0 2 0 1 2 2 0 3 2 3 2 0 2 1 1 1 1 3 0 1 -1 2 0 -2 3 1 -1 2 0 1
12 0 0 2 2 1 3 0 3 0 1 1 1 0 2 1 2 2 1 2 0 2 3 1 0 0 1 -1 2 0 -2 3 1 -1 2 0 1
3 0 1 0 0 0 2 0 1 2
2 1 0 0 0 0 -1
3 0 1 0 0 0 2 0 1 2
2 0 0 0 1 0 1
4 1 0 1 1 0 1 0 0 0 1 -1 0
3 0 2 0 0 0 1 0 1 2
2 1 0 0 0 0 -1
1 0 0 0
1 0 0 0
3 0 2 0 0 0 1 0 1 2
4 1 0 1 1 0 1 0 0 0 1 -1 0
4 1 0 0 1 1 1 0 0 0 1 -1 0
3 0 0 0 2 0 1 0 1 2
3 0 2 0 0 0 1 0 1 2
12 1 0 2 3 2 0 0 0 2 2 1 2 1 1 0 1 0 3 1 3 0 2 2 1 0 1 -1 2 0 -2 3 1 -1 2 0 1
8 0 0 0 3 1 2 1 1 0 1 0 2 1 0 1 3 0 1 -1 2 0 3 1 2
4 3 0 0 0 2 0 1 0 0 -1 -2 -3
3 1 0 2 0 0 0 0 -1 -2
2 0 0 1 0 0 -1
2 1 0 0 0 0 -1
8 1 1 0 3 1 2 1 3 1 0 0 0 0 2 0 1 0 1 -1 2 0 3 1 2
12 2 3 1 1 2 2 1 3 0 1 1 2 0 3 2 1 0 2 0 0 2 0 1 0 0 1 -1 2 0 -2 3 1 -1 2 0 1
2 1 0 0 0 0 -1
2 0 1 0 0 0 1
2 0 1 0 0 0 1
16 0 2 3 3 0 1 1 1 1 3 0 0 1 0 3 2 3 0 2 3 3 1 0 3 1 2 2 2 2 0 2 1 0 1 -1 2 0 -2 3 1 -1 -3 2 0 -2 1 -1 0
4 0 1 0 3 0 0 0 2 0 1 2 3
9 0 2 0 1 2 1 2 0 1 0 2 2 0 0 1 1 1 2 0 1 -1 2 0 -2 1 -1 0
2 0 0 1 0 0 -1
4 0 2 0 0 0 1 0 3 0 1 2 3
1 0 0 0
3 0 0 2 0 1 0 0 -1 -2
2 0 0 0 1 0 1
2 0 0 1 0 0 -1
16 3 2 0 0 3 1 1 3 0 3 0 2 2 2 2 0 2 1 2 3 3 3 0 1 1 2 1 1 3 0 1 0 0 1 -1 2 0 -2 3 1 -1 -3 2 0 -2 1 -1 0
4 0 3 0 1 0 2 0 0 0 1 2 3
9 0 0 0 1 1 1 1 0 1 2 2 2 2 0 0 2 2 1 0 1 -1 2 0 -2 1 -1 0
16 3 0 0 1 3 1 2 2 0 0 2 3 1 1 3 2 2 0 1 0 3 3 1 3 0 2 0 3 1 2 2 1 0 1 -1 2 0 -2 3 1 -1 -3 2 0 -2 1 -1 0
4 0 1 1 0 0 0 1 1 0 1 -1 0
2 0 0 0 1 0 1
8 3 1 3 0 0 0 1 0 2 0 0 1 1 1 2 1 0 1 -1 0 -2 -1 -3 -2
4 1 0 0 0 3 0 2 0 0 -1 -2 -3
6 0 0 1 0 2 0 2 1 1 1 0 1 0 1 -1 0 -2 -1
8 0 1 1 2 0 2 0 3 1 0 0 0 1 1 1 3 0 1 -1 2 0 3 1 2
4 0 0 2 0 1 0 3 0 0 -1 -2 -3
12 0 3 1 0 0 2 0 1 0 0 1 2 2 1 2 0 1 1 1 3 2 2 2 3 0 1 -1 2 0 -2 3 1 -1 2 0 1
16 3 2 2 1 2 2 1 2 1 0 3 0 0 3 3 3 0 1 0 2 3 1 2 3 2 0 1 3 1 1 0 0 0 1 -1 2 0 -2 3 1 -1 -3 2 0 -2 1 -1 0
8 1 0 0 0 1 1 0 1 3 0 3 1 2 1 2 0 0 1 -1 0 -2 -1 -3 -2
4 0 2 0 0 0 1 0 3 0 1 2 3
8 2 0 3 1 3 0 2 1 1 1 0 0 1 0 0 1 0 1 -1 0 -2 -1 -3 -2
12 1 1 2 2 1 2 0 3 2 1 0 2 2 3 0 1 0 0 1 3 2 0 1 0 0 1 -1 2 0 -2 3 1 -1 2 0 1
2 0 0 0 1 0 1
16 3 2 0 1 2 2 2 3 1 2 0 3 3 0 3 3 1 1 3 1 1 3 2 0 1 0 0 0 0 2 2 1 0 1 -1 2 0 -2 3 1 -1 -3 2 0 -2 1 -1 0
8 1 1 3 0 0 1 0 0 1 0 2 0 2 1 3 1 0 1 -1 0 -2 -1 -3 -2
9 0 2 1 2 1 0 0 1 2 1 0 0 2 2 2 0 1 1 0 1 -1 2 0 -2 1 -1 0
3 0 1 0 2 0 0 0 1 2
8 1 0 3 1 0 1 1 1 2 1 3 0 0 0 2 0 0 1 -1 0 -2 -1 -3 -2
9 0 0 2 1 1 1 2 0 2 2 0 2 1 0 1 2 0 1 0 1 -1 2 0 -2 1 -1 0
8 3 0 1 1 0 0 2 0 1 0 0 1 2 1 3 1 0 1 -1 0 -2 -1 -3 -2
8 1 0 0 0 1 3 0 1 0 2 1 1 1 2 0 3 0 1 -1 2 0 3 1 2
4 0 0 0 1 1 1 1 0 0 1 -1 0
12 1 0 3 2 0 1 0 0 0 2 3 1 2 1 2 0 1 2 2 2 3 0 1 1 0 1 -1 2 0 -2 1 -1 -3 0 -2 -1
16 0 3 2 3 3 2 2 1 2 0 0 0 2 2 1 2 3 1 0 2 0 1 3 0 1 3 3 3 1 0 1 1 0 1 -1 2 0 -2 3 1 -1 -3 2 0 -2 1 -1 0
4 0 3 0 1 0 2 0 0 0 1 2 3
8 1 2 0 0 1 3 0 2 0 1 1 1 0 3 1 0 0 1 -1 2 0 3 1 2
4 2 0 3 0 0 0 1 0 0 -1 -2 -3
6 2 1 0 1 0 0 2 0 1 0 1 1 0 1 -1 0 -2 -1
4 1 1 0 0 1 0 0 1 0 1 -1 0
9 2 1 1 2 2 0 1 1 1 0 2 2 0 1 0 2 0 0 0 1 -1 2 0 -2 1 -1 0
3 0 0 2 0 1 0 0 -1 -2
2 0 0 1 0 0 -1
16 0 0 0 2 1 0 3 3 1 2 1 1 0 3 1 3 3 2 2 3 0 1 3 0 2 1 3 1 2 0 2 2 0 1 -1 2 0 -2 3 1 -1 -3 2 0 -2 1 -1 0
4 2 0 3 0 1 0 0 0 0 -1 -2 -3
2 0 1 0 0 0 1
12 0 2 3 0 1 0 3 2 3 1 2 0 2 2 1 2 0 1 1 1 0 0 2 1 0 1 -1 2 0 -2 1 -1 -3 0 -2 -1
12 2 1 1 1 0 2 1 2 2 3 1 0 2 0 2 2 0 1 0 3 1 3 0 0 0 1 -1 2 0 -2 3 1 -1 2 0 1
2 0 1 0 0 0 1
6 2 1 2 0 0 0 1 1 1 0 0 1 0 1 -1 0 -2 -1
2 0 1 0 0 0 1
8 3 1 2 1 3 0 1 0 1 1 0 1 2 0 0 0 0 1 -1 0 -2 -1 -3 -2
2 1 0 0 0 0 -1
4 1 0 0 1 0 0 1 1 0 1 -1 0
12 0 3 1 3 0 1 1 0 0 0 2 1 2 0 1 1 2 2 2 3 1 2 0 2 0 1 -1 2 0 -2 3 1 -1 2 0 1
12 2 2 2 0 2 1 1 3 1 1 2 3 0 0 0 1 0 2 1 2 0 3 1 0 0 1 -1 2 0 -2 3 1 -1 2 0 1
4 3 0 0 0 2 0 1 0 0 -1 -2 -3
12 1 0 2 3 2 1 1 2 0 1 1 3 2 2 1 1 0 2 0 3 0 0 2 0 0 1 -1 2 0 -2 3 1 -1 2 0 1
4 0 1 0 3 0 0 0 2 0 1 2 3
9 0 2 1 1 2 1 1 2 0 0 2 0 1 0 0 1 2 2 0 1 -1 2 0 -2 1 -1 0
3 2 0 1 0 0 0 0 -1 -2
9 1 2 0 2 0 0 1 0 0 1 2 1 1 1 2 0 2 2 0 1 -1 2 0 -2 1 -1 0
16 1 1 0 0 3 2 0 2 3 3 2 2 0 1 2 0 3 1 0 3 3 0 1 0 2 1 1 3 1 2 2 3 0 1 -1 2 0 -2 3 1 -1 -3 2 0 -2 1 -1 0`

type point struct {
	x, y   int
	w      int
	indeg  int
	n1, n2 int
}

func solve596C(n int, coords [][2]int, wseq []int) (string, bool) {
	points := make([]point, n)
	pos := make(map[[2]int]int, n)
	for i := 0; i < n; i++ {
		points[i] = point{
			x: coords[i][0],
			y: coords[i][1],
			w: coords[i][1] - coords[i][0],
			n1: -1, n2: -1,
		}
		pos[[2]int{coords[i][0], coords[i][1]}] = i
	}

	for i := 0; i < n; i++ {
		x, y := points[i].x, points[i].y
		if j, ok := pos[[2]int{x + 1, y}]; ok {
			points[i].n1 = j
			points[j].indeg++
		}
		if j, ok := pos[[2]int{x, y + 1}]; ok {
			points[i].n2 = j
			points[j].indeg++
		}
	}

	avail := make(map[int][]int)
	for i := 0; i < n; i++ {
		if points[i].indeg == 0 {
			w := points[i].w
			avail[w] = append(avail[w], i)
		}
	}

	ans := make([][2]int, n)
	for i := 0; i < n; i++ {
		w := wseq[i]
		lst := avail[w]
		if len(lst) == 0 {
			return "NO", false
		}
		idx := lst[len(lst)-1]
		avail[w] = lst[:len(lst)-1]
		ans[i] = [2]int{points[idx].x, points[idx].y}

		if points[idx].n1 != -1 {
			j := points[idx].n1
			points[j].indeg--
			if points[j].indeg == 0 {
				avail[points[j].w] = append(avail[points[j].w], j)
			}
		}
		if points[idx].n2 != -1 {
			j := points[idx].n2
			points[j].indeg--
			if points[j].indeg == 0 {
				avail[points[j].w] = append(avail[points[j].w], j)
			}
		}
	}

	var sb strings.Builder
	sb.WriteString("YES\n")
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", ans[i][0], ans[i][1]))
	}
	return strings.TrimRight(sb.String(), "\n"), true
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

func parseCase(line string) (n int, coords [][2]int, wseq []int, ok bool) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return
	}
	expected := 1 + 3*n
	if len(fields) != expected {
		return
	}
	coords = make([][2]int, n)
	idx := 1
	for i := 0; i < n; i++ {
		x, err1 := strconv.Atoi(fields[idx])
		y, err2 := strconv.Atoi(fields[idx+1])
		if err1 != nil || err2 != nil {
			return
		}
		coords[i] = [2]int{x, y}
		idx += 2
	}
	wseq = make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[idx+i])
		if err != nil {
			return
		}
		wseq[i] = v
	}
	ok = true
	return
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesC), "\n")
	for idx, line := range lines {
		n, coords, wseq, ok := parseCase(line)
		if !ok {
			fmt.Fprintf(os.Stderr, "case %d malformed\n", idx+1)
			os.Exit(1)
		}
		wantStr, success := solve596C(n, coords, wseq)
		if !success {
			wantStr = "NO"
		}

		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		for _, p := range coords {
			fmt.Fprintf(&input, "%d %d\n", p[0], p[1])
		}
		for i, w := range wseq {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(w))
		}
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(wantStr) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\n", idx+1, wantStr, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
