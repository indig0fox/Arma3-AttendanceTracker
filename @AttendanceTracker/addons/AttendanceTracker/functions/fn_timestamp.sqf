// (parseSimpleArray ("AttendanceTracker" callExtension "getTimestamp")) select 0;

// need date for MySQL in format 2006-01-02 15:04:05

systemTimeUTC apply {if (_x < 10) then {"0" + str _x} else {str _x}} params [
	"_year",
	"_month",
	"_day",
	"_hour",
	"_minute",
	"_second",
	"_millisecond"
];

format[
	"%1-%2-%3T%4:%5:%6.000Z",
	_year,
	_month,
	_day,
	_hour,
	_minute,
	_second	
];

