package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Name       string `json:"Name"`
	Model      string `json:"Model"`
}

type ResponseBody struct {
	Status string  `json:"status"`
	Body string `json:"body"`
}

const BUCKET = "xxxx"
const region = "us-east-1"

func CreatePost(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	client := Client{}
	responseBody:= ResponseBody{}
	_ = json.NewDecoder(r.Body).Decode(&client)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
    if err != nil{
    	responseBody.Body = "Internal Error,"+err.Error()
    	responseBody.Status = "Error"
		json.NewEncoder(w).Encode(&responseBody)
	}
	uid, err := uuid.NewV4()
	if err != nil{
		responseBody.Body = "Internal Error,"+err.Error()
		responseBody.Status = "Error"
		json.NewEncoder(w).Encode(&responseBody)
	}
	err = WriteFiles(client,uid.String()+".json",sess)
    if err != nil{
		responseBody.Body = "Can't post client info to s3,"+err.Error()
		responseBody.Status = "Error"
		json.NewEncoder(w).Encode(&responseBody)
	}else{
		responseBody.Body = "Upload successfully"
		responseBody.Status = "Success"
		json.NewEncoder(w).Encode(&responseBody)
	}
}

// write files to s3
func WriteFiles(client Client, path string, svc *session.Session) error {
	var result []byte
	data, dataErr := json.Marshal(client)
	if dataErr != nil {
		return dataErr
	}
	result = append(result, data...)

	_, err := s3manager.NewUploader(svc).Upload(&s3manager.UploadInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(path),
		Body:   bytes.NewReader(result),
	})
	if err != nil {
		return err
	}
	return nil
}


func ReadGet(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	responseBody:= ResponseBody{}
	params := mux.Vars(r)
	client := Client{}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil{
		responseBody.Body = "Internal Error,"+err.Error()
		responseBody.Status = "Error"
		json.NewEncoder(w).Encode(&responseBody)
	}
	res, resErr := ReadFile(params["path"],sess)
	if resErr != nil{
		responseBody.Body = "Can't get info from s3,"+resErr.Error()
		responseBody.Status = "Error"
		json.NewEncoder(w).Encode(&responseBody)
	}else {
		_= json.Unmarshal([]byte(res),client)
		responseBody.Body = res
		responseBody.Status = "Success"
		json.NewEncoder(w).Encode(&responseBody)
	}
}

// read file content from s3 based on path
func ReadFile(filePath string, sess *session.Session)(string,error){

	// Create s3 connection
	svc := s3.New(sess)
	param := &s3.GetObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(filePath),
	}
	result, err := svc.GetObject(param)
	if err != nil{
		return "",err
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil{
		return "",err
	}
	return string(body),nil
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/model", CreatePost).Methods("POST")
	router.HandleFunc("/read/{path}", ReadGet).Methods("GET")
	http.ListenAndServe(":8000", router)
	fmt.Println("The server is starting...")
}
