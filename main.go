package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	ui "github.com/gizak/termui"
	"github.com/nsf/termbox-go"
)

var view = make([]float64, 0, 10)
var data = make([]float64, 0, 10)
var offset = 0

var lc *ui.LineChart
var par *ui.Par

func readStdIn() {
	r := bufio.NewReader(os.Stdin)
	for s, err := r.ReadString('\n'); err == nil; s, err = r.ReadString('\n') {
		f, e := strconv.ParseFloat(strings.TrimSpace(s), 64)
		if e != nil {
			par.Text = e.Error()
		}
		data = append(data, f)
	}
	//par.Text = fmt.Sprintf("%v", data)
	lc.Data = data
}

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	lc = ui.NewLineChart()
	lc.Data = data
	lc.Height = 60

	par = ui.NewPar("RESULTS")
	par.Height = 15

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, lc),
		),
		ui.NewRow(
			ui.NewCol(12, 0, par),
		),
	)
	ui.Body.Align()

	go readStdIn()

	//Set up input catching
	evt := make(chan termbox.Event)
	go func() {
		for {
			evt <- termbox.PollEvent()
		}
	}()
	go draw()
	for {
		select {
		case e := <-evt:
			switch e.Type {
			case termbox.EventKey:
				switch e.Ch {
				case 'q':
					return
				case 0:
					switch e.Key {
					case termbox.KeyArrowRight:
						offset += 1
					case termbox.KeyArrowLeft:
						offset -= 1
						if offset < 0 {
							offset = 0
						}
					}
				}
			}
		}

	}

}

func draw() {
	for {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		lc.Data = data[offset:]
		ui.Render(ui.Body)
		time.Sleep(time.Second / 10)
	}
}
