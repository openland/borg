package commands

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/kyokomi/emoji.v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/urfave/cli"
)

func doMapboxUpload(c *cli.Context) error {
	token := c.String("token")
	user := c.String("user")
	src := c.String("src")
	tileset := c.String("tileset")
	name := c.String("name")
	if token == "" {
		return cli.NewExitError("You should provide token", 1)
	}
	if name == "" {
		return cli.NewExitError("You should provide name", 1)
	}
	if user == "" {
		return cli.NewExitError("You should provide user", 1)
	}
	if src == "" {
		return cli.NewExitError("You should provide source file", 1)
	}
	if tileset == "" {
		return cli.NewExitError("You should provide destination tileset", 1)
	}

	// Request access token
	emoji.Println(":hammer: Requesting Upload Credentials")
	req, err := http.NewRequest("POST", "https://api.mapbox.com/uploads/v1/"+user+"/credentials?access_token="+token, bytes.NewBuffer(make([]byte, 0)))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		log.Println("Mapbox response Status:", resp.Status)
		log.Println("Mapbox response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("Mapbox response Body:", string(body))
		return cli.NewExitError("Unable to retreive upload credentials", 1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	response := make(map[string]string)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	emoji.Println(":hammer: Uploading to S3")
	// Credentials
	accessKeyID := response["accessKeyId"]
	secretAccessKey := response["secretAccessKey"]
	bucket := response["bucket"]
	key := response["key"]
	sessionToken := response["sessionToken"]
	// url := response["url"]
	conf := aws.Config{Region: aws.String("us-east-1"), Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, sessionToken)}

	//
	// File
	//

	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	//
	// Uploading
	//

	sess := session.New(&conf)
	svc := s3manager.NewUploader(sess)
	_, err = svc.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String("application/vnd.geo+json"),
		Body:        file,
	})
	if err != nil {
		return err
	}

	emoji.Println(":hammer: Commit changes")
	content := "{\"url\": \"http://" + bucket + ".s3.amazonaws.com/" + key + "\",\"tileset\": \"" + user + "." + tileset + "\",\"name\":\"" + name + "\"}"
	req, err = http.NewRequest("POST", "https://api.mapbox.com/uploads/v1/"+user+"?access_token="+token, bytes.NewBuffer([]byte(content)))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Status != "201 Created" {
		log.Println("Mapbox response Status:", resp.Status)
		log.Println("Mapbox response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("Mapbox response Body:", string(body))
		return cli.NewExitError("Unable to retreive commit upload", 1)
	}

	// curl -X POST -H "Content-Type: application/json" -H "Cache-Control: no-cache" -d '{
	// 	"url": "http://{bucket}.s3.amazonaws.com/{key}",
	// 	"tileset": "{username}.{tileset-name}"
	//   }' 'https://api.mapbox.com/uploads/v1/{username}?access_token=secret_access_token'

	return nil
}

func CreateMapboxCommands() []cli.Command {
	return []cli.Command{
		cli.Command{
			Name:  "mapbox",
			Usage: "Integration with MapBox",
			Subcommands: []cli.Command{
				cli.Command{
					Name:  "upload",
					Usage: "Upload tileset",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "token",
							Usage: "Mapbox Access Token",
						},
						cli.StringFlag{
							Name:  "user",
							Usage: "Mapbox User Name",
						},
						cli.StringFlag{
							Name:  "source, src",
							Usage: "Source geojson file",
						},
						cli.StringFlag{
							Name:  "tileset",
							Usage: "Tileset key",
						},
						cli.StringFlag{
							Name:  "name",
							Usage: "Tileset name",
						},
					},
					Action: doMapboxUpload,
				},
			},
		},
	}
}
