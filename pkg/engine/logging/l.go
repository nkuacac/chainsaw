package logging

import (
	"fmt"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/pkg/ext/output/color"
	"k8s.io/utils/clock"
)

const eraser = "\b\b\b\b\b\b\b\b\b"

type TestReport interface {
	SetOutput(in ...string)
}

type logger struct {
	t        TLogger
	clock    clock.PassiveClock
	test     string
	step     string
	resource client.Object
	reporter TestReport
}

func NewLogger(t TLogger, clock clock.PassiveClock, test string, step string, reporter TestReport) Logger {
	t.Helper()
	return &logger{
		t:        t,
		clock:    clock,
		test:     test,
		step:     step,
		reporter: reporter,
	}
}

func (l *logger) Log(operation Operation, status Status, color *color.Color, args ...fmt.Stringer) {
	a := l.formatLog(operation, status, color, args)
	l.t.Log(fmt.Sprint(a...))
	l.report(operation, status, args...)
}

func (l *logger) formatLog(operation Operation, status Status, color *color.Color, args []fmt.Stringer) []any {
	sprint := fmt.Sprint
	opLen := 9
	stLen := 5
	if color != nil {
		sprint = color.Sprint
		opLen += 14
		stLen += 14
	}
	a := make([]any, 0, len(args)+2)
	prefix := fmt.Sprintf("%s| %s | %s | %s | %-*s | %-*s |", eraser, l.clock.Now().Format("15:04:05"), sprint(l.test), sprint(l.step), opLen, sprint(operation), stLen, sprint(status))
	if l.resource != nil {
		gvk := l.resource.GetObjectKind().GroupVersionKind()
		key := client.Key(l.resource)
		prefix = fmt.Sprintf("%s %s/%s @ %s", prefix, gvk.GroupVersion(), gvk.Kind, client.Name(key))
	}
	a = append(a, prefix)
	for _, arg := range args {
		a = append(a, "\n")
		a = append(a, arg)
	}
	return a
}

func (l *logger) report(operation Operation, status Status, args ...fmt.Stringer) {
	a := l.formatLog(operation, status, nil, args)
	if l.reporter != nil {
		l.reporter.SetOutput(fmt.Sprint(a...))
	}
}

func (l *logger) WithResource(resource client.Object) Logger {
	return &logger{
		t:        l.t,
		clock:    l.clock,
		test:     l.test,
		step:     l.step,
		reporter: l.reporter,
		resource: resource,
	}
}
