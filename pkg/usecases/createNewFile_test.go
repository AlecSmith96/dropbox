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

func TestCreateNewFile_HappyPath(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	mockFileWriter := mock_adapters.NewMockFileModifier(ctrl)
	router := drivers.NewRouter("./dest", mockFileWriter)

	type createRequestBody struct {
		Path        string `json:"path"`
		Data        []byte `json:"data"`
		IsDirectory bool   `json:"isDirectory"`
	}
	requestBody := createRequestBody{
		Path: "/some/path.go",
		Data: []byte("some content"),
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	g.Expect(err).ToNot(HaveOccurred())

	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/v1/file", bytes.NewReader(requestBodyBytes))

	mockFileWriter.EXPECT().CreateFile("./dest/some/path.go", requestBody.Data, false).
		Return(nil).Times(1)

	router.ServeHTTP(w, req)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(w.Code).To(Equal(http.StatusOK))
}

func TestCreateNewFile_ValidationError(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	mockFileWriter := mock_adapters.NewMockFileModifier(ctrl)
	router := drivers.NewRouter("./dest", mockFileWriter)

	type createRequestBody struct {
		Path        string `json:"path"`
		Data        []byte `json:"data"`
		IsDirectory bool   `json:"isDirectory"`
	}
	requestBody := createRequestBody{
		Path: "",
		Data: []byte("some content"),
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	g.Expect(err).ToNot(HaveOccurred())

	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/v1/file", bytes.NewReader(requestBodyBytes))
	g.Expect(err).ToNot(HaveOccurred())

	mockFileWriter.EXPECT().CreateFile("./dest/some/path.go", requestBody.Data, false).
		Return(nil).Times(0)

	router.ServeHTTP(w, req)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(w.Code).To(Equal(http.StatusBadRequest))
	g.Expect(w.Body.String()).To(Equal(`{"message":"a bad request error occurred"}`))
}

func TestCreateNewFile_FileModifierReturnsError(t *testing.T) {
	g := NewGomegaWithT(t)
	w := httptest.NewRecorder()
	ctrl := gomock.NewController(t)
	mockFileWriter := mock_adapters.NewMockFileModifier(ctrl)
	router := drivers.NewRouter("./dest", mockFileWriter)

	type createRequestBody struct {
		Path        string `json:"path"`
		Data        []byte `json:"data"`
		IsDirectory bool   `json:"isDirectory"`
	}
	requestBody := createRequestBody{
		Path: "/some/path.go",
		Data: []byte("some content"),
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	g.Expect(err).ToNot(HaveOccurred())

	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/v1/file", bytes.NewReader(requestBodyBytes))
	g.Expect(err).ToNot(HaveOccurred())

	mockFileWriter.EXPECT().CreateFile("./dest/some/path.go", requestBody.Data, false).
		Return(errors.New("an error occurred")).Times(1)

	router.ServeHTTP(w, req)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(w.Code).To(Equal(http.StatusInternalServerError))
	g.Expect(w.Body.String()).To(Equal(`{"message":"an internal server error occurred"}`))
}
