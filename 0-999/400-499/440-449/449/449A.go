package main

import (
   "fmt"
)

func main() {
   var n, m, k int64
   if _, err := fmt.Scan(&n, &m, &k); err != nil {
       return
   }
   // Maximum possible cuts is (n-1)+(m-1)
   if k > n+m-2 {
       fmt.Println(-1)
       return
   }
   var ans int64 = 0
   // Case: use all cuts horizontally (if possible)
   if k <= n-1 {
       // x = k horizontal cuts, y = 0 vertical
       h := n / (k + 1)
       area := h * m
       if area > ans {
           ans = area
       }
   }
   // Case: use all cuts vertically
   if k <= m-1 {
       w := m / (k + 1)
       area := w * n
       if area > ans {
           ans = area
       }
   }
   // Case: saturate horizontal cuts (n-1), rest vertical
   if k > n-1 {
       // use x = n-1 horizontal cuts, y = k-(n-1)
       y := k - (n - 1)
       // strips in horizontal: n, so height = 1
       w := m / (y + 1)
       area := w
       if area > ans {
           ans = area
       }
   }
   // Case: saturate vertical cuts (m-1), rest horizontal
   if k > m-1 {
       x := k - (m - 1)
       h := n / (x + 1)
       area := h
       if area > ans {
           ans = area
       }
   }
   fmt.Println(ans)
}
