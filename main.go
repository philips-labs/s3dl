package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/cloudfoundry-community/gautocloud"
    "github.com/philips-software/gautocloud-connectors/hsdp"
    log "github.com/sirupsen/logrus"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/spf13/viper"
    "net/http"
    "os"

    "time"
)

var GitCommit = "deadbeaf"

func main() {
    viper.SetEnvPrefix("s3dl")
    viper.SetDefault("expire", 15)
    viper.AutomaticEnv()

    expire := viper.GetInt("expire")

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
    e.GET("/download", downloader(svc, expire))
    e.GET("/object/*", downloader(svc, expire))

    // Listen
    usePort := os.Getenv("PORT")
    if usePort == "" {
        usePort = "8080"
    }
    e.Start(":" + usePort)
}

// downloader generates HTTP 307 redirects to pre-signed bucket key URLs
func downloader(svc *hsdp.S3Client, expire int) echo.HandlerFunc {
    return func(e echo.Context) error {
        key := e.QueryParam("key")
        if key == "" {
            key = e.Param("_*")
        }
        log.Printf("Downloading: %s\n", key)
        req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
            Bucket: aws.String(svc.Bucket),
            Key:    aws.String(key),
        })
        str, err := req.Presign(time.Duration(expire) * time.Minute)
        if err != nil {
            return e.String(http.StatusBadRequest, err.Error())
        }
        return e.Redirect(http.StatusTemporaryRedirect, str)
    }
}