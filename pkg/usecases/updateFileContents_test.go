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

func TestUpdateFileContents_HappyPath(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	mockFileWriter := mock_adapters.NewMockFileModifier(ctrl)
	router := drivers.NewRouter("./dest", mockFileWriter)

	type updateFileContentsRequestBody struct {
		Path string `json:"path" binding:"required"`
		Data []byte `json:"data"`
	}
	requestBody := updateFileContentsRequestBody{
		Path: "/some/path.go",
		Data: []byte("some content"),
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	g.Expect(err).ToNot(HaveOccurred())

	req := httptest.NewRequest(http.MethodPut, "http://localhost:8080/v1/file", bytes.NewReader(requestBodyBytes))

	mockFileWriter.EXPECT().UpdateFile("./dest/some/path.go", requestBody.Data).
		Return(nil).Times(1)

	router.ServeHTTP(w, req)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(w.Code).To(Equal(http.StatusOK))
}

func TestUpdateFileContents_ValidationErr(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	mockFileWriter := mock_adapters.NewMockFileModifier(ctrl)
	router := drivers.NewRouter("./dest", mockFileWriter)

	type updateFileContentsRequestBody struct {
		Path string `json:"path" binding:"required"`
		Data []byte `json:"data"`
	}
	requestBody := updateFileContentsRequestBody{
		Data: []byte("some content"),
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	g.Expect(err).ToNot(HaveOccurred())

	req := httptest.NewRequest(http.MethodPut, "http://localhost:8080/v1/file", bytes.NewReader(requestBodyBytes))

	mockFileWriter.EXPECT().UpdateFile("./dest/some/path.go", requestBody.Data).
		Return(nil).Times(0)

	router.ServeHTTP(w, req)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(w.Code).To(Equal(http.StatusBadRequest))
	g.Expect(w.Body.String()).To(Equal(`{"message":"a bad request error occurred"}`))
}

func TestUpdateFileContents_ReturnsErr(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	mockFileWriter := mock_adapters.NewMockFileModifier(ctrl)
	router := drivers.NewRouter("./dest", mockFileWriter)

	type updateFileContentsRequestBody struct {
		Path string `json:"path" binding:"required"`
		Data []byte `json:"data"`
	}
	requestBody := updateFileContentsRequestBody{
		Path: "/some/path.go",
		Data: []byte("some content"),
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	g.Expect(err).ToNot(HaveOccurred())

	req := httptest.NewRequest(http.MethodPut, "http://localhost:8080/v1/file", bytes.NewReader(requestBodyBytes))

	mockFileWriter.EXPECT().UpdateFile("./dest/some/path.go", requestBody.Data).
		Return(errors.New("an error occurred")).Times(1)

	router.ServeHTTP(w, req)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(w.Code).To(Equal(http.StatusInternalServerError))
	g.Expect(w.Body.String()).To(Equal(`{"message":"an internal server error occurred"}`))
}
