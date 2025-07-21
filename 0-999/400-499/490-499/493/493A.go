package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var home, away string
   fmt.Fscan(reader, &home)
   fmt.Fscan(reader, &away)
   var n int
   fmt.Fscan(reader, &n)
   // Track yellow cards and red card status for each team
   yellowHome := make(map[int]int)
   yellowAway := make(map[int]int)
   redHome := make(map[int]bool)
   redAway := make(map[int]bool)
   for i := 0; i < n; i++ {
       var t int
       var teamType string
       var number int
       var cardType string
       fmt.Fscan(reader, &t, &teamType, &number, &cardType)
       if teamType == "h" {
           if redHome[number] {
               continue
           }
           if cardType == "r" {
               redHome[number] = true
               fmt.Println(home, number, t)
           } else if cardType == "y" {
               yellowHome[number]++
               if yellowHome[number] == 2 {
                   redHome[number] = true
                   fmt.Println(home, number, t)
               }
           }
       } else {
           if redAway[number] {
               continue
           }
           if cardType == "r" {
               redAway[number] = true
               fmt.Println(away, number, t)
           } else if cardType == "y" {
               yellowAway[number]++
               if yellowAway[number] == 2 {
                   redAway[number] = true
                   fmt.Println(away, number, t)
               }
           }
       }
   }
}
