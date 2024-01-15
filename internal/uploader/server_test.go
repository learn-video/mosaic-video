package uploader_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mauricioabreu/mosaic-video/internal/mocks"
	"github.com/mauricioabreu/mosaic-video/internal/uploader"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("error reading")
}

type FileUploadHandlerTestSuite struct {
	suite.Suite
	ctrl          *gomock.Controller
	storageClient *mocks.MockStorage
	handler       *uploader.FileUploadHandler
	logger        *zap.SugaredLogger
}

func (suite *FileUploadHandlerTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.storageClient = mocks.NewMockStorage(suite.ctrl)
	suite.logger = zap.NewNop().Sugar()
	suite.handler = uploader.NewHandler(suite.storageClient, suite.logger)
}

func (suite *FileUploadHandlerTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func TestFileUploadHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(FileUploadHandlerTestSuite))
}

func (suite *FileUploadHandlerTestSuite) TestSuccessfulUpload() {
	expectedPath := "test_folder/test_file.txt"
	suite.storageClient.EXPECT().Upload(expectedPath, gomock.Any()).Return(nil)

	r := mux.NewRouter()
	r.Handle("/upload/{folder}/{filename}", suite.handler)

	req := httptest.NewRequest("PUT", "/upload/test_folder/test_file.txt", bytes.NewReader([]byte("test content")))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	suite.Assert().Equal(http.StatusOK, resp.StatusCode)
}

func (suite *FileUploadHandlerTestSuite) TestReadError() {
	r := mux.NewRouter()
	r.Handle("/upload/{folder}/{filename}", suite.handler)

	req := httptest.NewRequest("PUT", "/upload/test_folder/test_file.txt", errorReader{})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	suite.Assert().Equal(http.StatusInternalServerError, resp.StatusCode)
}

func (suite *FileUploadHandlerTestSuite) TestUploadError() {
	expectedPath := "test_folder/test_file.txt"
	uploadErr := errors.New("upload failed")
	suite.storageClient.EXPECT().Upload(expectedPath, gomock.Any()).Return(uploadErr)

	r := mux.NewRouter()
	r.Handle("/upload/{folder}/{filename}", suite.handler)

	req := httptest.NewRequest("PUT", "/upload/test_folder/test_file.txt", bytes.NewReader([]byte("test content")))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	suite.Assert().Equal(http.StatusInternalServerError, resp.StatusCode)
}
