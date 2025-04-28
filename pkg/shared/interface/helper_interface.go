package shared_interface

type HelperService interface {
	LoadJsonData(filePath string, ch chan<- string)
}
