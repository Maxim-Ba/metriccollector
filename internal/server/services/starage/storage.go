package starage

import "context"

type Storage interface {
	Ping(ctx context.Context) error
}

func Ping(ctx context.Context, storage Storage)  error{
	err := storage.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}
