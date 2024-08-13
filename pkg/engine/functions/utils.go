package functions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jmespath/go-jmespath"
	"github.com/kyverno/chainsaw/pkg/engine/functions/tracectx"
	"io"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

func stable(in string) string {
	return in
}

func experimental(in string) string {
	return "x_" + in
}

func getArgAt(arguments []any, index int) (any, error) {
	if index >= len(arguments) {
		return nil, fmt.Errorf("index out of range (%d / %d)", index, len(arguments))
	}
	return arguments[index], nil
}

func getArg[T any](arguments []any, index int, out *T) error {
	arg, err := getArgAt(arguments, index)
	if err != nil {
		return err
	}
	if value, ok := arg.(T); !ok {
		return errors.New("invalid type")
	} else {
		*out = value
		return nil
	}
}

// Prettify returns the string representation of a value.
func Prettify(i interface{}) string {
	var buf bytes.Buffer
	prettify(reflect.ValueOf(i), 0, &buf)
	return buf.String()
}

// prettify will recursively walk value v to build a textual
// representation of the value.
func prettify(v reflect.Value, indent int, buf *bytes.Buffer) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		strtype := v.Type().String()
		if strtype == "time.Time" {
			_, err := fmt.Fprintf(buf, "%s", v.Interface())
			if err != nil {
				panic(err)
			}
			break
		} else if strings.HasPrefix(strtype, "io.") {
			buf.WriteString("<buffer>")
			break
		} else if strtype == "resource.Quantity" {
			value := v.Interface().(resource.Quantity)
			buf.WriteString(value.String())
			break
		}

		buf.WriteString("{\n")

		var names []string
		for i := 0; i < v.Type().NumField(); i++ {
			name := v.Type().Field(i).Name
			f := v.Field(i)
			if name[0:1] == strings.ToLower(name[0:1]) {
				continue // ignore unexported fields
			}
			if f.IsZero() {
				continue // ignore empty fields
			}
			if (f.Kind() == reflect.Ptr || f.Kind() == reflect.Slice || f.Kind() == reflect.Map) && f.IsNil() {
				continue // ignore unset fields
			}
			names = append(names, name)
		}

		for i, n := range names {
			val := v.FieldByName(n)

			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(n + ": ")
			prettify(val, indent+2, buf)

			if i < len(names)-1 {
				buf.WriteString(",\n")
			}
		}

		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")
	case reflect.Slice:
		strtype := v.Type().String()
		if strtype == "[]uint8" {
			_, err := fmt.Fprintf(buf, "<binary> len %d", v.Len())
			if err != nil {
				panic(err)
			}
			break
		}

		nl, id, id2 := "", "", ""
		if v.Len() > 3 {
			nl, id, id2 = "\n", strings.Repeat(" ", indent), strings.Repeat(" ", indent+2)
		}
		buf.WriteString("[" + nl)
		for i := 0; i < v.Len(); i++ {
			buf.WriteString(id2)
			prettify(v.Index(i), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString("," + nl)
			}
		}

		buf.WriteString(nl + id + "]")
	case reflect.Map:
		buf.WriteString("{\n")

		for i, k := range v.MapKeys() {
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(k.String() + ": ")
			prettify(v.MapIndex(k), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString(",\n")
			}
		}

		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")
	default:
		if !v.IsValid() {
			_, err := fmt.Fprint(buf, "<invalid value>")
			if err != nil {
				panic(err)
			}
			return
		}
		format := "%v"
		switch v.Interface().(type) {
		case string:
			format = "%q"
		case io.ReadSeeker, io.Reader:
			format = "buffer(%p)"
		}
		_, err := fmt.Fprintf(buf, format, v.Interface())
		if err != nil {
			panic(err)
		}
	}
}

func ParallelRun(ctx *tracectx.Context, max int, parallel int, toRun func(int) error) error {
	if max == 0 {
		return nil
	}
	errChan := make(chan string)
	var errors []string
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		for err := range errChan {
			errors = append(errors, err)
		}
		wg1.Done()
	}()

	waitChan := make(chan struct{}, parallel)
	var wg sync.WaitGroup

	//stepTimer := prometheus.NewTimer(costStep)
	for i := 0; i < max; i++ {
		if len(errors) > 0 {
			break
		}
		waitChan <- struct{}{}
		wg.Add(1)
		go func(idx int) {
			defer func() {
				msg := recover()
				if msg != nil {
					stack := []string{}
					for x := 3; x < 6; x++ {
						pc, file, line, _ := runtime.Caller(x)
						pcName := runtime.FuncForPC(pc).Name()
						stack = append(stack, fmt.Sprintf("%s\n%s:%d\n", pcName, file, line))
					}
					errChan <- ctx.Do().Step("panic", fmt.Sprintf("%v\n%s", msg, strings.Join(stack, ""))).String()
				}
				wg.Done()
				<-waitChan
			}()
			err := toRun(idx)
			if err != nil {
				errChan <- ctx.Do().Step("error", err.Error()).String()
			}
		}(i)
	}
	wg.Wait()
	close(errChan)
	wg1.Wait()

	if len(errors) != 0 {
		return fmt.Errorf("%s", strings.Join(errors, "\n"))
	}
	return nil
}

