package controllers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"os"
	"strings"
	"time"
	"bytes"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/bson"
	database "github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/db"
	"github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/models"
	"github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/strategies"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

// HabitsController holds important data for the habits controller
type HabitsController struct {
	habitsDB *database.HabitsDB
	currSession *s3.S3
	strategy *strategies.MainStrategy
}
// Get serves as a GET request for either all or one specific habit
func (controller *HabitsController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) < 1 {
		habits := controller.getAll()
		if habits == nil {
			empty := make([]*models.Habit, 0)
			w.Header().Add("Content-Type", "application/json")
			encoder := json.NewEncoder(w)
			encoder.Encode(empty)
		}
		response := struct {
			Habits []*models.Habit
		} {
			habits,
		}
		w.Header().Add("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.Encode(response.Habits)
        return
	}
	habitID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the habit ID.")
		return
	}

	filter := bson.D{{"_id", habitID}}

	var habit models.Habit

	err2 := controller.habitsDB.GetByID(filter, &habit)

	if(err2 != nil) {
		println(err2.Error())
		w.WriteHeader(500)
		fmt.Fprint(w, "Error retrieving habit from database.")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(habit)
}
func (controller *HabitsController) getAll() []*models.Habit {
	var results []*models.Habit
	filter := bson.D{}
	results, err := controller.habitsDB.Get(filter)
	if(err != nil) {
		results = nil
	}
	return results
}
func (controller *HabitsController) initializeController(r *mux.Router) {
	r.HandleFunc("/{id}", controller.Get).Methods(http.MethodGet)
	r.HandleFunc("/", controller.Get).Methods(http.MethodGet)
	r.HandleFunc("/", controller.CreateHabit).Methods(http.MethodPost)
	r.HandleFunc("", controller.CreateHabit).Methods(http.MethodPost)
	r.HandleFunc("/{id}", controller.EditHabit).Methods(http.MethodPut)
	r.HandleFunc("/{id}", controller.DeleteHabit).Methods(http.MethodDelete)
	r.HandleFunc("/{id}/complete", controller.CompleteHabit).Methods(http.MethodPost)
	r.HandleFunc("/check", controller.HealthCheck).Methods(http.MethodPost)
}

// SetHabitsController sets the controller for the sets up the habits controllet
func SetHabitsController(r *mux.Router, db *mongo.Database) {
	habit := database.HabitsDB{Habits: db.Collection(os.Getenv("MONGO_DB"))}

	creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), "") 
	_, err := creds.Get() 
	if err != nil { 
	// handle error
	} 
	cfg := aws.NewConfig().WithRegion("us-west-1").WithCredentials(creds)
	
	habitsController := HabitsController{habitsDB: &habit}
	habitsController.currSession = s3.New(session.New(), cfg)
	habitsController.strategy = &strategies.MainStrategy{}
	habitsController.initializeController(r)
}
// CreateHabit serves to receive a POST request to create a new habit
func (controller *HabitsController) CreateHabit(w http.ResponseWriter, r *http.Request) {
	var maxBytes int64 = 64 * 1024 * 1024
	validImageFormats := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
		"image/gif": true,
	}

	// Parse multipart form data
	err := r.ParseMultipartForm(maxBytes)

	if err != nil {
		
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the data provided.")
		return
	}
	// Ensure that both the image and age are contained in multipartform
	if len(r.MultipartForm.File["image"]) == 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, "Missing image.")
		return
	}
	if len(r.MultipartForm.Value["title"]) == 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, "Missing title.")
		return
	}
	if len(r.MultipartForm.Value["type"]) == 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, "Missing type.")
		return
	}
	if len(r.MultipartForm.Value["difficulty"]) == 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, "Missing difficulty.")
		return
	}
	if len(r.MultipartForm.Value["userId"]) == 0 {
		w.WriteHeader(400)
		fmt.Fprint(w, "Missing user ID.")
		return
	}

	userId := r.MultipartForm.Value["userId"][0]

	// get age from image
	hType, err := strconv.Atoi(r.MultipartForm.Value["type"][0])

	if (err != nil || hType < -1 || hType > 1){
		w.WriteHeader(400)
		fmt.Fprint(w, "Either type was not a number or is not within the permitted range.")
		return
	}

	difficulty, err := strconv.Atoi(r.MultipartForm.Value["difficulty"][0])

	if (err != nil){
		w.WriteHeader(400)
		fmt.Fprint(w, "Difficulty was not a number.")
		return
	}
	// get image file header
	imFileHeader := r.MultipartForm.File["image"][0]
	im, err := imFileHeader.Open()
	
	if err != nil {
		
		w.WriteHeader(500)
		fmt.Fprint(w, "Error opening the image file.")
		return
	}
	defer im.Close()
	if _, exists := validImageFormats[imFileHeader.Header["Content-Type"][0]]; !exists {
		w.WriteHeader(400)
		fmt.Fprint(w, "Image file format not valid.")
		return
	}
	if imFileHeader.Size/1000000 > 5 {
		w.WriteHeader(413)
		fmt.Fprint(w, "La imagen pesa más de 5 MB.")
		return
	}

	imageURL := "/habits/"

	filename := imFileHeader.Filename
	filename = strings.Replace(filename, " ", "", -1)

	additionalNum := strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond), 10)

	imageURL += additionalNum + filename

	file, err := imFileHeader.Open()

	size := imFileHeader.Size
	buffer := make([]byte, size) // read file content to buffer 

	file.Read(buffer) 
	fileBytes := bytes.NewReader(buffer) 
	fileType := http.DetectContentType(buffer) 

	params := &s3.PutObjectInput{ 
		Bucket: aws.String(os.Getenv("AWS_BUCKET")), 
		Key: aws.String(imageURL), 
		Body: fileBytes, 
		ContentLength: aws.Int64(size), 
		ContentType: aws.String(fileType), 
	}

	_, err = controller.currSession.PutObject(params) 
	if err != nil { 
		w.WriteHeader(500)
		fmt.Fprint(w, "Error contacting AWS.")
		return
	} 

	habit := models.Habit{
		Title:     	r.MultipartForm.Value["title"][0],
		Type: 		hType,
		Difficulty: difficulty,
		UserId:		userId,
		Image:		os.Getenv("AWS_OBJECT_PREFIX")+"/"+os.Getenv("AWS_BUCKET")+imageURL,
		Score:      0,
	}

	result, err := controller.habitsDB.Insert(habit)

	if err != nil {
		
		_ = os.Remove(("/static" + imageURL))
		w.WriteHeader(500)
		fmt.Fprint(w, "Error uploading image to database.")
		return
	}
	habit.ID = result.InsertedID.(primitive.ObjectID)
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(habit)
}

