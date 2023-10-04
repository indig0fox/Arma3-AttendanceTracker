addMissionEventHandler ["ExtensionCallback", {
	params ["_extension", "_function", "_data"];
	if !(_extension isEqualTo "AttendanceTracker") exitWith {};
	if !(_function isEqualTo ":MISSION:HASH:") exitWith {};

	_dataArr = parseSimpleArray _data;
	if (count _dataArr < 1) exitWith {};

	_dataArr params ["_startTime", "_hash"];
	ATNamespace setVariable ["missionStartTime", call attendanceTracker_fnc_timestamp];
	ATNamespace setVariable ["missionHash", _hash];

	removeMissionEventHandler [
		"ExtensionCallback", 
		_thisEventHandler
	];
}];

"AttendanceTracker" callExtension ":MISSION:HASH:";