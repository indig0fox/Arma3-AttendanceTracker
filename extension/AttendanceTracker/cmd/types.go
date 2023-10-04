package main

import (
	"database/sql"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type World struct {
	gorm.Model
	Author            string  `json:"author"`
	WorkshopID        string  `json:"workshopID"`
	DisplayName       string  `json:"displayName"`
	WorldName         string  `json:"worldName"`
	WorldNameOriginal string  `json:"worldNameOriginal"`
	WorldSize         float32 `json:"worldSize"`
	Latitude          float32 `json:"latitude"`
	Longitude         float32 `json:"longitude"`
	Missions          []Mission
}

type Mission struct {
	gorm.Model
	MissionName       string    `json:"missionName"`
	BriefingName      string    `json:"briefingName"`
	MissionNameSource string    `json:"missionNameSource"`
	OnLoadName        string    `json:"onLoadName"`
	Author            string    `json:"author"`
	ServerName        string    `json:"serverName"`
	ServerProfile     string    `json:"serverProfile"`
	MissionStart      time.Time `json:"missionStart" gorm:"index"`
	MissionHash       string    `json:"missionHash" gorm:"index"`
	WorldName         string    `json:"worldName" gorm:"-"`
	WorldID           uint
	World             World `gorm:"foreignkey:WorldID"`
	Attendees         []Session
}

func (m *Mission) UnmarshalJSON(data []byte) error {
	type Alias Mission
	aux := &struct {
		*Alias
		MissionStart string `json:"missionStart"`
	}{Alias: (*Alias)(m)}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	m.MissionStart, err = time.Parse(time.RFC3339, aux.MissionStart)
	if err != nil {
		return err
	}
	return nil
}

type Session struct {
	ID                uint         `json:"id" gorm:"primaryKey"`
	PlayerUID         string       `json:"playerUID" gorm:"index;primaryKey"`
	MissionHash       string       `json:"missionHash"`
	PlayerId          string       `json:"playerId"`
	JoinTimeUTC       sql.NullTime `json:"joinTimeUTC" gorm:"index"`
	DisconnectTimeUTC sql.NullTime `json:"disconnectTimeUTC" gorm:"index"`
	ProfileName       string       `json:"profileName"`
	SteamName         string       `json:"steamName"`
	IsJIP             bool         `json:"isJIP" gorm:"column:is_jip"`
	RoleDescription   string       `json:"roleDescription"`
	MissionID         uint
	Mission           Mission `gorm:"foreignkey:MissionID"`
}

func (s *Session) UnmarshalJSON(data []byte) error {
	type Alias Session
	aux := &struct {
		*Alias
		JoinTimeUTC       string `json:"joinTimeUTC"`
		DisconnectTimeUTC string `json:"disconnectTimeUTC"`
	}{Alias: (*Alias)(s)}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	if aux.JoinTimeUTC != "" {
		s.JoinTimeUTC.Time, err = time.Parse(time.RFC3339, aux.JoinTimeUTC)
		if err != nil {
			return err
		}
		s.JoinTimeUTC.Valid = true
	}
	if aux.DisconnectTimeUTC != "" {
		s.DisconnectTimeUTC.Time, err = time.Parse(time.RFC3339, aux.DisconnectTimeUTC)
		if err != nil {
			return err
		}
		s.DisconnectTimeUTC.Valid = true
	}
	return nil
}
