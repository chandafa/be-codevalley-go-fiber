package services

import (
	"log"
	"time"

	"code-valley-api/internal/models"
	"code-valley-api/internal/repositories"
	"code-valley-api/internal/websocket"
)

type GameClockService struct {
	worldRepo *repositories.WorldRepository
	ticker    *time.Ticker
	stopChan  chan bool
}

func NewGameClockService() *GameClockService {
	return &GameClockService{
		worldRepo: repositories.NewWorldRepository(),
		stopChan:  make(chan bool),
	}
}

func (s *GameClockService) Start() {
	// Update game time every 10 seconds (1 game minute = 10 real seconds)
	s.ticker = time.NewTicker(10 * time.Second)
	
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.updateGameTime()
			case <-s.stopChan:
				return
			}
		}
	}()
	
	log.Println("Game clock service started")
}

func (s *GameClockService) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.stopChan <- true
	log.Println("Game clock service stopped")
}

func (s *GameClockService) updateGameTime() {
	clock, err := s.worldRepo.GetGameClock()
	if err != nil {
		log.Printf("Failed to get game clock: %v", err)
		return
	}

	if clock.IsPaused {
		return
	}

	// Advance game time by 10 minutes
	clock.GameMinute += 10
	
	if clock.GameMinute >= 60 {
		clock.GameMinute = 0
		clock.GameHour++
		
		if clock.GameHour >= 24 {
			clock.GameHour = 0
			clock.GameDay++
			
			// Handle daily events
			s.handleDailyEvents(clock)
			
			if clock.GameDay > 28 {
				clock.GameDay = 1
				s.advanceSeason(clock)
			}
		}
	}

	clock.UpdatedAt = time.Now()
	
	if err := s.worldRepo.UpdateGameClock(clock); err != nil {
		log.Printf("Failed to update game clock: %v", err)
		return
	}

	// Broadcast time update to all clients
	websocket.GlobalHub.SendToAll(websocket.Message{
		Type: "time_update",
		Data: map[string]interface{}{
			"game_year":   clock.GameYear,
			"game_season": clock.GameSeason,
			"game_day":    clock.GameDay,
			"game_hour":   clock.GameHour,
			"game_minute": clock.GameMinute,
		},
	})

	// Handle hourly events
	if clock.GameMinute == 0 {
		s.handleHourlyEvents(clock)
	}
}

func (s *GameClockService) handleDailyEvents(clock *models.GameClock) {
	log.Printf("New day: Year %d, %s %d", clock.GameYear, clock.GameSeason, clock.GameDay)
	
	// Reset daily tasks, update code farms, etc.
	// This would trigger other services to handle daily resets
}

func (s *GameClockService) handleHourlyEvents(clock *models.GameClock) {
	// Handle NPC movement based on schedules
	// Update world objects
	// Process code farm growth
}

func (s *GameClockService) advanceSeason(clock *models.GameClock) {
	seasons := []string{"spring", "summer", "fall", "winter"}
	currentIndex := 0
	
	for i, season := range seasons {
		if season == clock.GameSeason {
			currentIndex = i
			break
		}
	}
	
	nextIndex := (currentIndex + 1) % len(seasons)
	clock.GameSeason = seasons[nextIndex]
	
	if clock.GameSeason == "spring" {
		clock.GameYear++
	}
	
	log.Printf("Season changed to %s, Year %d", clock.GameSeason, clock.GameYear)
	
	// Broadcast season change
	websocket.GlobalHub.SendToAll(websocket.Message{
		Type: "season_change",
		Data: map[string]interface{}{
			"new_season": clock.GameSeason,
			"game_year":  clock.GameYear,
		},
	})
}