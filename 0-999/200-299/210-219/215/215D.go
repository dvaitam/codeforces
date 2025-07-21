package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m int64
   fmt.Fscan(reader, &n, &m)
   var ans int64
   for i := int64(0); i < n; i++ {
       var t, T, x, cost int64
       fmt.Fscan(reader, &t, &T, &x, &cost)
       K := T - t
       switch {
       case K <= 0:
           // Always too hot: compensation for all children in one bus
           ans += cost + m*x
       case K >= m:
           // Never too hot in one bus
           ans += cost
       default:
           // Three strategies:
           // A: one bus, pay compensation for all
           costA := cost + m*x
           // B: minimal buses to eliminate full compensation but may have some partial
           b2 := m/(K+1) + 1
           k1 := m / b2
           r := m - k1*b2
           costB := b2*cost + r*(k1+1)*x
           // C: buses to avoid any compensation
           b3 := (m + K - 1) / K
           costC := b3 * cost
           // pick minimum
           best := costA
           if costB < best {
               best = costB
           }
           if costC < best {
               best = costC
           }
           ans += best
       }
   }
   fmt.Fprintln(writer, ans)
