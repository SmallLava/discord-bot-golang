package botinit

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/SmallLava/discord-leaving-bot/internal/database"
	"github.com/SmallLava/discord-leaving-bot/internal/service"
	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "測試機器人是否存活",
		},
		{
			Name:        "status",
			Description: "查詢個人請假紀錄",
		},
		{
			Name:        "set-announcement",
			Description: "設定此頻道為開會請假通知頻道",
		},
		{
			Name:        "leave",
			Description: "申請請假",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "type",
					Description: "請假類型",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "單次開會", Value: string(database.LeaveTypeSingle)},
						{Name: "時段請假", Value: string(database.LeaveTypeInterval)},
					},
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "start_date",
					Description:  "開始日期 (格式: YYYY-MM-DD，單次請假可輸入或選擇建議時間)",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "end_date",
					Description: "結束日期 (格式: YYYY-MM-DD, 單次請假可不填)",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "請假原因",
					Required:    false,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong! 🏓",
				},
			})
		},
	}
)

func slashCommand(bot *discordgo.Session) func() {
	// Initialize Service
	leaveRepo := database.NewLeaveRepository(database.DB)
	scheduleRepo := database.NewScheduleRepository(database.DB)
	configRepo := database.NewConfigRepository(database.DB)
	leaveService := service.NewLeaveService(leaveRepo, scheduleRepo, configRepo)

	commandHandlers["set-announcement"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := leaveService.SetAnnouncementChannel(i.ChannelID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "設定失敗: " + err.Error(),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "成功！此頻道已設定為開會請假通知頻道。",
			},
		})
	}

	commandHandlers["status"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		records, err := leaveService.GetUserLeaves(i.Member.User.ID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "查詢失敗: " + err.Error(),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		if len(records) == 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "您目前沒有請假紀錄。",
				},
			})
			return
		}

		content := "您的請假紀錄如下：\n"
		for _, r := range records {
			content += fmt.Sprintf("- [%s] %s ~ %s (%s)\n", 
				r.Type, 
				r.StartDate.Format("2006-01-02"), 
				r.EndDate.Format("2006-01-02"), 
				r.Reason)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})
	}

	commandHandlers["leave"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		leaveType := database.LeaveType(optionMap["type"].StringValue())
		startStr := optionMap["start_date"].StringValue()
		reason := ""
		if opt, ok := optionMap["reason"]; ok {
			reason = opt.StringValue()
		}

		start, err := service.ParseDate(startStr)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "日期格式錯誤，請使用 YYYY-MM-DD",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		end := start
		if opt, ok := optionMap["end_date"]; ok {
			end, err = service.ParseDate(opt.StringValue())
			if err != nil {
				// Handle end date error if needed
			}
		}

		err = leaveService.AddLeaveRecord(i.Member.User.ID, i.Member.User.Username, leaveType, start, end, reason)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "儲存失敗: " + err.Error(),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "請假申請成功！",
			},
		})
	}

	bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionApplicationCommandAutocomplete:
			data := i.ApplicationCommandData()
			if data.Name == "leave" {
				// Find focused option
				var focusedOpt *discordgo.ApplicationCommandInteractionDataOption
				for _, opt := range data.Options {
					if opt.Focused {
						focusedOpt = opt
						break
					}
				}

				if focusedOpt != nil && focusedOpt.Name == "start_date" {
					meetings, _ := leaveService.GetNextNMeetings(4)
					choices := []*discordgo.ApplicationCommandOptionChoice{}
					for _, m := range meetings {
						dateStr := m.Format("2006-01-02")
						choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
							Name:  fmt.Sprintf("%s (%s)", dateStr, m.Weekday().String()),
							Value: dateStr,
						})
					}
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: 8, // InteractionResponseApplicationCommandAutocompleteResult
						Data: &discordgo.InteractionResponseData{
							Choices: choices,
						},
					})
				}
			}
		}
	})

	var registeredCommands []*discordgo.ApplicationCommand
	guildID := os.Getenv("GUILD_ID")

	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		for _, v := range commands {
			cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, v)
			if err != nil {
				log.Printf("Cannot create '%v' command: %v", v.Name, err)
				continue
			}
			registeredCommands = append(registeredCommands, cmd)
		}
		log.Println("Slash commands registered.")

		// Start background scheduler
		go startScheduler(s, leaveService)
	})

	return func() {
		for _, v := range registeredCommands {
			err := bot.ApplicationCommandDelete(bot.State.User.ID, guildID, v.ID)
			if err != nil {
				log.Printf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}

func startScheduler(s *discordgo.Session, leaveService *service.LeaveService) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("Background scheduler started (Optimized).")

	// Cache schedules in memory and refresh every hour
	var cachedSchedules []database.MeetingSchedule
	var lastRefresh time.Time

	refreshSchedules := func() {
		scheds, err := leaveService.GetAllSchedules()
		if err == nil {
			cachedSchedules = scheds
			lastRefresh = time.Now()
			// Only log during refresh to keep console clean
			log.Printf("Meeting schedules cached (%d entries).", len(cachedSchedules))
		}
	}

	// Initial load
	refreshSchedules()

	for range ticker.C {
		now := time.Now()
		// Refresh cache once an hour
		if now.Sub(lastRefresh) > time.Hour {
			refreshSchedules()
		}

		currentTimeStr := now.Format("15:04")
		// Log a heartbeat every 30 minutes to show scheduler is alive
		if now.Minute() % 30 == 0 && now.Second() < 60 {
			log.Printf("Scheduler heartbeat: current time %s, day %s", currentTimeStr, now.Weekday())
		}

		for _, sched := range cachedSchedules {
			if now.Weekday() == sched.DayOfWeek && currentTimeStr == sched.StartTime {
				log.Printf("[MATCH] Meeting time detected: %s. Fetching data...", currentTimeStr)
				
				channelID, err := leaveService.GetAnnouncementChannel()
				if err != nil {
					log.Printf("Error getting announcement channel: %v", err)
					continue
				}
				if channelID == "" {
					log.Println("No announcement channel set, skipping notification.")
					continue
				}
				log.Printf("Using announcement channel: %s", channelID)

				userIDs, err := leaveService.GetMembersOnLeave(now)
				if err != nil {
					log.Printf("Error fetching members on leave: %v", err)
					continue
				}
				log.Printf("Found %d members on leave for today.", len(userIDs))

				if len(userIDs) > 0 {
					mentions := []string{}
					for _, id := range userIDs {
						mentions = append(mentions, fmt.Sprintf("<@%s>", id))
					}
					message := fmt.Sprintf("@silent %s 以上成員請假", strings.Join(mentions, " "))
					log.Printf("Sending message: %s", message)
					_, err = s.ChannelMessageSend(channelID, message)
					if err != nil {
						log.Printf("Error sending announcement: %v", err)
					} else {
						log.Println("Announcement sent successfully.")
					}
				} else {
					log.Println("No members on leave today, no message sent.")
				}
			}
		}
	}
}
