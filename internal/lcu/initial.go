package lcu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	lcu "github.com/its-haze/lcu-gopher"
	"github.com/its-haze/league-rpc/internal/discord"
	"github.com/its-haze/league-rpc/internal/state"
	"github.com/its-haze/league-rpc/pkg/constants"
	"github.com/its-haze/league-rpc/pkg/types"
)

// gatherInitialData fetches all initial state from LCU when first connecting.
func (c *Client) gatherInitialData() error {
	c.logger.Debug().Msg("Gathering initial data from League Client...")

	// Get current summoner info
	if err := c.fetchSummonerInfo(); err != nil {
		c.logger.Warn().Err(err).Msg("Failed to fetch summoner info")
	}

	// Get chat status (online/away)
	if err := c.fetchChatStatus(); err != nil {
		c.logger.Warn().Err(err).Msg("Failed to fetch chat status")
	}

	// Get TFT companion data
	if err := c.fetchTFTCompanion(); err != nil {
		c.logger.Debug().Err(err).Msg("Failed to fetch TFT companion (may not be selected)")
	}

	// Get ranked stats
	if err := c.fetchRankedStats(); err != nil {
		c.logger.Warn().Err(err).Msg("Failed to fetch ranked stats")
	}

	// Get current gameflow phase
	if err := c.fetchGameFlowPhase(); err != nil {
		c.logger.Warn().Err(err).Msg("Failed to fetch gameflow phase")
	}

	// Get lobby info (if in lobby)
	if err := c.fetchLobbyInfo(); err != nil {
		c.logger.Debug().Err(err).Msg("Not in lobby or failed to fetch lobby info")
		c.recoverGameModeFromLiveClientData()
	}

	c.logger.Debug().Msg("Initial data gathering complete")
	return nil
}

