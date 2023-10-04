addMissionEventHandler ["ExtensionCallback", {
	params ["_extension", "_function", "_data"];
	if !(_extension isEqualTo "AttendanceTracker") exitWith {};
	if !(_function isEqualTo ":GET:SETTINGS:") exitWith {};

	_dataArr = parseSimpleArray _data;
	diag_log format ["AT: Settings received: %1", _dataArr];
	if (count _dataArr < 1) exitWith {};

	private _settingsJSON = _dataArr select 0;
	private _settingsNamespace = [_settingsJSON] call CBA_fnc_parseJSON;
	{
		ATNamespace setVariable [_x, _settingsNamespace getVariable _x];
	} forEach (allVariables _settingsNamespace);
	ATDebug = ATNamespace getVariable "debug";
	ATUpdateDelay = ATNamespace getVariable "dbUpdateInterval";
	// remove last character (unit of time) and parse to number
	ATUpdateDelay = parseNumber (ATUpdateDelay select [0, count ATUpdateDelay - 1]);
	

	removeMissionEventHandler [
		"ExtensionCallback", 
		_thisEventHandler
	];
}];

"AttendanceTracker" callExtension ":GET:SETTINGS:";