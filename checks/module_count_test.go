package checks

import (
	"fmt"
	"github.com/veracode/scan_health/v2/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veracode/scan_health/v2/report"
)

func generateModule(name string) report.Module {

	return report.Module{
		Name:         name,
		IsIgnored:    false,
		IsThirdParty: false,
		Instances: []report.ModuleInstance{
			{IsDependency: false},
			{IsSelected: true},
		}}
}

func TestModuleCount(t *testing.T) {

	t.Run("Small Number of modules", func(t *testing.T) {
		t.Parallel()
		mockReport := report.Report{
			Modules: []report.Module{
				generateModule("module1.exe"),
				generateModule("module2.dll"),
			},
		}

		moduleCount(&mockReport)
		assert.Empty(t, mockReport.Issues)
	})

	t.Run("Large Number of modules", func(t *testing.T) {
		t.Parallel()

		var lotsOfModules []report.Module

		for i := 0; i < utils.MaximumModuleSelectedCountThreshold+1; i++ {
			moduleName := fmt.Sprintf("module%d", i)
			lotsOfModules = append(lotsOfModules, generateModule(moduleName))
		}

		mockReport := report.Report{
			Modules: lotsOfModules,
		}

		moduleCount(&mockReport)

		if !assert.Equal(t, 1, len(mockReport.Issues)) {
			t.FailNow()
		}

		assert.Contains(t, mockReport.Issues[0].Description, fmt.Sprintf("%d modules were selected", utils.MaximumModuleSelectedCountThreshold+1))
	})

}
