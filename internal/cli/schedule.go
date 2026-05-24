package cli

import (
	"fmt"
	"log"

	"github.com/SmallLava/discord-leaving-bot/internal/database"
	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage meeting schedules",
}

var scheduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all meeting schedules",
	Run: func(cmd *cobra.Command, args []string) {
		repo := database.NewScheduleRepository(database.DB)
		schedules, err := repo.GetAll()
		if err != nil {
			log.Fatalf("Error fetching schedules: %v", err)
		}

		fmt.Printf("%-5s | %-10s | %-10s | %-10s\n", "ID", "Day", "Start", "End")
		fmt.Println("-------------------------------------------")
		for _, s := range schedules {
			fmt.Printf("%-5d | %-10s | %-10s | %-10s\n", s.ID, s.DayOfWeek.String(), s.StartTime, s.EndTime)
		}
	},
}

func init() {
	scheduleCmd.AddCommand(scheduleListCmd)
	rootCmd.AddCommand(scheduleCmd)
}
