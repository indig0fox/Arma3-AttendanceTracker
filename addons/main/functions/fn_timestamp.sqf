// (parseSimpleArray ("AttendanceTracker" callExtension "getTimestamp")) select 0;

// const time.RFC3339 untyped string = "2006-01-02T15:04:05Z07:00"

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
	"%1-%2-%3T%4:%5:%6Z",
	_year,
	_month,
	_day,
	_hour,
	_minute,
	_second	
];

