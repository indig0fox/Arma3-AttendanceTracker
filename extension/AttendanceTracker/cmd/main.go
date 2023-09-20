package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C" // This is required to import the C code

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/indig0fox/Arma3-AttendanceTracker/internal/db"
	"github.com/indig0fox/Arma3-AttendanceTracker/internal/logger"
	"github.com/indig0fox/Arma3-AttendanceTracker/internal/util"
	"github.com/indig0fox/a3go/a3interface"
	"github.com/indig0fox/a3go/assemblyfinder"
)

const EXTENSION_NAME string = "AttendanceTracker"
const ADDON_NAME string = "AttendanceTracker"
const EXTENSION_VERSION string = "0.9.0.1"

// file paths
const ATTENDANCE_TABLE string = "attendance"
const MISSIONS_TABLE string = "missions"
const WORLDS_TABLE string = "worlds"

var currentMissionID uint = 0

var RVExtensionChannels = map[string]chan string{
	":START:":        make(chan string),
	":MISSION:HASH:": make(chan string),
	":GET:SETTINGS:": make(chan string),
}

var RVExtensionArgsChannels = map[string]chan []string{
	":LOG:MISSION:":  make(chan []string),
	":LOG:PRESENCE:": make(chan []string),
}

var (
	modulePath    string
	modulePathDir string

	initSuccess bool // default false
)

// configure log output
func init() {

	a3interface.SetVersion(EXTENSION_VERSION)
	a3interface.RegisterRvExtensionChannels(RVExtensionChannels)
	a3interface.RegisterRvExtensionArgsChannels(RVExtensionArgsChannels)

	go func() {
		var err error

		modulePath = assemblyfinder.GetModulePath()
		// get absolute path of module path
		modulePathAbs, err := filepath.Abs(modulePath)
		if err != nil {
			panic(err)
		}
		modulePathDir = filepath.Dir(modulePathAbs)

		result, configErr := util.LoadConfig(modulePathDir)
		logger.InitLoggers(&logger.LoggerOptionsType{
			Path: filepath.Join(
				modulePathDir,
				fmt.Sprintf(
					"%s_v%s.log",
					EXTENSION_NAME,
					EXTENSION_VERSION,
				)),
			AddonName:     ADDON_NAME,
			ExtensionName: EXTENSION_NAME,
			Debug:         util.ConfigJSON.GetBool("armaConfig.debug"),
			Trace:         util.ConfigJSON.GetBool("armaConfig.traceLogToFile"),
		})
		if configErr != nil {
			logger.Log.Error().Err(configErr).Msgf(`Error loading config`)
			return
		} else {
			logger.Log.Info().Msgf(result)
		}

		logger.ArmaOnly.Info().Msgf(`%s v%s started`, EXTENSION_NAME, "0.0.0")
		logger.ArmaOnly.Info().Msgf(`Log path: %s`, logger.ActiveOptions.Path)

		db.SetConfig(db.ConfigStruct{
			MySQLHost:     util.ConfigJSON.GetString("sqlConfig.mysqlHost"),
			MySQLPort:     util.ConfigJSON.GetInt("sqlConfig.mysqlPort"),
			MySQLUser:     util.ConfigJSON.GetString("sqlConfig.mysqlUser"),
			MySQLPassword: util.ConfigJSON.GetString("sqlConfig.mysqlPassword"),
			MySQLDatabase: util.ConfigJSON.GetString("sqlConfig.mysqlDatabase"),
		})
		err = db.Connect()
		if err != nil {
			logger.Log.Error().Err(err).Msgf(`Error connecting to database`)
			return
		}
		err = db.Client().Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
			&World{},
			&Mission{},
			&Session{},
		)
		if err != nil {
			logger.Log.Error().Err(err).Msgf(`Error migrating database schema`)
		}

		startA3CallHandlers()

		initSuccess = true
		logger.RotateLogs()
		a3interface.WriteArmaCallback(
			EXTENSION_NAME,
			":READY:",
		)

		go finalizeUnendedSessions()
	}()
}

