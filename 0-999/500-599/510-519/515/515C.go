package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   var s string
   fmt.Fscan(reader, &s)

   // Mapping each digit to its replacement digits
   mapping := map[rune]string{
       '0': "",
       '1': "",
       '2': "2",
       '3': "3",
       '4': "322",
       '5': "5",
       '6': "53",
       '7': "7",
       '8': "7222",
       '9': "7332",
   }

   var result []byte
   for _, ch := range s {
       if rep, ok := mapping[ch]; ok {
           for i := 0; i < len(rep); i++ {
               result = append(result, rep[i])
           }
       }
   }

   // Sort in descending order to form the largest possible number
   sort.Slice(result, func(i, j int) bool {
       return result[i] > result[j]
   })

   // Output result
   if len(result) == 0 {
       fmt.Println()
   } else {
       fmt.Println(string(result))
   }
}
