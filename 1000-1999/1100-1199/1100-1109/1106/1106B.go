package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Dish represents a menu item with available quantity a, cost c, and original id
type Dish struct {
   a  int64
   c  int64
   id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   dishes := make([]Dish, n)
   // read quantities
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &dishes[i].a)
   }
   // read costs and assign ids
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &dishes[i].c)
       dishes[i].id = i + 1
   }
   // sort by cost ascending, then by id
   sort.Slice(dishes, func(i, j int) bool {
       if dishes[i].c == dishes[j].c {
           return dishes[i].id < dishes[j].id
       }
       return dishes[i].c < dishes[j].c
   })
   // map original id to position in sorted slice
   pos := make([]int, n+1)
   for idx, dsh := range dishes {
       pos[dsh.id] = idx
   }
   // pointer for cheapest available
   p := 0
   // process queries
   for qi := 0; qi < m; qi++ {
       var t int
       var need int64
       fmt.Fscan(reader, &t, &need)
       idx := pos[t]
       var ans int64
       // take from requested dish first
       if dishes[idx].a >= need {
           dishes[idx].a -= need
           ans = need * dishes[idx].c
           fmt.Fprintln(writer, ans)
           continue
       }
       // take all remaining from this dish
       ans = dishes[idx].a * dishes[idx].c
       need -= dishes[idx].a
       dishes[idx].a = 0
       // take from cheapest dishes
       for p < n && need > 0 {
           if dishes[p].a == 0 {
               p++
               continue
           }
           if dishes[p].a >= need {
               ans += need * dishes[p].c
               dishes[p].a -= need
               need = 0
               break
           }
           ans += dishes[p].a * dishes[p].c
           need -= dishes[p].a
           dishes[p].a = 0
       }
       if need > 0 {
           // not enough quantity
           fmt.Fprintln(writer, 0)
       } else {
           fmt.Fprintln(writer, ans)
       }
   }
}
