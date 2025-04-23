package user

import (
	"log"

	dto "github.com/mrdolles/wallet-app/src/app/dto/user"

	repo "github.com/mrdolles/wallet-app/src/infra/persistence/mysql/user"
)

type UserUCInterface interface {
	Register(data *dto.RegisterReqDTO) (*dto.RegisterRespDTO, error)
	Login(data *dto.LoginReqDTO) (*dto.RegisterRespDTO, error)
}

type userUseCase struct {
	Repo repo.UserRepository
}

func NewUserUseCase(repo repo.UserRepository) UserUCInterface {
	return &userUseCase{
		Repo: repo,
	}
}

func (u *userUseCase) Register(data *dto.RegisterReqDTO) (*dto.RegisterRespDTO, error) {
	result, err := u.Repo.Register(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func (u *userUseCase) Login(data *dto.LoginReqDTO) (*dto.RegisterRespDTO, error) {

	result, err := u.Repo.Login(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}
