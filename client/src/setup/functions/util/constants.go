package util

//BinaryName represents the name of the generated binary file
const BinaryName = "producer-handler"

//BytesToMB transforms bytes into megabytes
func BytesToMB(sizeBytes int64) float64 {
	return float64(sizeBytes) / 1024. / 1024.
}

//MBToBytes transforms megabytes into bytes
func MBToBytes(sizeMB float64) int64 {
	return int64(sizeMB) * 1024 * 1024
}
