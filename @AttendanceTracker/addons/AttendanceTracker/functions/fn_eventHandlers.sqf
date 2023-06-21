[
	["OnUserConnected", {
		params ["_networkId", "_clientStateNumber", "_clientState"];

		[format ["(EventHandler) OnUserConnected fired: %1", _this], "DEBUG"] call attendanceTracker_fnc_log;

		private _userInfo = (getUserInfo _networkId);
		if (isNil "_userInfo") exitWith {
			[format ["(EventHandler) OnUserConnected: No user info found for %1", _networkId], "DEBUG"] call attendanceTracker_fnc_log;
		};

		_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit"];
		if (_isHC) exitWith {
			[format ["(EventHandler) OnUserConnected: %1 is HC, skipping", _playerID], "DEBUG"] call attendanceTracker_fnc_log;
		};

		(AttendanceTracker getVariable ["allUsers", createHashMap]) set [_networkId, _userInfo];
		(AttendanceTracker getVariable ["rowIds", createHashMap]) set [_networkId, [nil, nil]]; // reset rowId on connect
		[
			"Server",
			_playerID,
			_playerUID,
			_profileName,
			_steamName,
			nil // send rowId on d/c only
		] call attendanceTracker_fnc_logServerEvent;

	}],
	["OnUserDisconnected", {
		params ["_networkId", "_clientStateNumber", "_clientState"];

		[format ["(EventHandler) OnUserDisconnected fired: %1", _this], "DEBUG"] call attendanceTracker_fnc_log;

		if !(call attendanceTracker_fnc_missionLoaded) exitWith {
			[format ["(EventHandler) OnUserDisconnected: Server is in Mission Asked, likely mission selection state. Skipping.."], "DEBUG"] call attendanceTracker_fnc_log;
		};

		private _userInfo = (AttendanceTracker getVariable ["allUsers", createHashMap]) get _networkId;
		if (isNil "_userInfo") exitWith {
			[format ["(EventHandler) OnUserDisconnected: No user info found for %1", _networkId], "DEBUG"] call attendanceTracker_fnc_log;
		};

		_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit"];
		if (_isHC) exitWith {
			[format ["(EventHandler) OnUserDisconnected: %1 is HC, skipping", _playerID], "DEBUG"] call attendanceTracker_fnc_log;
		};

		private _rowId = ((AttendanceTracker getVariable ["rowIds", createHashMap]) getOrDefault [
			_networkId,
			[nil, nil]
		]) select 0;

		[
			"Server",
			_playerID,
			_playerUID,
			_profileName,
			_steamName,
			(if (!isNil "_rowId") then {_rowId} else {nil}) // send rowId on d/c only
		] call attendanceTracker_fnc_logServerEvent;
	}],
	["PlayerConnected", {
		params ["_id", "_uid", "_name", "_jip", "_owner", "_idstr"];

		[format ["(EventHandler) PlayerConnected fired: %1", _this], "DEBUG"] call attendanceTracker_fnc_log;

		if !(call attendanceTracker_fnc_missionLoaded) exitWith {
			[format ["(EventHandler) PlayerConnected: Server is in Mission Asked, likely mission selection state. Skipping.."], "DEBUG"] call attendanceTracker_fnc_log;
		};

		private _userInfo = (getUserInfo _idstr);
		if (isNil "_userInfo") exitWith {
			[format ["(EventHandler) PlayerConnected: No user info found for %1", _idstr], "DEBUG"] call attendanceTracker_fnc_log;
		};

		_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit"];
		if (_isHC) exitWith {
			[format ["(EventHandler) PlayerConnected: %1 is HC, skipping", _playerID], "DEBUG"] call attendanceTracker_fnc_log;
		};

		(AttendanceTracker getVariable ["allUsers", createHashMap]) set [_playerID, _userInfo];
		[
			"Mission",
			_playerID,
			_playerUID,
			_profileName,
			_steamName,
			_jip,
			roleDescription _unit,
			nil // send rowId on d/c only
		] call attendanceTracker_fnc_logMissionEvent;
	}],
	["PlayerDisconnected", {
		// NOTE: HandleDisconnect returns a DIFFERENT _id than PlayerDisconnected and above handlers, so we can't use it here
		params ["_id", "_uid", "_name", "_jip", "_owner", "_idstr"];

		[format ["(EventHandler) HandleDisconnect fired: %1", _this], "DEBUG"] call attendanceTracker_fnc_log;

		if !(call attendanceTracker_fnc_missionLoaded) exitWith {
			[format ["(EventHandler) HandleDisconnect: Server is in Mission Asked, likely mission selection state. Skipping.."], "DEBUG"] call attendanceTracker_fnc_log;
		};

		private _userInfo = (AttendanceTracker getVariable ["allUsers", createHashMap]) get _idstr;
		if (isNil "_userInfo") exitWith {
			[format ["(EventHandler) HandleDisconnect: No user info found for %1", _idstr], "DEBUG"] call attendanceTracker_fnc_log;
		};

		_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit", "_rowId"];
		if (_isHC) exitWith {
			[format ["(EventHandler) HandleDisconnect: %1 is HC, skipping", _playerID], "DEBUG"] call attendanceTracker_fnc_log;
		};

		private _rowId = ((AttendanceTracker getVariable ["rowIds", createHashMap]) getOrDefault [
			_idstr,
			[nil, nil]
		]) select 1;
		
		[
			"Mission",
			_playerID,
			_playerUID,
			_profileName,
			_steamName,
			_jip,
			nil,
			(if (!isNil "_rowId") then {_rowId} else {nil}) // send rowId on d/c only
		] call attendanceTracker_fnc_logMissionEvent;
		

		false;
	}],
	["OnUserKicked", {
		params ["_networkId", "_kickTypeNumber", "_kickType", "_kickReason", "_kickMessageIncReason"];

		[format ["(EventHandler) OnUserKicked fired: %1", _this], "DEBUG"] call attendanceTracker_fnc_log;

		if !(call attendanceTracker_fnc_missionLoaded) exitWith {
			[format ["(EventHandler) OnUserKicked: Server is in Mission Asked, likely mission selection state. Skipping.."], "DEBUG"] call attendanceTracker_fnc_log;
		};

		private _userInfo = (AttendanceTracker getVariable ["allUsers", createHashMap]) get _networkId;
		if (isNil "_userInfo") exitWith {
			[format ["(EventHandler) OnUserKicked: No user info found for %1", _networkId], "DEBUG"] call attendanceTracker_fnc_log;
		};

		_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit"];

		if (_isHC) exitWith {
			[format ["(EventHandler) OnUserKicked: %1 is HC, skipping", _playerID], "DEBUG"] call attendanceTracker_fnc_log;
		};

		private _rowId = ((AttendanceTracker getVariable ["rowIds", createHashMap]) getOrDefault [
			_networkId,
			[nil, nil]
		]) select 0;

		[
			"Server",
			_playerID,
			_playerUID,
			_profileName,
			_steamName,
			(if (!isNil "_rowId") then {_rowId} else {nil}) // send rowId on d/c only
		] call attendanceTracker_fnc_logServerEvent;


		private _rowId = ((AttendanceTracker getVariable ["rowIds", createHashMap]) getOrDefault [
			_networkId,
			[nil, nil]
		]) select 1;
		[
			"Mission",
			_playerID,
			_playerUID,
			_profileName,
			_steamName,
			nil // send rowId on d/c only
		] call attendanceTracker_fnc_logMissionEvent;
	}]
];