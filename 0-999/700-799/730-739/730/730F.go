package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, b int
   fmt.Fscan(reader, &n, &b)
   a := make([]int, n+1)
   p := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
       p[i+1] = a[i] / 10
   }
   p[1] = b

   ans := make([]int, n+1)
   var f []int
   var lck [11][]int

   // spendFre consumes from f
   var spendFre func(tot *int)
   spendFre = func(tot *int) {
       for *tot > 0 && len(f) > 0 {
           idx := f[len(f)-1]
           f = f[:len(f)-1]
           ans[idx]++
           *tot--
       }
   }
   // spend consumes using f and lck
   var spend func(tot *int)
   spend = func(tot *int) {
       spendFre(tot)
       for i := 10; i >= 1; i-- {
           for *tot > 0 && len(lck[i]) > 0 {
               *tot--
               idx := lck[i][len(lck[i])-1]
               lck[i] = lck[i][:len(lck[i])-1]
               // insert i copies of idx into f
               for k := 0; k < i; k++ {
                   f = append(f, idx)
               }
               spendFre(tot)
           }
       }
   }

   // main processing
   for i := n; i >= 1; i-- {
       lckFre := p[i+1]
       maxPoints := a[i] / 2
       fre := min(a[i]%10, maxPoints)
       // insert fre copies of i
       for j := 0; j < fre; j++ {
           f = append(f, i)
       }
       maxPoints -= fre
       // distribute groups of 10
       cnt10 := maxPoints / 10
       for j := 0; j < cnt10; j++ {
           if lckFre > 0 {
               lckFre--
               for k := 0; k < 10; k++ {
                   f = append(f, i)
               }
           } else {
               lck[10] = append(lck[10], i)
           }
       }
       rem := maxPoints % 10
       if rem > 0 {
           if lckFre > 0 {
               lckFre--
               for k := 0; k < rem; k++ {
                   f = append(f, i)
               }
           } else {
               lck[rem] = append(lck[rem], i)
           }
       }
       // spend p[i]
       spend(&p[i])
   }

   total := 0
   for i := 1; i <= n; i++ {
       total += a[i] - ans[i]
   }
   fmt.Fprintln(writer, total)
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, ans[i])
   }
   fmt.Fprintln(writer)
}
