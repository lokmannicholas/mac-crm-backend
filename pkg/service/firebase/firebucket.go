package firebase

import (
	"io"
	"path/filepath"

	"cloud.google.com/go/storage"
)

func (app *FirebaseApp) UploadFile(path, fileName string, src io.Reader) error {
	bucket, err := app.Storage.DefaultBucket()
	if err != nil {
		return err
	}
	object := bucket.Object(filepath.Join(path, fileName))
	writer := object.NewWriter(app.Context)
	//Set the attribute
	// writer.ObjectAttrs.Metadata = map[string]string{"id": fileName}
	defer writer.Close()

	if _, err := io.Copy(writer, src); err != nil {
		return err
	}

	object.ACL().Set(app.Context, storage.AllUsers, storage.RoleReader)
	return nil
}
func (app *FirebaseApp) DownloadFile(fileName string) (io.Reader, error) {
	bucket, err := app.Storage.DefaultBucket()
	if err != nil {
		return nil, err
	}
	object := bucket.Object(fileName)
	r, err := object.NewReader(app.Context)
	if err != nil {
		return nil, err
	}
	return r, nil

}
