package writer

import (
	"encoding/csv"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

var GatewaysWriterSingleton *GatewaysWriter

type GatewaysWriter struct {
	Writer *csv.Writer
	mux    sync.Mutex
}

func InitializeGatewaysWriter(file *os.File) {
	GatewaysWriterSingleton = &GatewaysWriter{Writer: csv.NewWriter(file)}
	GatewaysWriterSingleton.WriteRowToFile(
		"Gateway ID",
	)
}

func (writer *GatewaysWriter) WriteRowToFile(ID string) {
	writer.mux.Lock()
	if err := writer.Writer.Write([]string{ID}); err != nil {
		log.Fatal(err)
	}
	writer.mux.Unlock()
}
