package main

import (
   "bufio"
   "fmt"
   "os"
)

// Jar represents a honey jar with remaining honey and times Pooh has eaten from it
type Jar struct {
   rem    int
   eaten  int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   jars := make([]Jar, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &jars[i].rem)
       jars[i].eaten = 0
   }
   totalPiglet := 0
   // simulate until all jars given to Piglet
   for len(jars) > 0 {
       // find jar with max remaining honey
       maxIdx := 0
       for i := 1; i < len(jars); i++ {
           if jars[i].rem > jars[maxIdx].rem {
               maxIdx = i
           }
       }
       jar := &jars[maxIdx]
       // if less than k or eaten 3 times, give to Piglet
       if jar.rem < k || jar.eaten >= 3 {
           totalPiglet += jar.rem
           // remove jar
           jars = append(jars[:maxIdx], jars[maxIdx+1:]...)
       } else {
           // Pooh eats k kilos
           jar.rem -= k
           jar.eaten++
       }
   }
   fmt.Println(totalPiglet)
}
