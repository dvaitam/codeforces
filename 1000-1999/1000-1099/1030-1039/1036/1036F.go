package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = int64(1) << 61

var (
   w   []pair
   asd [64][]int64
)

type pair struct{ x, t int }

func mul(a, b int64) int64 {
   if a == 0 || b == 0 {
       return 0
   }
   if INF/b <= a {
       return INF
   }
   prod := a * b
   if prod < INF {
       return prod
   }
   return INF
}

func pw(a int64, b int) int64 {
   var r int64 = 1
   for i := 0; i < b; i++ {
       r = mul(r, a)
   }
   return r
}

func do2(n int64) int64 {
   var s, e int64 = 0, 1e9 + 7
   for e-s > 1 {
       m := (s + e) / 2
       if m*m <= n {
           s = m
       } else {
           e = m
       }
   }
   return s
}

func do3(n int64) int64 {
   var s, e int64 = 0, 1e6 + 7
   for e-s > 1 {
       m := (s + e) / 2
       if m*m*m <= n {
           s = m
       } else {
           e = m
       }
   }
   return s
}

func get(n int64, x int) int64 {
   if x == 2 {
       return do2(n)
   }
   if x == 3 {
       return do3(n)
   }
   // count of asd[x] values <= n
   arr := asd[x]
   // upper bound
   lo, hi := 0, len(arr)
   for lo < hi {
       mid := (lo + hi) / 2
       if arr[mid] <= n {
           lo = mid + 1
       } else {
           hi = mid
       }
   }
   return int64(lo)
}

func f(n int64) int64 {
   r := n - 1
   for _, p := range w {
       z := get(n, p.x)
       if z > 0 {
           z--
       }
       if p.t == 1 {
           r -= z
       } else {
           r += z
       }
   }
   return r
}

func init() {
   // prepare w and asd
   for i := 2; i < 61; i++ {
       n := i
       t := 0
       val := true
       j := 2
       for n > 1 {
           if n%j == 0 {
               if n%(j*j) == 0 {
                   val = false
                   break
               }
               t ^= 1
               n /= j
           } else {
               j++
           }
       }
       if val {
           w = append(w, pair{i, t})
       }
       if i > 3 {
           // build powers
           for tt := int64(1); ; tt++ {
               v := pw(tt, i)
               if v >= INF {
                   break
               }
               asd[i] = append(asd[i], v)
           }
       }
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var tn int
   if _, err := fmt.Fscan(in, &tn); err != nil {
       return
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for ; tn > 0; tn-- {
       var n int64
       fmt.Fscan(in, &n)
       fmt.Fprintln(out, f(n))
   }
}
