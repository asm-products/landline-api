package handlers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"os"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/gin-gonic/gin"
)

type UploadJSON struct {
	Name        string `json:"name" binding:"required"`
	Size        int    `json:"size" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

type Form struct {
	Key            string `json:"key"`
	ACL            string `json:"acl"`
	AWSAccessKeyId string `json:"AWSAccessKeyId"`
	CacheControl   string `json:"Cache-Control"`
	ContentType    string `json:"Content-Type"`
	Policy         string `json:"policy"`
	Signature      string `json:"signature"`
}

var (
	mimeTypesLower = map[string]string{
		".gif":  "image/gif",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".pdf":  "application/pdf",
		".png":  "image/png",
		".psd":  "image/psd",
	}
	extensions = reverseMap(mimeTypesLower)
)

func SignUpload(c *gin.Context) {
	user, err := GetUserFromContext(c)

	if err != nil {
		c.Fail(500, err)
	}

	var json UploadJSON
	c.Bind(&json)

	key := `/attachments/` + user.Id + `/` + uuid.NewUUID().String() + extensions[json.ContentType]

	f := &Form{
		Key:            key,
		ACL:            "public-read",
		AWSAccessKeyId: os.Getenv("AWS_ACCESS_KEY_ID"),
		CacheControl:   "max-age=31557600",
		ContentType:    json.ContentType,
	}

	f.build()

	href := "https://s3.amazonaws.com/" + os.Getenv("S3_BUCKET") + key

	c.JSON(200, gin.H{"form": f, "href": href})
}

func (f *Form) build() {
	f.Policy = f.policy()
	f.Signature = f.signature()
}

func (f *Form) policy() string {
	p := `{
    "expiration":"` + time.Now().Add(time.Minute*30).UTC().String() + `",
    "conditions": [
      { "bucket":"` + os.Getenv("S3_BUCKET") + `" },
      { "acl": "public-read" },
      { "cache-control": "max-age=31557600" },
      { "content-type":"` + f.ContentType + `" }
    ]
  }`

	return base64.StdEncoding.EncodeToString([]byte(p))
}

func (f *Form) signature() string {
	mac := hmac.New(sha1.New, []byte(os.Getenv("AWS_SECRET_ACCESS_KEY")))
	mac.Write([]byte(f.policy()))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func reverseMap(m map[string]string) map[string]string {
	reversedMap := make(map[string]string)
	for k, v := range m {
		reversedMap[v] = k
	}
	return reversedMap
}
