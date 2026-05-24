package cli

import (
	"fmt"
	"log"

	"github.com/SmallLava/discord-leaving-bot/internal/database"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all leave records",
	Run: func(cmd *cobra.Command, args []string) {
		repo := database.NewLeaveRepository(database.DB)
		records, err := repo.GetAll()
		if err != nil {
			log.Fatalf("Error fetching records: %v", err)
		}

		fmt.Printf("%-5s | %-15s | %-10s | %-20s | %-20s | %s\n", "ID", "User", "Type", "Start", "End", "Reason")
		fmt.Println("----------------------------------------------------------------------------------------------------")
		for _, r := range records {
			fmt.Printf("%-5d | %-15s | %-10s | %-20s | %-20s | %s\n", 
				r.ID, r.UserName, r.Type, 
				r.StartDate.Format("2006-01-02 15:04"), 
				r.EndDate.Format("2006-01-02 15:04"), 
				r.Reason)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
