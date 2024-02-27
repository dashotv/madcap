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
}

type FilesResponse struct {
	Count  int     `json:"count"`
	Result []*File `json:"files"`
}
