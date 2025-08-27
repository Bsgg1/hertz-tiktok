package user

import (
	"context"
	"sync"
	"tiktok/biz/dal/db"
	user "tiktok/biz/model/basic/user"
	"tiktok/biz/model/common"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{ctx: ctx, c: c}
}

func (s *UserService) UserRegister(req *user.DouyinUserRegisterRequest) (user_id int64, err error) {
	user, err := db.QueryUser(req.Username)
	if err != nil {
		return 0, err
	}
	if *user != (db.User{}) {
		return 0, errno.UserAlreadyExistErr
	}
	psd, err := utils.Crypt(req.Password)
	if err != nil {
		return 0, err
	}
	user_id, err = db.CreateUser(&db.User{
		UserName:        req.Username,
		Password:        psd,
		Avatar:          constants.TestAva,
		BackgroundImage: constants.TestBackground,
		Signature:       constants.TestSign,
	})
	return
}

func (s *UserService) UserInfo(req *user.DouyinUserRequest) (*common.User, error) {
	query_user_id := req.UserId
	current_user_id, exists := s.c.Get("current_user_id")
	if !exists {
		current_user_id = 0
	}
	return s.GetUserInfo(query_user_id, current_user_id.(int64))
}

func (s *UserService) GetUserInfo(query_user_id, user_id int64) (*common.User, error) {
	u := &common.User{}
	errChan := make(chan error, 7)
	defer close(errChan)
	var wg sync.WaitGroup
	wg.Add(7)
	go func() {
		dbUser, err := db.QueryUserById(query_user_id)
		if err != nil {
			errChan <- err
		} else {
			u.Name = dbUser.UserName
			u.Avatar = dbUser.Avatar
			u.BackgroundImage = dbUser.BackgroundImage
			u.Signature = dbUser.Signature
		}
		wg.Done()
	}()
	return nil, nil
}
