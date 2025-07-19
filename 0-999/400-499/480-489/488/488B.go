package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func chk(b []int) bool {
   dif := (b[3] - b[0]) * 4
   mid := (b[1] + b[2]) * 2
   if mid != dif {
       return false
   }
   sum := 0
   for i := 0; i < 4; i++ {
       sum += b[i]
   }
   if sum != dif {
       return false
   }
   return true
}

func verdict(a []int, b []int, x, y, z, r, cc int) bool {
   b[x] = a[0]
   b[y] = a[1]
   b[z] = a[2]
   var s, e int
   // cc==2 or cc==3
   if cc == 2 {
       // case cc==2
       var lim, lim2 int
       if y < z {
           s = b[y]
           e = 2000
           lim = z
           lim2 = 2000
       } else if x < z && z < y && y < r {
           s = b[x]
           e = b[y]
           lim = y
           lim2 = 2000
       } else if x < z && r < y {
           s = b[x]
           e = b[y]
           lim = z
           lim2 = y
       } else if z < x && x < r && r < y {
           s = 1
           e = b[x]
           lim = x
           lim2 = y
       } else {
           s = 1
           e = b[r]
           lim = z
           lim2 = x
       }
       for cur := s; cur <= e; cur++ {
           b[z] = cur
           tmp := (b[y] + b[z]) * 2
           b[r] = tmp - (b[x] + b[y] + b[z])
           if lim2 == 2000 {
               if b[r] >= b[lim] && chk(b) {
                   return true
               }
           } else {
               if b[r] >= b[lim] && b[r] <= b[lim2] && chk(b) {
                   return true
               }
           }
       }
   } else {
       // cc==3
       var lim, lim2 int
       if z < r {
           s = b[z]
           e = 2000
           lim = z
           lim2 = 2000
       } else if y < r && r < z {
           s = b[y]
           e = b[z]
           lim = y
           lim2 = z
       } else if x < r && r < y {
           s = b[x]
           e = b[y]
           lim = x
           lim2 = y
       } else {
           s = 1
           e = b[x]
           lim = 0 // unused
           lim2 = x
       }
       for cur := s; cur <= e; cur++ {
           b[r] = cur
           if lim == 0 {
               if b[r] <= b[lim2] && chk(b) {
                   return true
               }
           } else if lim2 == 2000 {
               if b[lim] <= b[r] && chk(b) {
                   return true
               }
           } else {
               if b[lim] <= b[r] && b[r] <= b[lim2] && chk(b) {
                   return true
               }
           }
       }
   }
   return false
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   b := make([]int, 4)
   hasAns := true
   switch n {
   case 0:
       fmt.Fprintln(out, "YES")
       fmt.Fprintln(out, 1)
       fmt.Fprintln(out, 1)
       fmt.Fprintln(out, 3)
       fmt.Fprintln(out, 3)
       return
   case 1:
       x := a[0]
       if x&1 == 1 {
           b[1] = x
           b[2] = x + 2
           b[3] = b[2] + ((x - 1) / 2)
           b[0] = b[1] - ((x - 1) / 2)
       } else {
           b[1], b[2] = x, x
           b[3] = b[2] + x/2
           b[0] = b[2] - x/2
       }
   case 2:
       sort.Ints(a)
       orders := [][5]int{{0, 1, 2, 3, 2}, {0, 2, 1, 3, 2}, {0, 3, 1, 2, 2}, {1, 2, 0, 3, 2}, {1, 3, 0, 2, 2}, {2, 3, 0, 1, 2}}
       hasAns = false
       for _, ord := range orders {
           if verdict(a, b, ord[0], ord[1], ord[2], ord[3], ord[4]) {
               hasAns = true
               break
           }
       }
   case 3:
       sort.Ints(a)
       orders := [][5]int{{0, 1, 2, 3, 3}, {0, 1, 3, 2, 3}, {0, 2, 3, 1, 3}, {1, 2, 3, 0, 3}}
       hasAns = false
       for _, ord := range orders {
           if verdict(a, b, ord[0], ord[1], ord[2], ord[3], ord[4]) {
               hasAns = true
               break
           }
       }
   case 4:
       sort.Ints(a)
       hasAns = chk(a)
       copy(b, a)
   }
   if hasAns {
       fmt.Fprintln(out, "YES")
       cnt := make(map[int]int)
       for i := 0; i < 4; i++ {
           cnt[b[i]]++
       }
       for i := 0; i < n; i++ {
           cnt[a[i]]--
       }
       keys := make([]int, 0, len(cnt))
       for k, v := range cnt {
           if v > 0 {
               keys = append(keys, k)
           }
       }
       sort.Ints(keys)
       for _, k := range keys {
           for i := 0; i < cnt[k]; i++ {
               fmt.Fprintln(out, k)
           }
       }
   } else {
       fmt.Fprintln(out, "NO")
   }
}
