package main
import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting 2 second pomodoro...")
	pomTimer := time.NewTimer(2 * time.Second)

	<- pomTimer.C

	fmt.Println("Pomodoro ended! Good Job.")
	time.Sleep(2 * time.Second)
	fmt.Println("Starting 1 second break.")
	breakTimer := time.NewTimer(1 * time.Second)

	<- breakTimer.C
	fmt.Println("Break over. Pomodoro session complete")
	
}

