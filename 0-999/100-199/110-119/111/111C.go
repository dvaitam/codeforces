package main

import "fmt"

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   if n > m {
       n, m = m, n
   }
   var res int
   switch n {
   case 1, 2:
       res = m*n - (m+2)/(4-n)
   case 3:
       m0 := m
       res = m0*n - (m0/4)*3
       r := m0 % 4
       res -= r
       if r == 0 {
           res--
       }
   case 4:
       res = m*n - m
       if m == 5 || m == 6 || m == 9 {
           res--
       }
   case 5:
       m0 := m
       res = m0*n - (m0/5)*6
       if m0 == 7 {
           res++
       }
       r := m0 % 5
       res -= r
       res--
       if r > 1 {
           res--
       }
   case 6:
       res = m*n - 10
   default:
       res = n * m
   }
   fmt.Println(res)
}
