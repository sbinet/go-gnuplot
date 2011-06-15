//
package gnuplot

import "os"
import "io/ioutil"
//import "bytes"
//import "container/vector"
import "exec"
import "fmt"

var g_gnuplot_cmd string
var g_gnuplot_prefix string = "go-gnuplot-"

func min(a,b int) int {
	if a<b {
		return a
	}
	return b
}

func init() {
	var err os.Error
	g_gnuplot_cmd, err = exec.LookPath("gnuplot")
	if err != nil {
		fmt.Printf("** could not find path to 'gnuplot':\n%v\n", err)
		panic("could not find 'gnuplot'")
	}
	fmt.Printf("-- found gnuplot command: %s\n", g_gnuplot_cmd)
}

type gnuplot_error struct {
	err string
}
func (e *gnuplot_error) String() string {
	return e.err
}
type plotter_process struct {
	handle  *exec.Cmd
}

func new_plotter_proc(persist bool) (*plotter_process, os.Error) {
	var proc_args []string
	if persist {
		proc_args = []string{"gnuplot", "-persist"}
	} else {
		proc_args = []string{"gnuplot"}
	}
	fmt.Printf("--> [%v] %v\n", g_gnuplot_cmd, proc_args)
	stdin  := exec.Pipe
	stdout := exec.Pipe
	//stdout := exec.PassThrough
	stderr := exec.MergeWithStdout

	cmd, err := exec.Run(g_gnuplot_cmd, proc_args, os.Environ(), "", 
		stdin, stdout, stderr)

	if err != nil {
		return nil, err
	}
	return &plotter_process{handle: cmd}, nil
}

type tmpfiles_db map[string]*os.File

type Plotter struct {
	proc *plotter_process
	debug bool
	plotcmd string
	nplots int // number of currently active plots
	style  string // current plotting style
	tmpfiles tmpfiles_db
}

func (self *Plotter) Cmd(format string, a ...interface{}) os.Error {
	cmd := fmt.Sprintf(format, a...) + "\n"
	n,err := self.proc.handle.Stdin.WriteString(cmd)
	
	if self.debug {
		//buf := new(bytes.Buffer)
		//io.Copy(buf, self.proc.handle.Stdout)
		fmt.Printf("cmd> %v", cmd)
		fmt.Printf("res> %v\n", n)
	}

	return err
}

func (self *Plotter) CheckedCmd(format string, a ...interface{}) {
	err := self.Cmd(format, a...)
	if err != nil {
		err_string := fmt.Sprintf("** err: %v\n", err)
		panic(err_string)
	}
}

func (self *Plotter) Close() (err os.Error) {
	if self.proc != nil && self.proc.handle != nil {
		err = self.proc.handle.Close()
	}
	self.ResetPlot()
	return err
}

func (self *Plotter) PlotNd(title string, data ...[]float64) os.Error {
	ndims := len(data)

	switch ndims {
	case 1: return self.PlotX(data[0], title)
	case 2: return self.PlotXY(data[0], data[1], title)
	case 3: return self.PlotXYZ(data[0], data[1], data[2], title)
	}

	return &gnuplot_error{fmt.Sprintf("invalid number of dims '%v'", ndims)}
}

func (self *Plotter) PlotX(data []float64, title string) os.Error {
	f, err := ioutil.TempFile(os.TempDir(), g_gnuplot_prefix)
	if err != nil {
		return err
	}
	fname := f.Name()
	self.tmpfiles[fname] = f
	for _,d := range data {
		f.WriteString(fmt.Sprintf("%v\n", d))
	}
	f.Close()
	cmd := "plot"
	if self.nplots > 0 {
		cmd = "replot"
	}
	
	var line string
	if title == "" {
		line = fmt.Sprintf("%s \"%s\" with %s", cmd, fname, self.style)
	} else {
		line = fmt.Sprintf("%s \"%s\" title \"%s\" with %s",
			cmd, fname, title, self.style)
	}
	self.nplots += 1
	return self.Cmd(line)
}

func (self *Plotter) PlotXY(x,y []float64, title string) os.Error {
	npoints := min(len(x), len(y))

	f, err := ioutil.TempFile(os.TempDir(), g_gnuplot_prefix)
	if err != nil {
		return err
	}
	fname := f.Name()
	self.tmpfiles[fname] = f

	for i:=0; i < npoints; i++ {
		f.WriteString(fmt.Sprintf("%v %v\n", x[i], y[i]))
	}

	f.Close()
	cmd := "plot"
	if self.nplots > 0 {
		cmd = "replot"
	}
	
	var line string
	if title == "" {
		line = fmt.Sprintf("%s \"%s\" with %s", cmd, fname, self.style)
	} else {
		line = fmt.Sprintf("%s \"%s\" title \"%s\" with %s",
			cmd, fname, title, self.style)
	}
	self.nplots += 1
	return self.Cmd(line)
}

