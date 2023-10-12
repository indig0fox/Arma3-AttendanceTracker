#include "script_component.hpp"

params ["_id", "_uid", "_name", "_jip", "_owner", "_idstr"];

["DEBUG", format ["(EventHandler) PlayerConnected fired: %1", _this]] call FUNC(log);

if !(call FUNC(missionLoaded)) exitWith {
	["DEBUG", format ["(EventHandler) PlayerConnected: Server is in Mission Asked, likely mission selection state. Skipping.."]] call FUNC(log);
};

private _userInfo = (getUserInfo _idstr);
if ((count _userInfo) isEqualTo 0) exitWith {
	["DEBUG", format ["(EventHandler) PlayerConnected: No user info found for %1", _idstr]] call FUNC(log);
};

_userInfo params ["_playerID", "_ownerId", "_playerUID", "_profileName", "_displayName", "_steamName", "_clientState", "_isHC", "_adminState", "_networkInfo", "_unit"];
if (_isHC) exitWith {
	[
		"DEBUG",
		format [
			"(EventHandler) PlayerConnected: %1 is HC, skipping",
			_playerID
		]
	] call FUNC(log);
};

// start CBA PFH
[
	"DEBUG",
	format [
		"(EventHandler) PlayerConnected: Starting CBA PFH for %1",
		_playerID
	]
] call FUNC(log);

[{
	params ["_args", "_handle"];

	// every dbUpdateInterval, queue a wait for the mission to be logged
	// times out after 30 seconds
	// used to ensure joins at start of mission (during db connect) are logged
	[{GVAR(missionLogged)}, {
			// check if player is still connected
			private _hash = _this;
			private _clientStateNumber = 0;

			private _userInfo = getUserInfo (_hash get "playerId");
			if (_userInfo isEqualTo []) exitWith {
				["DEBUG", format ["(EventHandler) PlayerConnected: %1 (UID) is no longer connected to the mission, exiting CBA PFH", _hash get "playerUID"]] call FUNC(log);
				[_handle] call CBA_fnc_removePerFrameHandler;
			};

			_clientStateNumber = _userInfo select 6;

			if (_clientStateNumber < 6) exitWith {
				["DEBUG", format ["(EventHandler) PlayerConnected: %1 (UID) is no longer connected to the mission, exiting CBA PFH", _hash get "playerUID"]] call FUNC(log);
				[_handle] call CBA_fnc_removePerFrameHandler;
			};

			["DEBUG", format [
				"(EventHandler) PlayerConnected: %1 (UID) is connected to the mission, logging. data: %2", 
				_hash get "playerUID",
				_hash
			]] call FUNC(log);
			GVAR(extensionName) callExtension [
				":LOG:PRESENCE:", [
				_hash
			]];
		},
		_args, // args
		30 // timeout
	] call CBA_fnc_waitUntilAndExecute;
	},
	GVAR(updateInterval),
	(createHashMapFromArray [ // args
		["playerId", _playerID],
		["playerUID", _playerUID],
		["profileName", _profileName],
		["steamName", _steamName],
		["isJIP", _jip],
		["roleDescription", if (roleDescription _unit isEqualTo "") then {"None"} else {roleDescription _unit}],
		["missionHash", GVAR(missionHash)]
	])
] call CBA_fnc_addPerFrameHandler;