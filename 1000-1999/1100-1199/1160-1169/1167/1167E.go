package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, x int
   fmt.Fscan(reader, &n, &x)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // positions for values 1..x
   INF := n + 5
   firstpos := make([]int, x+2)
   lastpos := make([]int, x+2)
   for v := 1; v <= x; v++ {
       firstpos[v] = INF
       lastpos[v] = 0
   }
   for i, v := range a {
       if firstpos[v] == INF {
           firstpos[v] = i + 1
       }
       lastpos[v] = i + 1
   }
   // sentinels
   lastpos[0] = 0
   firstpos[x+1] = INF

   // prefix good: values 1..v are in order
   prefixGood := make([]bool, x+2)
   prefixGood[0] = true
   for v := 1; v <= x; v++ {
       prefixGood[v] = prefixGood[v-1] && (lastpos[v-1] <= firstpos[v])
   }
   // suffix good: values v..x are in order
   goodSuffix := make([]bool, x+3)
   goodSuffix[x+1] = true
   for v := x; v >= 1; v-- {
       goodSuffix[v] = goodSuffix[v+1] && (lastpos[v] <= firstpos[v+1])
   }
   // find minimal p such that goodSuffix[p] == true
   p := 1
   for p <= x+1 && !goodSuffix[p] {
       p++
   }

   ans := int64(0)
   kptr := p
   // iterate l from 1..x
   for l := 1; l <= x; l++ {
       if !prefixGood[l-1] {
           break
       }
       A := lastpos[l-1]
       // r0plus1 = max(l+1, p)
       want := l + 1
       if p > want {
           want = p
       }
       if kptr < want {
           kptr = want
       }
       // find kptr minimal with firstpos[kptr] > A
       for kptr <= x+1 && firstpos[kptr] <= A {
           kptr++
       }
       // rmin = kptr - 1
       rmin := kptr - 1
       if rmin < l {
           rmin = l
       }
       // add count of r from rmin..x
       if rmin <= x {
           ans += int64(x - rmin + 1)
       }
   }
   fmt.Fprintln(writer, ans)
}
