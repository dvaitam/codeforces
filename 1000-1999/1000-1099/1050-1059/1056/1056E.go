package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var s, t string
    if _, err := fmt.Fscan(reader, &s, &t); err != nil {
        return
    }
    n := len(s)
    m := len(t)
    cnt := [2]int{}
    for i := 0; i < n; i++ {
        if s[i] == '0' {
            cnt[0]++
        } else {
            cnt[1]++
        }
    }
    const base = 228
    pow := make([]uint64, m+1)
    h := make([]uint64, m+1)
    pow[0] = 1
    for i := 0; i < m; i++ {
        pow[i+1] = pow[i] * base
        h[i+1] = h[i]*base + uint64(t[i]-'a'+1)
    }
    ans := 0
    // try length for r0 = a
    if cnt[0] > 0 {
        maxA := (m - 1) / cnt[0]
        for a := 1; a <= maxA; a++ {
            rem := m - a*cnt[0]
            if rem <= 0 || rem%cnt[1] != 0 {
                continue
            }
            b := rem / cnt[1]
            // check mapping
            it := 0
            var h0, h1 uint64
            seen0, seen1 := false, false
            ok := true
            for i := 0; i < n; i++ {
                if s[i] == '0' {
                    cur := h[it+a] - h[it]*pow[a]
                    if !seen0 {
                        h0 = cur
                        seen0 = true
                    } else if h0 != cur {
                        ok = false
                        break
                    }
                    it += a
                } else {
                    cur := h[it+b] - h[it]*pow[b]
                    if !seen1 {
                        h1 = cur
                        seen1 = true
                    } else if h1 != cur {
                        ok = false
                        break
                    }
                    it += b
                }
            }
            if ok && (h0 != h1 || a != b) {
                ans++
            }
        }
    }
    fmt.Println(ans)
}
