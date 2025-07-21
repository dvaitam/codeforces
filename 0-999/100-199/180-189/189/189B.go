package main

import (
    "fmt"
)

func main() {
    var w, h int64
    if _, err := fmt.Scan(&w, &h); err != nil {
        return
    }
    // Number of choices for half-diagonal lengths along width and height
    halfW := w / 2
    ceilHalfW := (w + 1) / 2
    halfH := h / 2
    ceilHalfH := (h + 1) / 2
    // Sum over a from 1..floor(w/2) of (w - 2a + 1) = floor(w/2) * ceil(w/2)
    countW := halfW * ceilHalfW
    // Similarly for height
    countH := halfH * ceilHalfH
    // Total rhombi is product of independent sums
    ans := countW * countH
    fmt.Println(ans)
}
