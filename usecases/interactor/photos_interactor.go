package interactor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"google-photos-sync/domain/model"
	"google-photos-sync/infrastructure/google"
	"google-photos-sync/interface/repository"
	"google-photos-sync/pkg/logger"
	"google-photos-sync/usecases/machine/states"
	"io"
	"net/http"
	"os"
	"strings"
)

type PhotosInteractor struct {
	credentialsRepository repository.CredentialsRepository
	stateRepository       repository.StateRepository
}

func MakePhotosInteractor(credentialsRepository repository.CredentialsRepository, stateRepository repository.StateRepository) (*PhotosInteractor, error) {
	service := &PhotosInteractor{
		credentialsRepository: credentialsRepository,
		stateRepository:       stateRepository,
	}

	return service, nil
}

func (service *PhotosInteractor) Statistic(ctx context.Context) error {
	// TODO statistics
	// [photos.count, albums.count, moved.count, photos_per_minutes, megabyte_per_minute]
	return errors.New("not implemented")
}

func (service *PhotosInteractor) SaveAccount(ctx context.Context, callback *model.GoogleCallback) error {
	state, err := service.stateRepository.GetLastState(ctx)
	if err != nil {
		return err
	}
	if state.To == states.Sync {
		return errors.New("sync in process")
	}

	authConfig, err := google.GetAuthConfig(callback.Status)
	if err != nil {
		return err
	}

	exchange, err := authConfig.Exchange(ctx, callback.Code)
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}

	encoder := json.NewEncoder(&buf)
	encoder.Encode(exchange)

	return service.credentialsRepository.Set(ctx, callback.Status, buf.Bytes())
}

func (service *PhotosInteractor) Complete(ctx context.Context) error {
	state, err := service.stateRepository.GetLastState(ctx)
	if err != nil {
		return err
	}
	if state.To != states.Sync {
		return errors.New("complete only for sync state")
	}

	// TODO complete
	return errors.New("not implemented")
}

func (service *PhotosInteractor) Sync(ctx context.Context) error {
	user, err := service.credentialsRepository.Get(ctx)
	if err != nil {
		return err
	}
	if user.CanSync() {
		return errors.New("can't sync. should add from and to accounts")
	}

	fromClient, fromGClient, err := service.getClient(ctx, user, model.From)
	if err != nil {
		return err
	}

	_, toGClient, err := service.getClient(ctx, user, model.To)
	if err != nil {
		return err
	}

	list, err := fromGClient.Albums.List(ctx)

	for _, album := range list {
		if album.Title == "Time-laps" {
			logger.Info(album.Title)

			listByAlbum, err := fromGClient.MediaItems.ListByAlbum(ctx, album.ID)
			if err != nil {
				return err
			}

			// todo if exists return previus album id
			create, err := toGClient.Albums.Create(ctx, album.Title)
			if err != nil {
				return err
			}

			toAlbumId := create.ID

			logger.Info(toAlbumId)

			for _, mediaItem := range listByAlbum {
				filename := mediaItem.Filename

				var url string
				if strings.HasPrefix(strings.ToLower(mediaItem.MimeType), "video") {
					url = mediaItem.BaseURL + "=dv"
				} else {
					url = fmt.Sprintf("%v=d", mediaItem.BaseURL)
				}

				resp, err := fromClient.Get(url)
				if err != nil {
					return err
				}

				file, err := os.Create(filename)
				if err != nil {
					return err
				}

				_, err = io.Copy(file, resp.Body)
				if err != nil {
					return err
				}

				resp.Body.Close()
				file.Close()

				// todo add implementation with reader / writer
				_, err = toGClient.UploadFileToAlbum(ctx, toAlbumId, filename)
				if err != nil {
					return err
				}

				os.Remove(filename)

				break
			}

		}
	}

	return err
}

func (service *PhotosInteractor) getClient(ctx context.Context, user *model.User, accountType model.AccountType) (*http.Client, *gphotos.Client, error) {
	configFrom, err := google.GetAuthConfig(accountType)
	if err != nil {
		return nil, nil, err
	}
	fromClient := configFrom.Client(ctx, user.Token(accountType))
	if err != nil {
		return nil, nil, err
	}
	fromGClient, err := gphotos.NewClient(fromClient)
	if err != nil {
		return nil, nil, err
	}

	return fromClient, fromGClient, err
}
