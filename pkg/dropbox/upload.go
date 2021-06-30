package dropbox

import (
	"bytes"
	"fmt"

	dbx "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

func Upload(client files.Client, path string, data []byte) error {
	commitInfo := files.NewCommitInfo(path)
	commitInfo.Mode = &files.WriteMode{Tagged: dbx.Tagged{Tag: "overwrite"}}

	_, err := client.Upload(commitInfo, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to upload: %s", err)
	}

	return nil
}
