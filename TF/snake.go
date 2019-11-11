package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Pair struct {
	a, b interface{}
}

const cycle int = 1

const size int = 20

var field [size][size]string

const players int = 2

func snakeInput(id int, up string, down string, left string, right string, inputs <-chan string) {

}

func snakeController(id int) {
	head := Pair{5, 5}
	field[head.a.(int)][head.b.(int)] = strconv.Itoa(id)
	//var body [5]Pair
	direction := 3
	for {
		switch direction {
		case 3:
			head.a = head.a.(int) - 1
			field[head.a.(int)][head.b.(int)] = strconv.Itoa(id)
			//for index := 0; index < len(body); index++ {
			//	if body[index].a != 0 {
			//		body[index].a = body[index].a.(int) - 1
			//		field[body[index].a.(int)][body[index].b.(int)] = strconv.Itoa(id)
			//	}
			//}

		}

	}
}

func drawGame() {
	for {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
		for rows := 0; rows < size; rows++ {
			for columns := 0; columns < size; columns++ {
				if field[rows][columns] == "" {
					field[rows][columns] = " "
				}

				fmt.Print(field[rows][columns])
			}
			fmt.Print("\n")
		}
		time.Sleep(time.Second * time.Duration(cycle))
	}
}

func sendInput(input string, inputs chan<- string) {
	for index := 0; index < players; index++ {
		inputs <- input
	}
}
func main() {
	go drawGame()
	inputs := make(chan string)
	go snakeInput(0, "w", "a", "s", "d", inputs)
	go snakeController(0)
	for {

	}
}
