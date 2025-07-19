package main

import "fmt"

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   n := len(s)
   cnt := make([]int, n+1)
   for i := 1; i <= n; i++ {
       cnt[i] = cnt[i-1]
       if s[i-1] == '[' {
           cnt[i]++
       }
   }
   // stack of positions (1-based), initialize with 0
   stack := []int{0}
   L, R, ans := 0, 0, 0
   for i := 1; i <= n; i++ {
       c := s[i-1]
       if c == '[' || c == '(' {
           stack = append(stack, i)
       } else {
           // determine if mismatch
           var topChar byte
           if len(stack) > 0 {
               pos := stack[len(stack)-1]
               if pos > 0 && pos <= n {
                   topChar = s[pos-1]
               }
           }
           if (c == ']' && topChar != '[') || (c == ')' && topChar != '(') {
               // reset on mismatch
               stack = []int{i}
               continue
           }
           // matching closing bracket: pop
           if len(stack) > 0 {
               stack = stack[:len(stack)-1]
           }
           // check valid sequence
           if len(stack) > 0 {
               l := stack[len(stack)-1]
               r := i
               if ans < cnt[r]-cnt[l] {
                   ans = cnt[r] - cnt[l]
                   L = l
                   R = r
               }
           }
       }
   }
   // output
   fmt.Println(ans)
   // substring from positions L+1..R in 1-based is s[L:R]
   if L < R && R <= n {
       fmt.Println(s[L:R])
   } else {
       fmt.Println()
   }
}
