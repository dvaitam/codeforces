package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n)
   fmt.Fscan(reader, &m)
   a := make([]int64, n)
   b := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i], &b[i])
   }

   type pair struct {
       diff int64
       idx  int
   }
   tab := make([]pair, n)
   for i := 0; i < n; i++ {
       tab[i] = pair{b[i] - a[i], i}
   }
   sort.Slice(tab, func(i, j int) bool {
       return tab[i].diff < tab[j].diff
   })

   pref := make([]int64, n+1)
   for i := 0; i < n; i++ {
       pref[i+1] = pref[i] + b[tab[i].idx]
   }
   suf := make([]int64, n+1)
   for i := n - 1; i >= 0; i-- {
       suf[i] = suf[i+1] + a[tab[i].idx]
   }

   wynik := make([]int64, n)
   for i := 0; i < n; i++ {
       idx := tab[i].idx
       wynik[idx] = b[idx]*int64(n-1-i) + a[idx]*int64(i)
       wynik[idx] += pref[i]
       wynik[idx] += suf[i+1]
   }

   for k := 0; k < m; k++ {
       var c, d int
       fmt.Fscan(reader, &c, &d)
       c--
       d--
       delta := a[c] + b[d]
       if a[d]+b[c] < delta {
           delta = a[d] + b[c]
       }
       wynik[c] -= delta
       wynik[d] -= delta
   }

   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(strconv.FormatInt(wynik[i], 10))
   }
   writer.WriteByte('\n')
}
