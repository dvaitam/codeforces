package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var month string
	if _, err := fmt.Fscan(reader, &month); err != nil {
		return
	}
	var season string
	switch month {
	case "December", "January", "February":
		season = "winter"
	case "March", "April", "May":
		season = "spring"
	case "June", "July", "August":
		season = "summer"
	case "September", "October", "November":
		season = "autumn"
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, season)
}
