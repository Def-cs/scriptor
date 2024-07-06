package event

import (
	"errors"
	"log"
	"scriptor.test/scriptor/configuration/constants"
	"scriptor.test/scriptor/errs"
)

type Event struct {
	name   string
	flag   int
	prefix string
}

func NewEvent(name string, prefix string, flags []string) *Event {
	return &Event{
		name:   name,
		prefix: prefix,
		flag:   generateFlag(flags),
	}
}

func generateFlag(flags []string) int {
	flag := 0
	for _, flagName := range flags {
		id, err := constants.FlagsList.MapElement(flagName)
		if err != nil {
			if errors.Is(err, errs.ErrElementNotFound(flagName)) {
				log.Fatal("This flag doesn't exist into package event.")
			} else {
				log.Fatal("Unexpected error:", err)
			}
		}
		flag |= id
	}
	return flag
}
