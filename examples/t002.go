package main

import "fmt"
//import "io/ioutil"
import "github.com/sbinet/go-gnuplot/pkg/gnuplot"

func main() {
	fname := ""
	persist := false
	debug := true

	p,err := gnuplot.NewPlotter(fname, persist, debug)
	if err != nil {
		err_string := fmt.Sprintf("** err: %v\n", err)
		panic(err_string)
	}
	defer p.Close()

	p.PlotX([]float64{0,1,2,3,4,5,6,7,8,9,10}, "some data")
	p.CheckedCmd("set terminal pdf")
	p.CheckedCmd("set output 'plot002.pdf'")
	p.CheckedCmd("replot")


	p.CheckedCmd("q")
	return
}

/* EOF */
