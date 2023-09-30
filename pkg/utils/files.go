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

func ImageUpload(fileHeader *multipart.FileHeader, uploadTo string) (string, error) {

	file, err := fileHeader.Open()
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
		return "", errors.New(errDir.Error())
	}

	var filename = ""
	// get file mime
	kind, _ := filetype.Match(buffer)

	if filetype.IsImage(buffer) {
		// if image
		filename, err = imageProcessing(buffer, filePath)
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

func FileUpload(c *fiber.Ctx, fileHeader *multipart.FileHeader, uploadTo string) (string, error) {

	file, err := fileHeader.Open()
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
		return "", errors.New(errDir.Error())
	}

	var filename = ""
	// get file mime
	kind, _ := filetype.Match(buffer)

	if filetype.IsImage(buffer) {
		// if image
		filename, err = imageProcessing(buffer, filePath)
		if err != nil {
			return "", errors.New(err.Error())
		}
	} else if filetype.IsVideo(buffer) {
		// if video
	} else {
		// if others file
		filename, err = fileProcessing(c, fileHeader, filePath)
		if err != nil {
			return "", errors.New(err.Error())
		}
		errSave := c.SaveFile(fileHeader, fmt.Sprintf("%s/%s", filePath, filename))
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
func imageProcessing(buffer []byte, dirname string) (string, error) {
	fn := strings.Replace(uuid.New().String(), "-", "", -1)
	filename := fn + ".webp"
	thumbnail_name := fn + "_thumbnail.webp"

	options := bimg.Options{
		Quality:       90,
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

func GetThumbnail(fileName string) *string {
	dirPath := filepath.Dir(fileName)
	fileExt := filepath.Ext(fileName)

	filename := filepath.Base(fileName)
	// fmt.Println("filename: ", filename)
	// fmt.Println("fileNameOnly: ", fileNameOnly)
	fileNameOnly := FileNameWithoutExtSliceNotation(filename)

	thumbnailName := fmt.Sprintf("%s/%s_thumbnail%s", dirPath, fileNameOnly, fileExt)

	return &thumbnailName
}

func FileNameWithoutExtSliceNotation(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func RemoveFileSilence(fileUrl string) error {
	originFile := fileUrl
	// fmt.Println("originFile: ", originFile)
	thumbnailName := GetThumbnail(fileUrl)
	// fmt.Println("thumbnailName: ", thumbnailName)

	// Using Remove() function
	// errF := os.Remove(originFile)
	// if errF != nil {
	// 	return errors.New(errF.Error())
	// }

	// Removing file from server
	os.Remove(originFile)
	// remove thumbnail
	os.Remove(*thumbnailName)

	return nil
}
