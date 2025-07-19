package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

var (
   n    int
   p    []int64
   qexp []int
   v    int64
   MA   int64
   M    int64
   a, b, c int64
   A, B, C int64
   ans  int64
)

func dfs(s int, now int64) {
   if now > MA {
       return
   }
   if s == n {
       rem := v / now
       // bound = rem + 2*a*sqrt(rem) with a=now
       bound := float64(rem) + 2*float64(now)*math.Sqrt(float64(rem))
       if bound < float64(ans) {
           a = now
           M = int64(math.Sqrt(float64(rem)) + 1e-8)
           cal(0, 1)
       }
       return
   }
   if qexp[s] > 0 {
       qexp[s]--
       dfs(s, now*p[s])
       qexp[s]++
   }
   dfs(s+1, now)
}

func cal(s int, now int64) {
   if now > M {
       return
   }
   if s == n {
       if now < a {
           return
       }
       b = now
       c = v / a / b
       sum := a*b + a*c + b*c
       if sum < ans {
           ans = sum
           A = a
           B = b
           C = c
       }
       return
   }
   if qexp[s] > 0 {
       qexp[s]--
       cal(s, now*p[s])
       qexp[s]++
   }
   cal(s+1, now)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var T int
   fmt.Fscan(in, &T)
   for T > 0 {
       T--
       fmt.Fscan(in, &n)
       p = make([]int64, n)
       qexp = make([]int, n)
       v = 1
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &p[i], &qexp[i])
           for j := 0; j < qexp[i]; j++ {
               v *= p[i]
           }
       }
       ans = math.MaxInt64
       // limit for a: cube root of v
       MA = int64(math.Pow(float64(v), 1.0/3.0) + 1e-8)
       dfs(0, 1)
       fmt.Fprintf(out, "%d %d %d %d\n", 2*ans, A, B, C)
   }
}
