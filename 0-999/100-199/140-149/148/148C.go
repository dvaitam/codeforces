package main

import (
  "fmt"
)

func main() {
  var n, m, k int
  if _, err := fmt.Scanf("%d%d%d", &n, &m, &k); err != nil {
    return
  }
  a := make([]int, n)
  a[0] = 1
  sum, mx, cnt := 1, 1, 1
  if m == 0 && k == 0 {
    for i := 0; i < n; i++ {
      fmt.Print(1, " ")
    }
    fmt.Println()
    return
  }
  if m+k+1 >= n && k == 0 {
    fmt.Println(-1)
    return
  }
  if m+k+1 < n {
    a[1] = 1
    sum = 2
    cnt = 2
  }
  for k > 0 {
    a[cnt] = sum + 1
    if a[cnt] > mx {
      mx = a[cnt]
    }
    sum += a[cnt]
    cnt++
    k--
  }
  for m > 0 {
    a[cnt] = mx + 1
    mx = a[cnt]
    cnt++
    m--
  }
  for i := 0; i < n; i++ {
    if a[i] > 50000 {
      fmt.Println(-1)
      return
    }
  }
  for i := 0; i < n; i++ {
    v := a[i]
    if v == 0 {
      v = 1
    }
    if i+1 < n {
      fmt.Print(v, " ")
    } else {
      fmt.Print(v)
    }
  }
  fmt.Println()
}
