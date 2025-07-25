package main

import "fmt"

func bigAdd(a, b string) string {
   i, j := len(a)-1, len(b)-1
   carry := 0
   var res []byte
   for i >= 0 || j >= 0 || carry > 0 {
       sum := carry
       if i >= 0 {
           sum += int(a[i] - '0')
           i--
       }
       if j >= 0 {
           sum += int(b[j] - '0')
           j--
       }
       res = append(res, byte('0'+sum%10))
       carry = sum / 10
   }
   // reverse res
   for l, r := 0, len(res)-1; l < r; l, r = l+1, r-1 {
       res[l], res[r] = res[r], res[l]
   }
   return string(res)
}

// compare numeric strings, return true if a < b
func less(a, b string) bool {
   if len(a) != len(b) {
       return len(a) < len(b)
   }
   return a < b
}

func main() {
   var l int
   var s string
   fmt.Scan(&l, &s)
   mid := l / 2
   left := -1
   for i := mid; i >= 1; i-- {
       if s[i] != '0' {
           left = i
           break
       }
   }
   right := -1
   for i := mid + 1; i < l; i++ {
       if s[i] != '0' {
           right = i
           break
       }
   }
   candidates := []int{}
   if left != -1 {
       candidates = append(candidates, left)
   }
   if right != -1 {
       candidates = append(candidates, right)
   }
   var best string
   for idx, pos := range candidates {
       a := s[:pos]
       b := s[pos:]
       sum := bigAdd(a, b)
       if idx == 0 || less(sum, best) {
           best = sum
       }
   }
   fmt.Println(best)
}