func startA3CallHandlers() error {
	go func() {
		for {
			select {
			case <-RVExtensionChannels[":START:"]:
				logger.Log.Trace().Msgf(`RVExtension :START: requested`)
				if !initSuccess {
					logger.Log.Warn().Msgf(`Received another :START: command before init was complete, ignoring.`)
					continue
				} else {
					logger.RotateLogs()
					a3interface.WriteArmaCallback(
						EXTENSION_NAME,
						":READY:",
					)
				}
			case <-RVExtensionChannels[":MISSION:HASH:"]:
				logger.Log.Trace().Msgf(`RVExtension :MISSION:HASH: requested`)
				timestamp, hash := getMissionHash()
				a3interface.WriteArmaCallback(
					EXTENSION_NAME,
					":MISSION:HASH:",
					timestamp,
					hash,
				)
			case <-RVExtensionChannels[":GET:SETTINGS:"]:
				logger.Log.Trace().Msg(`Settings requested`)
				armaConfig, err := util.ConfigArmaFormat()
				if err != nil {
					logger.Log.Error().Err(err).Msg(`Error when marshaling arma config`)
					continue
				}
				logger.Log.Trace().Str("armaConfig", armaConfig).Send()
				a3interface.WriteArmaCallback(
					EXTENSION_NAME,
					":GET:SETTINGS:",
					armaConfig,
				)
			case v := <-RVExtensionArgsChannels[":LOG:MISSION:"]:
				go func(data []string) {
					writeWorldInfo(v[1])
					writeMission(v[0])
				}(v)
			case v := <-RVExtensionArgsChannels[":LOG:PRESENCE:"]:
				go writeAttendance(v[0])
			}
		}
	}()

	return nil
}

// getMissionHash will return the current time in UTC and an md5 hash of that time
func getMissionHash() (sqlTime, hashString string) {
	// get md5 hash of string
	// https://stackoverflow.com/questions/2377881/how-to-get-a-md5-hash-from-a-string-in-golang

	nowTime := time.Now().UTC()
	// mysql format
	sqlTime = nowTime.Format("2006-01-02 15:04:05")
	hash := md5.Sum([]byte(sqlTime))
	hashString = fmt.Sprintf(`%x`, hash)

	return
}

// finalizeUnendedSessions will fill in the disconnect time for any sessions that have not been ended with a time 1 update interval after the join time
func finalizeUnendedSessions() {
	logger.Log.Debug().Msg("Filling missing disconnect events due to server restart.")
	// get all events with null DisconnectTime & set DisconnectTime
	var events []*Session
	db.Client().Model(&Session{}).
		Where("join_time_utc IS NOT NULL AND disconnect_time_utc IS NULL").
		Find(&events)
	for _, event := range events {
		// if difference between JoinTime and current time is greater than threshold, set to threshold
		if event.JoinTimeUTC.Time.Before(
			time.Now().Add(-1 * util.ConfigJSON.GetDuration("armaConfig.dbUpdateInterval")),
		) {
			// if more than the update interval has passed, set disconnect time as 1 interval after join
			event.DisconnectTimeUTC = sql.NullTime{
				Time:  event.JoinTimeUTC.Time.Add(util.ConfigJSON.GetDuration("armaConfig.dbUpdateInterval")),
				Valid: true,
			}
		} else {
			// otherwise, set disconnect time as now
			event.DisconnectTimeUTC = sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}
		}

		db.Client().Save(&event)
		if db.Client().Error != nil {
			logger.Log.Error().Err(db.Client().Error).Msgf(`Error when updating disconnect time for event %d`, event.ID)
		}
	}

	// log how many
	logger.Log.Info().Msgf(`Filled disconnect time of %d events.`, len(events))
}

