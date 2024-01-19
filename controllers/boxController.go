package controllers

import (
	"backend/schemas"
	"backend/utils"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ClaimBox(c *gin.Context) {
	var requestData struct {
		BoxID  int
		UserID string
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//return
	}

	cur, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").Find(context.Background(), bson.D{{}})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//return
	}

	ownerID, err := primitive.ObjectIDFromHex(requestData.UserID)
	emptyID, err := primitive.ObjectIDFromHex("000000000000")
	for cur.Next(context.TODO()) {
		var elem schemas.Box
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
			//log.Fatal(err)
		}

		if elem.BoxId == requestData.BoxID {
			if elem.OwnerId != emptyID {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Someone already owns this box"})
				//return
			}

			_, error := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").UpdateOne(context.Background(),
				bson.D{{Key: "boxid", Value: requestData.BoxID}},
				bson.D{{Key: "$set", Value: bson.D{
					{Key: "ownerid", Value: ownerID},
				}}},
				options.Update().SetUpsert(true))

			if error != nil {
				fmt.Println(error)
			}

			c.IndentedJSON(http.StatusOK, gin.H{"message": "Box ownership successfully updated!"})
			//return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Box not found"})

}

func AddUserBox(c *gin.Context) {
	var requestData struct {
		UserID     string
		SmartBoxID string
		Lat        float64
		Lon        float64
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//return
	}
	fmt.Println(requestData)

	cur, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").Find(context.Background(), bson.D{{}})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//return
	}

	str, _ := primitive.ObjectIDFromHex(requestData.UserID)
	emptyId, _ := primitive.ObjectIDFromHex("000000000000")
	for cur.Next(context.TODO()) {
		var elem schemas.Box
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
			//log.Fatal(err)
		}

		tmp, _ := strconv.Atoi(requestData.SmartBoxID)
		if elem.BoxId == tmp && elem.OwnerId != emptyId && elem.OwnerId != str {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Someone already own this box"})
			//return
		}
		if elem.BoxId == tmp && elem.OwnerId == str {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Someone already own this box"})
			//return
		}
		if elem.BoxId == tmp && elem.OwnerId == emptyId {
			_, error := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").UpdateOne(context.Background(),
				bson.D{{Key: "boxid", Value: elem.BoxId}},
				bson.D{{Key: "$set", Value: bson.D{
					{Key: "ownerid", Value: str},
					{Key: "latitude", Value: requestData.Lat},
					{Key: "longitude", Value: requestData.Lon},
				}}},
				options.Update().SetUpsert(true))

			if error != nil {
				fmt.Println(error)
			}

			c.IndentedJSON(http.StatusOK, gin.H{"message": "Box successfully updated!"})
			//return
		}
	}

	var box schemas.Box
	var temp []primitive.ObjectID
	box.BoxId, _ = strconv.Atoi(requestData.SmartBoxID)
	box.OwnerId = emptyId
	box.AccessIds = temp
	box.Latitude = requestData.Lat
	box.Longitude = requestData.Lon
	fmt.Println(box)

	result, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").InsertOne(context.Background(), box)
	fmt.Println(err)
	fmt.Println(result.InsertedID)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Box ownership successful!"})
}

func RemoveBox(c *gin.Context) {
	var boxid = c.Param("id")
	boxIdInt, _ := strconv.Atoi(boxid)

	res, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").DeleteOne(context.TODO(), bson.D{{Key: "boxid", Value: boxIdInt}})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Error while deleting this box")
	}

	fmt.Print(res.DeletedCount)
	c.IndentedJSON(http.StatusOK, "successfully deleted")
}

func ClearBoxOwner(c *gin.Context) {
	boxid := c.Param("id")
	boxIdInt, _ := strconv.Atoi(boxid)
	noOwner, _ := primitive.ObjectIDFromHex("000000000000")
	filter := bson.D{{Key: "boxid", Value: boxIdInt}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "ownerid", Value: noOwner}}}}

	res, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error while clearing the owner of this box"})
		//return
	}

	if res.ModifiedCount == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Box not found"})
		//return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Owner of the box successfully cleared"})
}

// preimenuj v get user boxes and acesses, naredi dodatno da samo boxe vrne
func GetUserBoxesAndAccesses(c *gin.Context) {
	var allBoxes []schemas.Box
	var usrid = c.Param("id")
	str, _ := primitive.ObjectIDFromHex(usrid)

	cur, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").Find(context.TODO(), bson.D{{Key: "ownerid", Value: str}})
	if err == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusInternalServerError, "Error")
	}

	var usernames [][]string
	for cur.Next(context.TODO()) {
		var elem schemas.Box
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
			//log.Fatal(err)
		}
		if len(elem.AccessIds) > 0 {
			var boxUsernames []string
			for _, id := range elem.AccessIds {
				user := schemas.User{}
				err := utils.CheckBase().Database("PametniPaketnik").Collection("users").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
				if err == mongo.ErrNoDocuments {
					continue // user not found, skip to next id
				} else if err != nil {
					fmt.Println(err)
					//log.Fatal(err)
				}
				boxUsernames = append(boxUsernames, user.Username)
			}
			usernames = append(usernames, boxUsernames)
		} else {
			usernames = append(usernames, []string{})
		}

		allBoxes = append(allBoxes, elem)
	}

	if len(allBoxes) == 0 {
		c.IndentedJSON(http.StatusBadRequest, "Error")
	}
	obj := bson.M{"allBoxes": allBoxes, "usernames": usernames}
	c.IndentedJSON(http.StatusOK, obj)
}

