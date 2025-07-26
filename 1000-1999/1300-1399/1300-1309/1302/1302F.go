package main

import (
   "fmt"
)

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   if len(s) != 5 {
       return
   }
   d := make([]int, 5)
   for i := 0; i < 5; i++ {
       d[i] = int(s[i] - '0')
   }
   rotate := func(pos, times int) {
       d[pos] = (d[pos] + times) % 10
   }

   if d[0]+d[3] > 10 {
       rotate(0, 3)
   } else {
       rotate(3, 8)
   }
   if d[2]+d[1] > 8 {
       rotate(3, 9)
   } else {
       rotate(4, 8)
   }
   if d[2]%2 == 1 {
       rotate(2, 3)
   } else {
       rotate(2, 4)
   }
   if d[4] > d[1] {
       rotate(3, 1)
   } else {
       rotate(1, 7)
   }
   if d[0]%2 == 1 {
       rotate(0, 3)
   } else {
       rotate(2, 5)
   }
   if d[3]%2 == 1 {
       rotate(3, 7)
   } else {
       rotate(0, 9)
   }
   if d[3] > d[0] {
       rotate(3, 9)
   } else {
       rotate(3, 2)
   }
   if d[0] > d[2] {
       rotate(1, 1)
   } else {
       rotate(2, 1)
   }
   if d[4] > d[2] {
       rotate(3, 5)
   } else {
       rotate(4, 8)
   }
   if d[0]+d[2] > 8 {
       rotate(3, 5)
   } else {
       rotate(1, 5)
   }
   if d[0] > d[3] {
       rotate(3, 3)
   } else {
       rotate(1, 3)
   }
   if d[2]+d[0] > 9 {
       rotate(1, 9)
   } else {
       rotate(1, 2)
   }
   if d[3]+d[2] > 10 {
       rotate(3, 7)
   } else {
       rotate(4, 7)
   }
   if d[2] > d[1] {
       rotate(2, 2)
   } else {
       rotate(3, 6)
   }
   if d[0] > d[2] {
       rotate(0, 9)
   } else {
       rotate(1, 9)
   }
   if d[2]%2 == 1 {
       rotate(2, 9)
   } else {
       rotate(0, 5)
   }
   if d[2]+d[4] > 9 {
       rotate(2, 4)
   } else {
       rotate(2, 9)
   }
   if d[2] > d[0] {
       rotate(4, 1)
   } else {
       rotate(4, 7)
   }
   if d[0] > d[2] {
       rotate(1, 9)
   } else {
       rotate(3, 6)
   }
   if d[1]+d[2] > 10 {
       rotate(1, 2)
   } else {
       rotate(2, 6)
   }

   for i := 0; i < 5; i++ {
       fmt.Print(d[i])
   }
   fmt.Println()
}
