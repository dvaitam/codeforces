package main

import (
   "bufio"
   "fmt"
   "os"
)

const MAXN = 100000

var d [MAXN + 1]int

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   // precompute divisor counts
   for i := 1; i <= MAXN; i++ {
       for j := i; j <= MAXN; j += i {
           d[j]++
       }
   }
   r := bufio.NewReader(os.Stdin)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()

   var t int
   if _, err := fmt.Fscan(r, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var a, b, c int
       fmt.Fscan(r, &a, &b, &c)
       // compute gcds
       x := gcd(gcd(a, b), c)
       m := gcd(a, b)
       n := gcd(b, c)
       p := gcd(c, a)
       // use int64 for results
       var ans int64
       // all combinations
       ans += int64(d[a]) * int64(d[b]) * int64(d[c])
       // subtract one divisor divides only one of others
       tmp := int64(d[a] - d[m] - d[p] + d[x])
       ans -= tmp * int64(d[n]) * int64(d[n]-1) / 2
       tmp = int64(d[b] - d[m] - d[n] + d[x])
       ans -= tmp * int64(d[p]) * int64(d[p]-1) / 2
       tmp = int64(d[c] - d[n] - d[p] + d[x])
       ans -= tmp * int64(d[m]) * int64(d[m]-1) / 2
       // all three divide all
       ans -= 5 * int64(d[x]) * int64(d[x]-1) * int64(d[x]-2) / 6
       ans -= 2 * int64(d[x]) * int64(d[x]-1)
       // other combinations
       tmp = int64(d[m] - d[x])
       ans -= 3 * tmp * int64(d[x]) * int64(d[x]-1) / 2
       tmp = int64(d[n] - d[x])
       ans -= 3 * tmp * int64(d[x]) * int64(d[x]-1) / 2
       tmp = int64(d[p] - d[x])
       ans -= 3 * tmp * int64(d[x]) * int64(d[x]-1) / 2
       tmp = int64(d[m] - d[x])
       ans -= tmp * int64(d[x])
       tmp = int64(d[n] - d[x])
       ans -= tmp * int64(d[x])
       tmp = int64(d[p] - d[x])
       ans -= tmp * int64(d[x])
       tmp = int64(d[m] - d[x])
       ans -= int64(d[x]) * tmp * (tmp - 1) / 2
       tmp = int64(d[n] - d[x])
       ans -= int64(d[x]) * tmp * (tmp - 1) / 2
       tmp = int64(d[p] - d[x])
       ans -= int64(d[x]) * tmp * (tmp - 1) / 2
       ans -= 2 * int64(d[x]) * int64(d[m]-d[x]) * int64(d[n]-d[x])
       ans -= 2 * int64(d[x]) * int64(d[m]-d[x]) * int64(d[p]-d[x])
       ans -= 2 * int64(d[x]) * int64(d[n]-d[x]) * int64(d[p]-d[x])
       tmp = int64(d[n] - d[x] + d[p] - d[x])
       ans -= tmp * int64(d[m]-d[x]) * int64(d[m]-d[x]-1) / 2
       tmp = int64(d[m] - d[x] + d[n] - d[x])
       ans -= tmp * int64(d[p]-d[x]) * int64(d[p]-d[x]-1) / 2
       tmp = int64(d[m] - d[x] + d[p] - d[x])
       ans -= tmp * int64(d[n]-d[x]) * int64(d[n]-d[x]-1) / 2
       ans -= int64(d[m]-d[x]) * int64(d[n]-d[x]) * int64(d[p]-d[x])
       // output
       fmt.Fprintln(w, ans)
   }
}
