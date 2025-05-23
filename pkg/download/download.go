package download

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"
	"strings"

	"github.com/timoknapp/soundcloud-cli/pkg/soundcloud"
	"github.com/timoknapp/soundcloud-cli/pkg/tag"
)

func start(downloadPath, format string, track soundcloud.Track) error {

	if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
		os.Mkdir(downloadPath, os.ModePerm)
	}

	url := track.PermalinkURL

	chunks := strings.Split(url, "/")
	l := len(chunks)
	normalizedTrackName := chunks [l-1]
	normalizedArtistName := chunks [l-2]

	normalizedTrackName = strings.Replace(normalizedTrackName, "-", "_", -1)
	normalizedArtistName = strings.Replace(normalizedArtistName, "-", "_", -1)

	fileName := normalizedArtistName + "-" + normalizedTrackName

	artistSubFolder := filepath.Join(downloadPath, normalizedArtistName)
	os.MkdirAll(artistSubFolder, os.ModePerm)

	var fileToWrite string
	fileToWrite = artistSubFolder + string(os.PathSeparator) + fileName + "." + format

	if _, err := os.Stat(fileToWrite); err == nil {
		return errors.New("soundcloud-cli: track already exists")
	}

	downloadedTrack, err := downloadFile(track.MediaURL)
	if err != nil {
		return err
	}
	err = os.WriteFile(fileToWrite, downloadedTrack, 0644)
	if err != nil {
		return err
	}
	//todo add more informations
	err = tag.AddTags(fileToWrite, fileName, track)
	if err != nil {
		return err
	}
	return nil
}

// Prepare use track and prepare the download
func Prepare(track soundcloud.Track, dlPath string, quality string) {
	downloadStartTime := time.Now()
	absoluteDownloadPath, err := filepath.Abs(dlPath)
	if err != nil {
		log.Printf("could not get absolute downloadPath")
		absoluteDownloadPath = "/"
	}
	log.Printf("start downloading track to %v", absoluteDownloadPath)
	err = start(dlPath, quality, track)
	if err != nil {
		log.Println(err)
	} else {
		downloadProgress := (float64(1) / float64(1)) * 100
		log.Printf("%3.f%% finished downloading: %v", downloadProgress, track.Title)
	}
	downloadDuration := time.Since(downloadStartTime)
	log.Printf("downloaded track in %s to %v", downloadDuration.Round(time.Second), absoluteDownloadPath)
}
