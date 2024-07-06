package scriptor

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"scriptor.test/scriptor/errs"
	"scriptor.test/scriptor/structs/event"
	"scriptor.test/scriptor/structs/record"
	"strconv"
	"sync"
	"time"
)

type Scriptor struct {
	m               sync.Mutex
	logMessages     chan *record.Record
	files           map[string]*os.File
	logType         map[string]*event.Event
	openedLoggers   map[string]map[string]*log.Logger
	paths           []string
	countScriptings int
}

var isRecordExist = false

func (scriptor *Scriptor) ThrowLog(file string, message string, t string) error {
	if (file != "stop") && (message != "stop") && (t != "stop") {
		scriptor.logMessages <- record.NewRecord(file, message, t)
		return nil
	} else {
		return errs.ErrReservedStopParametr
	}
}

func (scriptor *Scriptor) NewLog(name string, prefix string, flags []string) string {
	scriptor.logType[name] = event.NewEvent(name, prefix, flags)
	return name
}

func NewScriptor() (*Scriptor, error) {
	// при инициализации, анализирует структуру папок и создает словарь из папок названий и путей до них.
	// открывает коннекты до файлов и ждет с открытыми коннектами у актульных файлов
	ch := make(chan *record.Record)
	if isRecordExist {
		return nil, errs.ErrScriptorIsSinglengton
	}
	res, err := dirExists("logs")
	if err != nil {
		log.Fatal(err)
	}

	if !res {

		if err = os.Mkdir("logs", os.ModePerm); err != nil {
			log.Fatal(err)
		}

		scriptor := &Scriptor{
			logMessages:     ch,
			files:           make(map[string]*os.File),
			logType:         make(map[string]*event.Event),
			openedLoggers:   make(map[string]map[string]*log.Logger),
			paths:           []string{},
			countScriptings: 0,
		}
		isRecordExist = true

		go scriptor.startUpdating()
		go scriptor.StartScripting()
		return scriptor, nil
	}

	dirs, err := os.ReadDir("logs")
	if err != nil {
		log.Fatal(err)
	}

	paths := []string{}
	files := make(map[string]*os.File)

	for _, dir := range dirs {

		if dir.IsDir() {
			files[dir.Name()] = actualFile(dir.Name())
			paths = append(paths, path.Join("logs", dir.Name()))
		} else {
			continue
		}
	}

	scriptor := &Scriptor{
		logMessages:     ch,
		files:           files,
		paths:           paths,
		logType:         make(map[string]*event.Event),
		openedLoggers:   make(map[string]map[string]*log.Logger),
		countScriptings: 0,
	}
	isRecordExist = true

	go scriptor.startUpdating()
	go scriptor.StartScripting()
	return scriptor, nil
}

func (scriptor *Scriptor) StopScripting() error {
	if scriptor.countScriptings > 0 {
		scriptor.countScriptings -= 1
		scriptor.logMessages <- record.NewRecord("stop", "stop", "stop")
		return nil
	}
	return errs.ErrNoWorkingScriptors
}

func (scriptor *Scriptor) StartScripting() error {
	scriptor.countScriptings += 1
	for {
		select {
		case val := <-scriptor.logMessages:
			if val.Message() == "stop" {
				return nil
			}
			if _, ok := scriptor.files[val.FileName()]; !ok {
				var err error
				if err = os.Mkdir(filepath.Join("logs", val.FileName()), os.ModePerm); err != nil {
					log.Fatal(err)
				}

				scriptor.files[val.FileName()], err = os.Create(filepath.Join("logs", val.FileName(), generateTimeFileName(val.FileName(), time.Now())))
				if err != nil {
					log.Fatal(err)
				}
			}

			if _, ok := scriptor.openedLoggers[val.FileName()]; !ok {
				scriptor.openedLoggers[val.FileName()] = make(map[string]*log.Logger)
			}
			if _, ok := scriptor.openedLoggers[val.FileName()][val.LogType()]; !ok {

				if _, ok := scriptor.logType[val.LogType()]; ok {
					scriptor.openedLoggers[val.FileName()][val.LogType()] = log.New(scriptor.files[val.FileName()], scriptor.logType[val.LogType()].Prefix(), scriptor.logType[val.LogType()].Flag())
				} else {
					return errs.ErrLoggerNotFound(val.LogType())
				}
			}
			scriptor.openedLoggers[val.FileName()][val.LogType()].Println(val.Message())
		}
	}
}

func (scriptor *Scriptor) startUpdating() {
	for {
		time.Sleep(time.Duration(timeToNextDay()) * time.Second)
		scriptor.m.Lock()

		for path := range scriptor.files {
			var err error
			scriptor.files[path], err = os.Create(filepath.Join("logs", path, generateTimeFileName(path, time.Now())))
			if err != nil {
				log.Fatal(err)
			}
		}
		scriptor.m.Unlock()
	}
}

func timeToNextDay() int {
	currentTime := time.Now()
	nextDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, 0, 0, 0, 0, currentTime.Location())
	secondsToNextDay := int(nextDay.Sub(currentTime).Seconds())
	return secondsToNextDay
}

func actualFile(dir string) *os.File {
	actualFileName := generateTimeFileName(dir, time.Now())
	files, err := os.ReadDir(filepath.Join("logs", dir))

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name() == actualFileName {
			if !file.IsDir() {
				fileToReturn, err := os.OpenFile(filepath.Join("logs", dir, actualFileName), os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}
				return fileToReturn
			}
		}
	}

	file, err := os.Create(filepath.Join("logs", dir, actualFileName))
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func generateTimeFileName(name string, time time.Time) string {
	return name + "_" + strconv.Itoa(time.Day()) + "_" + time.Month().String() + "_" + strconv.Itoa(time.Year()) + ".txt"
}

func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
