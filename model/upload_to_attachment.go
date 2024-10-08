package model

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"

	"google.golang.org/appengine/v2/mail"
)

// MaxFileSize is Attachment max size (MB)
var MaxFileSize = os.Getenv("MAX_FILE_SIZE")

// UploadToAttachment to change upload file to attachment
type UploadToAttachment struct {
	UploadFile   multipart.File
	UploadHeader *multipart.FileHeader
	Attachment   mail.Attachment
}

// Change is main code
func (u *UploadToAttachment) Change() (err error) {
	if u.UploadFile == nil {
		return
	}

	var data []byte
	buffer := make([]byte, 1024)
	totalBytesRead := 0

	for {
		// Read a kilobyte of data from the file into our buffer
		n, readErr := u.UploadFile.Read(buffer)

		totalBytesRead += n

		// Alert the user if there as error reading the file
		if readErr != nil && readErr.Error() != `EOF` {
			err = errors.New("error to read file")
			return
		}
		MaxFileSizeInt, _ := strconv.Atoi(MaxFileSize)
		if totalBytesRead > MaxFileSizeInt*1024*1024 {
			err = fmt.Errorf("%s is larger than %d MB please resize then re-upload", u.UploadHeader.Filename, MaxFileSizeInt)
			return
		}

		// Copy the bytes into our data array
		data = append(data, buffer[:n]...)

		// Stop reading file if we reach the end or there's no data to copy.
		if readErr != nil && readErr == io.EOF {
			break
		}
	}
	// Add the file to attachments
	u.Attachment = mail.Attachment{
		Name: u.UploadHeader.Filename,
		Data: data,
	}
	return
}
