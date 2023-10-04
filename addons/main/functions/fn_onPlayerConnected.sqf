params ["_id", "_uid", "_name", "_jip", "_owner", "_idstr"];

[format ["(EventHandler) PlayerConnected fired: %1", _this], "DEBUG"] call attendanceTracker_fnc_log;

if !(call attendanceTracker_fnc_missionLoaded) exitWith {
	[format ["(EventHandler) PlayerConnected: Server is in Mission Asked, likely mission selection state. Skipping.."], "DEBUG"] call attendanceTracker_fnc_log;
};

private _userInfo = (getUserInfo _idstr);
if ((count _userInfo) isEqualTo 0) exitWith {
	[format ["(EventHandler) PlayerConnected: No user info found for %1", _idstr], "DEBUG"] call attendanceTracker_fnc_log;
};

_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit"];
if (_isHC) exitWith {
	[
		format [
			"(EventHandler) PlayerConnected: %1 is HC, skipping",
			_playerID
		],
		"DEBUG"
	] call attendanceTracker_fnc_log;
};

// start CBA PFH
[
	format [
		"(EventHandler) PlayerConnected: Starting CBA PFH for %1",
		_playerID
	], 
	"DEBUG"
] call attendanceTracker_fnc_log;

[
	{
		params ["_args", "_handle"];
		// check if player is still connected
		_args params ["_playerID", "_playerUID", "_profileName", "_steamName", "_jip", "_roleDescription"];
		private _userInfo = getUserInfo _playerID;
		private _clientStateNumber = 0;
		if (_userInfo isEqualTo []) exitWith {
			[_handle] call CBA_fnc_removePerFrameHandler;
		};

		_clientStateNumber = _userInfo select 6;

		if (_clientStateNumber < 6) exitWith {
			[format ["(EventHandler) PlayerConnected: %1 (UID) is no longer connected to the mission, exiting CBA PFH", _playerID], "DEBUG"] call attendanceTracker_fnc_log;
			[_handle] call CBA_fnc_removePerFrameHandler;
		};

		_args call attendanceTracker_fnc_writePlayer;
	},
	ATUpdateDelay,
	[
		_playerID,
		_playerUID,
		_profileName,
		_steamName,
		_jip,
		roleDescription _unit
	]
] call CBA_fnc_addPerFrameHandler;