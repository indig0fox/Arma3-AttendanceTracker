freeExtension "AttendanceTracker";
"AttendanceTracker" callExtension "connectDB";
sleep 2;
"attendanceTracker" callExtension ["logAttendance", ["{""playerUID"": ""76561197991996737"", ""roleDescription"": ""NULL"", ""missionNameSource"": ""aaaltisaiatk"", ""eventType"": ""ConnectedMission"", ""briefingName"": ""aaaltisaiatk"", ""profileName"": ""IndigoFox"", ""serverName"": ""IndigoFox on DESKTOP-6B2U0AT"", ""steamName"": ""IndigoFox"", ""onLoadName"": ""NULL"", ""missionName"": ""aaaltisaiatk"", ""isJIP"": false, ""author"": ""IndigoFox"", ""missionStart"": ""1682549469590908300""}"]]
"attendanceTracker" callExtension ["logAttendance", ["{""playerUID"": ""76561197991996737"", ""roleDescription"": ""NULL"", ""missionNameSource"": ""aaaltisaiatk"", ""eventType"": ""ConnectedServer"", ""briefingName"": ""aaaltisaiatk"", ""profileName"": ""IndigoFox"", ""serverName"": ""IndigoFox on DESKTOP-6B2U0AT"", ""steamName"": ""IndigoFox"", ""onLoadName"": ""NULL"", ""missionName"": ""aaaltisaiatk"", ""isJIP"": false, ""author"": ""IndigoFox"", ""missionStart"": ""1682549469590908300""}"]]
sleep 15;
exit;