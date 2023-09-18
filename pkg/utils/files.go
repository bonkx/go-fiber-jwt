package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/h2non/bimg"
	"github.com/h2non/filetype"
)

// A new folder is created at the root of the project.
func createFolder(dirname string) error {
	_, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirname, 0755)
		if errDir != nil {
			return errDir
		}
	}
	return nil
}

func FileUpload(c *fiber.Ctx, fileParam string, uploadTo string) (string, error) {

	fileheader, err := c.FormFile(fileParam)
	if err != nil {
		return "", errors.New(err.Error())
	}

	file, err := fileheader.Open()
	if err != nil {
		return "", errors.New(err.Error())
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		return "", errors.New(err.Error())
	}

	year, month, day := time.Now().Date()
	// initial dir path
	filePath := fmt.Sprintf("media/%s/%v/%v/%v", uploadTo, year, int(month), day)

	// create dir
	errDir := createFolder(filePath)
	if errDir != nil {
		return "", errors.New(err.Error())
	}

	var filename = ""
	// get file mime
	kind, _ := filetype.Match(buffer)

	if filetype.IsImage(buffer) {
		// if image
		filename, err = imageProcessing(buffer, 90, filePath)
		if err != nil {
			return "", errors.New(err.Error())
		}
	} else if filetype.IsVideo(buffer) {
		// if video
	} else {
		// if others file
		filename, err = fileProcessing(c, fileheader, filePath)
		if err != nil {
			return "", errors.New(err.Error())
		}
		errSave := c.SaveFile(fileheader, fmt.Sprintf("%s/%s", filePath, filename))
		if errSave != nil {
			return "", errors.New(err.Error())
		}
	}

	imageUrl := fmt.Sprintf("%s/%s", filePath, filename)

	log.Printf("Successfully uploaded %s to %s", kind.MIME.Value, imageUrl)

	//here we save our file to our path
	return imageUrl, nil
}

// The mime type of the image is changed, it is compressed and then saved in the specified folder.
func imageProcessing(buffer []byte, quality int, dirname string) (string, error) {
	fn := strings.Replace(uuid.New().String(), "-", "", -1)
	filename := fn + ".webp"
	thumbnail_name := fn + "_thumbnail.webp"

	options := bimg.Options{
		Quality:       quality,
		StripMetadata: false,
	}

	rorated, err := bimg.NewImage(buffer).AutoRotate()
	if err != nil {
		return filename, err
	}

	resized, err := resizeImage(rorated)
	if err != nil {
		return filename, err
	}

	converted, err := bimg.NewImage(resized).Convert(bimg.WEBP)
	if err != nil {
		return filename, err
	}

	processed, err := bimg.NewImage(converted).Process(options)
	if err != nil {
		return filename, err
	}

	// write media
	writeError := bimg.Write(fmt.Sprintf("./"+dirname+"/%s", filename), processed)
	if writeError != nil {
		return filename, writeError
	}

	// create thumbnail
	_thumbnail, err := bimg.NewImage(converted).Thumbnail(200)
	thumbWriteError := bimg.Write(fmt.Sprintf("./"+dirname+"/%s", thumbnail_name), _thumbnail)
	if thumbWriteError != nil {
		return filename, writeError
	}

	return filename, nil
}

func resizeImage(b []byte) ([]byte, error) {

	//[#1] Create image from bytes
	origImage := bimg.NewImage(b)

	//[#2] calculate relative height using aspect ratio
	origSize, _ := origImage.Size()
	width := 1920
	height := 0
	if origSize.Width > origSize.Height {
		aspectRatio := float64(origSize.Height / origSize.Width)
		height = int(float64(width) * aspectRatio)
	} else {
		aspectRatio := float64(origSize.Width / origSize.Height)
		height = int(float64(width) * aspectRatio)
	}

	//[#3] Apply resize operation with given width and height
	newImage, err := origImage.Resize(width, height)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	return newImage, nil
}

func fileProcessing(c *fiber.Ctx, file *multipart.FileHeader, dirname string) (string, error) {

	// rename file
	uniqueId := uuid.New()
	name := strings.Replace(uniqueId.String(), "-", "", -1)
	fileExt := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", name, fileExt)

	// save file
	errSave := c.SaveFile(file, fmt.Sprintf("./"+dirname+"/%s", filename))
	if errSave != nil {
		return "", errors.New(errSave.Error())
	}

	return filename, nil
}

func ImageUpload(c *fiber.Ctx, fileParam string, uploadTo string) (string, error) {

	fileheader, err := c.FormFile(fileParam)
	if err != nil {
		return "", errors.New(err.Error())
	}

	file, err := fileheader.Open()
	if err != nil {
		return "", errors.New(err.Error())
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		return "", errors.New(err.Error())
	}

	year, month, day := time.Now().Date()
	// initial dir path
	filePath := fmt.Sprintf("media/%s/%v/%v/%v", uploadTo, year, int(month), day)

	// create dir
	errDir := createFolder(filePath)
	if errDir != nil {
		return "", errors.New(err.Error())
	}

	var filename = ""
	// get file mime
	kind, _ := filetype.Match(buffer)

	if filetype.IsImage(buffer) {
		// if image
		filename, err = imageProcessing(buffer, 90, filePath)
		if err != nil {
			return "", errors.New(err.Error())
		}
	} else {
		// return "", errors.New("The file must be a file of type: jpeg, jpg, png")
		return "", errors.New("The file under validation must be an image (jpg, jpeg, png, bmp, gif, svg, or webp).")
	}

	imageUrl := fmt.Sprintf("%s/%s", filePath, filename)

	log.Printf("Successfully uploaded %s to %s", kind.MIME.Value, imageUrl)

	//here we save our file to our path
	return imageUrl, nil
}