package report

import (
	"encoding/json"
	"fmt"
	"os"
)

type ret struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
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
			for i, step := range test.steps {
				for j, op := range step.reports {
					if op.err != nil {
						row.Message = fmt.Sprintf("step %d op %d - %s: %s", i, j, op.operationType, op.operationType)
					}
				}
			}
			rows = append(rows, row)
		}
	}

	marshal, err := json.Marshal(rows)
	if err != nil {
		return err
	}
	return os.WriteFile(file, marshal, 0o600)
}
