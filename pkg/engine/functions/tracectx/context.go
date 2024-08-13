package tracectx

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/mohae/deepcopy"
)

type Context struct {
	InvokeTraces []*InvokeTrace
	StepTraces   []*StepTrace
	Parent       *Context
}

func (c *Context) String() string {
	var ret = []string{"InvokeTrace:"}
	for _, trace := range c.InvokeTraces {
		ret = append(ret, trace.String())
	}

	var steps []string
	cur := c
	prefix := ""
	for cur != nil {
		if len(cur.StepTraces) > 0 {
			for _, step := range cur.StepTraces {
				steps = append(steps, prefix+step.String())
			}
		}
		cur = cur.Parent
		prefix += " "
	}
	if len(steps) > 0 {
		ret = append(ret, "Steps:")
		ret = append(ret, steps...)
	}
	return strings.Join(ret, "\n")
}

type InvokeTrace struct {
	File string
	Line int
	Name string
}

func (c *InvokeTrace) String() string {
	return fmt.Sprintf("%s:%d %s", c.File, c.Line, c.Name)
}

type StepTrace struct {
	File string
	Line int
	Step string
	Args []interface{}
}

func (c *StepTrace) String() string {
	return fmt.Sprintf("%s:%d %s%v", c.File, c.Line, c.Step, c.Args)
}

func (c *Context) Step(name string, args ...interface{}) *Context {
	_, file, line, _ := runtime.Caller(1)
	c.StepTraces = append(c.StepTraces, &StepTrace{file, line, name, args})
	return c
}
func (c *Context) By(name string, f func()) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("%s\n %s:%d\n", name, file, line)
	f()
}
func (c *Context) trace(file string, line int, funcName string) {
	c.InvokeTraces = append(c.InvokeTraces, &InvokeTrace{file, line, funcName})
}

func (c *Context) Do() *Context {
	newc := deepcopy.Copy(c).(*Context)
	//var newStepTrace []*StepTrace
	//newStepTrace = append(newStepTrace, c.StepTraces...)
	newc.StepTraces = []*StepTrace{}
	newc.Parent = c
	pc, file, line, _ := runtime.Caller(1)
	funcName := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	if len(newc.InvokeTraces) > 0 {
		newc.InvokeTraces[len(newc.InvokeTraces)-1].Name = funcName[len(funcName)-1]
	}
	newc.trace(file, line, "self")
	return newc
}

func nowStamp() string {
	return time.Now().Format(time.StampMicro)
}
func (c *Context) InfoF(format string, args ...interface{}) {
	fmt.Printf("Info "+nowStamp()+": "+format+"\n", args...)
	dep := len(c.InvokeTraces)
	if dep > 0 {
		file := c.InvokeTraces[dep-1].File
		line := c.InvokeTraces[dep-1].Line
		fmt.Printf("%v:%v\n", file, line)
	}
	if dep > 1 {
		file := c.InvokeTraces[0].File
		line := c.InvokeTraces[0].Line
		fmt.Printf("%v:%v\n", file, line)
	}
}
