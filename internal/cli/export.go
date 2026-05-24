package cli

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/SmallLava/discord-leaving-bot/internal/database"
	"github.com/spf13/cobra"
)

var exportFormat string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export leave records to CSV or Markdown",
	Run: func(cmd *cobra.Command, args []string) {
		repo := database.NewLeaveRepository(database.DB)
		records, err := repo.GetAll()
		if err != nil {
			log.Fatalf("Error fetching records: %v", err)
		}

		switch exportFormat {
		case "csv":
			exportCSV(records)
		case "md":
			exportMarkdown(records)
		default:
			fmt.Println("Invalid format. Use 'csv' or 'md'.")
		}
	},
}

func exportCSV(records []database.LeaveRecord) {
	file, err := os.Create("export.csv")
	if err != nil {
		log.Fatalf("Cannot create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ID", "UserID", "UserName", "Type", "StartDate", "EndDate", "Reason"})
	for _, r := range records {
		writer.Write([]string{
			fmt.Sprintf("%d", r.ID),
			r.UserID,
			r.UserName,
			string(r.Type),
			r.StartDate.Format("2006-01-02"),
			r.EndDate.Format("2006-01-02"),
			r.Reason,
		})
	}
	fmt.Println("Exported to export.csv")
}

func exportMarkdown(records []database.LeaveRecord) {
	file, err := os.Create("export.md")
	if err != nil {
		log.Fatalf("Cannot create Markdown file: %v", err)
	}
	defer file.Close()

	fmt.Fprintln(file, "| ID | User | Type | Start | End | Reason |")
	fmt.Fprintln(file, "|----|------|------|-------|-----|--------|")
	for _, r := range records {
		fmt.Fprintf(file, "| %d | %s | %s | %s | %s | %s |\n",
			r.ID, r.UserName, r.Type,
			r.StartDate.Format("2006-01-02"),
			r.EndDate.Format("2006-01-02"),
			r.Reason)
	}
	fmt.Println("Exported to export.md")
}

func init() {
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "csv", "Export format (csv or md)")
	rootCmd.AddCommand(exportCmd)
}
