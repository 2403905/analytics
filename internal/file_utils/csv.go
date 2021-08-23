package file_utils

import (
	"io/ioutil"

	"github.com/jszwec/csvutil"
)

func ReadModel(model interface{}, dest string) error {
	b, err := ioutil.ReadFile(dest)
	if err != nil {
		return err
	}
	if err := csvutil.Unmarshal(b, model); err != nil {
		return err
	}
	return nil
}
