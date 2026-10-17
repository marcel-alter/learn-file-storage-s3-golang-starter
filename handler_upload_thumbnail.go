package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here
	const maxMemory = 10 * 1024 * 1024
	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ParseMultipartForm failed!", err)
	}
	mpFile, mpHeader, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "FormFle failed!", err)
	}
	mediaT := mpHeader.Header.Get("Content-Type")
	imageData, err := io.ReadAll(mpFile)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ReadAll failed!", err)
	}
	videoMetaData, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find Video in Database!", err)
	}
	if userID != videoMetaData.UserID {
		respondWithError(w, http.StatusUnauthorized, "Video owned by different User", err)
		return
	}
	rawData := base64.StdEncoding.EncodeToString(imageData)
	thumbURL := "data:" + mediaT + ";base64," + rawData
	videoMetaData.ThumbnailURL = &thumbURL
	cfg.db.UpdateVideo(videoMetaData)

	respondWithJSON(w, http.StatusOK, videoMetaData)
}
