package firebase

import (
	"context"

	fb "firebase.google.com/go"
	"google.golang.org/api/option"
)

type firebaseAppOption func(*firebaseAppOptions)

type firebaseApp struct {
	app *fb.App
}

type firebaseAppOptions struct {
	credentialFile string
}

func NewFirebaseApp(ctx context.Context, opts ...firebaseAppOption) (*firebaseApp, error) {
	var options firebaseAppOptions
	for _, opt := range opts {
		opt(&options)
	}

	if options.credentialFile == "" {
		// TODO: Replace to logger Fatal
		panic("firebase credential file is required")
	}

	app, err := fb.NewApp(ctx, nil, option.WithCredentialsFile(options.credentialFile))
	if err != nil {
		return nil, err
	}

	return &firebaseApp{app}, nil
}

func (f *firebaseApp) Auth(ctx context.Context) (*firebaseAuth, error) {
	return NewFirebaseAuth(ctx, f.app)
}
