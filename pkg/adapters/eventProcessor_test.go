package adapters

import (
	"errors"
	mock_adapters "github.com/AlecSmith96/dopbox/mocks"
	"github.com/AlecSmith96/dopbox/pkg/entities"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestNewEventProcessor_ProcessEvent_PathTrimmingReturnsErr(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)
	mockHTTPClient := mock_adapters.NewMockRequestSender(ctrl)

	eventProcessor := NewEventProcessor(mockHTTPClient, "./source/path")

	event := entities.FilesystemEvent{
		Name:      "./source/path/file.go/source/path",
		Operation: entities.OperationCreated,
		FileContents: entities.FileContents{
			IsDirectory: false,
			Inode:       0,
			Data:        nil,
		},
	}

	mockHTTPClient.EXPECT().SendCreateRequest("/file.go", nil, false).
		Return(nil).Times(0)

	err := eventProcessor.ProcessEvent(event)
	g.Expect(err).To(MatchError("invalid trimmed path produced"))
}

func TestNewEventProcessor_ProcessEvent_CreateHappyPath(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)
	mockHTTPClient := mock_adapters.NewMockRequestSender(ctrl)

	eventProcessor := NewEventProcessor(mockHTTPClient, "./source/path")

	event := entities.FilesystemEvent{
		Name:      "./source/path/file.go",
		Operation: entities.OperationCreated,
		FileContents: entities.FileContents{
			IsDirectory: false,
			Inode:       0,
			Data:        nil,
		},
	}

	mockHTTPClient.EXPECT().SendCreateRequest("/file.go", nil, false).
		Return(nil)

	err := eventProcessor.ProcessEvent(event)
	g.Expect(err).ToNot(HaveOccurred())
}

func TestNewEventProcessor_ProcessEvent_CreateReturnsErr(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)
	mockHTTPClient := mock_adapters.NewMockRequestSender(ctrl)

	eventProcessor := NewEventProcessor(mockHTTPClient, "./source/path")

	event := entities.FilesystemEvent{
		Name:      "./source/path/file.go",
		Operation: entities.OperationCreated,
		FileContents: entities.FileContents{
			IsDirectory: false,
			Inode:       0,
			Data:        nil,
		},
	}

	mockHTTPClient.EXPECT().SendCreateRequest("/file.go", nil, false).
		Return(errors.New("an error occurred"))

	err := eventProcessor.ProcessEvent(event)
	g.Expect(err).To(MatchError("an error occurred"))
}

func TestNewEventProcessor_ProcessEvent_RenameHappyPath(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)
	mockHTTPClient := mock_adapters.NewMockRequestSender(ctrl)

	eventProcessor := NewEventProcessor(mockHTTPClient, "./source/path")

	event := entities.FilesystemEvent{
		Name:         "./source/path/new-file.go",
		PreviousPath: "./source/path/old-file.go",
		Operation:    entities.OperationRenamed,
		FileContents: entities.FileContents{
			IsDirectory: false,
			Inode:       0,
			Data:        nil,
		},
	}

	mockHTTPClient.EXPECT().SendRenameRequest("/old-file.go", "/new-file.go").
		Return(nil)

	err := eventProcessor.ProcessEvent(event)
	g.Expect(err).ToNot(HaveOccurred())
}

func TestNewEventProcessor_ProcessEvent_RenamePathTrimmingReturnsErr(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)
	mockHTTPClient := mock_adapters.NewMockRequestSender(ctrl)

	eventProcessor := NewEventProcessor(mockHTTPClient, "./source/path")

	event := entities.FilesystemEvent{
		Name:      "./source/path/new-file.go",
		Operation: entities.OperationRenamed,
		FileContents: entities.FileContents{
			IsDirectory: false,
			Inode:       0,
			Data:        nil,
		},
	}

	mockHTTPClient.EXPECT().SendRenameRequest("/old-file.go", "/new-file.go").
		Return(nil).Times(0)

	err := eventProcessor.ProcessEvent(event)
	g.Expect(err).To(MatchError("invalid trimmed path produced"))
}

