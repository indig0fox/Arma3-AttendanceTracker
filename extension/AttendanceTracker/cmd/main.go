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
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/indig0fox/Arma3-AttendanceTracker/internal/db"
	"github.com/indig0fox/Arma3-AttendanceTracker/internal/logger"
	"github.com/indig0fox/Arma3-AttendanceTracker/internal/util"
	"github.com/indig0fox/a3go/a3interface"
	"github.com/indig0fox/a3go/assemblyfinder"
	"github.com/rs/zerolog"
)

const EXTENSION_NAME string = "AttendanceTracker"
const ADDON_NAME string = "AttendanceTracker"

// file paths
const ATTENDANCE_TABLE string = "attendance"
const MISSIONS_TABLE string = "missions"
const WORLDS_TABLE string = "worlds"

var (
	EXTENSION_VERSION string = "DEVELOPMENT"

	modulePath    string
	modulePathDir string

	loadedMission *Mission
	loadedWorld   *World
)

// configure log output
func init() {

	a3interface.SetVersion(EXTENSION_VERSION)
	a3interface.NewRegistration(":START:").
		SetFunction(onStartCommand).
		SetRunInBackground(false).
		Register()

	a3interface.NewRegistration(":MISSION:HASH:").
		SetFunction(onMissionHashCommand).
		SetRunInBackground(false).
		Register()

	a3interface.NewRegistration(":GET:SETTINGS:").
		SetFunction(onGetSettingsCommand).
		SetRunInBackground(false).
		Register()

	a3interface.NewRegistration(":LOG:MISSION:").
		SetDefaultResponse(`Logging mission data`).
		SetArgsFunction(onLogMissionArgsCommand).
		SetRunInBackground(true).
		Register()

	a3interface.NewRegistration(":LOG:PRESENCE:").
		SetDefaultResponse(`Logging presence data`).
		SetArgsFunction(onLogPresenceArgsCommand).
		SetRunInBackground(true).
		Register()

	go func() {
		var err error

		modulePath = assemblyfinder.GetModulePath()
		modulePathDir = filepath.Dir(modulePath)

		result, configErr := util.LoadConfig(modulePathDir)
		logger.InitLoggers(&logger.LoggerOptionsType{
			Path: filepath.Join(
				modulePathDir,
				fmt.Sprintf(
					"%s_v%s.log",
					EXTENSION_NAME,
					EXTENSION_VERSION,
				)),
			AddonName:        ADDON_NAME,
			ExtensionName:    EXTENSION_NAME,
			ExtensionVersion: EXTENSION_VERSION,
			Debug:            util.ConfigJSON.GetBool("armaConfig.debug"),
			Trace:            util.ConfigJSON.GetBool("armaConfig.trace"),
		})
		logger.RotateLogs()
		if configErr != nil {
			logger.Log.Error().Err(configErr).Msgf(`Error loading config`)
			return
		} else {
			logger.Log.Info().Msgf(result)
		}

		logger.Log.Info().Msgf(`%s v%s started`, EXTENSION_NAME, EXTENSION_VERSION)
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

		logger.Log.Info().
			Str("dialect", db.Client().Dialector.Name()).
			Str("database", db.Client().Migrator().CurrentDatabase()).
			Str("host", util.ConfigJSON.GetString("sqlConfig.mysqlHost")).
			Int("port", util.ConfigJSON.GetInt("sqlConfig.mysqlPort")).
			Msgf(`Connected to database`)

		err = db.Client().Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
			&World{},
			&Mission{},
			&Session{},
		)
		if err != nil {
			logger.Log.Error().Err(err).Msgf(`Error migrating database schema`)
		} else {
			logger.Log.Info().Msgf(`Database schema migrated`)
		}

		a3interface.WriteArmaCallback(
			EXTENSION_NAME,
			":READY:",
		)

		go finalizeUnendedSessions()
	}()
}

func onStartCommand(
	ctx a3interface.ArmaExtensionContext,
	data string,
) (string, error) {
	logger.Log.Debug().Msgf(`RVExtension :START: requested`)
	loadedWorld = nil
	loadedMission = nil
	return fmt.Sprintf(
		`["%s v%s started"]`,
		EXTENSION_NAME,
		EXTENSION_VERSION,
	), nil
}

