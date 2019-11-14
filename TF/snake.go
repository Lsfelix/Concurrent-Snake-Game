package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	term "github.com/nsf/termbox-go"
)

const SIZE int = 20

var field [SIZE][SIZE]string

const PLAYERS int = 2

func treatInput(inputs chan<- rune) {
	term.Init()
	defer term.Close()

	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {

			case term.KeyEsc:
				os.Exit(3)
			}

			//fmt.Println(ev.Ch)
			inputs <- ev.Ch

		case term.EventError:
			fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
			panic(ev.Err)
		}

	}

}

func snakeController(id int, inputs <-chan rune) {
	name := " " + strconv.Itoa(id) + " "

	headX := 5
	headY := 5

	field[headY][headX] = name
	//body := list.New()

	for {
		key := <-inputs
		switch key {
		case rune(119):

			field[headY][headX] = " _ "

			headY--

			field[headY][headX] = name

		case rune(97):

			field[headY][headX] = " _ "

			headX--

			field[headY][headX] = name

		case rune(115):

			field[headY][headX] = " _ "

			headY++

			field[headY][headX] = name

		case rune(100):

			field[headY][headX] = " _ "

			headX++

			field[headY][headX] = name

		}
		drawGame()
	}

}

func drawGame() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
	for rows := 0; rows < SIZE; rows++ {
		for columns := 0; columns < SIZE; columns++ {
			if field[rows][columns] == "" {
				field[rows][columns] = " _ "
			}

			fmt.Print(field[rows][columns])
		}
		fmt.Print("\n")
	}
}

func sendInput(input string, inputs chan<- string) {
	for index := 0; index < PLAYERS; index++ {
		inputs <- input
	}
}
func main() {
	inputs := make(chan rune, 100)
	go treatInput(inputs)

	drawGame()

	go snakeController(0, inputs)
	for {

	}
}
