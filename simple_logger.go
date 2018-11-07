// Copyright 2018 Foo Coders (www.foocoders.io).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"fmt"
	"log"
	"os"
)

// SimpleLogger is the standar logger implementation
type SimpleLogger struct {
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	level         int
	slackWeebHook string
	errorsToSlack bool
}

// NewSimpleLogger creates a new instance of SimpleLogger.
func NewSimpleLogger() (sl *SimpleLogger) {
	sl = new(SimpleLogger)
	sl.level = LevelInfo

	sl.debugLogger = log.New(os.Stdout, "[DBG] ", log.LstdFlags)
	sl.infoLogger = log.New(os.Stdout, "[INF] ", log.LstdFlags)
	sl.errorLogger = log.New(os.Stdout, "[ERR] ", log.LstdFlags)
	return
}

// Debug prints the arguments to the debug logger.
func (sl *SimpleLogger) Debug(v ...interface{}) {
	if sl.GetLevel() >= LevelDebug {
		sl.debugLogger.Print(v...)
	}
}

// Debugf prints the arguments to the debug logger. Arguments are handled like in fmt.Printf.
func (sl *SimpleLogger) Debugf(format string, v ...interface{}) {
	if sl.GetLevel() >= LevelDebug {
		sl.debugLogger.Printf(format, v...)
	}
}

// Info prints the arguments to the info logger.
func (sl *SimpleLogger) Info(v ...interface{}) {
	if sl.GetLevel() >= LevelInfo {
		sl.infoLogger.Print(v...)
	}
}

// Infof prints the arguments to the info logger. Arguments are handled like in fmt.Printf.
func (sl *SimpleLogger) Infof(format string, v ...interface{}) {
	if sl.GetLevel() >= LevelInfo {
		sl.infoLogger.Printf(format, v...)
	}
}

// Error prints the arguments to the error logger.
func (sl *SimpleLogger) Error(v ...interface{}) {
	sl.errorLogger.Print(v...)

	if sl.errorsToSlack && sl.slackWeebHook != "" {
		text := fmt.Sprint(v...)
		go SendAlert(sl.slackWeebHook, "", "Error", ColorWarning, fmt.Sprintf("`%s`", text))
	}
}

// Errorf prints the arguments to the error logger. Arguments are handled like in fmt.Printf.
func (sl *SimpleLogger) Errorf(format string, v ...interface{}) {
	sl.errorLogger.Printf(format, v...)

	if sl.errorsToSlack && sl.slackWeebHook != "" {
		text := fmt.Sprintf(format, v...)
		go SendAlert(sl.slackWeebHook, "", "Error", ColorWarning, fmt.Sprintf("`%s`", text))
	}
}

// Fatal prints the arguments to the error logger, followed by a call to os.Exit(1).
func (sl *SimpleLogger) Fatal(v ...interface{}) {
	sl.errorLogger.Fatal(v...)

	if sl.errorsToSlack && sl.slackWeebHook != "" {
		text := fmt.Sprint(v...)
		go SendAlert(sl.slackWeebHook, "", "Fatal", ColorDanger, fmt.Sprintf("`%s`", text))
	}
}

// Fatalf prints the arguments to the error logger, followed by a call to os.Exit(1).
// Arguments are handled like in fmt.Printf.
func (sl *SimpleLogger) Fatalf(format string, v ...interface{}) {
	sl.infoLogger.Fatalf(format, v...)

	if sl.errorsToSlack && sl.slackWeebHook != "" {
		text := fmt.Sprintf(format, v...)
		go SendAlert(sl.slackWeebHook, "", "Fatal", ColorDanger, fmt.Sprintf("`%s`", text))
	}
}

// Print prints the arguemtns to the info logger. It's good for standar logger compatibility
func (sl *SimpleLogger) Print(v ...interface{}) {
	sl.infoLogger.Print(v...)
}

// Printf prints the arguemtns to the info logger. It's good for standar logger compatibility
func (sl *SimpleLogger) Printf(format string, v ...interface{}) {
	sl.infoLogger.Printf(format, v...)
}

// SetLevel sets the log level (0=ERROR, 1=INFO, 2=DEBUG)
func (sl *SimpleLogger) SetLevel(level int) {
	sl.level = level
}

// GetLevel returns the current log level
func (sl *SimpleLogger) GetLevel() int {
	return sl.level
}

// SetSlackWeebhook configures the slack webhook used to report errors
func (sl *SimpleLogger) SetSlackWeebhook(s string) {
	sl.slackWeebHook = s
}

// EnableErrorsToSlack enable/disable the option to automatically
// log the errors to the configured slack channel
func (sl *SimpleLogger) EnableErrorsToSlack(b bool) {
	sl.errorsToSlack = b
}

// LogToSlack sends a message to the configured channel, if it's enabled
func (sl *SimpleLogger) LogToSlack(webHook, title, text string) {
	go func() {
		if err := SendAlert(webHook, "", title, ColorGood, text); err != nil {
			sl.Errorf("Found an error sending notification to Slack: %s", err)
		}
	}()
}
