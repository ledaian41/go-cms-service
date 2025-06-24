package shared_interface

type HelperService interface {
	LoadSchema(filePath string, ch chan<- string)
	LoadJsonData(filePath string, ch chan<- string)
}
