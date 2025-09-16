package pub

import (
	"errors"
	"fmt"
)

var (
	ErrAssetMissingObjects        = errors.New("asset: could not find objects/descriptors")
	ErrAssetDescriptorMissingName = errors.New("asset descriptor: missing Name")
)

type ErrAssetMismatchedDescriptorType struct {
	AssetType      string
	DescriptorName string
	DescriptorType string
}

func (e ErrAssetMismatchedDescriptorType) Error() string {
	return fmt.Sprintf("asset descriptor \"%s\" mismatched file type: got \"%s\" but expected \"%s\" (the top-level asset element type)", e.DescriptorName, e.DescriptorType, e.AssetType)
}

// Asset represents a media element such as an image or video. Supports specifying multiple [AssetDescriptor]s which will be used as fallback formats (in the specified order) when the asset is not supported by the application.
type Asset struct {
	Objects         []AssetDescriptor
	AlternativeText string
	Caption         Content
}

func (a *Asset) MainDescriptor() AssetDescriptor {
	return a.Objects[0]
}

func (a *Asset) OtherDescriptors() []AssetDescriptor {
	if len(a.Objects) == 1 {
		return nil
	}

	var descriptors []AssetDescriptor

	for i := 1; i < len(a.Objects); i++ {
		object := &a.Objects[i]

		descriptors = append(descriptors, *object)
	}

	return descriptors
}

func (a *Asset) EnsureValid() error {
	if len(a.Objects) == 0 {
		return ErrAssetMissingObjects
	}

	main := a.MainDescriptor()
	if err := main.CheckValid(); err != nil {
		return err
	}
	mainType := main.Type

	otherDescriptors := a.OtherDescriptors()
	for i := range otherDescriptors {
		o := &otherDescriptors[i]

		if mainType != "" && o.Type != "" && o.Type != mainType {
			return ErrAssetMismatchedDescriptorType{
				AssetType:      mainType,
				DescriptorType: o.Type,
				DescriptorName: o.Name,
			}
		}
	}

	return nil
}

// AssetDescriptor represents an individual file format of a media element.
type AssetDescriptor struct {
	Name   string
	Type   string
	Format string
}

func (a *AssetDescriptor) CheckValid() error {
	if a.Name == "" {
		return ErrAssetDescriptorMissingName
	}

	return nil
}
