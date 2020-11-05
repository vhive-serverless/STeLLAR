package writer

import (
	"encoding/csv"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

//GatewaysWriterSingleton is the object used to write IDs of created gateways to a file. It is safe to use
//concurrently (mutual exclusion lock).
var GatewaysWriterSingleton *gatewaysWriter

type gatewaysWriter struct {
	Writer *csv.Writer
	mux    sync.Mutex
}

//InitializeGatewaysWriter create a new writer for given file and writes the header `Gateway ID`.
func InitializeGatewaysWriter(file *os.File) {
	GatewaysWriterSingleton = &gatewaysWriter{Writer: csv.NewWriter(file)}
	GatewaysWriterSingleton.WriteGatewayID(
		"Gateway ID",
	)
}

//WriteGatewayID is used to write the specified gateway ID to the initialized file.
func (writer *gatewaysWriter) WriteGatewayID(ID string) {
	writer.mux.Lock()
	if err := writer.Writer.Write([]string{ID}); err != nil {
		log.Fatal(err)
	}
	writer.Writer.Flush()
	writer.mux.Unlock()
}
