package mock_user

import (
	dto "github.com/mrdolles/wallet-app/src/app/dto/user"
	repo "github.com/mrdolles/wallet-app/src/infra/persistence/mysql/user"

	"github.com/stretchr/testify/mock"
)

type MockUsersRepo struct {
	mock.Mock
}

func NewMockUsersRepo() *MockUsersRepo {
	return &MockUsersRepo{}
}

var _ repo.UserRepository = &MockUsersRepo{}

func (m *MockUsersRepo) Register(data *dto.RegisterReqDTO) (*dto.RegisterRespDTO, error) {
	args := m.Called(data)
	var (
		res *dto.RegisterRespDTO
		err error
	)

	if n, ok := args.Get(0).(*dto.RegisterRespDTO); ok {

		res = n
	}

	if n, ok := args.Get(1).(error); ok {

		err = n
	}

	return res, err
}

func (m *MockUsersRepo) Login(data *dto.LoginReqDTO) (*dto.RegisterRespDTO, error) {
	args := m.Called(data)
	var (
		res *dto.RegisterRespDTO
		err error
	)

	if n, ok := args.Get(0).(*dto.RegisterRespDTO); ok {

		res = n
	}

	if n, ok := args.Get(1).(error); ok {

		err = n
	}

	return res, err
}
