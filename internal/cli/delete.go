package cli

import (
	"fmt"
	"log"
	"strconv"

	"github.com/SmallLava/discord-leaving-bot/internal/database"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a leave record by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Fatalf("Invalid ID: %v", err)
		}

		repo := database.NewLeaveRepository(database.DB)
		err = repo.Delete(uint(id))
		if err != nil {
			log.Fatalf("Failed to delete record: %v", err)
		}

		fmt.Printf("Record %d deleted successfully.\n", id)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