func TestNewEventProcessor_ProcessEvent_RenameReturnsErr(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)
	mockHTTPClient := mock_adapters.NewMockRequestSender(ctrl)

	eventProcessor := NewEventProcessor(mockHTTPClient, "./source/path")

	event := entities.FilesystemEvent{
		Name:         "./source/path/new-file.go",
		PreviousPath: "./source/path/old-file.go",
		Operation:    entities.OperationRenamed,
		FileContents: entities.FileContents{
			IsDirectory: false,
			Inode:       0,
			Data:        nil,
		},
	}

	mockHTTPClient.EXPECT().SendRenameRequest("/old-file.go", "/new-file.go").
		Return(errors.New("an error occurred"))

	err := eventProcessor.ProcessEvent(event)
	g.Expect(err).To(MatchError("an error occurred"))
}

func TestNewEventProcessor_ProcessEvent_DeleteHappyPath(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)
	mockHTTPClient := mock_adapters.NewMockRequestSender(ctrl)

	eventProcessor := NewEventProcessor(mockHTTPClient, "./source/path")

	event := entities.FilesystemEvent{
		Name:      "./source/path/file.go",
		Operation: entities.OperationDeleted,
		FileContents: entities.FileContents{
			IsDirectory: false,
			Inode:       0,
			Data:        nil,
		},
	}

	mockHTTPClient.EXPECT().SendDeleteRequest("/file.go").
		Return(nil)

	err := eventProcessor.ProcessEvent(event)
	g.Expect(err).ToNot(HaveOccurred())
}

func TestNewEventProcessor_ProcessEvent_DeleteReturnsErr(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)
	mockHTTPClient := mock_adapters.NewMockRequestSender(ctrl)

	eventProcessor := NewEventProcessor(mockHTTPClient, "./source/path")

	event := entities.FilesystemEvent{
		Name:      "./source/path/file.go",
		Operation: entities.OperationDeleted,
		FileContents: entities.FileContents{
			IsDirectory: false,
			Inode:       0,
			Data:        nil,
		},
	}

	mockHTTPClient.EXPECT().SendDeleteRequest("/file.go").
		Return(errors.New("an error occurred"))

	err := eventProcessor.ProcessEvent(event)
	g.Expect(err).To(MatchError("an error occurred"))
}

func TestNewEventProcessor_ProcessEvent_UpdateHappyPath(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)
	mockHTTPClient := mock_adapters.NewMockRequestSender(ctrl)

	eventProcessor := NewEventProcessor(mockHTTPClient, "./source/path")

	event := entities.FilesystemEvent{
		Name:      "./source/path/file.go",
		Operation: entities.OperationModified,
		FileContents: entities.FileContents{
			IsDirectory: false,
			Inode:       0,
			Data:        []byte("some content"),
		},
	}

	mockHTTPClient.EXPECT().SendUpdateRequest("/file.go", event.FileContents.Data).
		Return(nil)

	err := eventProcessor.ProcessEvent(event)
	g.Expect(err).ToNot(HaveOccurred())
}

func TestNewEventProcessor_ProcessEvent_UpdateReturnsErr(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)
	mockHTTPClient := mock_adapters.NewMockRequestSender(ctrl)

	eventProcessor := NewEventProcessor(mockHTTPClient, "./source/path")

	event := entities.FilesystemEvent{
		Name:      "./source/path/file.go",
		Operation: entities.OperationModified,
		FileContents: entities.FileContents{
			IsDirectory: false,
			Inode:       0,
			Data:        []byte("some content"),
		},
	}

	mockHTTPClient.EXPECT().SendUpdateRequest("/file.go", event.FileContents.Data).
		Return(errors.New("an error occurred"))

	err := eventProcessor.ProcessEvent(event)
	g.Expect(err).To(MatchError("an error occurred"))
}

func TestNewEventProcessor_ProcessEvent_InvalidOperation(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)
	mockHTTPClient := mock_adapters.NewMockRequestSender(ctrl)

	eventProcessor := NewEventProcessor(mockHTTPClient, "./source/path")

	event := entities.FilesystemEvent{
		Name:      "./source/path/file.go",
		Operation: "INVALID_OPERATION",
		FileContents: entities.FileContents{
			IsDirectory: false,
			Inode:       0,
			Data:        []byte("some content"),
		},
	}

	err := eventProcessor.ProcessEvent(event)
	g.Expect(err).To(MatchError("unknown event operation: INVALID_OPERATION"))
}
