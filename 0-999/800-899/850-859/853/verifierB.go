package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type flight struct {
	day  int
	city int
	cost int64
	from bool // true if arrival (to 0), false if departure (from 0)
}

func solveCase(line string) string {
	fields := strings.Fields(line)
	idx := 0
	if len(fields) < 3 {
		return ""
	}
	n, _ := strconv.Atoi(fields[idx])
	idx++
	m, _ := strconv.Atoi(fields[idx])
	idx++
	k, _ := strconv.Atoi(fields[idx])
	idx++
	flights := make([]flight, m)
	for i := 0; i < m; i++ {
		d, _ := strconv.Atoi(fields[idx])
		idx++
		f, _ := strconv.Atoi(fields[idx])
		idx++
		t, _ := strconv.Atoi(fields[idx])
		idx++
		c, _ := strconv.ParseInt(fields[idx], 10, 64)
		idx++
		if t == 0 {
			flights[i] = flight{day: d, city: f, cost: c, from: true}
		} else if f == 0 {
			flights[i] = flight{day: d, city: t, cost: c, from: false}
		} else {
			// ignore other flights
			flights[i] = flight{day: d, city: 0, cost: c, from: true}
		}
	}

	const INF int64 = 1 << 60
	maxDay := 0
	for _, f := range flights {
		if f.day > maxDay {
			maxDay = f.day
		}
	}
	maxIndex := maxDay + k + 5
	arrCost := make([]int64, maxIndex+2)
	depCost := make([]int64, maxIndex+2)
	for i := range arrCost {
		arrCost[i] = INF
		depCost[i] = INF
	}

	arrivals := make([]flight, 0)
	departures := make([]flight, 0)
	for _, f := range flights {
		if f.from {
			arrivals = append(arrivals, f)
		} else {
			departures = append(departures, f)
		}
	}
	sort.Slice(arrivals, func(i, j int) bool { return arrivals[i].day < arrivals[j].day })
	bestArr := make([]int64, n+1)
	for i := range bestArr {
		bestArr[i] = INF
	}
	sum := int64(0)
	cntMissing := n
	idxA := 0
	for day := 0; day <= maxIndex; day++ {
		for idxA < len(arrivals) && arrivals[idxA].day == day {
			f := arrivals[idxA]
			if bestArr[f.city] == INF {
				bestArr[f.city] = f.cost
				sum += f.cost
				cntMissing--
			} else if f.cost < bestArr[f.city] {
				sum += f.cost - bestArr[f.city]
				bestArr[f.city] = f.cost
			}
			idxA++
		}
		if cntMissing == 0 {
			arrCost[day] = sum
		}
	}
	for day := 1; day <= maxIndex; day++ {
		if arrCost[day] > arrCost[day-1] {
			arrCost[day] = arrCost[day-1]
		}
	}

	sort.Slice(departures, func(i, j int) bool { return departures[i].day > departures[j].day })
	bestDep := make([]int64, n+1)
	for i := range bestDep {
		bestDep[i] = INF
	}
	sum = 0
	cntMissing = n
	idxD := 0
	for day := maxIndex; day >= 0; day-- {
		for idxD < len(departures) && departures[idxD].day == day {
			f := departures[idxD]
			if bestDep[f.city] == INF {
				bestDep[f.city] = f.cost
				sum += f.cost
				cntMissing--
			} else if f.cost < bestDep[f.city] {
				sum += f.cost - bestDep[f.city]
				bestDep[f.city] = f.cost
			}
			idxD++
		}
		if cntMissing == 0 {
			depCost[day] = sum
		}
	}
	for day := maxIndex - 1; day >= 0; day-- {
		if depCost[day] > depCost[day+1] {
			depCost[day] = depCost[day+1]
		}
	}

	ans := INF
	for day := 0; day <= maxDay; day++ {
		dIdx := day + k + 1
		if dIdx > maxIndex {
			break
		}
		if arrCost[day] == INF || depCost[dIdx] == INF {
			continue
		}
		total := arrCost[day] + depCost[dIdx]
		if total < ans {
			ans = total
		}
	}
	if ans == INF {
		return "-1"
	}
	return fmt.Sprintf("%d", ans)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesRaw = `1 4 3 2 1 0 7 3 1 0 2 4 1 0 9 10 0 1 2
1 2 3 4 1 0 7 11 0 1 1
2 6 2 1 1 0 6 2 0 1 4 3 1 0 6 3 2 0 1 4 0 1 2 8 0 2 6
1 4 1 1 0 1 3 1 1 0 1 3 1 0 6 5 0 1 6
3 7 2 1 1 0 4 3 3 0 3 4 2 0 5 6 0 1 6 6 0 3 2 9 0 1 7 9 0 2 9
1 2 1 3 1 0 9 6 0 1 4
1 3 2 1 1 0 3 4 1 0 1 8 0 1 10
2 8 2 3 0 2 2 3 2 0 2 4 1 0 3 4 2 0 2 6 0 1 7 7 0 2 8 7 2 0 8 10 0 1 5
3 11 3 2 0 3 5 2 1 0 8 4 2 0 7 4 3 0 6 4 3 0 8 8 0 1 8 9 0 1 3 9 0 2 4 9 0 2 8 9 0 3 6 10 1 0 5
1 3 3 1 1 0 1 2 1 0 9 9 0 1 5
2 8 2 2 2 0 5 2 2 0 5 3 2 0 5 4 1 0 8 7 0 1 5 7 0 2 10 9 1 0 3 9 1 0 9
1 4 2 2 0 1 7 4 1 0 5 6 0 1 3 9 0 1 8
3 7 2 2 2 0 8 3 1 0 6 3 3 0 8 5 2 0 1 6 0 3 8 8 0 2 9 9 0 1 7
3 6 1 1 2 0 2 3 1 0 8 3 3 0 1 5 0 3 10 6 0 2 1 8 0 1 1
1 4 2 3 1 0 4 5 1 0 8 6 1 0 7 7 0 1 10
3 9 3 2 1 0 8 3 2 0 3 4 3 0 6 7 0 1 9 7 3 0 2 8 0 2 3 8 0 3 1 8 1 0 4 11 0 3 7
1 2 2 2 1 0 4 7 0 1 10
1 4 3 1 0 1 10 3 1 0 6 8 1 0 4 10 0 1 8
2 4 2 1 1 0 1 1 2 0 5 4 0 2 3 5 0 1 6
1 4 2 2 0 1 6 2 1 0 10 4 1 0 10 8 0 1 9
2 7 3 1 2 0 4 3 1 0 5 4 0 2 4 4 2 0 5 7 0 1 9 7 0 2 2 9 0 2 8
2 5 3 1 1 0 5 4 2 0 8 5 0 1 1 7 0 2 8 11 0 2 1
1 2 1 2 1 0 2 5 0 1 7
1 4 2 1 1 0 7 3 1 0 3 5 0 1 7 7 0 1 1
1 3 3 2 1 0 5 9 0 1 9 10 1 0 5
1 2 1 1 1 0 3 6 0 1 4
1 4 2 1 1 0 9 2 1 0 2 7 0 1 1 10 1 0 4
3 6 1 1 3 0 9 3 1 0 6 4 2 0 5 5 0 3 6 6 0 1 6 7 0 2 3
3 6 3 1 2 0 3 2 1 0 7 2 3 0 4 6 0 2 8 7 0 1 3 7 0 3 4
2 8 2 1 2 0 9 2 1 0 8 3 0 1 6 3 0 2 8 5 0 1 10 6 0 2 8 8 1 0 9 9 1 0 5
2 4 2 2 1 0 1 3 2 0 8 7 0 1 4 7 0 2 5
2 5 3 1 2 0 3 3 1 0 1 7 0 2 1 8 0 1 4 9 0 1 10
1 2 1 3 1 0 2 5 0 1 6
3 7 2 2 1 0 8 3 2 0 5 4 3 0 7 5 0 3 8 7 0 2 6 7 0 3 7 8 0 1 3
1 2 3 4 1 0 7 9 0 1 1
3 10 1 2 0 2 3 2 1 0 2 2 2 0 4 2 3 0 2 2 3 0 3 4 0 2 3 6 0 1 4 7 0 3 10 7 1 0 8 10 0 2 6
2 8 1 3 1 0 3 3 2 0 8 5 0 2 9 6 1 0 5 8 0 1 6 8 0 2 1 8 0 2 1 8 1 0 7
3 8 1 3 1 0 4 3 2 0 4 3 3 0 1 3 3 0 4 5 0 2 10 7 0 1 7 7 0 3 3 10 0 3 3
2 6 2 1 0 2 5 3 1 0 6 4 2 0 3 7 0 2 3 9 0 1 7 10 1 0 6
3 7 3 1 1 0 9 2 2 0 9 2 3 0 10 6 0 3 9 7 0 2 2 8 0 1 5 9 0 2 5
2 4 2 3 2 0 6 4 1 0 5 7 0 2 7 9 0 1 9
1 4 1 2 1 0 10 4 0 1 9 7 0 1 6 8 0 1 10
1 2 3 3 1 0 8 7 0 1 4
3 8 1 1 3 0 4 2 0 3 7 2 1 0 7 4 0 1 10 4 0 3 8 4 2 0 4 5 2 0 10 9 0 2 9
2 6 1 1 1 0 8 2 2 0 7 3 0 1 5 7 0 1 10 7 0 2 7 8 1 0 10
3 8 2 1 1 0 7 2 2 0 1 3 3 0 6 4 0 3 9 5 0 2 10 6 0 3 2 7 0 1 9 8 1 0 6
2 6 1 3 0 2 1 3 2 0 9 4 1 0 3 5 0 2 3 7 0 1 7 9 2 0 3
2 4 3 1 1 0 9 2 2 0 5 6 0 1 5 8 0 2 5
2 7 3 2 1 0 8 2 2 0 9 2 2 0 10 3 1 0 9 4 0 1 6 6 0 2 6 7 0 1 1
1 3 3 1 1 0 2 5 0 1 7 10 1 0 6
2 6 3 1 1 0 3 2 2 0 3 5 0 1 10 6 0 2 6 7 0 1 1 10 2 0 1
3 8 1 1 1 0 7 2 1 0 1 3 0 1 8 3 2 0 10 4 3 0 7 7 0 2 5 7 0 3 1 8 0 1 7
2 5 1 2 0 2 5 4 1 0 10 4 2 0 9 6 0 1 6 9 0 2 8
1 4 3 1 1 0 2 2 1 0 3 7 0 1 6 9 0 1 1
1 3 2 2 1 0 9 3 0 1 4 6 0 1 3
3 6 2 2 2 0 5 2 3 0 6 4 1 0 7 6 0 3 3 7 0 1 4 7 0 2 3
2 6 2 1 1 0 10 1 2 0 7 1 2 0 9 2 1 0 7 5 0 1 3 5 0 2 3
3 12 1 2 0 2 9 2 1 0 7 2 3 0 6 4 0 1 7 4 0 2 10 4 1 0 7 4 2 0 10 5 0 3 9 6 0 2 7 6 0 3 10 7 2 0 10 8 0 2 4
2 6 1 1 1 0 10 1 2 0 9 3 0 2 3 5 0 1 7 7 0 2 10 7 2 0 7
2 5 1 3 1 0 8 3 2 0 9 5 0 2 6 6 2 0 8 8 0 1 3
3 12 3 1 1 0 4 3 1 0 4 3 2 0 5 3 3 0 9 6 1 0 10 6 2 0 10 7 0 1 3 7 0 3 3 7 3 0 2 7 3 0 4 10 0 2 1 10 0 3 1
1 2 3 4 1 0 6 10 0 1 6
3 12 1 2 0 3 4 3 2 0 9 3 3 0 4 3 3 0 9 4 1 0 5 5 0 3 3 6 0 2 7 7 0 1 3 7 0 1 6 7 0 3 5 8 0 3 6 10 0 3 2
2 4 1 1 1 0 2 3 2 0 2 5 0 1 8 8 0 2 6
3 8 3 1 2 0 9 2 3 0 9 4 1 0 4 4 3 0 7 5 1 0 7 6 0 2 1 9 0 1 9 9 0 3 10
1 4 2 1 1 0 7 4 0 1 10 7 1 0 1 8 0 1 9
2 7 3 2 1 0 2 3 1 0 8 3 2 0 3 4 1 0 5 7 0 2 9 9 0 1 4 9 0 2 10
2 6 1 1 1 0 3 2 2 0 1 3 0 1 3 4 0 2 7 8 0 2 2 9 0 2 2
3 7 1 2 1 0 8 2 2 0 1 3 3 0 2 5 0 3 5 6 0 1 5 7 0 2 6 10 3 0 6
2 5 1 3 1 0 4 3 2 0 7 6 0 2 2 7 2 0 4 8 0 1 6
3 8 3 3 1 0 7 3 3 0 2 4 2 0 5 6 0 3 3 6 1 0 7 8 0 2 8 9 0 1 6 10 0 3 3
2 8 2 1 0 1 7 2 2 0 3 3 2 0 8 4 1 0 10 6 1 0 7 7 0 2 6 8 0 1 7 9 0 1 8
1 4 1 1 1 0 3 2 0 1 2 3 1 0 1 6 0 1 1
3 7 2 1 1 0 9 2 2 0 4 4 1 0 10 4 3 0 8 6 0 2 7 7 0 1 2 7 0 3 7
2 8 1 2 0 2 1 2 2 0 10 3 0 1 7 3 1 0 6 4 0 2 4 4 0 2 5 7 0 1 8 7 0 1 9
1 3 3 1 1 0 9 3 1 0 3 8 0 1 8
2 7 2 1 1 0 10 4 2 0 4 6 0 1 4 7 0 2 4 7 1 0 9 7 2 0 8 9 0 2 8
1 3 1 3 0 1 4 3 1 0 6 8 0 1 10
2 4 3 1 1 0 3 3 2 0 10 7 0 1 1 8 0 2 8
2 5 3 4 0 1 9 4 1 0 8 4 2 0 1 9 0 1 3 10 0 2 9
3 12 3 2 2 0 10 2 3 0 8 4 0 1 9 4 1 0 3 4 1 0 3 4 3 0 10 6 0 2 8 6 1 0 3 7 0 3 5 8 0 3 9 9 0 2 1 11 0 1 10
2 6 2 1 1 0 7 1 2 0 3 4 1 0 2 7 0 2 8 8 0 1 4 9 2 0 1
3 12 3 1 0 3 6 2 0 1 10 2 2 0 9 3 1 0 10 4 1 0 9 4 2 0 5 4 3 0 9 8 0 3 2 8 1 0 7 9 0 1 1 9 0 2 5 10 0 2 1
2 8 3 1 2 0 3 1 2 0 4 3 1 0 4 6 1 0 2 7 0 1 3 7 0 2 10 9 0 1 2 10 2 0 7
1 3 3 3 1 0 7 5 1 0 9 7 0 1 6
3 8 1 1 1 0 3 1 3 0 10 3 2 0 6 5 0 3 2 6 0 1 1 6 0 2 6 6 0 3 2 10 1 0 9
1 2 2 2 1 0 3 6 0 1 7
1 4 3 4 1 0 6 6 0 1 4 8 0 1 5 9 0 1 5
3 12 3 1 1 0 2 1 3 0 5 2 0 2 6 3 2 0 1 3 3 0 7 5 0 3 3 5 0 3 9 5 1 0 5 6 0 1 7 7 0 1 5 7 0 2 10 9 0 2 2
1 3 3 3 1 0 9 4 0 1 5 9 0 1 1
3 12 1 2 2 0 5 4 1 0 7 4 3 0 8 5 1 0 3 5 2 0 8 6 0 2 9 6 0 3 7 6 1 0 4 7 1 0 8 8 0 1 10 8 1 0 2 8 2 0 9
3 7 1 2 2 0 5 3 1 0 4 3 3 0 9 4 0 3 7 5 0 2 3 5 0 3 3 6 0 1 3
3 7 1 1 3 0 7 2 1 0 10 3 2 0 1 4 0 3 7 5 0 1 4 8 0 2 4 9 2 0 9
3 6 2 1 2 0 4 2 1 0 2 3 3 0 7 5 0 1 2 7 0 2 8 8 0 3 7
2 4 3 1 2 0 1 2 1 0 10 5 0 2 8 7 0 1 4
1 4 2 2 1 0 10 5 0 1 3 7 0 1 1 10 1 0 6
1 4 2 2 1 0 10 5 0 1 6 8 0 1 6 8 1 0 8
3 9 3 3 0 3 9 3 2 0 10 4 1 0 10 4 3 0 4 6 2 0 8 8 0 2 4 9 0 3 8 10 0 3 9 11 0 1 9
3 6 3 1 2 0 4 1 3 0 4 2 1 0 5 5 0 2 5 6 0 3 10 9 0 1 2
1 3 3 3 1 0 1 7 0 1 1 10 0 1 5`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		expected := solveCase(line)
		got, err := run(bin, line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
