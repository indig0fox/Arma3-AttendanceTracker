#include "..\script_mod.hpp"

if (!isServer) exitWith {};

ATNamespace = false call CBA_fnc_createNamespace;
ATDebug = true;
"AttendanceTracker" callExtension ":START:";


// we'll wait for the asynchronous init steps of the extension to finish, to confirm we have a DB connection and our config was loaded. If there are errors with either, the extension won't reply and initiate further during this mission.
addMissionEventHandler ["ExtensionCallback", {
	params ["_name", "_function", "_data"];
	if !(_name isEqualTo "AttendanceTracker") exitWith {};
	if !(_function isEqualTo ":READY:") exitWith {};

	call attendanceTracker_fnc_getMissionHash;
	call attendanceTracker_fnc_getSettings;

	[
		{// wait until settings have been loaded from extension
			!isNil {ATNamespace getVariable "missionHash"} &&
			!isNil {ATDebug}
		},
		{

			// get world and mission context
			ATNamespace setVariable [
				"worldContext",
				call attendanceTracker_fnc_getWorldInfo
			];
			ATNamespace setVariable [
				"missionContext", 
				call attendanceTracker_fnc_getMissionInfo
			];

			// write them to establish DB rows
			"AttendanceTracker" callExtension [
				":LOG:MISSION:",
				[
					[ATNamespace getVariable "missionContext"] call CBA_fnc_encodeJSON,
					[ATNamespace getVariable "worldContext"] call CBA_fnc_encodeJSON
				]
			];

			// add player connected (to mission) handler
			addMissionEventHandler ["PlayerConnected", {
				_this call attendanceTracker_fnc_onPlayerConnected;
			}];
		},
		[],
		10, // 10 second timeout
		{ // timeout code
			["Failed to load settings", "ERROR"] call attendanceTracker_fnc_log;
		}
	] call CBA_fnc_waitUntilAndExecute;

	removeMissionEventHandler [
		"ExtensionCallback",
		_thisEventHandler
	];
}];