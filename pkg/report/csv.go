package report

import (
	"fmt"
	"github.com/kyverno/chainsaw/pkg/utils/table"
	"os"
)

func saveCsv(report *Report, file string) error {
	var rows [][]interface{}

	perFolder := map[string][]*TestReport{}
	for _, test := range report.tests {
		perFolder[test.test.BasePath] = append(perFolder[test.test.BasePath], test)
	}
	for _, tests := range perFolder {
		for _, test := range tests {
			var row []interface{}
			row = append(row, test.test.Test.Name)
			if test.failed {
				row = append(row, "fail")
			} else if test.skipped {
				row = append(row, "skip")
			} else {
				row = append(row, "pass")
			}
			for i, step := range test.steps {
				for j, op := range step.reports {
					if op.err != nil {
						row = append(row, fmt.Sprintf("step %d op %d - %s: %s", i, j, op.operationType, op.err))
					}
				}
			}
			rows = append(rows, row)
		}
	}
	csv := table.RenderCSV(rows)
	return os.WriteFile(file, []byte(csv), 0o600)
}