func onMissionHashCommand(
	ctx a3interface.ArmaExtensionContext,
	data string,
) (string, error) {
	logger.Log.Debug().Msgf(`RVExtension :MISSION:HASH: requested`)
	timestamp, hash := getMissionHash()
	return fmt.Sprintf(
		`[%q, %q]`,
		timestamp,
		hash,
	), nil
}

func onGetSettingsCommand(
	ctx a3interface.ArmaExtensionContext,
	data string,
) (string, error) {
	logger.Log.Debug().Msg(`RVExtension :GET:SETTINGS: requested`)
	// get arma config
	c := util.ConfigJSON.Get("armaConfig")
	armaConfig := a3interface.ToArmaHashMap(c)
	return fmt.Sprintf(
		`[%s]`,
		armaConfig,
	), nil
}

func onLogMissionArgsCommand(
	ctx a3interface.ArmaExtensionContext,
	command string,
	args []string,
) (string, error) {
	thisLogger := logger.Log.With().Str("command", command).Interface("ctx", ctx).Logger()
	thisLogger.Debug().Msgf(`RVExtension :LOG:MISSION: requested`)
	var err error
	world, err := writeWorldInfo(args[0], thisLogger)
	if err != nil {
		return ``, err
	}
	loadedWorld = &world

	mission, err := writeMission(args[1], thisLogger)
	if err != nil {
		return ``, err
	}
	loadedMission = &mission

	a3interface.WriteArmaCallback(
		EXTENSION_NAME,
		":LOG:MISSION:SUCCESS:",
	)

	return ``, nil
}

func onLogPresenceArgsCommand(
	ctx a3interface.ArmaExtensionContext,
	command string,
	args []string,
) (string, error) {
	thisLogger := logger.Log.With().Str("command", command).Interface("ctx", ctx).Logger()
	thisLogger.Debug().Msgf(`RVExtension :LOG:PRESENCE: requested`)
	writeAttendance(args[0], thisLogger)
	return ``, nil
}

