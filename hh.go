package main

import (
	tm "github.com/buger/goterm"
	"github.com/gobuffalo/packr"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type CoOrdinate struct {
	x int
	y int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}


func turnOffCursor() {
	fmt.Print("\033[?25l")
}

func turnOnCursor() {
	fmt.Print("\033[?25h")
}


var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["darwin"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok { //if we defined a clear func for that platform:
		value()  //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}


func main() {
	turnOffCursor()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	stop := make(chan bool, 1)
	box := packr.NewBox("./resources")
	data, err := box.FindString("hh1.txt")
	check(err)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- true
	}()

	go func() {
		rand.Seed(time.Now().Unix())
		colours := []int{
			tm.RED,
			tm.GREEN,
			tm.BLUE,
			tm.CYAN,
			tm.BLACK,
			tm.YELLOW,
			tm.MAGENTA,
		}

		for {
			select {
			case <-stop:
				tm.Clear()
				tm.Flush()
			default:
				for i := 0; i < 20; i ++ {
					tm.Clear()
					tm.MoveCursor(1, 1)
					for _, c := range data {
						if c == '%' {
							n := rand.Int() %  len(colours)
							clr := colours[n]
							tm.Print(tm.Color("%", clr))
							//fmt.Printf("%s", red(string(c)))
						}else{
							tm.Print(string(c))

						}
					}
					tm.Flush()
					time.Sleep(400 * time.Millisecond)
				}
			}
		}
	}()

	<-done
	stop <- true
	tm.MoveCursor(1, 1)
	tm.Clear()
	tm.Flush()
	CallClear()
	fmt.Println("Based on https://github.com/hvishwanath/happy-holidays")
	turnOnCursor()
}
