class CfgPatches {
	class AttendanceTracker {
		units[] = {};
		weapons[] = {};
		requiredVersion = 2.10;
		requiredAddons[] = {};
		author[] = {"IndigoFox"};
		authorUrl = "http://example.com";
	};
};

class CfgFunctions {
	class attendanceTracker {
		class functions {
			file = "\AttendanceTracker\functions";
			class postInit {postInit = 1;};
			class eventHandlers {};
			class callbackHandler {postInit = 1;};
			class log {};
			class logMissionEvent {};
			class logServerEvent {};
			class timestamp {};
		};
	};
};