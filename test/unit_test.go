package matchServiceTest

import (
	"context"
	"testing"

	mock_adapters "github.com/akshaybt001/DatingApp_MatchMaking_Service/internal/adapters/mockAdapters"
	"github.com/akshaybt001/DatingApp_MatchMaking_Service/internal/service"
	"github.com/akshaybt001/DatingApp_proto_files/pb"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUnMatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapter := mock_adapters.NewMockAdapterInterface(ctrl)
	matchService := service.NewMatchService(mockAdapter, "", "")

	testUserID := "test-user-id"

	tests := []struct {
		name           string
		request        *pb.GetByUserId
		setupMocks     func()
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:    "Success",
			request: &pb.GetByUserId{Id: testUserID},
			setupMocks: func() {
				mockAdapter.EXPECT().IsMatchExist(testUserID).Return(true, nil).Times(1)
				mockAdapter.EXPECT().UnMatch(testUserID).Return(nil).Times(1)
			},
			expectedError: false,
		},
		{
			name:    "Fail - Not Matched",
			request: &pb.GetByUserId{Id: testUserID},
			setupMocks: func() {
				mockAdapter.EXPECT().IsMatchExist(testUserID).Return(false, nil).Times(1)
			},
			expectedError:  true,
			expectedErrMsg: "cannot unmatch as it is not matched user",
		},
		{
			name:    "Fail - IsMatchExist Error",
			request: &pb.GetByUserId{Id: testUserID},
			setupMocks: func() {
				mockAdapter.EXPECT().IsMatchExist(testUserID).Return(false, assert.AnError).Times(1)
			},
			expectedError:  true,
			expectedErrMsg: assert.AnError.Error(),
		},
		{
			name:    "Fail - UnMatch Error",
			request: &pb.GetByUserId{Id: testUserID},
			setupMocks: func() {
				mockAdapter.EXPECT().IsMatchExist(testUserID).Return(true, nil).Times(1)
				mockAdapter.EXPECT().UnMatch(testUserID).Return(assert.AnError).Times(1)
			},
			expectedError:  true,
			expectedErrMsg: assert.AnError.Error(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setupMocks()

			_, err := matchService.UnMatch(context.Background(), test.request)

			if test.expectedError {
				assert.Error(t, err)
				assert.Equal(t, test.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
