addMissionEventHandler ["ExtensionCallback", {
	params ["_name", "_function", "_data"];
	if !(_name == "AttendanceTracker") exitWith {};

	// Validate data param
	if (isNil "_data") then {_data = ""};

	if (_data isEqualTo "") exitWith {
		[
			format ["Callback empty data: %1", _function],
			"WARN"
		] call attendanceTracker_fnc_log;
		false;
	};

	if (typeName _data != "STRING") exitWith {
		[
			format ["Callback invalid data: %1: %2", _function, _data],
			"WARN"
		] call attendanceTracker_fnc_log;
		false;
	};

	// Parse response from string array
	private "_response";
	try {
		// diag_log format ["Raw callback: %1: %2", _function, _data];
		_response = parseSimpleArray _data;
	} catch {
		[
			format ["Callback invalid data: %1: %2: %3", _function, _data, _exception],
			"WARN"
		] call attendanceTracker_fnc_log;
	};

	if (isNil "_response") exitWith {false};

	switch (_function) do {
		case "connectDB": {
			systemChat format ["AttendanceTracker: %1", _response#0];
			[_response#0, _response#1, _function] call attendanceTracker_fnc_log;
		};
		default {
			_response call attendanceTracker_fnc_log;
		};
	};
	true;
}];