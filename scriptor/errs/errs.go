package errs

import "errors"

var (
	ErrReservedStopParametr = errors.New("It's reserved parameters. Use method StopScripting() instead")
	ErrNoWorkingScriptors   = errors.New("No scriptings working")
	ErrLoggerNotFound       = func(name string) error {
		return errors.New("Logger not found: " + name)
	}
	ErrElementNotFound = func(name string) error {
		return errors.New("Element not found: " + name)
	}
	ErrScriptorIsSinglengton = errors.New("You can create only one scriptor")
)
