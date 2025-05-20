package entities

const (
	OperationCreated  = "CREATE"
	OperationModified = "MODIFIED"
	OperationDeleted  = "DELETED"
	OperationRenamed  = "RENAMED"
)

// FilesystemEvent is a struct that represents a file event. It stores the name of the file and the operation that
// happened to it, along with the metadata.
type FilesystemEvent struct {
	Name         string
	Operation    string
	PreviousPath string
	FileContents FileContents
}

// FileContents is a struct that represents a file's metadata. It does whether it is a directory, its inode and the
// contents of the file.
type FileContents struct {
	IsDirectory bool
	Inode       uint64
	Data        []byte
}
