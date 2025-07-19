package main

import (
   "bufio"
   "fmt"
   "os"
)

const M = 1000000007
const maxA = 1000000

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   ff := make([]int, maxA+1)
   for i := 2; i <= maxA; i++ {
       if ff[i] == 0 {
           for j := i; j <= maxA; j += i {
               if ff[j] == 0 {
                   ff[j] = i
               }
           }
       }
   }

   g := make([]int, maxA+1)
   g[0] = 1
   for _, val := range a {
       // generate divisors of val
       divisors := []int{1}
       tmp := val
       for tmp > 1 {
           p := ff[tmp]
           cnt := 0
           for tmp%p == 0 {
               tmp /= p
               cnt++
           }
           oldSize := len(divisors)
           mult := 1
           for k := 0; k < cnt; k++ {
               mult *= p
               for i := 0; i < oldSize; i++ {
                   divisors = append(divisors, divisors[i]*mult)
               }
           }
       }
       // record previous g[d-1]
       values := make([]int, len(divisors))
       for i, d := range divisors {
           values[i] = g[d-1]
       }
       // update g
       for i, d := range divisors {
           g[d] += values[i]
           if g[d] >= M {
               g[d] -= M
           }
       }
   }

   var ans int64
   for i := 1; i <= maxA; i++ {
       ans += int64(g[i])
   }
   ans %= M
   fmt.Fprintln(writer, ans)
}