// recoverGameModeFromLiveClientData sets GameMode (and, for Practice Tool,
func (c *Client) recoverGameModeFromLiveClientData() {
	if c.liveGame == nil || c.state.Get().GameFlowPhase != types.GameFlowInProgress {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rawMode, err := c.liveGame.GameMode(ctx)
	if err != nil || rawMode == "" {
		c.logger.Debug().Err(err).Msg("Could not recover game mode from Live Client Data")
		return
	}

	isPractice := rawMode == "PRACTICETOOL"

	c.state.UpdateField(func(s *state.State) {
		s.GameMode = types.GameMode(rawMode)
		if isPractice {
			s.QueueName = "Practice Tool"
			s.IsPractice = true
		}
	})

	c.logger.Info().Str("game_mode", rawMode).Bool("is_practice", isPractice).
		Msg("Recovered game mode from Live Client Data (no lobby available)")
}

// fetchSummonerInfo gets the current summoner's profile
func (c *Client) fetchSummonerInfo() error {
	summoner, err := c.lcu.GetCurrentSummoner()
	if err != nil {
		return fmt.Errorf("get current summoner: %w", err)
	}

	c.state.UpdateSummonerIcon(summoner.ProfileIconID)
	c.logger.Debug().Int("icon_id", summoner.ProfileIconID).Msg("Updated summoner icon")

	return nil
}

// fetchChatStatus gets the current online/away status
func (c *Client) fetchChatStatus() error {
	resp, err := c.lcu.Get(constants.EndpointChatStatus)
	if err != nil {
		return fmt.Errorf("get chat status: %w", err)
	}
	defer resp.Body.Close()

	var chatStatus struct {
		Availability string `json:"availability"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&chatStatus); err != nil {
		return fmt.Errorf("decode chat status: %w", err)
	}

	c.state.UpdateAvailability(chatStatus.Availability)
	c.logger.Debug().Str("availability", chatStatus.Availability).Msg("Updated availability")

	return nil
}

// fetchTFTCompanion gets the selected TFT companion
func (c *Client) fetchTFTCompanion() error {
	resp, err := c.lcu.Get(constants.EndpointTFTCompanion)
	if err != nil {
		return fmt.Errorf("get TFT companion: %w", err)
	}
	defer resp.Body.Close()

	var companionData struct {
		SelectedLoadoutItem struct {
			ItemID       int    `json:"itemId"`
			LoadoutsIcon string `json:"loadoutsIcon"`
			Name         string `json:"name"`
			Description  string `json:"description"`
		} `json:"selectedLoadoutItem"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&companionData); err != nil {
		return fmt.Errorf("decode TFT companion: %w", err)
	}

	item := companionData.SelectedLoadoutItem

	// Extract icon path from full path
	// Example: "ASSETS/Characters/TFT/..." → "characters/tft/..."
	fullPath := item.LoadoutsIcon
	parts := strings.Split(fullPath, "ASSETS/")
	if len(parts) < 2 {
		return fmt.Errorf("invalid icon path format: %s", fullPath)
	}
	finalPath := strings.ToLower(parts[1])

	iconURL := discord.GetTFTCompanionURL(finalPath)

	c.state.UpdateTFTCompanion(item.ItemID, iconURL, item.Name, item.Description)
	c.logger.Debug().Str("name", item.Name).Msg("Updated TFT companion")

	return nil
}

// fetchRankedStats gets the current ranked statistics
func (c *Client) fetchRankedStats() error {
	resp, err := c.lcu.Get(constants.EndpointRankedStats)
	if err != nil {
		return fmt.Errorf("get ranked stats: %w", err)
	}
	defer resp.Body.Close()

	var rankedData struct {
		Queues []struct {
			QueueType    string `json:"queueType"`
			Tier         string `json:"tier"`
			Division     string `json:"division"`
			LeaguePoints int    `json:"leaguePoints"`
			RatedTier    string `json:"ratedTier"`    // For Arena
			RatedRating  int    `json:"ratedRating"`  // For Arena
		} `json:"queues"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rankedData); err != nil {
		return fmt.Errorf("decode ranked stats: %w", err)
	}

	// Parse ranked stats for each queue
	currentState := c.state.Get()
	for _, queue := range rankedData.Queues {
		switch queue.QueueType {
		case "RANKED_SOLO_5x5":
			currentState.SummonerRank = types.RankedStats{
				Tier:         types.Tier(queue.Tier),
				Division:     types.Division(queue.Division),
				LeaguePoints: queue.LeaguePoints,
			}

		case "RANKED_FLEX_SR":
			currentState.SummonerRankFlex = types.RankedStats{
				Tier:         types.Tier(queue.Tier),
				Division:     types.Division(queue.Division),
				LeaguePoints: queue.LeaguePoints,
			}

		case "RANKED_TFT":
			currentState.TFTRank = types.RankedStats{
				Tier:         types.Tier(queue.Tier),
				Division:     types.Division(queue.Division),
				LeaguePoints: queue.LeaguePoints,
			}

		case "CHERRY":
			currentState.ArenaRank = types.ArenaStats{
				RatedTier:   types.ArenaRatedTier(queue.RatedTier),
				Tier:        types.MapRatedTierToTier(types.ArenaRatedTier(queue.RatedTier)),
				RatedRating: queue.RatedRating,
			}
		}
	}

	c.state.Update(currentState)
	c.logger.Debug().Msg("Updated ranked stats")

	return nil
}

// fetchGameFlowPhase gets the current game flow phase
func (c *Client) fetchGameFlowPhase() error {
	resp, err := c.lcu.Get(constants.EndpointGameflow)
	if err != nil {
		return fmt.Errorf("get gameflow phase: %w", err)
	}
	defer resp.Body.Close()

	var phase string
	if err := json.NewDecoder(resp.Body).Decode(&phase); err != nil {
		return fmt.Errorf("decode gameflow phase: %w", err)
	}

	c.state.UpdateGameFlowPhase(phase)
	c.logger.Debug().Str("phase", phase).Msg("Updated gameflow phase")

	return nil
}

// fetchLobbyInfo gets current lobby information (if in lobby). The LCU
func (c *Client) fetchLobbyInfo() error {
	resp, err := c.lcu.Get(constants.EndpointLobby)
	if err != nil {
		return fmt.Errorf("get lobby: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("not in a lobby (status %d)", resp.StatusCode)
	}

	var lobbyData lcu.Lobby
	if err := json.NewDecoder(resp.Body).Decode(&lobbyData); err != nil {
		return fmt.Errorf("decode lobby: %w", err)
	}

	// Update lobby state
	c.updateLobbyState(&lobbyData)

	return nil
}
