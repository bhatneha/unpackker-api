package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"mime/multipart"

	con "github.com/bhatneha/unpackker-api/config"
	"github.com/gin-gonic/gin"
)

type Form struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func packker(c *gin.Context) {
	var set con.PackkerInput
	if err := c.ShouldBindJSON(&set); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stat, err := set.Pack()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": stat})
}

func unpackker(c *gin.Context) {
	var set con.UnPackkerInput
	if err := c.ShouldBindJSON(&set); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stat, err := set.Unpack()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": stat})
}

func upload(c *gin.Context) {
	uid := con.GetUinqueID()
	folder := "/tmp"
	if err := con.CreateDir(folder, uid); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	var form Form
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	path := filepath.Join("/tmp", uid, form.File.Filename)
	if err := c.SaveUploadedFile(form.File, path); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"assetpath": path, "name": form.File.Filename, "assetversion": "0.0.1"})
}

func delete(c *gin.Context) {
	var del con.DeleteInput
	if err := c.ShouldBindJSON(&del); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	dir, _ := filepath.Split(del.Assetpath)
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("file %s removal failed", del.Assetpath)})
	}
	c.JSON(http.StatusOK, gin.H{"response": "Deleted"})
}

func download(c *gin.Context) {
	Filename := c.Request.URL.Query().Get("asset")
	if Filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Get 'asset' not specified in url"})
		return
	}

	Openfile, err := os.Open(Filename)
	if err != nil {
		//File not found, send 404
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	defer Openfile.Close() //Close after function return

	FileHeader := make([]byte, 512)
	_, err = Openfile.Read(FileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	FileContentType := http.DetectContentType(FileHeader)

	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+Filename)
	c.Writer.Header().Set("Content-Type", FileContentType)
	c.Writer.Header().Set("Content-Length", FileSize)

	_, err = Openfile.Seek(0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	_, err = io.Copy(c.Writer, Openfile) //'Copy' the file to the client
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
