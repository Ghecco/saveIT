package controllers

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	filesFolder    = "usersFiles/"
	picturesFolder = "/pictures/"
	videosFolder   = "/videos/"
)

// Check if user files folder exist
func CheckUserFolder(username string) bool {
	_, err := os.Stat(fmt.Sprintf("%s%s", filesFolder, username))
	if os.IsNotExist(err) {
		fmt.Printf("Folder %s%s not founded.", filesFolder, username)
		return false
	}
	return true
}

func AddUserFolder(username string) bool {
	err := os.Mkdir(fmt.Sprintf("%s%s", filesFolder, username), 0755)
	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Printf("Folder %s%s create.", filesFolder, username)
	return true
}

func RemoveUserFolder(username string) bool {
	err := os.Remove(fmt.Sprintf("%s%s", filesFolder, username))
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func CopyMedia(username string, photo *tb.Photo) bool {

	srcFile, _ := os.Open("4.jpg")
	endFile, _ := os.Create(srcFile.Name())
	_, err := io.Copy(endFile, photo.FileReader)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
func AddPicture(username string, File *tb.File, b *tb.Bot) bool {

	dir := fmt.Sprintf("%s%s/%s", filesFolder, username, picturesFolder)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0777)
		fmt.Print("prova a creare la dir")
	}
	endFile, _ := os.Create(fmt.Sprintf("%s/%s.jpg", dir, File.UniqueID))
	tempFile, _ := b.GetFile(File)
	_, err := io.Copy(endFile, tempFile)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

func AddVideo(username string, File *tb.File, b *tb.Bot) bool {

	dir := fmt.Sprintf("%s%s/%s", filesFolder, username, videosFolder)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
	endFile, _ := os.Create(fmt.Sprintf("%s/%s.mp4", dir, File.UniqueID))
	tempFile, _ := b.GetFile(File)
	_, err := io.Copy(endFile, tempFile)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

func getAllUserFiles(username string) (photos, videos []string) {

	root := fmt.Sprintf("%s%s/%s", filesFolder, username, picturesFolder)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		photos = append(photos, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	root = fmt.Sprintf("%s%s/%s", filesFolder, username, videosFolder)
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		videos = append(videos, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return photos, videos
}

func SendFilesToUser(username string, b *tb.Bot, m *tb.Message) {
	photos, videos := getAllUserFiles(username)

	b.Send(m.Sender, "Your Pictures:")
	for _, p := range photos {
		temp := &tb.Photo{File: tb.FromDisk(p)}
		b.Send(m.Sender, temp)
	}

	b.Send(m.Sender, "Your Videos:")
	for _, v := range videos {
		temp := &tb.Photo{File: tb.FromDisk(v)}
		b.Send(m.Sender, temp)
	}
	b.Send(m.Sender, "No more pictures & videos")
}
