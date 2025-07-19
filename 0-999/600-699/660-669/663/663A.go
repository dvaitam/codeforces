package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       // ignore error
   }
   tokens := strings.Fields(line)
   // Parse signs: first term assumed '+'
   sgns := []string{"+"}
   i := 0
   for ; i < len(tokens); i++ {
       t := tokens[i]
       if t == "=" {
           break
       }
       if t == "+" || t == "-" {
           sgns = append(sgns, t)
       }
   }
   // Next token is the target n
   if i+1 >= len(tokens) {
       return
   }
   n, err := strconv.Atoi(tokens[i+1])
   if err != nil {
       return
   }
   // Initial sum of ones with signs
   st := 1
   for j := 1; j < len(sgns); j++ {
       if sgns[j] == "+" {
           st++
       } else {
           st--
       }
   }
   // Initialize all variables to 1
   sol := make([]int, len(sgns))
   for j := range sol {
       sol[j] = 1
   }
   // Greedily adjust to reach n
   for j := 0; j < len(sgns); j++ {
       if st > n && sgns[j] == "-" {
           if st-n > n-1 {
               sol[j] = n
               st -= n - 1
           } else {
               sol[j] += st - n
               st = n
           }
       }
       if st < n && sgns[j] == "+" {
           if n-st > n-1 {
               sol[j] = n
               st += n - 1
           } else {
               sol[j] += n - st
               st = n
           }
       }
   }
   // Output result
   if st == n {
       fmt.Println("Possible")
       out := fmt.Sprintf("%d", sol[0])
       for j := 1; j < len(sgns); j++ {
           out += fmt.Sprintf(" %s %d", sgns[j], sol[j])
       }
       out += fmt.Sprintf(" = %d", n)
       fmt.Println(out)
   } else {
       fmt.Println("Impossible")
   }
}