func writeWorldInfo(worldInfo string) {
	// worldInfo is json, parse it
	var wi World
	fixedString := unescapeArmaQuotes(worldInfo)
	err := json.Unmarshal([]byte(fixedString), &wi)
	if err != nil {
		logger.Log.Error().Err(err).Msgf(`Error when unmarshalling world info`)
		return
	}

	// write world if not exist
	var dbWorld World
	db.Client().Where("world_name = ?", wi.WorldName).First(&dbWorld)
	if dbWorld.ID == 0 {
		db.Client().Create(&wi)
		if db.Client().Error != nil {
			logger.Log.Error().Err(db.Client().Error).Msgf(`Error when creating world`)
			return
		}
		logger.Log.Info().Msgf(`World %s created.`, wi.WorldName)
	} else {
		// don't do anything if exists
		logger.Log.Debug().Msgf(`World %s exists with ID %d.`, wi.WorldName, dbWorld.ID)
	}
}

func writeMission(missionJSON string) {
	var err error
	// writeLog(functionName, fmt.Sprintf(`["%s", "DEBUG"]`, Mission))
	// Mission is json, parse it
	var mi Mission
	fixedString := fixEscapeQuotes(trimQuotes(missionJSON))
	err = json.Unmarshal([]byte(fixedString), &mi)
	if err != nil {
		logger.Log.Error().Err(err).Msgf(`Error when unmarshalling mission`)
		return
	}

	// get world from WorldName
	var dbWorld World
	db.Client().Where("world_name = ?", mi.WorldName).First(&dbWorld)
	if dbWorld.ID == 0 {
		logger.Log.Error().Msgf(`World %s not found.`, mi.WorldName)
		return
	}

	mi.WorldID = dbWorld.ID

	// write mission to database
	db.Client().Create(&mi)
	if db.Client().Error != nil {
		logger.Log.Error().Err(db.Client().Error).Msgf(`Error when creating mission`)
		return
	}
	logger.Log.Info().Msgf(`Mission %s created with ID %d`, mi.MissionName, mi.ID)
	currentMissionID = mi.ID
}

func writeAttendance(data string) {
	var err error
	// data is json, parse it
	stringjson := unescapeArmaQuotes(data)
	var event Session
	err = json.Unmarshal([]byte(stringjson), &event)
	if err != nil {
		logger.Log.Error().Err(err).Msgf(`Error when unmarshalling attendance`)
		return
	}

	// search existing event
	var dbEvent Session
	db.Client().
		Where(
			"player_uid = ? AND mission_hash = ?",
			event.PlayerUID,
			event.MissionHash,
		).
		Order("join_time_utc desc").
		First(&dbEvent)
	if dbEvent.ID != 0 {
		// update disconnect time
		dbEvent.DisconnectTimeUTC = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		err = db.Client().Save(&dbEvent).Error
		if err != nil {
			logger.Log.Error().Err(err).
				Msgf(`Error when updating disconnect time for event %d`, dbEvent.ID)
			return
		}
		logger.Log.Debug().Msgf(`Attendance updated for %s (%s)`,
			dbEvent.ProfileName,
			dbEvent.PlayerUID,
		)
	} else {
		// insert new row
		event.JoinTimeUTC = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		if currentMissionID == 0 {
			logger.Log.Error().Msgf(`Current mission ID not set, cannot create attendance event`)
			return
		}
		event.MissionID = currentMissionID
		err = db.Client().Create(&event).Error
		if err != nil {
			logger.Log.Error().Err(err).Msgf(`Error when creating attendance event`)
			return
		}
		logger.Log.Debug().Msgf(`Attendance created for %s (%s)`,
			event.ProfileName,
			event.PlayerUID,
		)
	}
}

func getTimestamp() string {
	// get the current unix timestamp in nanoseconds
	// return time.Now().Local().Unix()
	return time.Now().Format("2006-01-02 15:04:05")
}

func trimQuotes(s string) string {
	// trim the start and end quotes from a string
	return strings.Trim(s, `"`)
}

func fixEscapeQuotes(s string) string {
	// fix the escape quotes in a string
	return strings.Replace(s, `""`, `"`, -1)
}

func unescapeArmaQuotes(s string) string {
	return fixEscapeQuotes(trimQuotes(s))
}

func main() {
	// loadConfig()
	// fmt.Println("Running DB connect/migrate to build schema...")
	// err := connectDB()
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("DB connect/migrate complete!")
	// }
	// fmt.Scanln()
}
