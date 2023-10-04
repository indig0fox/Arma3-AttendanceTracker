addMissionEventHandler ["ExtensionCallback", {
	params ["_name", "_function", "_data"];
	if !(_name isEqualTo "AttendanceTracker") exitWith {};

	if (ATDebug && _function isNotEqualTo ":LOG:") then {
		diag_log format ["Raw callback: %1 _ %2", _function, _data];
	};

	_dataArr = parseSimpleArray _data;
	if (count _dataArr < 1) exitWith {};

	switch (_function) do {
		case ":LOG:": {
			diag_log formatText[
				"[Attendance Tracker] %1",
				_dataArr select 0
			];
		};
		default {
			[format["%1", _dataArr]] call attendanceTracker_fnc_log;
		};
	};
	true;
}];