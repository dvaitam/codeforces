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
)

// Embedded solver from the accepted solution

type Point struct {
	x, y float64
}

type Interval struct {
	L, R     float64
	A, B     float64
	w_id     int
	is_entry bool
}

type Event struct {
	angle  float64
	typ    int
	iv_idx int
}

func oracleSolve(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	scanInt := func() int {
		if !scanner.Scan() {
			return 0
		}
		res, _ := strconv.Atoi(scanner.Text())
		return res
	}

	scanFloat := func() float64 {
		if !scanner.Scan() {
			return 0
		}
		res, _ := strconv.ParseFloat(scanner.Text(), 64)
		return res
	}

	N := scanInt()

	var intervals []Interval

	for i := 0; i < N; i++ {
		k := scanInt()
		pts := make([]Point, k)
		for j := 0; j < k; j++ {
			pts[j].x = scanFloat()
			pts[j].y = scanFloat()
		}

		for j := 0; j < k/2; j++ {
			pts[j], pts[k-1-j] = pts[k-1-j], pts[j]
		}

		thetas := make([]float64, k+1)
		for j := 0; j < k; j++ {
			thetas[j] = math.Atan2(pts[j].y, pts[j].x)
			if thetas[j] < 0 {
				thetas[j] += 2 * math.Pi
			}
		}
		for j := 1; j < k; j++ {
			for thetas[j]-thetas[j-1] > math.Pi {
				thetas[j] -= 2 * math.Pi
			}
			for thetas[j]-thetas[j-1] < -math.Pi {
				thetas[j] += 2 * math.Pi
			}
		}
		thetas[k] = math.Atan2(pts[0].y, pts[0].x)
		if thetas[k] < 0 {
			thetas[k] += 2 * math.Pi
		}
		for thetas[k]-thetas[k-1] > math.Pi {
			thetas[k] -= 2 * math.Pi
		}
		for thetas[k]-thetas[k-1] < -math.Pi {
			thetas[k] += 2 * math.Pi
		}

		min_t := thetas[0]
		for j := 1; j <= k; j++ {
			if thetas[j] < min_t {
				min_t = thetas[j]
			}
		}
		shift := math.Floor(min_t / (2 * math.Pi))
		for j := 0; j <= k; j++ {
			thetas[j] -= shift * 2 * math.Pi
		}

		for j := 0; j < k; j++ {
			u := thetas[j]
			v := thetas[j+1]
			if math.Abs(u-v) < 1e-9 {
				continue
			}

			is_entry := u > v
			L, R := u, v
			if is_entry {
				L, R = v, u
			}

			p1, p2 := pts[j], pts[(j+1)%k]
			C := p1.x*p2.y - p2.x*p1.y
			if math.Abs(C) < 1e-9 {
				continue
			}
			A := (p2.y - p1.y) / C
			B := -(p2.x - p1.x) / C

			for L >= 2*math.Pi {
				L -= 2 * math.Pi
				R -= 2 * math.Pi
			}

			if R > 2*math.Pi {
				intervals = append(intervals, Interval{L: L, R: 2 * math.Pi, A: A, B: B, w_id: i, is_entry: is_entry})
				intervals = append(intervals, Interval{L: 0, R: R - 2*math.Pi, A: A, B: B, w_id: i, is_entry: is_entry})
			} else {
				intervals = append(intervals, Interval{L: L, R: R, A: A, B: B, w_id: i, is_entry: is_entry})
			}
		}
	}

	if len(intervals) == 0 {
		return fmt.Sprintf("%.4f", 0.0)
	}

	var events []Event
	for i := range intervals {
		events = append(events, Event{angle: intervals[i].L, typ: 1, iv_idx: i})
		events = append(events, Event{angle: intervals[i].R, typ: -1, iv_idx: i})
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].angle < events[j].angle
	})

	var total_area float64
	last_angle := events[0].angle

	active_entries := make([]int, 0, 1000)
	active_exits := make([]int, N)
	for i := range active_exits {
		active_exits[i] = -1
	}

	n_events := len(events)
	idx := 0

	for idx < n_events {
		j := idx
		for j < n_events && events[j].angle == events[idx].angle {
			j++
		}

		curr_angle := events[idx].angle
		if curr_angle > last_angle {
			theta_mid := (last_angle + curr_angle) / 2.0
			cos_t := math.Cos(theta_mid)
			sin_t := math.Sin(theta_mid)

			best_iv_idx := -1
			max_inv_d := -math.MaxFloat64

			for _, ev_idx := range active_entries {
				iv := &intervals[ev_idx]
				inv_d := iv.A*cos_t + iv.B*sin_t
				if inv_d > max_inv_d {
					max_inv_d = inv_d
					best_iv_idx = ev_idx
				}
			}

			if best_iv_idx != -1 {
				best_iv := &intervals[best_iv_idx]
				w_id := best_iv.w_id
				exit_iv_idx := active_exits[w_id]

				if exit_iv_idx != -1 {
					exit_iv := &intervals[exit_iv_idx]

					inv_d1_en := best_iv.A*math.Cos(last_angle) + best_iv.B*math.Sin(last_angle)
					inv_d2_en := best_iv.A*math.Cos(curr_angle) + best_iv.B*math.Sin(curr_angle)
					area_en := 0.5 * math.Sin(curr_angle-last_angle) / (inv_d1_en * inv_d2_en)

					inv_d1_ex := exit_iv.A*math.Cos(last_angle) + exit_iv.B*math.Sin(last_angle)
					inv_d2_ex := exit_iv.A*math.Cos(curr_angle) + exit_iv.B*math.Sin(curr_angle)
					area_ex := 0.5 * math.Sin(curr_angle-last_angle) / (inv_d1_ex * inv_d2_ex)

					total_area += area_ex - area_en
				}
			}
			last_angle = curr_angle
		}

		for k := idx; k < j; k++ {
			ev := events[k]
			iv := &intervals[ev.iv_idx]
			if ev.typ == 1 {
				if iv.is_entry {
					active_entries = append(active_entries, ev.iv_idx)
				} else {
					active_exits[iv.w_id] = ev.iv_idx
				}
			} else {
				if iv.is_entry {
					for fi, v := range active_entries {
						if v == ev.iv_idx {
							active_entries[fi] = active_entries[len(active_entries)-1]
							active_entries = active_entries[:len(active_entries)-1]
							break
						}
					}
				} else {
					if active_exits[iv.w_id] == ev.iv_idx {
						active_exits[iv.w_id] = -1
					}
				}
			}
		}

		idx = j
	}

	return fmt.Sprintf("%.4f", total_area)
}

// End of embedded solver

type Test struct{ input string }

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []Test {
	rand.Seed(2)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(3) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			c := 4
			x0 := (j + 1) * 10
			y0 := (j + 1) * 10
			w := rand.Intn(3) + 1
			h := rand.Intn(3) + 1
			fmt.Fprintf(&sb, "%d\n", c)
			fmt.Fprintf(&sb, "%d %d\n", x0, y0)
			fmt.Fprintf(&sb, "%d %d\n", x0+w, y0)
			fmt.Fprintf(&sb, "%d %d\n", x0+w, y0+h)
			fmt.Fprintf(&sb, "%d %d\n", x0, y0+h)
		}
		tests = append(tests, Test{sb.String()})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	tests := genTests()
	for i, tc := range tests {
		exp := oracleSolve(tc.input)
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
