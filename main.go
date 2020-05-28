package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/cloudfoundry-community/gautocloud"
    "github.com/philips-software/gautocloud-connectors/hsdp"
    log "github.com/sirupsen/logrus"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "net/http"
    "os"

    "time"
)

var GitCommit = "deadbeaf"

func main() {
    // S3 Bucket
    var svc *hsdp.S3Client
    err := gautocloud.Inject(&svc)

    if err != nil {
        log.Printf("error: %v\n", err)
        return
    }
    
    // Web server
    e := echo.New()
    e.Use(middleware.Logger())
    e.GET("/download", downloader(svc))

    // Listen
    usePort := os.Getenv("PORT")
    if usePort == "" {
        usePort = "8080"
    }
    e.Start(":" + usePort)
}

// downloader generates HTTP 307 redirects to pre-signed bucket key URLs
func downloader(svc *hsdp.S3Client) echo.HandlerFunc {
    return func(e echo.Context) error {
        key := e.QueryParam("key")
        log.Printf("Downloading: %s\n", key)
        req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
            Bucket: aws.String(svc.Bucket),
            Key:    aws.String(key),
        })
        str, err := req.Presign(15 * time.Minute)
        if err != nil {
            return e.String(http.StatusBadRequest, err.Error())
        }
        return e.Redirect(http.StatusTemporaryRedirect, str)
    }
}