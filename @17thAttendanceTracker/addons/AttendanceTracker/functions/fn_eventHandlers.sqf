[
	["OnUserConnected", {
		params ["_networkId", "_clientStateNumber", "_clientState"];
		private _userInfo = (getUserInfo _networkId);
		_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit"];
		if (_isHC) exitWith {};

		[
			"ConnectedServer",
			_playerUID,
			_profileName,
			_steamName
		] call attendanceTracker_fnc_logServerEvent;

		(AttendanceTracker getVariable ["allUsers", createHashMap]) set [_networkId, _userInfo];

		[format ["(EventHandler) OnUserConnected fired: %1", _this], "DEBUG"] call attendanceTracker_fnc_log;
	}],
	["OnUserDisconnected", {
		params ["_networkId", "_clientStateNumber", "_clientState"];
		private _userInfo = (AttendanceTracker getVariable ["allUsers", createHashMap]) get [_networkId, nil];
		if (isNil "_userInfo") exitWith {};

		_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit"];
		if (_isHC) exitWith {};

		[
			"DisconnectedServer",
			_playerUID,
			_profileName,
			_steamName
		] call attendanceTracker_fnc_logServerEvent;

		[format ["(EventHandler) OnUserDisconnected fired: %1", _this], "DEBUG"] call attendanceTracker_fnc_log;
	}],
	["PlayerConnected", {
		params ["_id", "_uid", "_name", "_jip", "_owner", "_idstr"];
		private _userInfo = (getUserInfo _idstr);
		if (isNil "_userInfo") exitWith {};

		_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit"];

		(AttendanceTracker getVariable ["allUsers", createHashMap]) set [_playerID, _userInfo];

		if (_isHC) exitWith {};
		
		[
			"ConnectedMission",
			_playerUID,
			_profileName,
			_steamName,
			_jip,
			roleDescription _unit
		] call attendanceTracker_fnc_logMissionEvent;

		[format ["(EventHandler) PlayerConnected fired: %1", _this], "DEBUG"] call attendanceTracker_fnc_log;
	}],
	["HandleDisconnect", {
		params ["_unit", "_id", "_uid", "_name"];
		private _userInfo = (AttendanceTracker getVariable ["allUsers", createHashMap]) get [_id toFixed 0, nil];
		if (isNil "_userInfo") exitWith {};

		_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit"];

		if (_isHC) exitWith {};
		
		[
			"DisconnectedMission",
			_playerUID,
			_profileName,
			_steamName,
			_jip
		] call attendanceTracker_fnc_logMissionEvent;
		
		[format ["(EventHandler) HandleDisconnect fired: %1", _this], "DEBUG"] call attendanceTracker_fnc_log;
		false;
	}]
];