package main

import (
	"bufio"
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	term "github.com/nsf/termbox-go"
)

type Pair struct {
	a, b interface{}
}

const SIZE int = 20

var field [SIZE][SIZE]string

var startGame sync.Mutex
var updateGame sync.Mutex

const PLAYERS int = 3

func treatInput(inputs chan<- rune) {
	term.Init()

	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {

			case term.KeyEsc:
				term.Close()
				os.Exit(3)
			}
			inputs <- ev.Ch

		case term.EventError:
			panic(ev.Err)
		}

	}

}

func spawnFruit() {
	r := rand.New(rand.NewSource(13))
	for {
		x := r.Intn(SIZE)
		y := r.Intn(SIZE)
		if field[y][x] != " _ " {
			continue
		}
		field[y][x] = " F "
		time.Sleep(3 * time.Second)
	}

}

func snakeController(id int, inputs chan rune, up rune, down rune, left rune, right rune) {

	name := " " + strconv.Itoa(id) + " "

	headX := (SIZE / (PLAYERS + 1)) * id
	headY := (SIZE / (PLAYERS + 1)) * id

	field[headY][headX] = name
	body := list.New()

	startGame.Lock()
	if id != PLAYERS {
		startGame.Unlock()
	}
	//fmt.Println("Cobra" + strconv.Itoa(id) + "Desbloquada")
	for {
		//update := false
		key := <-inputs

		oldHeadX := headX
		oldHeadY := headY

		switch key {
		case up:

			field[headY][headX] = " _ "

			headY--

		case left:

			field[headY][headX] = " _ "

			headX--

		case down:

			field[headY][headX] = " _ "

			headY++

		case right:

			field[headY][headX] = " _ "

			headX++

		default:
			inputs <- key
			//updateGame.Lock()
			continue
		}

		updateGame.Lock()

		if (field[headY][headX] != " _ " && field[headY][headX] != " F ") || headX < 0 || headX > SIZE || headY < 0 || headY > SIZE {
			for e := body.Front(); e != nil; e = e.Next() {
				field[e.Value.(Pair).a.(int)][e.Value.(Pair).b.(int)] = " _ "
			}
			updateGame.Unlock()
			break
		}

		if body.Len() > 0 {
			change := body.Back()
			field[change.Value.(Pair).a.(int)][change.Value.(Pair).b.(int)] = " _ "
			body.Remove(change)
			aux := Pair{oldHeadY, oldHeadX}
			body.PushFront(aux)
		}

		if field[headY][headX] == " F " {
			point := Pair{headY, headX}
			body.PushFront(point)
		}

		field[headY][headX] = name

		for e := body.Front(); e != nil; e = e.Next() {
			field[e.Value.(Pair).a.(int)][e.Value.(Pair).b.(int)] = name
		}

		drawGame()
	}

}

func drawGame() {
	cmd := exec.Command("clear") //Linux
	//cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
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
	updateGame.Unlock()
}

func sendInput(input string, inputs chan<- string) {
	for index := 0; index < PLAYERS; index++ {
		inputs <- input
	}
}

func main() {
	inputs := make(chan rune)

	startGame.Lock()

	reader := bufio.NewReader(os.Stdin)
	for index := 1; index < PLAYERS+1; index++ {
		fmt.Println("Jogador " + strconv.Itoa(index))
		fmt.Println("Digite a tecla que usará para ir para cima.")
		up, _, _ := reader.ReadRune()
		reader.Reset(os.Stdin)

		fmt.Println("Digite a tecla que usará para ir para a esquerda.")
		left, _, _ := reader.ReadRune()
		reader.Reset(os.Stdin)

		fmt.Println("Digite a tecla que usará para ir para baixo.")
		down, _, _ := reader.ReadRune()
		reader.Reset(os.Stdin)

		fmt.Println("Digite a tecla que usará para ir para a direita.")
		right, _, _ := reader.ReadRune()
		reader.Reset(os.Stdin)

		go snakeController(index, inputs, up, down, left, right)
	}

	go treatInput(inputs)

	updateGame.Lock()
	drawGame()
	//Alterar isso para ocorrer apenas depois que todas cobrinhas se "Desbloquarem"
	go spawnFruit()

	startGame.Unlock()

	for {

	}
}
