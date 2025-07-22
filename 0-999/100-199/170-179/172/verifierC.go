package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Student struct {
	t   int64
	x   int64
	idx int
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n, m int, students []Student) []int64 {
	ans := make([]int64, n)
	queue := []Student{}
	nextI := 0
	var curTime int64
	processed := 0
	for processed < n {
		if len(queue) == 0 && nextI < n && curTime < students[nextI].t {
			curTime = students[nextI].t
		}
		for nextI < n && students[nextI].t <= curTime && len(queue) < m {
			queue = append(queue, students[nextI])
			nextI++
		}
		for len(queue) < m && nextI < n {
			curTime = students[nextI].t
			queue = append(queue, students[nextI])
			nextI++
		}
		departTime := curTime
		k := len(queue)
		batch := make([]Student, k)
		copy(batch, queue)
		sort.Slice(batch, func(i, j int) bool { return batch[i].x < batch[j].x })
		var unloadTime int64
		lastX := int64(-1)
		var cntAtX int64
		for _, st := range batch {
			if st.x != lastX {
				if lastX != -1 {
					unloadTime += 1 + cntAtX/2
				}
				lastX = st.x
				cntAtX = 1
			} else {
				cntAtX++
			}
		}
		if lastX != -1 {
			unloadTime += 1 + cntAtX/2
		}
		maxX := batch[k-1].x
		for _, st := range batch {
			ans[st.idx] = departTime + st.x
		}
		curTime = departTime + maxX + unloadTime + maxX
		processed += k
		queue = queue[:0]
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+2*n {
			fmt.Fprintf(os.Stderr, "test %d bad count\n", idx)
			os.Exit(1)
		}
		students := make([]Student, n)
		for i := 0; i < n; i++ {
			ti, _ := strconv.Atoi(parts[2+2*i])
			xi, _ := strconv.Atoi(parts[3+2*i])
			students[i] = Student{t: int64(ti), x: int64(xi), idx: i}
		}
		want := expected(n, m, students)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			input.WriteString(fmt.Sprintf("%d %d\n", students[i].t, students[i].x))
		}
		gotStr, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		outFields := strings.Fields(gotStr)
		if len(outFields) != n {
			fmt.Fprintf(os.Stderr, "case %d wrong output length\n", idx)
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			val, _ := strconv.ParseInt(outFields[i], 10, 64)
			if val != want[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at student %d: expected %d got %d\n", idx, i+1, want[i], val)
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