// EditHabit serves to receive a PUT request to modify an existing habit
func (controller *HabitsController) EditHabit(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) < 1 {
		w.WriteHeader(500)
		fmt.Fprint(w, "No ID was provided.")
		return
	}

	habitID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the habit ID.")
		return
	}

	// Parse url-encoded data
	err = r.ParseForm()

	if err != nil {
		
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the data provided.")
		return
	}

	filter := bson.D{{"_id", habitID}}

	var existing models.Habit

	err = controller.habitsDB.GetByID(filter, &existing)

	if err != nil {
		
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while contacting database.")
		return
	}

	newTitle := existing.Title

	newType := existing.Type

	newDiff := existing.Difficulty

	newValues := new(models.Habit) // Person being a struct type
    decoder := schema.NewDecoder()

    err = decoder.Decode(newValues, r.Form)
    if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the input form.")
		return
    }
	
	if newValues.Title != existing.Title && len(r.Form["title"]) > 0 {
		newTitle = newValues.Title
	}
	if newValues.Type != existing.Type && len(r.Form["type"]) > 0 {
		newType = newValues.Type
	}
	if newValues.Difficulty != existing.Difficulty && len(r.Form["difficulty"]) > 0{
		newDiff = newValues.Difficulty
	}

	updateDoc := bson.D{{"$set", bson.D{{"title", newTitle},
		{"type", newType},
		{"difficulty", newDiff}}}}

	_, err = controller.habitsDB.UpdateOne(filter, updateDoc)

	if err != nil {
		
		w.WriteHeader(500)
		fmt.Fprint(w, "Error updating data in the database.")
		return
	}

	var result models.Habit

	err = controller.habitsDB.GetByID(filter, &result)

	if err != nil {
		
		w.WriteHeader(500)
		fmt.Fprint(w, "Error retrieving new values.")
		return
	}
	
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(result)
}

func (controller *HabitsController) DeleteHabit(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) < 1 {
		w.WriteHeader(500)
		fmt.Fprint(w, "No ID was provided.")
		return
	}

	habitID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the habit ID.")
		return
	}

	filter := bson.D{{"_id", habitID}}

	var existing models.Habit

	err = controller.habitsDB.GetByID(filter, &existing)

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while retrieving habit from database.")
		return
	}

	_, err = controller.habitsDB.DeleteOne(filter)

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while deleting habit from database.")
		return
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(strings.Replace(existing.Image, os.Getenv("AWS_OBJECT_PREFIX")+"/"+os.Getenv("AWS_BUCKET"), "", -1)),
	}

	_, err = controller.currSession.DeleteObject(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				w.WriteHeader(500)
				fmt.Fprint(w, aerr.Error())
				return
			}
		} else {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, "OK")
}

func (controller * HabitsController) CompleteHabit(w http.ResponseWriter, r * http.Request){
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) < 1 {
		w.WriteHeader(500)
		fmt.Fprint(w, "No ID was provided.")
		return
	}
	
	habitID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the habit ID.")
		return
	}

	err = r.ParseForm()

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Error parsing form.")
		return
	}

	if len(r.Form["completionStatus"]) < 1 {
		w.WriteHeader(500)
		fmt.Fprint(w, "No completion status was provided.")
		return
	}

	completionStatus, err := strconv.Atoi(r.Form["completionStatus"][0])

	if err != nil || completionStatus > 1 || completionStatus < 0{
		
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while intepreting the completion status: either not a number or not within range [0,1].")
		return
	}

	var habit models.Habit

	filter := bson.D{{"_id", habitID}}

	err = controller.habitsDB.GetByID(filter, &habit)

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Error retrieving habit from the database.")
		return
	}

	if (habit.Type == 1 && completionStatus == 0) || (habit.Type == -1 && completionStatus == 1) {
		w.WriteHeader(400)
		fmt.Fprint(w, "Incompatible habit type and completion status.")
		return
	}

	newScore := controller.strategy.CalculateScore(habit.Score, completionStatus, habit.Difficulty)

	updateDoc := bson.D{{"$set", bson.D{{"score", newScore}}}}

	_, err = controller.habitsDB.UpdateOne(filter, updateDoc)

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Error updating data in the database.")
		return
	}

	var result models.Habit

	err = controller.habitsDB.GetByID(filter, &result)

	if err != nil {
		
		w.WriteHeader(500)
		fmt.Fprint(w, "Error retrieving new values.")
		return
	}
	
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(result)

}

func (controller * HabitsController) HealthCheck(w http.ResponseWriter, r * http.Request){
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, "OK")
}