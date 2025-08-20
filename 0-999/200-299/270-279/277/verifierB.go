package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"
)

type Point struct{ x, y int64 }

func runBinary(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errb bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errb
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("%v\n%s", err, errb.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
    m := rng.Intn(98) + 3 // 3..100
    n := rng.Intn(m+1) + m // m..2m
    return fmt.Sprintf("%d %d\n", n, m)
}

func parsePoints(out string, n int) ([]Point, error) {
    fields := strings.Fields(out)
    if len(fields) == 1 && fields[0] == "-1" {
        return nil, fmt.Errorf("output -1 for feasible case")
    }
    if len(fields) < 2*n {
        return nil, fmt.Errorf("expected %d integers, got %d", 2*n, len(fields))
    }
    if len(fields) > 2*n {
        // be strict to catch extra output
        return nil, fmt.Errorf("got extra output: %d integers (need %d)", len(fields), 2*n)
    }
    pts := make([]Point, 0, n)
    for i := 0; i < 2*n; i += 2 {
        xi, err := strconv.ParseInt(fields[i], 10, 64)
        if err != nil {
            return nil, fmt.Errorf("invalid int: %v", err)
        }
        yi, err := strconv.ParseInt(fields[i+1], 10, 64)
        if err != nil {
            return nil, fmt.Errorf("invalid int: %v", err)
        }
        pts = append(pts, Point{xi, yi})
    }
    return pts, nil
}

func withinBounds(pts []Point) error {
    const LIM = 100000000
    for i, p := range pts {
        if p.x < -LIM || p.x > LIM || p.y < -LIM || p.y > LIM {
            return fmt.Errorf("point %d out of bounds: %d %d", i+1, p.x, p.y)
        }
    }
    return nil
}

func uniquePoints(pts []Point) error {
    seen := make(map[Point]struct{}, len(pts))
    for i, p := range pts {
        if _, ok := seen[p]; ok {
            return fmt.Errorf("duplicate point at index %d: %d %d", i+1, p.x, p.y)
        }
        seen[p] = struct{}{}
    }
    return nil
}

func collinear(a, b, c Point) bool {
    // cross((b-a),(c-a)) == 0
    x1 := b.x - a.x
    y1 := b.y - a.y
    x2 := c.x - a.x
    y2 := c.y - a.y
    return x1*y2-y1*x2 == 0
}

func noThreeCollinear(pts []Point) error {
    n := len(pts)
    for i := 0; i < n; i++ {
        for j := i + 1; j < n; j++ {
            for k := j + 1; k < n; k++ {
                if collinear(pts[i], pts[j], pts[k]) {
                    return fmt.Errorf("three collinear: (%d,%d), (%d,%d), (%d,%d)",
                        pts[i].x, pts[i].y, pts[j].x, pts[j].y, pts[k].x, pts[k].y)
                }
            }
        }
    }
    return nil
}

func cross(o, a, b Point) int64 {
    return (a.x-o.x)*(b.y-o.y) - (a.y-o.y)*(b.x-o.x)
}

func convexHullSize(pts []Point) int {
    n := len(pts)
    if n <= 1 {
        return n
    }
    // sort by x, then y
    ps := make([]Point, n)
    copy(ps, pts)
    // simple insertion sort due to small n (<= 200)
    for i := 1; i < n; i++ {
        j := i
        for j > 0 && (ps[j].x < ps[j-1].x || (ps[j].x == ps[j-1].x && ps[j].y < ps[j-1].y)) {
            ps[j], ps[j-1] = ps[j-1], ps[j]
            j--
        }
    }
    lower := make([]Point, 0, n)
    for _, p := range ps {
        for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], p) <= 0 {
            lower = lower[:len(lower)-1]
        }
        lower = append(lower, p)
    }
    upper := make([]Point, 0, n)
    for i := n - 1; i >= 0; i-- {
        p := ps[i]
        for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
            upper = upper[:len(upper)-1]
        }
        upper = append(upper, p)
    }
    // concatenate lower and upper, removing duplicate endpoints
    hull := len(lower) + len(upper) - 2
    if hull < 0 {
        hull = 0
    }
    return hull
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    seed := time.Now().UnixNano()
    fmt.Fprintf(os.Stderr, "seed: %d\n", seed)
    rng := rand.New(rand.NewSource(seed))
    for i := 1; i <= 100; i++ {
        input := generateCase(rng)
        out, err := runBinary(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:%s", i, err, input)
            os.Exit(1)
        }
        // parse input to get n,m
        var n, m int
        fmt.Sscanf(input, "%d %d", &n, &m)
        pts, err := parsePoints(out, n)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed to parse output: %v\ninput:%soutput:%s\n", i, err, input, out)
            os.Exit(1)
        }
        if err := withinBounds(pts); err != nil {
            fmt.Fprintf(os.Stderr, "case %d bounds error: %v\ninput:%soutput:%s\n", i, err, input, out)
            os.Exit(1)
        }
        if err := uniquePoints(pts); err != nil {
            fmt.Fprintf(os.Stderr, "case %d duplicate error: %v\ninput:%soutput:%s\n", i, err, input, out)
            os.Exit(1)
        }
        if err := noThreeCollinear(pts); err != nil {
            fmt.Fprintf(os.Stderr, "case %d collinearity error: %v\ninput:%soutput:%s\n", i, err, input, out)
            os.Exit(1)
        }
        h := convexHullSize(pts)
        if h != m {
            fmt.Fprintf(os.Stderr, "case %d failed convexity check: hull=%d, want=%d\ninput:%soutput:%s\n", i, h, m, input, out)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}
