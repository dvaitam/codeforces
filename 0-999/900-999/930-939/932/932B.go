package main

import (
   "bufio"
   "fmt"
   "os"
)

const maxN = 1000000

func digitProduct(n int) int {
   prod := 1
   for n > 0 {
      d := n % 10
      if d > 0 {
         prod *= d
      }
      n /= 10
   }
   return prod
}

func main() {
   g := make([]int, maxN+1)
   for i := 1; i <= maxN; i++ {
      if i < 10 {
         g[i] = i
      } else {
         p := digitProduct(i)
         g[i] = g[p]
      }
   }
   pref := make([][]int, 10)
   for i := range pref {
      pref[i] = make([]int, maxN+1)
   }
   for i := 1; i <= maxN; i++ {
      for k := 1; k <= 9; k++ {
         pref[k][i] = pref[k][i-1]
      }
      val := g[i]
      pref[val][i]++
   }

   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var q int
   fmt.Fscan(reader, &q)
   for ; q > 0; q-- {
      var l, r, k int
      fmt.Fscan(reader, &l, &r, &k)
      if k < 1 || k > 9 {
         fmt.Fprintln(writer, 0)
      } else {
         fmt.Fprintln(writer, pref[k][r]-pref[k][l-1])
      }
   }
}
