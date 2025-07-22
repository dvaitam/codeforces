package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type event struct {
	angle float64
	delta int
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveE(n int, d float64, circles [][3]float64) string {
	events := make([]event, 0, n*8)
	twoPi := 2 * math.Pi
	for i := 0; i < n; i++ {
		xi, yi, ri := circles[i][0], circles[i][1], circles[i][2]
		Di := math.Hypot(xi, yi)
		phi := math.Atan2(yi, xi)
		kmin := int(math.Ceil((Di - ri) / d))
		if kmin < 1 {
			kmin = 1
		}
		kmax := int(math.Floor((Di + ri) / d))
		for k := kmin; k <= kmax; k++ {
			kd := float64(k) * d
			A := (kd - ri) / Di
			B := (kd + ri) / Di
			if A > 1 || B < -1 {
				continue
			}
			if A < -1 {
				A = -1
			}
			if B > 1 {
				B = 1
			}
			al := math.Acos(A)
			au := math.Acos(B)
			intervals := [][2]float64{{phi + au, phi + al}, {phi + twoPi - al, phi + twoPi - au}}
			for _, iv := range intervals {
				s, e := iv[0], iv[1]
				if s < 0 {
					s += twoPi * (math.Floor(-s/twoPi) + 1)
				}
				if e < 0 {
					e += twoPi * (math.Floor(-e/twoPi) + 1)
				}
				s = math.Mod(s, twoPi)
				e = math.Mod(e, twoPi)
				if e < s {
					events = append(events, event{s, 1})
					events = append(events, event{twoPi, -1})
					events = append(events, event{0, 1})
					events = append(events, event{e, -1})
				} else {
					events = append(events, event{s, 1})
					events = append(events, event{e, -1})
				}
			}
		}
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].angle == events[j].angle {
			return events[i].delta > events[j].delta
		}
		return events[i].angle < events[j].angle
	})
	maxCnt := 0
	cur := 0
	for _, ev := range events {
		cur += ev.delta
		if cur > maxCnt {
			maxCnt = cur
		}
	}
	return fmt.Sprintf("%d", maxCnt)
}

func generateCaseE(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	d := float64(rng.Intn(6) + 5)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %.0f\n", n, d)
	for i := 0; i < n; i++ {
		xi := float64(rng.Intn(21) - 10)
		yi := float64(rng.Intn(21) - 10)
		ri := float64(rng.Intn(5) + 1)
		for math.Hypot(xi, yi) <= ri {
			xi = float64(rng.Intn(21) - 10)
			yi = float64(rng.Intn(21) - 10)
		}
		fmt.Fprintf(&sb, "%.0f %.0f %.0f\n", xi, yi, ri)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseE(rng)
		scanner := bufio.NewScanner(strings.NewReader(tc))
		scanner.Split(bufio.ScanWords)
		var fields []string
		for scanner.Scan() {
			fields = append(fields, scanner.Text())
		}
		n, _ := strconv.Atoi(fields[0])
		dVal, _ := strconv.ParseFloat(fields[1], 64)
		circles := make([][3]float64, n)
		idx := 2
		for j := 0; j < n; j++ {
			x, _ := strconv.ParseFloat(fields[idx], 64)
			y, _ := strconv.ParseFloat(fields[idx+1], 64)
			r, _ := strconv.ParseFloat(fields[idx+2], 64)
			circles[j] = [3]float64{x, y, r}
			idx += 3
		}
		expect := solveE(n, dVal, circles)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
