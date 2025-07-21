package main
import (
    "bufio"
    "fmt"
    "math"
    "os"
)
func main() {
    reader := bufio.NewReader(os.Stdin)
    var s string
    if _, err := fmt.Fscan(reader, &s); err != nil {
        return
    }
    // build P' with breaks between same foot
    n := len(s)
    p := make([]byte, 0, n*2)
    for i := 0; i < n; i++ {
        p = append(p, s[i])
        if s[i] != 'X' && i+1 < n && s[i+1] == s[i] {
            p = append(p, 'X')
        }
    }
    // ensure even length
    if len(p)%2 == 1 {
        p = append(p, 'X')
    }
    m := len(p)
    // count matches for two alignments
    var m0, m1 int
    for i := 0; i < m; i++ {
        c := p[i]
        if c == 'L' {
            if i%2 == 0 {
                m0++
            }
            if i%2 == 1 {
                m1++
            }
        } else if c == 'R' {
            if i%2 == 1 {
                m0++
            }
            if i%2 == 0 {
                m1++
            }
        }
    }
    best := m0
    if m1 > best {
        best = m1
    }
    res := float64(best) * 100.0 / float64(m)
    // truncate to 6 decimals
    res = math.Floor(res*1e6) / 1e6
    fmt.Printf("%.6f", res)
}
