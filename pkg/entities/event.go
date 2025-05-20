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

type FileContents struct {
	Inode uint64
	Data  []byte
}
