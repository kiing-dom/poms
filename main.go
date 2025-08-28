package main
import (
	"fmt"
	"time"
)

func main() {
	workDuration := 2 * time.Second
	breakDuration := 1 * time.Second

	fmt.Println("Starting", workDuration, "second", "pomodoro...")
	pomTimer := time.NewTimer(workDuration)

	workSeconds := int(workDuration.Seconds())
	for i := workSeconds; i > 0; i-- {
		fmt.Println(i, "seconds left...")
		time.Sleep(1 * time.Second)
	}

	<- pomTimer.C

	fmt.Println("Pomodoro ended! Good Job.")
	time.Sleep(2 * time.Second)

	fmt.Println("Starting ", breakDuration, "second ", "break.")
	breakTimer := time.NewTimer(breakDuration)

	breakSeconds := int(breakDuration.Seconds())
	for i := breakSeconds; i > 0; i-- {
		fmt.Println(i, "seconds left")
	}

	<- breakTimer.C
	fmt.Println("Break over. Pomodoro session complete")
	
}

