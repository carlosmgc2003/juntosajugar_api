package models

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

var (
	InvalidMeetingPlace    = errors.New("Gamemeeting Model: Invalid meeting place")
	InvalidMeetingSchedule = errors.New("Gamemeeting Model: Invalid meeting time")
	InvalidMeetingOwner    = errors.New("Gamemeeting Model: Invalid Owner ID")
	InvalidBoardgame       = errors.New("Gamemeeting Model: Invalid Boardgame ID")
	InvalidMaxPlayers      = errors.New("Gamemeeting Model: Invalid Max Players Quantity")
)

func (G *Gamemeeting) FromJson(requestBody []byte, db *gorm.DB) error {
	type tempStruct struct {
		Id         uint   `json:"id"`
		Place      string `json:"place"`
		Scheduled  string `json:"scheduled"`
		Owner      uint   `json:"owner"`
		Game       uint   `json:"game"`
		MaxPlayers uint   `json:"max_players"`
	}
	var temp tempStruct
	err := json.Unmarshal(requestBody, &temp)
	if err != nil {
		return err
	}
	G.ID = temp.Id
	G.Place = temp.Place
	G.Scheduled, err = time.Parse(time.RFC3339, temp.Scheduled)
	G.OwnerId = temp.Owner
	G.GameId = temp.Game
	G.MaxPlayers = temp.MaxPlayers
	return err
}

func validMeetingPlace(meetingplace string) bool {
	var characters = len(meetingplace)
	return characters <= 100 && characters >= 5
}

func validScheduledTime(scheduledtime time.Time) bool {
	return time.Now().Before(scheduledtime)
}

func validOwner(ownerid int, db *gorm.DB) bool {
	var owner User
	result := db.First(&owner, ownerid)
	return errors.Is(result.Error, gorm.ErrRecordNotFound)
}

func validGame(gameid int, db *gorm.DB) bool {
	var game Boardgame
	result := db.First(&game, gameid)
	return errors.Is(result.Error, gorm.ErrRecordNotFound)
}

func validMaxPlayers(maxplayers int) bool {
	return maxplayers <= 12
}