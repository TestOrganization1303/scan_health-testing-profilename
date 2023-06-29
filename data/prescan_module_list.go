package data

import (
	"encoding/xml"
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
	"net/http"
	"strconv"
	"strings"
)

type prescanModuleList struct {
	XMLName xml.Name        `xml:"prescanresults"`
	Modules []prescanModule `xml:"module"`
}

type prescanModule struct {
	XMLName        xml.Name             `xml:"module"`
	Id             int                  `xml:"id,attr"`
	Name           string               `xml:"name,attr"`
	Status         string               `xml:"status,attr"`
	Platform       string               `xml:"platform,attr"`
	Size           string               `xml:"size,attr"`
	MD5            string               `xml:"checksum,attr"`
	HasFatalErrors bool                 `xml:"has_fatal_errors,attr"`
	IsDependency   bool                 `xml:"is_dependency,attr"`
	Issues         []prescanModuleIssue `xml:"issue"`
	SizeBytes      int
}

type prescanModuleIssue struct {
	XMLName xml.Name `xml:"issue"`
	Details string   `xml:"details,attr"`
}

func (api API) getPrescanModuleList(r *report.Report) {
	var url = fmt.Sprintf("https://analysiscenter.veracode.com/api/5.0/getprescanresults.do?app_id=%d&build_id=%d", r.Scan.ApplicationId, r.Scan.BuildId)
	response := api.makeApiRequest(url, http.MethodGet)

	moduleList := prescanModuleList{}

	err := xml.Unmarshal(response, &moduleList)

	if err != nil {
		utils.ErrorAndExit("Could not get prescan results", err)
	}

	// Sort modules by name for consistency
	// We will sort later actually
	//sort.Slice(moduleList.Modules, func(i, j int) bool {
	//	return moduleList.Modules[i].Name < moduleList.Modules[j].Name
	//})

	for _, module := range moduleList.Modules {
		var issues []string

		for _, issue := range module.Issues {
			issues = append(issues, issue.Details)
		}

		r.Modules = append(
			r.Modules,
			report.Module{
				Id:             module.Id,
				Name:           module.Name,
				Status:         module.Status,
				Platform:       module.Platform,
				Size:           module.Size,
				MD5:            module.MD5,
				HasFatalErrors: module.HasFatalErrors,
				IsDependency:   module.IsDependency,
				Issues:         issues,
				//SizeBytes:      calculateModuleSize(module.Size),
			},
		)
	}
}

func calculateModuleSize(size string) int {
	var totalModuleSize = 0
	totalModuleSize += convertSize(size, "GB", 1e+9)
	totalModuleSize += convertSize(size, "MB", 1e+6)
	totalModuleSize += convertSize(size, "KB", 1000)
	return totalModuleSize
}

func convertSize(size, measurement string, multiplier int) int {
	if !strings.HasSuffix(size, measurement) {
		return 0
	}

	formattedSize := strings.TrimSuffix(size, measurement)
	sizeInt, err := strconv.Atoi(formattedSize)

	if err != nil {
		panic(err)
	}

	return sizeInt * multiplier

}

func (module prescanModule) getFatalReason() string {
	for _, issue := range strings.Split(module.Status, ",") {
		if strings.HasPrefix(issue, "(Fatal)") {
			return strings.Replace(strings.Replace(issue, "(Fatal)", "", 1), " - 1 File", "", 1)
		}
	}

	return ""
}