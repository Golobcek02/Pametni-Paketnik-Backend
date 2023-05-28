package controllers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func LoginFaceID(c *gin.Context) {
	userId := strings.TrimSpace(c.Param("id"))

	// Parse the multipart form
	err := c.Request.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//return
	}

	// Get the files from the "File" field
	files := make([]*multipart.FileHeader, 0)
	for i := 0; true; i++ {
		fileHeader, exists := c.Request.MultipartForm.File[fmt.Sprintf("image%d", i)]
		if !exists {
			break
		}
		files = append(files, fileHeader...)
	}

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		//return
	}

	// Create a directory to store the images
	err = os.MkdirAll("images/"+userId, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//return
	}

	// Iterate over the files and save them
	for _, file := range files {
		// Open the file
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			//return
		}

		// Create the destination file
		dstPath := fmt.Sprintf("images/%s/%s", userId, file.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			src.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			//return
		}

		// Copy the file contents to the destination
		_, err = io.Copy(dst, src)
		dst.Close()
		src.Close()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			//return
		}
	}

	cmd := exec.Command("python", "scripts/LoginFaceId.py", userId)
	out, err := cmd.Output()
	println(string(out))
	if err != nil {
		println(err.Error())
		////return
	}

	removeErr := os.RemoveAll("images/" + userId)
	if removeErr != nil {
		fmt.Println(removeErr.Error())
		//return
	}
	stringOut := string(out)
	if stringOut[len(stringOut)-1] != 'T' {
		c.JSON(http.StatusOK, false)
	}
	//res = true
	c.JSON(http.StatusOK, true)
}

func RegisterFaceID(c *gin.Context) {

	userId := c.Param("id")

	// Parse the multipart form
	err := c.Request.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//return
	}

	// Get the files from the "File" field
	files := make([]*multipart.FileHeader, 0)
	for i := 0; true; i++ {
		fileHeader, exists := c.Request.MultipartForm.File[fmt.Sprintf("image%d", i)]
		if !exists {
			break
		}
		files = append(files, fileHeader...)
	}

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		//return
	}

	// Create a directory to store the images
	err = os.MkdirAll("images/"+userId, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//return
	}

	// Iterate over the files and save them
	for _, file := range files {
		// Open the file
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			//return
		}

		// Create the destination file
		dstPath := fmt.Sprintf("images/%s/%s", userId, file.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			src.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			//return
		}

		// Copy the file contents to the destination
		_, err = io.Copy(dst, src)
		dst.Close()
		src.Close()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			//return
		}
	}

	cmd := exec.Command("python", "scripts/RegisterFaceId.py", userId)
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
	}

	stringOut := string(out)

	if stringOut[len(stringOut)-1] != 'T' {
		c.JSON(http.StatusOK, false)
	}

	c.JSON(http.StatusOK, true)
}
