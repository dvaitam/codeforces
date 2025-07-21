package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000000

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Precompute Fibonacci f[0..n]
   f := make([]int, n+1)
   if n >= 0 {
       f[0] = 1
   }
   if n >= 1 {
       f[1] = 1
   }
   for i := 2; i <= n; i++ {
       f[i] = f[i-1] + f[i-2]
       if f[i] >= mod {
           f[i] -= mod
       }
   }
   // Process operations
   for k := 0; k < m; k++ {
       var t int
       fmt.Fscan(reader, &t)
       switch t {
       case 1:
           var x, v int
           fmt.Fscan(reader, &x, &v)
           if x >= 1 && x <= n {
               a[x] = v
           }
       case 2:
           var l, r int
           fmt.Fscan(reader, &l, &r)
           sum := 0
           // sum f[x-l] * a[x]
           for x := l; x <= r; x++ {
               idx := x - l
               sum += f[idx] * a[x]
               if sum >= mod {
                   sum %= mod
               }
           }
           fmt.Fprintln(writer, sum % mod)
       default:
           // For E1, no type 3; ignore
           // If type 3 present, skip inputs
           var l, r, d int
           fmt.Fscan(reader, &l, &r, &d)
       }
   }
}
