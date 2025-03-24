package types

var ValidProjectTypes = map[string]bool{
	"ongoing":  true,
	"one-time": true,
}

var ValidProjectRates = map[string]bool{
	"hourly": true,
	"fixed":  true,
}

var ValidProjectLengths = map[string]bool{
	"<1":   true,
	"1-3":  true,
	"3-6":  true,
	"6-12": true,
	"12+":  true,
}

var ValidProjectHoursPerWeek = map[string]bool{
	"<10":   true,
	"10-20": true,
	"20-40": true,
	"40-60": true,
	"80+":   true,
}
