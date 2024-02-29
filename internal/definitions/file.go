package definitions

// File represents a file in the file system.
// model: true
type File struct {
	Path       string `json:"path"`
	Size       int64  `json:"size"`
	ModifiedAt string `json:"modified_at"`
}

type FileService interface {
	Index(IndexRequest) FilesResponse
	Update(KeyRequest) EmptyResponse
	Stat(KeyRequest) EmptyResponse
}

type FilesResponse struct {
	Count  int64   `json:"count"`
	Result []*File `json:"files"`
}

type EmptyResponse struct{}