// getMissionHash will return the current time in UTC and an md5 hash of that time
func getMissionHash() (sqlTime, hashString string) {
	// get md5 hash of string
	// https://stackoverflow.com/questions/2377881/how-to-get-a-md5-hash-from-a-string-in-golang

	nowTime := time.Now().UTC()
	// mysql format
	sqlTime = nowTime.Format(time.RFC3339)
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

func writeWorldInfo(worldInfo string, thisLogger zerolog.Logger) (World, error) {

	parsedInterface, err := a3interface.ParseSQF(worldInfo)
	if err != nil {
		thisLogger.Error().Err(err).Msgf(`Error when parsing world info`)
		return World{}, err
	}

	parsedMap, err := a3interface.ParseSQFHashMap(parsedInterface)
	if err != nil {
		thisLogger.Error().Err(err).Msgf(`Error when parsing world info`)
		return World{}, err
	}

	thisLogger.Trace().Msgf(`parsedMap: %+v`, parsedMap)

	// create world object from map[string]interface{}
	var wi = World{}
	worldBytes, err := json.Marshal(parsedMap)
	if err != nil {
		thisLogger.Error().Err(err).Msgf(`Error when marshalling world info`)
		return World{}, err
	}
	err = json.Unmarshal(worldBytes, &wi)
	if err != nil {
		thisLogger.Error().Err(err).Msgf(`Error when unmarshalling world info`)
		return World{}, err
	}

	thisLogger.Trace().Msgf(`World info: %+v`, wi)

	var dbWorld World
	db.Client().Where("world_name = ?", wi.WorldName).First(&dbWorld)
	// if world exists, use it
	if dbWorld.ID > 0 {
		thisLogger.Debug().Msgf(`World %s exists with ID %d.`, wi.WorldName, dbWorld.ID)
		return dbWorld, nil
	}

	// write world if not exist
	db.Client().Create(&wi)
	if db.Client().Error != nil {
		thisLogger.Error().Err(db.Client().Error).Msgf(`Error when creating world`)
		return World{}, db.Client().Error
	}
	thisLogger.Info().Msgf(`World %s created.`, wi.WorldName)

	return wi, nil
}

func writeMission(data string, thisLogger zerolog.Logger) (Mission, error) {
	var err error
	parsedInterface, err := a3interface.ParseSQF(data)
	if err != nil {
		thisLogger.Error().Err(err).Msgf(`Error when parsing mission info`)
		return Mission{}, err
	}

	parsedMap, err := a3interface.ParseSQFHashMap(parsedInterface)
	if err != nil {
		thisLogger.Error().Err(err).Msgf(`Error when parsing mission info`)
		return Mission{}, err
	}

	thisLogger.Trace().Msgf(`parsedMap: %+v`, parsedMap)

	var mi Mission
	// create mission object from map[string]interface{}
	missionBytes, err := json.Marshal(parsedMap)
	if err != nil {
		thisLogger.Error().Err(err).Msgf(`Error when marshalling mission info`)
		return Mission{}, err
	}
	err = json.Unmarshal(missionBytes, &mi)
	if err != nil {
		thisLogger.Error().Err(err).Msgf(`Error when unmarshalling mission info`)
		return Mission{}, err
	}

	if loadedWorld == nil {
		thisLogger.Error().Msgf(`Current world ID not set, cannot create mission`)
		return Mission{}, err
	}
	if loadedWorld.ID == 0 {
		thisLogger.Error().Msgf(`Current world ID is 0, cannot create mission`)
		return Mission{}, err
	}
	mi.WorldID = loadedWorld.ID

	// write mission to database
	db.Client().Create(&mi)
	if db.Client().Error != nil {
		thisLogger.Error().Err(db.Client().Error).Msgf(`Error when creating mission`)
		return Mission{}, db.Client().Error
	}
	thisLogger.Info().Msgf(`Mission %s created with ID %d`, mi.MissionName, mi.ID)

	a3interface.WriteArmaCallback(
		EXTENSION_NAME,
		":LOG:MISSION:SUCCESS:",
		"World and mission logged successfully.",
	)

	return mi, nil
}

func writeAttendance(data string, thisLogger zerolog.Logger) {
	var err error

	parsedInterface, err := a3interface.ParseSQF(data)
	if err != nil {
		thisLogger.Error().Err(err).Str("data", data).Msgf(`Error when parsing attendance info`)
		return
	}

	parsedMap, err := a3interface.ParseSQFHashMap(parsedInterface)
	if err != nil {
		thisLogger.Error().Err(err).Str("data", data).Msgf(`Error when parsing attendance info`)
		return
	}

	thisLogger.Trace().Msgf(`parsedMap: %+v`, parsedMap)

	var thisSession Session
	// create session object from map[string]interface{}
	sessionBytes, err := json.Marshal(parsedMap)
	if err != nil {
		thisLogger.Error().Err(err).Str("data", data).Msgf(`Error when marshalling attendance info`)
		return
	}

	err = json.Unmarshal(sessionBytes, &thisSession)
	if err != nil {
		thisLogger.Error().Err(err).Str("data", data).Msgf(`Error when unmarshalling attendance info`)
		return
	}

	thisLogger2 := thisLogger.With().
		Str("playerId", thisSession.PlayerId).
		Str("playerUID", thisSession.PlayerUID).
		Str("profileName", thisSession.ProfileName).
		Logger()

	// search existing event
	var dbEvent Session

	db.Client().
		Where(
			"player_id = ? AND mission_hash = ?",
			thisSession.PlayerId,
			thisSession.MissionHash,
		).
		Order("join_time_utc desc").
		First(&dbEvent)

	if dbEvent.ID > 0 {
		// update disconnect time
		dbEvent.DisconnectTimeUTC = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		err = db.Client().Save(&dbEvent).Error
		if err != nil {
			thisLogger2.Error().Err(err).
				Msgf(`Error when updating disconnect time for event %d`, dbEvent.ID)
			return
		}
		thisLogger2.Debug().Msgf(`Attendance updated with ID %d`,
			dbEvent.ID,
		)
	} else {
		// insert new row
		thisSession.JoinTimeUTC = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		if loadedMission == nil {
			thisLogger2.Error().Msgf(`Current mission ID not set, cannot create attendance event`)
			return
		}
		thisSession.MissionID = loadedMission.ID
		err = db.Client().Create(&thisSession).Error
		if err != nil {
			thisLogger2.Error().Err(err).Msgf(`Error when creating attendance event`)
			return
		}
		thisLogger2.Info().Msgf(`Attendance created with ID %d`,
			thisSession.ID,
		)
	}
}

func getTimestamp() string {
	// get the current unix timestamp in nanoseconds
	// return time.Now().Local().Unix()
	return time.Now().Format("2006-01-02 15:04:05")
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
