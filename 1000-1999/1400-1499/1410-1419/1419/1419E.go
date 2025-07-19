package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   fmt.Fscan(in, &t)
   for tt := 0; tt < t; tt++ {
       var n int64
       fmt.Fscan(in, &n)
       // factorize n
       p := make([]int64, 0)
       e := make([]int, 0)
       tmp := n
       for i := int64(2); i*i <= tmp; i++ {
           if tmp%i == 0 {
               p = append(p, i)
               e = append(e, 0)
               for tmp%i == 0 {
                   tmp /= i
                   e[len(e)-1]++
               }
           }
       }
       if tmp > 1 {
           p = append(p, tmp)
           e = append(e, 1)
       }
       k := len(p)
       // divisors grouped
       d := make([][]int64, k)
       if k >= 2 {
           prod := p[0] * p[k-1]
           d[k-1] = append(d[k-1], prod)
       }
       // dfs to generate divisors
       var dfs func(cur int, x int64, mp int)
       dfs = func(cur int, x int64, mp int) {
           if cur == k {
               if mp == -1 {
                   return
               }
               if k >= 2 && x == p[0]*p[k-1] {
                   return
               }
               d[mp] = append(d[mp], x)
               return
           }
           // skip prime[cur]
           dfs(cur+1, x, mp)
           // include prime[cur] one or more times
           xi := x
           for i := 1; i <= e[cur]; i++ {
               xi *= p[cur]
               if mp == -1 {
                   dfs(cur+1, xi, cur)
               } else {
                   dfs(cur+1, xi, mp)
               }
           }
       }
       dfs(0, 1, -1)
       // special answer if two primes both exponent 1
       ans := 0
       if k == 2 && e[0] == 1 && e[1] == 1 {
           ans = 1
       }
       // output divisors in order
       for i := 0; i < k; i++ {
           next := p[(i+1)%k]
           // move any divisible by next to end
           for j := 0; j < len(d[i]); j++ {
               if d[i][j]%next == 0 {
                   last := len(d[i]) - 1
                   d[i][j], d[i][last] = d[i][last], d[i][j]
               }
           }
           for _, x := range d[i] {
               fmt.Fprint(out, x, " ")
           }
       }
       fmt.Fprintln(out)
       fmt.Fprintln(out, ans)
   }
}
