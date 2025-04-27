package shared_utils

import "fmt"

const (
	KB = 1024
	MB = 1024 * KB
)

type FileSize int64

func (fs FileSize) Bytes() int64 {
	return int64(fs)
}

func (fs FileSize) MB() float64 {
	return float64(fs) / float64(MB)
}

func (fs FileSize) String() string {
	mb := fs.MB()
	return fmt.Sprintf("%.2f MB", mb)
}
