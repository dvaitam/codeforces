package main

import (
   "fmt"
)

func main() {
   var k int
   if _, err := fmt.Scan(&k); err != nil {
       return
   }
   // columns: driver + 10 rows of 3 seats + last row of 4 seats
   cols := 12
   width := 1 + cols*2
   height := 6
   // prepare bus canvas
   bus := make([][]rune, height)
   // top and bottom borders
   border := make([]rune, width)
   border[0] = '+'
   for i := 1; i < width-1; i++ {
       border[i] = '-'
   }
   border[width-1] = '+'
   bus[0] = border
   bus[height-1] = border
   // function to return seat height (number of seats across width) for column c
   seatH := func(c int) int {
       switch c {
       case 0:
           return 1 // driver
       case cols - 1:
           return 4 // last row
       default:
           return 3 // regular rows
       }
   }
   // generate seat positions in boarding order: last row first
   var posList [][2]int
   for c := cols - 1; c >= 1; c-- {
       for p := 1; p <= seatH(c); p++ {
           posList = append(posList, [2]int{p, c})
       }
   }
   // mark taken seats
   taken := make(map[[2]int]bool)
   for i := 0; i < k && i < len(posList); i++ {
       taken[posList[i]] = true
   }
   // build seat rows
   for r := 1; r <= height-2; r++ {
       row := make([]rune, width)
       // fill default underscores and pipes
       for j := 0; j < width; j++ {
           row[j] = '_'
       }
       row[0] = '|'
       row[width-1] = '|'
       // each column
       for c := 0; c < cols; c++ {
           pos := 1 + c*2
           // separator
           row[pos+1] = '|'
           // seat or empty area
           if r <= seatH(c) {
               if c == 0 {
                   // driver
                   row[pos] = 'D'
               } else {
                   // passenger seat
                   key := [2]int{r, c}
                   if taken[key] {
                       row[pos] = 'O'
                   } else {
                       row[pos] = '#'
                   }
               }
           } else {
               // area below seats stays underscore
               row[pos] = '_'
           }
       }
       bus[r] = row
   }
   // print bus
   for i := 0; i < height; i++ {
       fmt.Println(string(bus[i]))
   }
}
