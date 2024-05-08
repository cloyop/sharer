package fs

import "fmt"

var KB int = 1024
var MB int = 1_048_576
var GB int = 1_073_741_824

func BytesToScalaBytes(bytesQ int) string {
	if bytesQ < KB {
		return fmt.Sprintf("%.d bytes", bytesQ)
	}
	if bytesQ < MB {
		return fmt.Sprintf("%.d Kilobytes", bytesQ/KB)
	}
	if bytesQ < GB {
		return fmt.Sprintf("%.d Megabytes", bytesQ/MB)
	}
	return fmt.Sprintf("%.d Gigabytes", bytesQ/GB)
}
