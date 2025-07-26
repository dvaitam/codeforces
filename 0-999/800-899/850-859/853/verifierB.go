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
	file, err := os.Open("testcasesB.txt")
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
