package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Driver struct {
	name   string
	points int
	pos    []int
}

func expected(races [][]string) (string, string) {
	ptsAward := []int{25, 18, 15, 12, 10, 8, 6, 4, 2, 1}
	const maxPos = 50
	
	driverMap := make(map[string]*Driver)
	
	for _, race := range races {
		for pos, name := range race {
			if _, exists := driverMap[name]; !exists {
				driverMap[name] = &Driver{
					name:   name,
					points: 0,
					pos:    make([]int, maxPos),
				}
			}
			
			driver := driverMap[name]
			if pos < maxPos {
				driver.pos[pos]++
			}
			if pos < len(ptsAward) {
				driver.points += ptsAward[pos]
			}
		}
	}
	
	drivers := make([]*Driver, 0, len(driverMap))
	for _, d := range driverMap {
		drivers = append(drivers, d)
	}
	
	originalDrivers := make([]*Driver, len(drivers))
	copy(originalDrivers, drivers)
	sort.Slice(originalDrivers, func(i, j int) bool {
		a, b := originalDrivers[i], originalDrivers[j]
		if a.points != b.points {
			return a.points > b.points
		}
		for k := 0; k < maxPos; k++ {
			if a.pos[k] != b.pos[k] {
				return a.pos[k] > b.pos[k]
			}
		}
		return false
	})
	
	alternativeDrivers := make([]*Driver, len(drivers))
	copy(alternativeDrivers, drivers)
	sort.Slice(alternativeDrivers, func(i, j int) bool {
		a, b := alternativeDrivers[i], alternativeDrivers[j]
		if a.pos[0] != b.pos[0] {
			return a.pos[0] > b.pos[0]
		}
		if a.points != b.points {
			return a.points > b.points
		}
		for k := 1; k < maxPos; k++ {
			if a.pos[k] != b.pos[k] {
				return a.pos[k] > b.pos[k]
			}
		}
		return false
	})
	
	originalChamp := ""
	alternativeChamp := ""
	if len(originalDrivers) > 0 {
		originalChamp = originalDrivers[0].name
	}
	if len(alternativeDrivers) > 0 {
		alternativeChamp = alternativeDrivers[0].name
	}
	
	return originalChamp, alternativeChamp
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

func generateDriverName(rng *rand.Rand) string {
	length := rng.Intn(10) + 3
	chars := make([]byte, length)
	for i := 0; i < length; i++ {
		if rng.Intn(2) == 0 {
			chars[i] = byte('a' + rng.Intn(26))
		} else {
			chars[i] = byte('A' + rng.Intn(26))
		}
	}
	return string(chars)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		t := rng.Intn(20) + 1
		
		driverPool := make([]string, 0)
		for j := 0; j < rng.Intn(20)+5; j++ {
			driverPool = append(driverPool, generateDriverName(rng))
		}
		
		races := make([][]string, t)
		for j := 0; j < t; j++ {
			n := rng.Intn(len(driverPool)) + 1
			if n > 50 {
				n = 50
			}
			
			selected := make(map[int]bool)
			race := make([]string, n)
			for k := 0; k < n; k++ {
				idx := rng.Intn(len(driverPool))
				for selected[idx] {
					idx = rng.Intn(len(driverPool))
				}
				selected[idx] = true
				race[k] = driverPool[idx]
			}
			races[j] = race
		}
		
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", t))
		for _, race := range races {
			input.WriteString(fmt.Sprintf("%d\n", len(race)))
			for _, driver := range race {
				input.WriteString(driver + "\n")
			}
		}
		
		expectedOrig, expectedAlt := expected(races)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		
		lines := strings.Split(got, "\n")
		if len(lines) != 2 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected 2 lines, got %d\ninput:\n%s", i+1, len(lines), input.String())
			os.Exit(1)
		}
		
		gotOrig := strings.TrimSpace(lines[0])
		gotAlt := strings.TrimSpace(lines[1])
		
		if gotOrig != expectedOrig {
			fmt.Fprintf(os.Stderr, "case %d failed: original champion expected %s got %s\ninput:\n%s", i+1, expectedOrig, gotOrig, input.String())
			os.Exit(1)
		}
		
		if gotAlt != expectedAlt {
			fmt.Fprintf(os.Stderr, "case %d failed: alternative champion expected %s got %s\ninput:\n%s", i+1, expectedAlt, gotAlt, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}