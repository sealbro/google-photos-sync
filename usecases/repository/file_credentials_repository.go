package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"google-photos-sync/domain/model"
	"google-photos-sync/infrastructure/web/auth"
	"google-photos-sync/interface/repository"
	"os"
	"time"
)

type FileCredentialsRepository struct {
}

func MakeFileCredentialsRepository() repository.CredentialsRepository {
	return &FileCredentialsRepository{}
}

func (service *FileCredentialsRepository) Exists(userId string, accountType model.AccountType) (bool, error) {
	_, err := os.Stat(getAccountKey(userId, accountType))

	return err != nil, err
}

func (service *FileCredentialsRepository) Set(ctx context.Context, accountType model.AccountType, token []byte) error {
	f, err := os.OpenFile(getAccountKey(auth.GetUserId(ctx), accountType), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	json.NewEncoder(f).Encode(token)

	return nil
}

func (service *FileCredentialsRepository) Get(ctx context.Context) (*model.User, error) {
	userId := auth.GetUserId(ctx)
	fromToken, err := getToken(userId, model.From)
	if err != nil {
		return nil, err
	}

	toToken, err := getToken(userId, model.To)
	if err != nil {
		return nil, err
	}

	return &model.User{
		Id:        userId,
		Created:   time.Now(),
		Modified:  time.Now(),
		FromToken: fromToken,
		ToToken:   toToken,
	}, nil
}

func getToken(userId string, accountType model.AccountType) ([]byte, error) {
	f, err := os.Open(getAccountKey(userId, accountType))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	filesize := fileInfo.Size()
	buffer := make([]byte, filesize)

	f.Read(buffer)

	return buffer, err
}

func getAccountKey(userId string, accountType model.AccountType) string {
	return fmt.Sprintf("%s_%s", userId, accountType)
}