func GetUserBoxes(c *gin.Context) {
	var allBoxes []schemas.Box
	var usrid = c.Param("id")
	str, _ := primitive.ObjectIDFromHex(usrid)

	cur, err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").Find(context.TODO(), bson.D{{Key: "ownerid", Value: str}})
	if err == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusInternalServerError, "Error")
	}

	for cur.Next(context.TODO()) {
		var elem schemas.Box
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
			//log.Fatal(err)
		}

		// Ignore the AccessIds for this request
		elem.AccessIds = nil

		allBoxes = append(allBoxes, elem)
	}

	if len(allBoxes) == 0 {
		c.IndentedJSON(http.StatusBadRequest, "Error")
	}
	obj := bson.M{"allBoxes": allBoxes}
	c.IndentedJSON(http.StatusOK, obj)
}

func AuthenticateUser(c *gin.Context) {
	var requestData struct {
		UserID string
		BoxID  int
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//return
	}
	fmt.Println(requestData)
	str, _ := primitive.ObjectIDFromHex(requestData.UserID)

	var res schemas.Box
	err := utils.CheckBase().Database("PametniPaketnik").Collection("boxes").FindOne(context.TODO(), bson.D{{Key: "boxid", Value: requestData.BoxID}}).Decode(&res)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, "Error")
	}
	var result = false

	if str == res.OwnerId {
		result = true
		c.JSON(http.StatusOK, result)
		//return
	}

	for _, v := range res.AccessIds {
		if v == str {
			result = true
			c.JSON(http.StatusOK, result)
			//return
		}
	}

	//c.JSON(http.StatusForbidden, result)
}

func AddImage(c *gin.Context) {
	boxID := c.Param("id")

	// Save the uploaded BMP picture
	ctx := context.WithValue(context.Background(), "boxID", boxID)
	err := SaveBMPImage(ctx, c.Request)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error saving BMP picture"})
		return
	}

	err = runCSharpCompression(boxID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error running C# compression"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Image successfully added"})
}

func SaveBMPImage(ctx context.Context, req *http.Request) error {
	boxID, ok := ctx.Value("boxID").(string)
	if !ok {
		return fmt.Errorf("boxID not found in context")
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()

	if err := os.MkdirAll("images/boxPictures", os.ModePerm); err != nil {
		return err
	}

	fileName := filepath.Join("images/boxPictures", boxID+".bmp")

	dstFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, file)
	if err != nil {
		return err
	}

	return nil
}

func runCSharpCompression(fileName string) error {
	// Specify the path to the C# binary
	csharpBinary := "scripts/compression.exe"

	// Log the command being executed
	fmt.Printf("Executing command: %s %s\n", csharpBinary, fileName)

	// Run the C# binary and pass the filename as an argument
	cmd := exec.Command(csharpBinary, fileName)

	// Capture and log the standard output and standard error
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		// Log the error message
		fmt.Printf("Error running C# compression: %s\n", err)

		// Log the standard output and standard error
		fmt.Printf("C# Binary Output:\n%s\n", stdout.String())
		fmt.Printf("C# Binary Error:\n%s\n", stderr.String())
	}

	bmpFilePath := filepath.Join("images/boxPictures", fileName+".bmp")
	if err := os.Remove(bmpFilePath); err != nil {
		fmt.Printf("Error deleting BMP file: %s\n", err)
	}

	return err
}

func CheckBinFile(c *gin.Context) {
	boxID := c.Param("id")
	binFilePath := filepath.Join("images/boxPictures/", fmt.Sprintf("out_%s.bin", boxID))

	_, err := os.Stat(binFilePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{"exists": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": true})
}

func GetBMPImage(c *gin.Context) {
	boxID := c.Param("id")

	// Check if the .bin file exists for the given boxID
	binFilePath := filepath.Join("images/boxPictures/", "out_"+boxID+".bin")
	if _, err := os.Stat(binFilePath); os.IsNotExist(err) {
		// .bin file does not exist, return an error or handle it as needed
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Decompress the image
	decompressedImagePath, err := runCSharpDecompression(binFilePath, boxID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error running decompression"})
		return
	}

	// Serve the decompressed BMP image
	c.File(decompressedImagePath)

	// Delete the decompressed BMP image after serving
	if err := os.Remove(decompressedImagePath); err != nil {
		fmt.Printf("Error deleting decompressed image: %s\n", err)
	}
}

func runCSharpDecompression(binFilePath string, boxID string) (string, error) {
	// Specify the path to the C# decompression binary
	csharpDecompressionBinary := "scripts/Vaja2.exe"

	// Generate the decompressed image path
	decompressedImagePath := filepath.Join("images/boxPictures/", "decompressed_"+boxID+".bmp")
	fmt.Printf("Executing command: %s %s %s\n", csharpDecompressionBinary, binFilePath, decompressedImagePath)

	cmd := exec.Command(csharpDecompressionBinary, boxID, decompressedImagePath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		// Log the error message
		fmt.Printf("Error running C# decompression: %s\n", err)

		// Log the standard output and standard error
		fmt.Printf("C# Binary Output:\n%s\n", stdout.String())
		fmt.Printf("C# Binary Error:\n%s\n", stderr.String())
		return "", err
	}

	return decompressedImagePath, nil
}
