package report

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ret struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Remark  string `json:"remark"`
}

func saveJson(report *Report, file string) error {
	var rows []ret

	perFolder := map[string][]*TestReport{}
	for _, test := range report.tests {
		perFolder[test.test.BasePath] = append(perFolder[test.test.BasePath], test)
	}
	for _, tests := range perFolder {
		for _, test := range tests {
			var row ret
			row.Name = test.test.Test.Name
			if test.failed {
				row.Status = "fail"
			} else if test.skipped {
				row.Status = "skip"
			} else {
				row.Status = "pass"
			}
			if test.err != nil {
				row.Remark = test.err.Error()
			}
			for i, step := range test.steps {
				if step.err != nil {
					row.Remark = step.err.Error()
				}
				for j, op := range step.reports {
					if op.err != nil {
						row.Remark = fmt.Sprintf("step %d: %s op %d - %s[%s]: %s", i, step.step.Description, j, op.operationType, op.name, op.err.Error())
					}
				}
			}
			row.Message = strings.Join(test.output, "\n")
			rows = append(rows, row)
		}
	}

	marshal, err := json.Marshal(rows)
	if err != nil {
		return err
	}
	return os.WriteFile(file, marshal, 0o600)
}
