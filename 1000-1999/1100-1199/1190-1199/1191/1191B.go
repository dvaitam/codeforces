package main

import (
   "fmt"
)

func main() {
   var a, b, c string
   if _, err := fmt.Scan(&a, &b, &c); err != nil {
       return
   }
   tiles := []string{a, b, c}
   type tile struct{ num int; suit byte }
   var t [3]tile
   for i, s := range tiles {
       t[i].num = int(s[0] - '0')
       t[i].suit = s[1]
   }
   // check already have mentsu among the three
   // koutsu: all identical
   if t[0] == t[1] && t[1] == t[2] {
       fmt.Println(0)
       return
   }
   // shuntsu: same suit, nums are x, x+1, x+2
   if t[0].suit == t[1].suit && t[1].suit == t[2].suit {
       nums := []int{t[0].num, t[1].num, t[2].num}
       // simple sort 3 elements
       for i := 0; i < 2; i++ {
           for j := i + 1; j < 3; j++ {
               if nums[i] > nums[j] {
                   nums[i], nums[j] = nums[j], nums[i]
               }
           }
       }
       if nums[1] == nums[0]+1 && nums[2] == nums[1]+1 {
           fmt.Println(0)
           return
       }
   }
   // check if one draw is enough
   // for any pair of tiles
   for i := 0; i < 3; i++ {
       for j := i + 1; j < 3; j++ {
           ti, tj := t[i], t[j]
           // same tile => need one more for koutsu
           if ti == tj {
               fmt.Println(1)
               return
           }
           // same suit for sequence
           if ti.suit == tj.suit {
               d := ti.num - tj.num
               if d < 0 {
                   d = -d
               }
               if d <= 2 {
                   fmt.Println(1)
                   return
               }
           }
       }
   }
   // otherwise need two tiles
   fmt.Println(2)
}
