package main

import (
   "fmt"
   "sort"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   ys := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&ys[i])
   }
   prev := 1000
   res := make([]int, n)
   for i, y := range ys {
       // generate candidates differing by at most one digit
       ystr := fmt.Sprintf("%04d", y)
       var cands []int
       // original
       if y >= 1000 && y <= 2011 {
           cands = append(cands, y)
       }
       // one-digit changes
       for j := 0; j < 4; j++ {
           for d := byte('0'); d <= '9'; d++ {
               if d == ystr[j] {
                   continue
               }
               if j == 0 && d == '0' {
                   continue
               }
               zbytes := []byte(ystr)
               zbytes[j] = d
               z := int(zbytes[0]-'0')*1000 + int(zbytes[1]-'0')*100 + int(zbytes[2]-'0')*10 + int(zbytes[3]-'0')
               if z < 1000 || z > 2011 {
                   continue
               }
               cands = append(cands, z)
           }
       }
       // unique and sort
       seen := make(map[int]bool, len(cands))
       uniq := make([]int, 0, len(cands))
       for _, z := range cands {
           if !seen[z] {
               seen[z] = true
               uniq = append(uniq, z)
           }
       }
       sort.Ints(uniq)
       // pick smallest >= prev
       pick := -1
       for _, z := range uniq {
           if z >= prev {
               pick = z
               break
           }
       }
       if pick < 0 {
           fmt.Println("No solution")
           return
       }
       res[i] = pick
       prev = pick
   }
   for _, z := range res {
       fmt.Println(z)
   }
}
