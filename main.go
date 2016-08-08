package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang.org/x/text/unicode/norm"
)

type Strings map[string]string

func Encode(s string) Strings {
	ret := map[string]string{}

	ret["origin"] = s

	ret["nfd"] = string(norm.NFD.Bytes([]byte(s)))
	ret["nfd-s"] = string(norm.NFD.Bytes([]byte(s + "-nfd")))

	ret["nfc"] = string(norm.NFC.Bytes([]byte(s)))
	ret["nfc-s"] = string(norm.NFC.Bytes([]byte(s + "-nfc")))

	ret["nfkd"] = string(norm.NFKD.Bytes([]byte(s)))
	ret["nfkd-s"] = string(norm.NFKD.Bytes([]byte(s + "-nfkd")))

	ret["nfkc"] = string(norm.NFKC.Bytes([]byte(s)))
	ret["nfkc-s"] = string(norm.NFKC.Bytes([]byte(s + "-nfkc")))

	return ret
}

func Upload(src, dst, bucket string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(dst),
		Body:   file,
	})

	return err
}

func main() {
	bucket := flag.String("b", "", "bucket name")
	prefixDir := flag.String("pd", "", "prefix dir")
	src := flag.String("s", "", "src file")
	flag.Parse()

	fmt.Println(*bucket)
	fmt.Println(*prefixDir)
	fmt.Println(*src)
	info, err := os.Stat(*src)
	if err != nil {
		panic(err)
	}

	if info.IsDir() {
		panic(errors.New("dir"))
	}

	eNames := Encode(info.Name())
	for k, v := range eNames {
		if strings.Contains(k, "-s") {
			if err := Upload(*src, *prefixDir+"/"+v, *bucket); err != nil {
				panic(err)
			}
		}
	}
}
