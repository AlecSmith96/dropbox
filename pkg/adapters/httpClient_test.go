package adapters

import (
	"errors"
	mock_adapters "github.com/AlecSmith96/dopbox/mocks"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
)

func TestIsServerLive_HappyPath(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)

	mockHTTPClient := mock_adapters.NewMockHttpClient(ctrl)
	client := NewHTTPClient(mockHTTPClient, "http://localhost:8080")

	mockHTTPClient.EXPECT().Do(gomock.AssignableToTypeOf(&http.Request{})).Return(&http.Response{
		Status:     "200 OK",
		StatusCode: 200,
	}, nil)

	isLive := client.IsServerLive()
	g.Expect(isLive).To(BeTrue())
}

func TestIsServerLive_HappyPathServerNotLive(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)

	mockHTTPClient := mock_adapters.NewMockHttpClient(ctrl)
	client := NewHTTPClient(mockHTTPClient, "http://localhost:8080")

	mockHTTPClient.EXPECT().Do(gomock.AssignableToTypeOf(&http.Request{})).Return(&http.Response{
		Status:     "400 Bad Request",
		StatusCode: 400,
	}, nil)

	isLive := client.IsServerLive()
	g.Expect(isLive).To(BeFalse())
}

func TestIsServerLive_RequestReturnsErr(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)

	mockHTTPClient := mock_adapters.NewMockHttpClient(ctrl)
	client := NewHTTPClient(mockHTTPClient, "http://localhost:8080")

	mockHTTPClient.EXPECT().Do(gomock.AssignableToTypeOf(&http.Request{})).Return(&http.Response{
		Status:     "500 Internal Server Error",
		StatusCode: 500,
	}, errors.New("an error occurred"))

	isLive := client.IsServerLive()
	g.Expect(isLive).To(BeFalse())
}

func TestSendCreateRequest_HappyPath(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)

	path := "/some/path/file.go"
	data := []byte("some content")
	isDirectory := false

	mockHTTPClient := mock_adapters.NewMockHttpClient(ctrl)
	client := NewHTTPClient(mockHTTPClient, "http://localhost:8080")

	mockHTTPClient.EXPECT().Do(gomock.AssignableToTypeOf(&http.Request{})).Return(&http.Response{
		Status:     "200 OK",
		StatusCode: 200,
	}, nil)

	err := client.SendCreateRequest(path, data, isDirectory)
	g.Expect(err).ToNot(HaveOccurred())
}

func TestSendCreateRequest_ServerReturnsFailedStatusCode(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)

	path := "/some/path/file.go"
	data := []byte("some content")
	isDirectory := false

	mockHTTPClient := mock_adapters.NewMockHttpClient(ctrl)
	client := NewHTTPClient(mockHTTPClient, "http://localhost:8080")

	mockHTTPClient.EXPECT().Do(gomock.AssignableToTypeOf(&http.Request{})).Return(&http.Response{
		Status:     "500 Internal Server Error",
		StatusCode: 500,
	}, nil)

	err := client.SendCreateRequest(path, data, isDirectory)
	g.Expect(err).To(MatchError("request failed with status code 500"))
}

func TestSendCreateRequest_RequestReturnsErr(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)

	path := "/some/path/file.go"
	data := []byte("some content")
	isDirectory := false

	mockHTTPClient := mock_adapters.NewMockHttpClient(ctrl)
	client := NewHTTPClient(mockHTTPClient, "http://localhost:8080")

	mockHTTPClient.EXPECT().Do(gomock.AssignableToTypeOf(&http.Request{})).Return(nil, errors.New("an error occurred"))

	err := client.SendCreateRequest(path, data, isDirectory)
	g.Expect(err).To(MatchError("an error occurred"))
}
