package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var s string
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   if n%2 == 1 {
       fmt.Println(0)
       return
   }
   a := make([]int, n)
   num := 0
   for i := 0; i < n; i++ {
       if s[i] == '(' {
           num++
       } else {
           num--
       }
       a[i] = num
   }
   switch a[n-1] {
   case 0:
       fmt.Println(0)
   case -2:
       ok := true
       for _, v := range a {
           if v < -2 {
               ok = false
               break
           }
       }
       if !ok {
           fmt.Println(0)
           return
       }
       right := n - 1
       for i, v := range a {
           if v < 0 {
               right = i
               break
           }
       }
       cnt := 0
       for i := 0; i <= right; i++ {
           if s[i] == ')' {
               cnt++
           }
       }
       fmt.Println(cnt)
   case 2:
       ok := true
       for _, v := range a {
           if v < 0 {
               ok = false
               break
           }
       }
       if !ok {
           fmt.Println(0)
           return
       }
       left := 0
       for i := n - 1; i >= 0; i-- {
           if a[i] < 2 {
               left = i + 1
               break
           }
       }
       cnt := 0
       for i := left; i < n; i++ {
           if s[i] == '(' {
               cnt++
           }
       }
       fmt.Println(cnt)
   default:
       fmt.Println(0)
   }
}
