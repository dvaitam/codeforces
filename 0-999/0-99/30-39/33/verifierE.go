package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func parseHM(s string) int {
	parts := strings.Split(s, ":")
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	return h*60 + m
}

func solveE(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var m, n, k int
	fmt.Fscan(reader, &m, &n, &k)
	subj := make([]string, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &subj[i])
	}
	tsolve := make(map[string]int)
	times := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &times[i])
		tsolve[subj[i]] = times[i]
	}
	type seg struct{ start, end, idx int }
	breaks := make([]seg, 4)
	for i := 0; i < 4; i++ {
		var s string
		fmt.Fscan(reader, &s)
		parts := strings.Split(s, "-")
		st := parseHM(parts[0])
		en := parseHM(parts[1])
		breaks[i] = seg{st, en, i + 1}
	}
	totalMin := k * 1440
	workTimes := make([]int, 0, totalMin)
	fb := make([]bool, 1440)
	for _, b := range breaks {
		for t := b.start; t <= b.end; t++ {
			fb[t] = true
		}
	}
	for d := 0; d < k; d++ {
		base := d * 1440
		for t := 0; t < 1440; t++ {
			if !fb[t] {
				workTimes = append(workTimes, base+t)
			}
		}
	}
	capMax := len(workTimes)
	type job struct{ p, c, cap, idx int }
	jobs := make([]job, 0, n)
	for i := 0; i < n; i++ {
		var ssub, tstr string
		var di, ci int
		fmt.Fscan(reader, &ssub, &di, &tstr, &ci)
		p, ok := tsolve[ssub]
		if !ok {
			continue
		}
		dd := parseHM(tstr)
		deadline := (di-1)*1440 + dd
		lo, hi := 0, len(workTimes)
		for lo < hi {
			mid := (lo + hi) / 2
			if workTimes[mid] < deadline {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
		cap := lo
		jobs = append(jobs, job{p: p, c: ci, cap: cap, idx: i + 1})
	}
	sort.Slice(jobs, func(i, j int) bool { return jobs[i].cap < jobs[j].cap })
	negInf := -1 << 60
	dp := make([]int, capMax+1)
	for i := 1; i <= capMax; i++ {
		dp[i] = negInf
	}
	picks := make([][]bool, len(jobs))
	for i := range picks {
		picks[i] = make([]bool, capMax+1)
	}
	curCap := capMax
	for i, jb := range jobs {
		for t := jb.cap + 1; t <= curCap; t++ {
			dp[t] = negInf
		}
		for t := jb.cap; t >= jb.p; t-- {
			if dp[t-jb.p]+jb.c > dp[t] {
				dp[t] = dp[t-jb.p] + jb.c
				picks[i][t] = true
			}
		}
		curCap = jb.cap
	}
	bestProfit := 0
	bestT := 0
	for t := 0; t <= curCap; t++ {
		if dp[t] > bestProfit {
			bestProfit = dp[t]
			bestT = t
		}
	}
	sel := make([]job, 0)
	t := bestT
	for i := len(jobs) - 1; i >= 0; i-- {
		if t >= 0 && picks[i][t] {
			sel = append(sel, jobs[i])
			t -= jobs[i].p
		}
	}
	for i, j := 0, len(sel)-1; i < j; i, j = i+1, j-1 {
		sel[i], sel[j] = sel[j], sel[i]
	}
	sort.Slice(breaks, func(i, j int) bool { return breaks[i].end < breaks[j].end })
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", bestProfit))
	sb.WriteString(fmt.Sprintf("%d\n", len(sel)))
	offset := 0
	for _, jb := range sel {
		startG := workTimes[offset]
		finishG := workTimes[offset+jb.p-1]
		offset += jb.p
		minuteInDay := startG % 1440
		ti := 4
		prev := -1
		for _, b := range breaks {
			if b.end < minuteInDay && b.end > prev {
				ti = b.idx
				prev = b.end
			}
		}
		h0 := minuteInDay / 60
		m0 := minuteInDay % 60
		finDay := finishG/1440 + 1
		finMin := finishG % 1440
		h1 := finMin / 60
		m1 := finMin % 60
		sb.WriteString(fmt.Sprintf("%d %d %02d:%02d %d %02d:%02d\n", jb.idx, ti, h0, m0, finDay, h1, m1))
	}
	return sb.String()
}

func genTestE() (string, string) {
	m := 2
	n := 2 + rand.Intn(2)
	k := 1
	subjects := []string{"math", "prog"}
	times := []int{rand.Intn(3) + 1, rand.Intn(3) + 1}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", m, n, k)
	fmt.Fprintf(&sb, "%s\n%s\n", subjects[0], subjects[1])
	fmt.Fprintf(&sb, "%d %d\n", times[0], times[1])
	sb.WriteString("00:00-00:10\n01:00-01:10\n02:00-02:10\n03:00-03:10\n")
	for i := 0; i < n; i++ {
		subj := subjects[rand.Intn(len(subjects))]
		day := 1
		exTime := fmt.Sprintf("0%d:30", 4+i)
		cost := rand.Intn(10) + 1
		fmt.Fprintf(&sb, "%s %d %s %d\n", subj, day, exTime, cost)
	}
	in := sb.String()
	out := solveE(in)
	return in, out
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go <binary>")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genTestE()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nInput:\n%s\nOutput:\n%s\n", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
