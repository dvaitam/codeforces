package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type para struct {
   a, b int
}

var staff []para

// find returns the para in staff representing the edge from x.b to the next city
func find(x para) para {
   l, r := 0, len(staff)-1
   var mid int
   for l <= r {
       mid = (l + r) / 2
       if staff[mid].a > x.b {
           r = mid - 1
       } else if staff[mid].a < x.b {
           l = mid + 1
       } else {
           if staff[mid].b != x.a {
               return staff[mid]
           }
           if mid+1 < len(staff) && staff[mid+1].a == x.b && staff[mid+1].b != x.a {
               return staff[mid+1]
           }
           return staff[mid-1]
       }
   }
   return staff[mid]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   staff = make([]para, 0, 2*n)
   for i := 0; i < n; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       staff = append(staff, para{a: u, b: v})
       staff = append(staff, para{a: v, b: u})
   }
   sort.Slice(staff, func(i, j int) bool {
       if staff[i].a != staff[j].a {
           return staff[i].a < staff[j].a
       }
       return staff[i].b < staff[j].b
   })
   // find starting city (odd degree)
   idx := 0
   for idx+1 < len(staff) && staff[idx].a == staff[idx+1].a {
       idx += 2
   }
   cur := staff[idx]
   // build path
   res := make([]int, 0, n+1)
   res = append(res, cur.a)
   for i := 0; i < n; i++ {
       if i == n-1 {
           res = append(res, cur.b)
       } else {
           cur = find(cur)
           res = append(res, cur.a)
       }
   }
   // output
   for i, v := range res {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
}
