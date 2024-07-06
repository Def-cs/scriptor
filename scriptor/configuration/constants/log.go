package constants

import (
	"log"
	"scriptor.test/scriptor/configuration"
)

//type constantsMapStrings map[string]string

var FlagsList = configuration.NewConstantsMapIntsStorage(map[string]int{
	"Ldate":         log.Ldate,
	"Ltime":         log.Ltime,
	"Lmicroseconds": log.Lmicroseconds,
	"Llongfile":     log.Llongfile,
	"Lshortfile":    log.Lshortfile,
	"LUTC":          log.LUTC,
	"Lmsgprefix":    log.Lmsgprefix,
	"LstdFlags":     log.LstdFlags,
},
)
