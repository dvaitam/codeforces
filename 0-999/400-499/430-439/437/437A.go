package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // labels in order
   labels := []string{"A", "B", "C", "D"}
   lengths := make([]int, 4)
   for i := 0; i < 4; i++ {
       line, err := reader.ReadString('\n')
       if err != nil && len(line) == 0 {
           // no more input
           break
       }
       // remove trailing newline/carriage return
       if len(line) > 0 && (line[len(line)-1] == '\n' || line[len(line)-1] == '\r') {
           line = line[:len(line)-1]
       }
       // description starts after "X."
       if len(line) >= 2 {
           desc := line[2:]
           lengths[i] = len(desc)
       } else {
           lengths[i] = 0
       }
   }
   greatCount := 0
   greatIdx := -1
   for i := 0; i < 4; i++ {
       shorter := true
       longer := true
       for j := 0; j < 4; j++ {
           if i == j {
               continue
           }
           if lengths[i]*2 > lengths[j] {
               shorter = false
           }
           if lengths[i] < 2*lengths[j] {
               longer = false
           }
       }
       if shorter || longer {
           greatCount++
           greatIdx = i
       }
   }
   if greatCount == 1 {
       fmt.Println(labels[greatIdx])
   } else {
       fmt.Println("C")
   }
}
