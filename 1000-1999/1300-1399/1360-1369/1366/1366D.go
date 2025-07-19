package main

import (
   "bufio"
   "fmt"
   "os"
)

const M = 10000005

var pre []int32
var primes []int32

// initSieve computes smallest prime factor for each number up to M-1
func initSieve() {
   pre = make([]int32, M)
   primes = make([]int32, 0, 700000)
   for i := 2; i < M; i++ {
       if pre[i] == 0 {
           pre[i] = int32(i)
           primes = append(primes, int32(i))
       }
       for _, p := range primes {
           t := int32(i) * p
           if t >= M {
               break
           }
           pre[t] = p
           if p == pre[i] {
               break
           }
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   initSieve()
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x%2 == 0 {
           // remove all factors of 2
           for x%2 == 0 {
               x /= 2
           }
           if x == 1 {
               a[i], b[i] = -1, -1
           } else {
               a[i], b[i] = 2, x
           }
       } else {
           var pa, pb int32 = -1, -1
           xx := int32(x)
           for xx > 1 {
               t := pre[int(xx)]
               if pa == -1 {
                   pa = t
               } else if pb == -1 {
                   pb = t
               }
               // remove all occurrences of t
               for xx%t == 0 {
                   xx /= t
               }
               if pa != -1 && pb != -1 {
                   break
               }
           }
           if pa == -1 || pb == -1 {
               a[i], b[i] = -1, -1
           } else {
               a[i], b[i] = int(pa), int(pb)
           }
       }
   }
   // output
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprint(a[i]))
   }
   writer.WriteByte('\n')
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprint(b[i]))
   }
   writer.WriteByte('\n')
}
