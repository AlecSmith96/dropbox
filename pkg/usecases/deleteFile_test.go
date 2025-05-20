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

func TestDeleteFile_HappyPath(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	mockFileWriter := mock_adapters.NewMockFileModifier(ctrl)
	router := drivers.NewRouter("./dest", mockFileWriter)

	type deleteRequestBody struct {
		Path string `json:"path"`
	}
	requestBody := deleteRequestBody{
		Path: "/some/path.go",
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	g.Expect(err).ToNot(HaveOccurred())

	req := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/v1/file", bytes.NewReader(requestBodyBytes))

	mockFileWriter.EXPECT().DeleteFile("./dest/some/path.go").
		Return(nil).Times(1)

	router.ServeHTTP(w, req)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(w.Code).To(Equal(http.StatusOK))
}

func TestDeleteFile_ReturnsErr(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	mockFileWriter := mock_adapters.NewMockFileModifier(ctrl)
	router := drivers.NewRouter("./dest", mockFileWriter)

	type deleteRequestBody struct {
		Path string `json:"path"`
	}
	requestBody := deleteRequestBody{
		Path: "/some/path.go",
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	g.Expect(err).ToNot(HaveOccurred())

	req := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/v1/file", bytes.NewReader(requestBodyBytes))

	mockFileWriter.EXPECT().DeleteFile("./dest/some/path.go").
		Return(errors.New("an error occurred")).Times(1)

	router.ServeHTTP(w, req)
	g.Expect(w.Code).To(Equal(http.StatusInternalServerError))
	g.Expect(w.Body.String()).To(Equal(`{"message":"an internal server error occurred"}`))
}
