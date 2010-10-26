package main

import "fmt"
import "gnuplot"
import "math"

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

	p.SetStyle("steps")
	p.PlotFunc([]float64{0,1,2,3,4,5,6,7,8,9,10}, 
		func (x float64) float64 {return math.Exp(float64(x) + 2.)},
		"test plot-func")
	p.SetXLabel("my x data")
	p.SetYLabel("my y data")
	p.CheckedCmd("set terminal pdf")
	p.CheckedCmd("set output 'plot004.pdf'")
	p.CheckedCmd("replot")


	p.CheckedCmd("q")
	return
}

/* EOF */
