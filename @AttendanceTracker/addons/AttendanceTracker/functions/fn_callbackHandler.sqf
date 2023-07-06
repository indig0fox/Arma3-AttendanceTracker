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

	if (missionNamespace getVariable ["AttendanceTracker_" + "debug", true]) then {
		diag_log format ["Raw callback: %1: %2", _function, _data];
	};

	// Parse response from string array
	private "_response";
	try {
		// diag_log format ["Raw callback: %1: %2", _function, _data];
		_response = parseSimpleArray _data;
		if (_response isEqualTo []) then {
			throw "Failed to parse response as array";
		};
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
			if (_response#0 == "SUCCESS") then {				
				// log world info
				private _response = "AttendanceTracker" callExtension [
					"logWorld",
					[
						[(call attendanceTracker_fnc_getWorldInfo)] call CBA_fnc_encodeJSON
					]
				];
				
				// log mission info and get back the row Id to send with future messages
				private _response = "AttendanceTracker" callExtension [
					"logMission",
					[
						[AttendanceTracker getVariable ["missionContext", createHashMap]] call CBA_fnc_encodeJSON
					]
				];

				missionNamespace setVariable ["AttendanceTracker_DBConnected", true];
			};
		};
		case "writeMission": {
			if (_response#0 == "MISSION_ID") then {
				AttendanceTracker_missionId = _response#1;
			};
		};
		default {
			_response call attendanceTracker_fnc_log;
		};
	};
	true;
}];