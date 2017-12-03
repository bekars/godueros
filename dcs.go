package godueros

import (
	"bytes"
	"mime/multipart"
	"fmt"
	"io"
	"io/ioutil"
	"net/textproto"
)

type DuDCS struct {

}

func (dcs *DuDCS)GetMultiPartData(event string, audio []byte, boundary string) (body *bytes.Buffer, err error) {
	body = &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	if boundary != "" {
		writer.SetBoundary(boundary)
	}

	var w io.Writer
	if w, err = createJsonPart(writer,"metadata"); err != nil {
		fmt.Printf("Create Event Form Field Err: %v\n", err)
	}

	if _, err = w.Write([]byte(event)); err != nil {
		fmt.Printf("Write Event Form Field Err: %v\n", err)
	}

	if w, err = createAudioPart(writer,"audio"); err != nil {
		fmt.Printf("Create Audio Form Field Err: %v\n", err)
	}

	if _, err = w.Write(audio); err != nil {
		fmt.Printf("Write Audio Form Field Err: %v\n", err)
	}

	return body, err
}

func createJsonPart(writer *multipart.Writer, fieldname string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"`, fieldname))
	h.Set("Content-Type", "application/json; charset=UTF-8")
	return writer.CreatePart(h)
}

func createAudioPart(writer *multipart.Writer, fieldname string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"`, fieldname))
	h.Set("Content-Type", "application/octet-stream")
	return writer.CreatePart(h)
}


func (dcs *DuDCS)ReadMultiPartData(r io.Reader, boundary string) {
	reader := multipart.NewReader(r, boundary)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			CheckErr(err)
		}
		content, err := ioutil.ReadAll(part)
		if err != nil {
			CheckErr(err)
		}
		fmt.Println("Form name: ", part.FormName(), "'s", "Content is: ", string(content))

		if part.FormName() != "metadata" {
			ioutil.WriteFile("duvoice.wav", content, 0666)
		}
	}
}
