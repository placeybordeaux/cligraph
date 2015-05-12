package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	ui "github.com/gizak/termui"
	"github.com/montanaflynn/stats"
	"github.com/nsf/termbox-go"
)

var view = make([]float64, 0, 10)
var data = make([]float64, 0, 10)
var offset = 0

var lc *ui.LineChart
var total_par *ui.Par
var window_par *ui.Par

func readStdIn() {
	r := bufio.NewReader(os.Stdin)
	for s, err := r.ReadString('\n'); err == nil; s, err = r.ReadString('\n') {
		f, e := strconv.ParseFloat(strings.TrimSpace(s), 64)
		if e != nil {
			total_par.Text = e.Error() //TODO make this a bit more reasonable
		}
		data = append(data, f)
		//total_par.Text = info(data)
	}
	lc.Data = data
}

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer func() {
		if r := recover(); r != nil {
			ui.Close()
			panic(r)
		}
	}()

	defer ui.Close()

	lc = ui.NewLineChart()
	lc.Data = data
	lc.Height = 50

	total_par = ui.NewPar("RESULTS")
	total_par.Height = 15

	window_par = ui.NewPar("RESULTS")
	window_par.Height = 15

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, lc),
		),
		ui.NewRow(
			ui.NewCol(6, 0, total_par),
			ui.NewCol(6, 0, window_par),
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
					case termbox.KeySpace:
						offset += (ui.Body.Width - 7) * 2
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

func update_window() {
	window_length := (ui.Body.Width - 7) * 2
	//                  /             |  |
	//    total width__/    borders__/    \_ two points per character
	if window_length < len(data) {
		lc.Data = data[offset : offset+window_length]
	} else if offset < len(data) {
		lc.Data = data[offset:]
	}
	window_par.Text = info(lc.Data)
}

func info(d []float64) string {
	defer func() {
		if r := recover(); r != nil {
			ui.Close()
			panic(r)
		}
	}()
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "%v total points\n", len(d))
	fmt.Fprintf(b, "Min: %v\tAvg: %v\tMax: %v\n", stats.Min(d), stats.Mean(d), stats.Max(d))
	//fmt.Fprintf(b, "25th %v\t50th %v\t75th %v\t90th\t%v\t99th %v\n", stats.Percentile(d, 25), stats.Percentile(d, 50), stats.Percentile(d, 75), stats.Percentile(d, 90), stats.Percentile(d, 99))
	fmt.Fprintf(b, "Standard deviation\n", stats.StdDevS(d))
	return b.String()
}

func draw() {
	for {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		update_window()
		ui.Render(ui.Body)
		time.Sleep(time.Second / 10)
	}
}
