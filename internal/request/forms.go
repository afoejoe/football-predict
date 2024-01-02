package request

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/go-playground/form/v4"
)

var decoder = form.NewDecoder()

func DecodeForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	return decodeURLValues(r.Form, dst)
}

func DecodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	return decodeURLValues(r.PostForm, dst)
}

func DecodeQueryString(r *http.Request, dst any) error {
	return decodeURLValues(r.URL.Query(), dst)
}

func decodeURLValues(v url.Values, dst any) error {
	decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return time.Parse("2006-01-02T15:04-0700", vals[0]+"+0300")
	}, time.Time{})

	err := decoder.Decode(dst, v)

	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
	}

	return err
}
