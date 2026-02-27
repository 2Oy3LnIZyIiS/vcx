package domains

import (
	"fmt"
	"vcx/pkg/logging"
)


var log = logging.GetLogger()


// LogError logs a domain operation error with consistent formatting
func LogError(domain, operation string, err error) {
	log.Error(fmt.Sprintf("%s %s Failed: %v", domain, operation, err))
}
