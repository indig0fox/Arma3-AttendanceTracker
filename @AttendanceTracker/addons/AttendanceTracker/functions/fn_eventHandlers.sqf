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


		[ // write d/c for past events
			"Server",
			_playerID,
			_playerUID,
			_profileName,
			_steamName
		] call attendanceTracker_fnc_writeDisconnect;

		// [
		// 	"Server",
		// 	_playerID,
		// 	_playerUID,
		// 	_profileName,
		// 	_steamName,
		// 	nil,
		// 	nil
		// ] call attendanceTracker_fnc_writeConnect;

		// start CBA PFH
		[format ["(EventHandler) OnUserConnected: Starting CBA PFH for %1", _playerID], "DEBUG"] call attendanceTracker_fnc_log;
		[
			{
				params ["_args", "_handle"];
				// check if player is still connected
				private _playerID = _args select 1;
				private _playerUID = _args select 2;
				if (allUsers find _playerID == -1) exitWith {
					[format ["(EventHandler) OnUserConnected: %1 (UID %2) is no longer connected, exiting CBA PFH", _playerUID], "DEBUG"] call attendanceTracker_fnc_log;
					_args call attendanceTracker_fnc_writeConnect;
					[_handle] call CBA_fnc_removePerFrameHandler;
				};

				_args call attendanceTracker_fnc_writeConnect;
			},
			missionNamespace getVariable ["AttendanceTracker_" + "dbupdateintervalseconds", 300],
			[
				"Server",
				_playerID,
				_playerUID,
				_profileName,
				_steamName,
				nil,
				nil
			]
		] call CBA_fnc_addPerFrameHandler;
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

		[
			"Server",
			_playerID,
			_playerUID,
			_profileName,
			_steamName
		] call attendanceTracker_fnc_writeConnect;
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


		[ // write d/c for past events
			"Mission",
			_playerID,
			_playerUID,
			_profileName,
			_steamName,
			_jip,
			nil
		] call attendanceTracker_fnc_writeDisconnect;

		// [
		// 	"Mission",
		// 	_playerID,
		// 	_playerUID,
		// 	_profileName,
		// 	_steamName,
		// 	_jip,
		// 	roleDescription _unit
		// ] call attendanceTracker_fnc_writeConnect;

		// start CBA PFH
		[format ["(EventHandler) PlayerConnected: Starting CBA PFH for %1", _playerID], "DEBUG"] call attendanceTracker_fnc_log;
		[
			{
				params ["_args", "_handle"];
				// check if player is still connected
				private _playerID = _args select 1;
				private _playerUID = _args select 2;
				private _userInfo = getUserInfo _playerID;
				private _clientStateNumber = _userInfo select 6;
				if (_clientStateNumber < 6) exitWith {
					[format ["(EventHandler) PlayerConnected: %1 (UID) is no longer connected to the mission, exiting CBA PFH", _playerID], "DEBUG"] call attendanceTracker_fnc_log;
					_args call attendanceTracker_fnc_writeConnect;
					[_handle] call CBA_fnc_removePerFrameHandler;
				};

				_args call attendanceTracker_fnc_writeConnect;
			},
			missionNamespace getVariable ["AttendanceTracker_" + "dbupdateintervalseconds", 300],
			[
				"Mission",
				_playerID,
				_playerUID,
				_profileName,
				_steamName,
				_jip,
				roleDescription _unit
			]
		] call CBA_fnc_addPerFrameHandler;
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
		
		[
			"Mission",
			_playerID,
			_playerUID,
			_profileName,
			_steamName,
			_jip,
			nil
		] call attendanceTracker_fnc_writeConnect;

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

		[
			"Server",
			_playerID,
			_playerUID,
			_profileName,
			_steamName,
			nil,
			nil
		] call attendanceTracker_fnc_writeConnect;

		[
			"Mission",
			_playerID,
			_playerUID,
			_profileName,
			_steamName,
			nil,
			nil
		] call attendanceTracker_fnc_writeConnect;
	}]
];