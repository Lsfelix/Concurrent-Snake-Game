//Concurrent Snake Game -  Lucas Félix
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

	"./MCCSemaforo"
	term "github.com/nsf/termbox-go"
)

//Pair : Estrutura de Par de valores para armazenar posições da cobra (X,Y)
type Pair struct {
	a, b interface{}
}

//SIZE : Tamanho do campo do jogo, configurável.
const SIZE int = 20

//PLAYERS : Definição do número de jogadores, configurável
const PLAYERS int = 3

//SPAWNPERIOD : Definição do período de geramento de frutas, configurável
const SPAWNPERIOD time.Duration = 3

//Definição global do campo do jogo.
var field [SIZE][SIZE]string

//Definição global de variáveis utilizadas para coordenação dos processos;.
var inicialization = 0
var waiting = 0

//Definição de Mutex utilizados
var startGame sync.Mutex
var updateGame sync.Mutex
var lockWaiting sync.Mutex

var inputSemaphore = MCCSemaforo.NewSemaphore(0)

//Rotina que detecta entradas e as envia para o canal de entradas.
func treatInput(snakes [PLAYERS + 1]chan rune) {
	//Inicializa o módulo
	term.Init()

	waiting = PLAYERS

	//Detecta entradas até o jogo acabar ou detectar uma entrada "ESC"
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {

			case term.KeyEsc:
				term.Close()
				os.Exit(3)
			}

			for index := 1; index <= PLAYERS; index++ {
				snakes[index] <- ev.Ch
			}

			for index := 0; index < waiting; index++ {
				inputSemaphore.Wait()
			}

		case term.EventError:
			panic(ev.Err)
		}

	}

}

//Rotina que gera as frutas no campo
func spawnFruit() {
	//Cria um gerador de números aleátorios a partir de uma seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		//Aleatóriamente encontra um local da grade.
		x := r.Intn(SIZE)
		y := r.Intn(SIZE)
		//Se o local já está ocupado, ele volta ao inicio da iteração.
		if field[y][x] != " _ " {
			continue
		}
		//Senão, coloca uma fruta no local.
		field[y][x] = " F "
		time.Sleep(SPAWNPERIOD * time.Second)
	}

}

//Rotina de controle das cobras.
func snakeController(id int, end chan int, inputs chan rune, up rune, down rune, left rune, right rune) {

	//Define a aparência da cobra no campo de acordo com seu id.
	name := " " + strconv.Itoa(id) + " "

	dead := false

	//Define a posição inicial da cobra de acordo com o tamanho do campo e número máximo de jogadoes.
	headX := (SIZE / (PLAYERS + 1)) * id
	headY := (SIZE / (PLAYERS + 1)) * id

	//Posiciona cabeça da cobra e cria estrutura do corpo
	field[headY][headX] = name
	body := list.New()

	//Espera a inicializaçao de todas as cobras.
	startGame.Lock()
	//Aumenta o contador
	inicialization++
	//Se todas cobras já passaram, não necessita dar unlock.
	if inicialization != PLAYERS {
		startGame.Unlock()
	}

	//Iteração de controle
	for {

		//Recebe Input
		key := <-inputs

		if dead {
			continue
		}

		oldHeadX := headX
		oldHeadY := headY

		//Se input corresponder a esta cobra, movimenta, senão, manda input para outras cobras e espera alguma aceitar o input.
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

			//Mutex para singal do semaforo, como é uma adição a uma variável pode dar problema
			lockWaiting.Lock()
			inputSemaphore.Signal()
			lockWaiting.Unlock()

			continue
		}

		//Toma controle das alterações no campo.
		updateGame.Lock()

		//Checa se ocorreu uma colisão.
		if headX < 0 || headX > SIZE-1 || headY < 0 || headY > SIZE-1 || (field[headY][headX] != " _ " && field[headY][headX] != " F ") {
			for e := body.Front(); e != nil; e = e.Next() {
				field[e.Value.(Pair).a.(int)][e.Value.(Pair).b.(int)] = " _ "
			}
			waiting--
			inputSemaphore.Signal()
			updateGame.Unlock()
			end <- id
			dead = true
			continue
		}

		//Movimenta a cobra
		if body.Len() > 0 {
			change := body.Back()
			field[change.Value.(Pair).a.(int)][change.Value.(Pair).b.(int)] = " _ "
			body.Remove(change)
			aux := Pair{oldHeadY, oldHeadX}
			body.PushFront(aux)
		}

		//Checa se comeu uma fruta
		if field[headY][headX] == " F " {
			point := Pair{headY, headX}
			body.PushFront(point)
		}

		field[headY][headX] = name

		//Desenha no campo
		for e := body.Front(); e != nil; e = e.Next() {
			field[e.Value.(Pair).a.(int)][e.Value.(Pair).b.(int)] = name
		}

		//Atualiza o terminal
		drawGame()

		//Libera as cobras aguardando
		inputSemaphore.Signal()

		//Libera o controle do campo
		updateGame.Unlock()

	}

}

//Rotina que desenha o campo do jogo.
func drawGame() {
	//cmd := exec.Command("clear") //Linux
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

func main() {

	//Cria canal de detecção de cobras mortas
	end := make(chan int)

	//Array para marcar quais cobras ainda estão no jogo
	var vivos [PLAYERS + 1]int
	vivos[0] = PLAYERS
	for index := 1; index <= PLAYERS; index++ {
		vivos[index] = 1
	}

	//Inicializa o mutex como trancado.
	startGame.Lock()

	reader := bufio.NewReader(os.Stdin)

	var snakesInputs [PLAYERS + 1]chan rune

	//Inicializa cobras com comandos costumizados
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

		//Cria o canal de inputs.
		inputs := make(chan rune)
		snakesInputs[index] = inputs

		go snakeController(index, end, inputs, up, down, left, right)
	}

	go treatInput(snakesInputs)
	go spawnFruit()

	drawGame()

	startGame.Unlock()

	for {
		morte := <-end
		vivos[morte] = 0
		vivos[0]--
		if vivos[0] == 1 {
			for index := 1; index <= PLAYERS; index++ {
				if vivos[index] == 1 {
					fmt.Println("Parabéns cobra " + strconv.Itoa(index) + "!! Você venceu o jogo.")
					os.Exit(0)
				}
			}
		}
	}
}