func NewLabelSelector(kvs ...string) labels.Selector {
	if len(kvs) == 0 || len(kvs[0]) == 0 {
		return labels.Everything()
	}
	if kvs[0] == "nil" {
		return nil
	}

	dict, toStruct := toMap(kvs)

	if toStruct {
		var requirements []labels.Requirement
		for k, v := range dict {
			if len(v) == 0 {
				if strings.Contains(k, "!") {
					require, err := labels.NewRequirement(strings.TrimSpace(strings.ReplaceAll(k, "!", "")), selection.DoesNotExist, nil)
					if err == nil {
						requirements = append(requirements, *require)
					}
				} else {
					require, err := labels.NewRequirement(k, selection.Exists, nil)
					if err == nil {
						requirements = append(requirements, *require)
					}
				}
			} else {
				if strings.Contains(k, "!") {
					require, err := labels.NewRequirement(strings.TrimSpace(strings.ReplaceAll(k, "!", "")), selection.NotIn, strings.Split(v, "|"))
					if err == nil {
						requirements = append(requirements, *require)
					}
				} else {
					require, err := labels.NewRequirement(k, selection.In, strings.Split(v, "|"))
					if err == nil {
						requirements = append(requirements, *require)
					}
				}
			}
		}
		return labels.NewSelector().Add(requirements...)
	}
	return labels.SelectorFromSet(dict)
}

func toMap(kvs []string) (map[string]string, bool) {
	labels := make(map[string]string)
	var toStruct = false
	for _, kv := range kvs {
		kv = strings.TrimSpace(kv)
		if len(kv) == 0 {
			continue
		}
		arr := strings.Split(kv, "=")
		if len(arr) == 1 {
			labels["app"] = arr[0]
		} else {
			labels[arr[0]] = arr[1]
			if len(arr[1]) == 0 {
				toStruct = true
			}
			if strings.Contains(arr[1], "|") {
				toStruct = true
			}
			if strings.Contains(arr[0], "!") {
				toStruct = true
			}
		}
	}
	return labels, toStruct
}

func tableFormat(arguments []any) (any, error) {
	var input map[string]interface{}
	var path string
	if err := getArg(arguments, 0, &input); err != nil {
		return "", err
	}
	if len(arguments) >= 2 {
		if err := getArg(arguments, 1, &path); err != nil {
			return "", err
		}
	}

	var desc [][]interface{}
	items, ok := input["items"].([]interface{})
	if ok {
		for _, item := range items {
			single := addSingleItem(item, path)
			if single != nil {
				desc = append(desc, single)
			}
		}
	} else {
		single := addSingleItem(input, path)
		if single != nil {
			desc = append(desc, single)
		}
	}

	fmt.Printf("\n%s\n", TableRenderWithName("", []string{"name", "status"}, desc))
	return "", nil
}

func addSingleItem(item interface{}, path string) []interface{} {
	itemMap, ok := item.(map[string]interface{})
	if !ok {
		return nil
	}
	obj := unstructured.Unstructured{Object: itemMap}
	obj.GetAPIVersion()
	name := obj.GetName()

	if path != "" {
		search, err := jmespath.Search(path, item)
		if err == nil {
			return []interface{}{name, search}
		}
	}
	marshal, _ := json.MarshalIndent(itemMap["status"], "", " ")
	return []interface{}{name, string(marshal)}
}

func TableRenderWithName(name string, header []string, rows [][]interface{}) string {
	t := table.NewWriter()
	t.Style().Options = table.OptionsNoBordersAndSeparators
	if len(header) > 0 {
		new_header := convertToInterface(header...)
		new_header[0] = name
		t.AppendHeader(new_header)
	}
	for i, row := range rows {
		new_row := []interface{}{i}
		new_row = append(new_row, row...)
		t.AppendRow(new_row)
	}
	return t.Render()
}

func convertToInterface(t ...string) []interface{} {
	s := make([]interface{}, len(t)+1)
	for i, v := range t {
		s[i+1] = v
	}
	return s
}
