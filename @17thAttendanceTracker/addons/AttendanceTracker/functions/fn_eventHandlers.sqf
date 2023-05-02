[
	["OnUserConnected", {
		params ["_networkId", "_clientStateNumber", "_clientState"];

		[format ["(EventHandler) OnUserConnected fired: %1", _this], "DEBUG"] call attendanceTracker_fnc_log;

		private _userInfo = (getUserInfo _networkId);
		_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit"];
		if (_isHC) exitWith {};

		[
			"ConnectedServer",
			_playerID,
			_playerUID,
			_profileName,
			_steamName
		] call attendanceTracker_fnc_logServerEvent;

		(AttendanceTracker getVariable ["allUsers", createHashMap]) set [_networkId, _userInfo];
	}],
	["OnUserDisconnected", {
		params ["_networkId", "_clientStateNumber", "_clientState"];

		[format ["(EventHandler) OnUserDisconnected fired: %1", _this], "DEBUG"] call attendanceTracker_fnc_log;

		private _userInfo = (AttendanceTracker getVariable ["allUsers", createHashMap]) get _networkId;
		if (isNil "_userInfo") exitWith {};

		_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit"];
		if (_isHC) exitWith {};

		[
			"DisconnectedServer",
			_playerID,
			_playerUID,
			_profileName,
			_steamName
		] call attendanceTracker_fnc_logServerEvent;
	}],
	["PlayerConnected", {
		params ["_id", "_uid", "_name", "_jip", "_owner", "_idstr"];

		[format ["(EventHandler) PlayerConnected fired: %1", _this], "DEBUG"] call attendanceTracker_fnc_log;

		private _userInfo = (getUserInfo _idstr);
		if (isNil "_userInfo") exitWith {};

		_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit"];

		(AttendanceTracker getVariable ["allUsers", createHashMap]) set [_playerID, _userInfo];

		if (_isHC) exitWith {};
		
		[
			"ConnectedMission",
			_playerID,
			_playerUID,
			_profileName,
			_steamName,
			_jip,
			roleDescription _unit
		] call attendanceTracker_fnc_logMissionEvent;
	}],
	["PlayerDisconnected", {
		// NOTE: HandleDisconnect returns a DIFFERENT _id than PlayerDisconnected and above handlers, so we can't use it here
		params ["_id", "_uid", "_name", "_jip", "_owner", "_idstr"];

		[format ["(EventHandler) HandleDisconnect fired: %1", _this], "DEBUG"] call attendanceTracker_fnc_log;

		private _userInfo = (AttendanceTracker getVariable ["allUsers", createHashMap]) get _idstr;
		if (isNil "_userInfo") exitWith {
			[format ["(EventHandler) HandleDisconnect: No user info found for %1", _idstr], "DEBUG"] call attendanceTracker_fnc_log;
		};

		_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit"];

		if (_isHC) exitWith {};
		
		[
			"DisconnectedMission",
			_playerID,
			_playerUID,
			_profileName,
			_steamName,
			_jip
		] call attendanceTracker_fnc_logMissionEvent;
		

		false;
	}]
];