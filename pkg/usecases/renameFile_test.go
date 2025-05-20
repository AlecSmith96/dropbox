package usecases_test

import (
	"bytes"
	"encoding/json"
	"errors"
	mock_adapters "github.com/AlecSmith96/dopbox/mocks"
	"github.com/AlecSmith96/dopbox/pkg/drivers"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRenameFile_HappyPath(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	mockFileWriter := mock_adapters.NewMockFileModifier(ctrl)
	router := drivers.NewRouter("./dest", mockFileWriter)

	type renameFileRequestBody struct {
		Path         string `json:"path" binding:"required"`
		PreviousPath string `json:"previousPath"`
	}
	requestBody := renameFileRequestBody{
		Path:         "/some/new-path.go",
		PreviousPath: "/some/old-path.go",
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	g.Expect(err).ToNot(HaveOccurred())

	req := httptest.NewRequest(http.MethodPatch, "http://localhost:8080/v1/file", bytes.NewReader(requestBodyBytes))

	mockFileWriter.EXPECT().RenameFile("./dest/some/old-path.go", "./dest/some/new-path.go").
		Return(nil).Times(1)

	router.ServeHTTP(w, req)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(w.Code).To(Equal(http.StatusOK))
}

func TestRenameFile_ValidationErr(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	mockFileWriter := mock_adapters.NewMockFileModifier(ctrl)
	router := drivers.NewRouter("./dest", mockFileWriter)

	type renameFileRequestBody struct {
		Path         string `json:"path" binding:"required"`
		PreviousPath string `json:"previousPath"`
	}
	requestBody := renameFileRequestBody{}

	requestBodyBytes, err := json.Marshal(&requestBody)
	g.Expect(err).ToNot(HaveOccurred())

	req := httptest.NewRequest(http.MethodPatch, "http://localhost:8080/v1/file", bytes.NewReader(requestBodyBytes))

	mockFileWriter.EXPECT().RenameFile("./dest/some/old-path.go", "./dest/some/new-path.go").
		Return(nil).Times(0)

	router.ServeHTTP(w, req)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(w.Code).To(Equal(http.StatusBadRequest))
}

func TestRenameFile_ReturnsErr(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	mockFileWriter := mock_adapters.NewMockFileModifier(ctrl)
	router := drivers.NewRouter("./dest", mockFileWriter)

	type renameFileRequestBody struct {
		Path         string `json:"path" binding:"required"`
		PreviousPath string `json:"previousPath"`
	}
	requestBody := renameFileRequestBody{
		Path:         "/some/new-path.go",
		PreviousPath: "/some/old-path.go",
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	g.Expect(err).ToNot(HaveOccurred())

	req := httptest.NewRequest(http.MethodPatch, "http://localhost:8080/v1/file", bytes.NewReader(requestBodyBytes))

	mockFileWriter.EXPECT().RenameFile("./dest/some/old-path.go", "./dest/some/new-path.go").
		Return(errors.New("an error occurred")).Times(1)

	router.ServeHTTP(w, req)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(w.Code).To(Equal(http.StatusInternalServerError))
	g.Expect(w.Body.String()).To(Equal(`{"message":"an internal server error occurred"}`))
}
