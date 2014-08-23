/*
	This package simply instantiates a logging 
	mechanism based on the go log package.  This does 
	not use log rolling.  It is expected that the log 
	rolling will be done in the leader implementation 
	where the logging is first instantiated.
*/
package log

import (
	"io"
	"log"
)

// The log package allows four types of logging, 
// the TRACE, INFO, WARNING, and ERROR logging.
var (
	TRACE   *log.Logger
	INFO    *log.Logger
	WARNING *log.Logger
	ERROR   *log.Logger
)

// Called in the node leader in order to set 
// the io.Writer's the logger should be writing 
// too.  The default is to stdout and stderr.  
// But this should be changed to a file to hold the 
// logs which then gets rolled over time such that 
// the memory doesn't get lost.
func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	TRACE = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	INFO = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	WARNING = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ERROR = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