func (self *Plotter) PlotXYZ(x,y,z []float64, title string) os.Error {
	npoints := min(len(x), len(y))
	npoints = min(npoints, len(z))
	f, err := ioutil.TempFile(os.TempDir(), g_gnuplot_prefix)
	if err != nil {
		return err
	}
	fname := f.Name()
	self.tmpfiles[fname] = f

	for i:=0; i < npoints; i++ {
		f.WriteString(fmt.Sprintf("%v %v %v\n", x[i], y[i], z[i]))
	}

	f.Close()
	cmd := "splot"
	if self.nplots > 0 {
		cmd = "replot"
	}
	
	var line string
	if title == "" {
		line = fmt.Sprintf("%s \"%s\" with %s", cmd, fname, self.style)
	} else {
		line = fmt.Sprintf("%s \"%s\" title \"%s\" with %s",
			cmd, fname, title, self.style)
	}
	self.nplots += 1
	return self.Cmd(line)
}

type Func func(x float64) float64

func (self *Plotter) PlotFunc(data []float64, fct Func, title string) os.Error {
	
	f, err := ioutil.TempFile(os.TempDir(), g_gnuplot_prefix)
	if err != nil {
		return err
	}
	fname := f.Name()
	self.tmpfiles[fname] = f

	for _,x := range data {
		f.WriteString(fmt.Sprintf("%v %v\n", x, fct(x)))
	}

	f.Close()
	cmd := "plot"
	if self.nplots > 0 {
		cmd = "replot"
	}
	
	var line string
	if title == "" {
		line = fmt.Sprintf("%s \"%s\" with %s", cmd, fname, self.style)
	} else {
		line = fmt.Sprintf("%s \"%s\" title \"%s\" with %s",
			cmd, fname, title, self.style)
	}
	self.nplots += 1
	return self.Cmd(line)
}

func (self *Plotter) SetStyle(style string) (err os.Error) {
	allowed := []string{
		"lines", "points", "linepoints",
		"impulses", "dots",
		"steps",
		"errorbars",
		"boxes",
		"boxerrorbars"}

	for _,s := range allowed {
		if s == style {
			self.style = style
			err = nil
			return err
		}
	}

	fmt.Printf("** style '%v' not in allowed list %v\n", style, allowed)
	fmt.Printf("** default to 'points'\n")
	self.style = "points"
	err = &gnuplot_error{fmt.Sprintf("invalid style '%s'", style)}

	return err
}

func (self *Plotter) SetXLabel(label string) os.Error {
	return self.Cmd(fmt.Sprintf("set xlabel '%s'", label))
}

func (self *Plotter) SetYLabel(label string) os.Error {
	return self.Cmd(fmt.Sprintf("set ylabel '%s'", label))
}

func (self *Plotter) SetZLabel(label string) os.Error {
	return self.Cmd(fmt.Sprintf("set zlabel '%s'", label))
}

func (self *Plotter) SetLabels(labels ...string) os.Error {
	ndims := len(labels)
	if ndims > 3 || ndims <= 0 {
		return &gnuplot_error{fmt.Sprintf("invalid number of dims '%v'", ndims)}
	}
	var err os.Error = nil

	for i,label := range labels {
		switch i {
		case 0: 
			ierr := self.SetXLabel(label)
			if ierr != nil {
				err = ierr
				return err
			}
		case 1:
			ierr := self.SetYLabel(label)
			if ierr != nil {
				err = ierr
				return err
			}
		case 2:
			ierr := self.SetZLabel(label)
			if ierr != nil {
				err = ierr
				return err
			}
		}
	}
	return nil
}

func (self *Plotter) ResetPlot() (err os.Error) {
	for fname, fhandle := range self.tmpfiles {
		ferr := fhandle.Close()
		if ferr != nil {
			err = ferr
		}
		os.Remove(fname)
	}
	self.nplots = 0
	return err
}

func NewPlotter(fname string, persist, debug bool) (*Plotter, os.Error) {
	p := &Plotter{proc: nil, debug: debug, plotcmd: "plot", 
	nplots:0, style:"points"}
	p.tmpfiles = make(tmpfiles_db)

	if fname != "" {
		panic("NewPlotter with fname is not yet supported")
	} else {
		proc, err := new_plotter_proc(persist)
		if err != nil {
			return nil, err
		}
		p.proc = proc
	}
	return p, nil
}

/* EOF */
