package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var f, m, s string
   if _, err := fmt.Fscan(reader, &f, &m, &s); err != nil {
       return
   }
   // mapping of gesture to what it beats
   beats := map[string]string{
       "rock":     "scissors",
       "scissors": "paper",
       "paper":    "rock",
   }
   // determine odd gesture out
   gestures := []string{f, m, s}
   idx := -1
   common := ""
   // two equal, one different
   switch {
   case gestures[0] == gestures[1] && gestures[0] != gestures[2]:
       idx = 2
       common = gestures[0]
   case gestures[0] == gestures[2] && gestures[0] != gestures[1]:
       idx = 1
       common = gestures[0]
   case gestures[1] == gestures[2] && gestures[1] != gestures[0]:
       idx = 0
       common = gestures[1]
   default:
       fmt.Println("?")
       return
   }
   // check if odd gesture beats the common one
   if beats[gestures[idx]] == common {
       switch idx {
       case 0:
           fmt.Println("F")
       case 1:
           fmt.Println("M")
       case 2:
           fmt.Println("S")
       }
   } else {
       fmt.Println("?")
   }
}
