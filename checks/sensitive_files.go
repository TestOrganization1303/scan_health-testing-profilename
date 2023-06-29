package checks

import (
	"fmt"
	"github.com/antfie/scan_health/v2/report"
	"github.com/antfie/scan_health/v2/utils"
)

// Test cases
// https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview::1656378:24066707:24037910:24053560::::

func sensitiveFiles(r *report.Report) {
	detectSecretFiles(r)
	detectBackupFiles(r)
	detectWordDocuments(r)
	detectSpreadsheets(r)
}

func detectSecretFiles(r *report.Report) {
	var sensitiveFilePatterns = []string{
		"*.asc",
		"*.crt",
		"*.gpg",
		"*.jks",
		"*.key",
		"*.p7b",
		"*.p7s",
		"*.pem",
		"*.pfx",
		"*.pgp",
	}

	var foundFiles = r.FancyListMatchUploadedFiles(sensitiveFilePatterns)

	if len(foundFiles) == 0 {
		return
	}

	issueDescription := "These files could contain secrets and should not be upload to Veracode for SAST scanning."
	issueText := fmt.Sprintf("A potentially sensitive or secret file was uploaded: \"%s\". %s", foundFiles[0], issueDescription)

	if len(foundFiles) > 1 {
		issueText = fmt.Sprintf(
			"%d potentially sensitive or secret files were uploaded: %s. %s",
			len(foundFiles),
			utils.Top5StringList(foundFiles), issueDescription)
	}

	r.ReportIssue(issueText, report.IssueSeverityHigh)
	r.MakeRecommendation("Do not upload any secrets, certificates, or key files.")
	r.MakeRecommendation("Do not upload unnecessary files.")
}

func detectBackupFiles(r *report.Report) {
	var sensitiveFilePatterns = []string{
		"*.bac", "*.back", "*.backup", "*.old", "*.orig", "*.bak",
	}

	var foundFiles = r.FancyListMatchUploadedFiles(sensitiveFilePatterns)

	if len(foundFiles) == 0 {
		return
	}

	issueDescription := "These files could contain secrets and should not be upload to Veracode for SAST scanning."
	issueText := fmt.Sprintf("A potentially sensitive backup/old/scratch file was uploaded: \"%s\". %s", foundFiles[0], issueDescription)

	if len(foundFiles) > 1 {
		issueText = fmt.Sprintf(
			"%d potentially sensitive backup/old/scratch files were uploaded: %s. %s",
			len(foundFiles),
			utils.Top5StringList(foundFiles), issueDescription)
	}

	r.ReportIssue(issueText, report.IssueSeverityHigh)
	r.MakeRecommendation("Do not upload backup, old or scratch files.")
	r.MakeRecommendation("Do not upload unnecessary files.")
}

func detectWordDocuments(r *report.Report) {
	var sensitiveFilePatterns = []string{
		"*.docx", "*.doc", "*.docm", "*.odt",
	}

	var foundFiles = r.FancyListMatchUploadedFiles(sensitiveFilePatterns)

	if len(foundFiles) == 0 {
		return
	}

	issueDescription := "These files could contain secrets and should not be upload to Veracode for SAST scanning."
	issueText := fmt.Sprintf("A Word document was uploaded: \"%s\". %s", foundFiles[0], issueDescription)

	if len(foundFiles) > 1 {
		issueText = fmt.Sprintf(
			"%d Word documents were uploaded: %s. %s",
			len(foundFiles),
			utils.Top5StringList(foundFiles), issueDescription)
	}

	r.ReportIssue(issueText, report.IssueSeverityHigh)
	r.MakeRecommendation("Office documents could contain sensitive information or secrets and should not be uploaded.")
	r.MakeRecommendation("Do not upload unnecessary files.")
}

func detectSpreadsheets(r *report.Report) {
	var sensitiveFilePatterns = []string{
		"*.xlsx", "*.xls", "*.xlsm", "*.ods",
	}

	var foundFiles = r.FancyListMatchUploadedFiles(sensitiveFilePatterns)

	if len(foundFiles) == 0 {
		return
	}

	issueDescription := "These files could contain secrets and should not be upload to Veracode for SAST scanning."
	issueText := fmt.Sprintf("A spreadsheet was uploaded: \"%s\". %s", foundFiles[0], issueDescription)

	if len(foundFiles) > 1 {
		issueText = fmt.Sprintf(
			"%d spreadsheets were uploaded: %s. %s",
			len(foundFiles),
			utils.Top5StringList(foundFiles), issueDescription)
	}

	r.ReportIssue(issueText, report.IssueSeverityHigh)
	r.MakeRecommendation("Office documents could contain sensitive information or secrets and should not be uploaded.")
	r.MakeRecommendation("Do not upload unnecessary files.")
}