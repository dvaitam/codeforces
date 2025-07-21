package main

import (
   "bufio"
   "fmt"
   "os"
   "unicode"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   weights := map[rune]int{
       'Q': 9,
       'R': 5,
       'B': 3,
       'N': 3,
       'P': 1,
   }
   white, black := 0, 0
   for i := 0; i < 8; i++ {
       if !scanner.Scan() {
           break
       }
       line := scanner.Text()
       for _, r := range line {
           if r == '.' {
               continue
           }
           w := weights[unicode.ToUpper(r)]
           if unicode.IsUpper(r) {
               white += w
           } else {
               black += w
           }
       }
   }
   if white > black {
       fmt.Println("White")
   } else if black > white {
       fmt.Println("Black")
   } else {
       fmt.Println("Draw")
   }
}
