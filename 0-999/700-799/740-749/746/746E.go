package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Card struct {
   v   int
   pos int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   if m > 200000 {
       m = 200000
   }
   c := make([]Card, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &c[i].v)
       c[i].pos = i
   }
   sort.Slice(c, func(i, j int) bool { return c[i].v < c[j].v })

   d := make([]bool, m+1)
   odd := make([]int, 0, n)
   even := make([]int, 0, n)
   del := make([]int, 0, n)

   // mark used and classify, collect duplicates
   for i := 0; i < n; i++ {
       if c[i].v <= m && c[i].v > 0 {
           d[c[i].v] = true
       }
       if c[i].v%2 != 0 {
           odd = append(odd, i)
       } else {
           even = append(even, i)
       }
       // duplicates
       j := i + 1
       for j < n && c[j].v == c[i].v {
           c[j].v = -1
           del = append(del, j)
           j++
       }
       i = j - 1
   }

   h := n / 2
   // remove excess
   for len(odd) > h {
       del = append(del, odd[0])
       odd = odd[1:]
   }
   for len(even) > h {
       del = append(del, even[0])
       even = even[1:]
   }

   res := len(del)
   oddp, evenp := 1, 2
   // fill odd
   for len(odd) < h {
       for oddp <= m && d[oddp] {
           oddp += 2
       }
       if oddp > m {
           fmt.Fprintln(out, -1)
           return
       }
       d[oddp] = true
       idx := del[len(del)-1]
       del = del[:len(del)-1]
       c[idx].v = oddp
       odd = append(odd, idx)
   }
   // fill even
   for len(even) < h {
       for evenp <= m && d[evenp] {
           evenp += 2
       }
       if evenp > m {
           fmt.Fprintln(out, -1)
           return
       }
       d[evenp] = true
       idx := del[len(del)-1]
       del = del[:len(del)-1]
       c[idx].v = evenp
       even = append(even, idx)
   }

   // build ans
   ans := make([]int, n)
   for i := 0; i < n; i++ {
       ans[c[i].pos] = c[i].v
   }
   // output
   fmt.Fprintln(out, res)
   for i, v := range ans {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v)
   }
   out.WriteByte('\n')
}
