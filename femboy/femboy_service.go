package femboy

import (
	"fmt"
	"furrybot/config"
	"math/rand"
	"sort"
	"time"
)

type RateLimitError struct {
	TimeLeftMs int64
}

func (err *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit. avaialable in %d ms", err.TimeLeftMs)
}

type NoPlayersError struct{}

func (err *NoPlayersError) Error() string {
	return "there are no players in this game"
}

type FemboyGameService struct {
	players      []*FemboyGamePlayer
	lastFemboyMs int64
}

type FemboyGamePlayer struct {
	UserId int64
	Wins   int
}

func NewFemboyGameService() *FemboyGameService {
	return &FemboyGameService{make([]*FemboyGamePlayer, 0), 0}
}

func (fs *FemboyGameService) RegisterPlayer(userId int64) bool {
	for _, p := range fs.players {
		if p.UserId == userId {
			return false
		}
	}

	fs.players = append(fs.players, &FemboyGamePlayer{userId, 0})
	return true
}

func (fs *FemboyGameService) RemovePlayerByUserId(userId int64) {
	players := make([]*FemboyGamePlayer, 0, len(fs.players)-1)

	for i := 0; i < len(fs.players); i++ {
		if fs.players[i].UserId == userId {
			continue
		}

		players = append(players, fs.players[i])
	}

	fs.players = players
}

func (fs *FemboyGameService) PickWinner() (int64, error) {
	if len(fs.players) == 0 {
		return 0, &NoPlayersError{}
	}

	timeElapsed := time.Now().UTC().UnixMilli() - fs.lastFemboyMs

	if timeElapsed < config.Settings.FemboyCooldownMs {
		return 0, &RateLimitError{config.Settings.FemboyCooldownMs - timeElapsed}
	}

	winner := fs.players[rand.Intn(len(fs.players))]
	winner.Wins++
	fs.lastFemboyMs = time.Now().UTC().UnixMilli()

	return winner.UserId, nil
}

func (fs *FemboyGameService) GetSortedPlayerSlice() []*FemboyGamePlayer {
	sort.Slice(fs.players, func(i, j int) bool {
		return fs.players[i].Wins > fs.players[j].Wins
	})

	return fs.players
}
