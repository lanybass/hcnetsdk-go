package hcnetsdk

import (
	"io/ioutil"
	"os"
	"testing"
)

func Test(t *testing.T) {
	sdk := new(HCNetSDK)
	sdk.Init()
	err := sdk.Login("192.168.0.109", 8000, "admin", "HikQPUWBS")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		ok := sdk.Logout()
		if !ok {
			t.Fatal("fail in logout")
		}
	}()
	//d, _ := os.Getwd()
	err, b := sdk.CaptureJPEGPictureNew(&JPEGParam{
		PicSize:    9,
		PicQuality: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("pic.jpg",b,os.FileMode(0644))
}
