package imgs

import (
)

type ProfileImage struct {
	ImageBytes		[]byte 		`bson:"encodedData"		json:"encodedData"`
}