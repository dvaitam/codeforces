package main

import (
   "fmt"
   "os"
   "bufio"
   "strings"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   t, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   s = strings.TrimSpace(s)
   t = strings.TrimSpace(t)
   sh := strings.Split(s, ":")
   th := strings.Split(t, ":")
   shh, _ := strconv.Atoi(sh[0])
   smm, _ := strconv.Atoi(sh[1])
   thh, _ := strconv.Atoi(th[0])
   tmm, _ := strconv.Atoi(th[1])
   // current minutes since midnight
   curr := shh*60 + smm
   sleep := thh*60 + tmm
   // compute bedtime minutes
   p := curr - sleep
   // wrap around 24h
   p = ((p % (24*60)) + (24 * 60)) % (24 * 60)
   ph := p / 60
   pm := p % 60
   fmt.Printf("%02d:%02d\n", ph, pm)
}
