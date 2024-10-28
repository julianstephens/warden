package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/julianstephens/warden/internal/backend/common"
	"github.com/julianstephens/warden/internal/warden"
)

type LocalHandler struct{}

func (h *LocalHandler) WriteConfig(data []byte) error { return nil }
func (h *LocalHandler) WriteKey(ctx context.Context, filename string, reader common.IReader) error {
	bReader, ok := reader.(*common.ByteReader)
	if !ok {
		return fmt.Errorf("invalid byte reader")
	}

	loc := ctx.Value("location").(string)

	err := warden.EnsureDir(path.Join(loc, "keys"))
	if err != nil {
		return fmt.Errorf("unable to create key dir: %+v", err)
	}

	_, err = os.Stat(filename)

	if os.IsNotExist(err) {
		keyfile, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("unable to create keyfile: %+v", err)
		}
		defer keyfile.Close()

		wroteBytes, err := io.Copy(keyfile, bReader.Reader)
		if err != nil {
			return err
		}

		if wroteBytes != bReader.Len {
			return fmt.Errorf("expected to write %d bytes, wrote %d", bReader.Len, wroteBytes)
		}

		err = keyfile.Close()
		if err != nil {
			return err
		}

		err = makeReadonly(filename)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("keyfile conflict: %s", filename)
	}

	return nil
}

func (h *LocalHandler) WritePack(data []byte) error { return nil }
