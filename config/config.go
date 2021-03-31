package config

import (
	"github.com/nikhilsbhat/unpackker/pkg/packer"
	"github.com/nikhilsbhat/unpackker/pkg/unpacker"
)

//PackkerInput holds fields required to perform packing of an asset
type PackkerInput struct {
	PackData *packer.PackkerInput `json:"packdata"`
}

//UnPackkerInput holds fields reuired to perform unpacking of an asset
type UnPackkerInput struct {
	UnpackData *unpacker.UnPackkerInput `json:"unpackdata"`
}

//DeleteInput holds fields required to delete an asset
type DeleteInput struct {
	Assetpath string `json:"assetpath"`
}
