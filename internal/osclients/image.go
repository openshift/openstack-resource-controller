/*
Copyright 2021 The ORC Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package osclients

import (
	"context"
	"fmt"
	"io"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/imagedata"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/imageimport"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
)

type ImageClient interface {
	ListImages(listOpts images.ListOptsBuilder) ([]images.Image, error)
	GetImage(id string) (*images.Image, error)
	CreateImage(ctx context.Context, createOpts images.CreateOptsBuilder) (*images.Image, error)
	DeleteImage(ctx context.Context, id string) error
	UploadData(ctx context.Context, id string, data io.Reader) error
	GetImportInfo(ctx context.Context) (*imageimport.ImportInfo, error)
	CreateImport(ctx context.Context, id string, createOpts imageimport.CreateOptsBuilder) error
}

type imageClient struct{ client *gophercloud.ServiceClient }

// NewImageClient returns a new glance client.
func NewImageClient(providerClient *gophercloud.ProviderClient, providerClientOpts *clientconfig.ClientOpts) (ImageClient, error) {
	images, err := openstack.NewImageV2(providerClient, gophercloud.EndpointOpts{
		Region:       providerClientOpts.RegionName,
		Availability: clientconfig.GetEndpointType(providerClientOpts.EndpointType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create image service client: %v", err)
	}

	return imageClient{images}, nil
}

func (c imageClient) ListImages(listOpts images.ListOptsBuilder) ([]images.Image, error) {
	pages, err := images.List(c.client, listOpts).AllPages(context.TODO())
	if err != nil {
		return nil, err
	}
	return images.ExtractImages(pages)
}

func (c imageClient) GetImage(id string) (*images.Image, error) {
	image := &images.Image{}
	err := images.Get(context.TODO(), c.client, id).ExtractInto(image)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (c imageClient) CreateImage(ctx context.Context, createOpts images.CreateOptsBuilder) (*images.Image, error) {
	image, err := images.Create(ctx, c.client, createOpts).Extract()
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (c imageClient) DeleteImage(ctx context.Context, id string) error {
	return images.Delete(ctx, c.client, id).ExtractErr()
}

func (c imageClient) UploadData(ctx context.Context, id string, data io.Reader) error {
	return imagedata.Upload(ctx, c.client, id, data).ExtractErr()
}

func (c imageClient) GetImportInfo(ctx context.Context) (*imageimport.ImportInfo, error) {
	return imageimport.Get(ctx, c.client).Extract()
}

func (c imageClient) CreateImport(ctx context.Context, id string, createOpts imageimport.CreateOptsBuilder) error {
	return imageimport.Create(ctx, c.client, id, createOpts).ExtractErr()
}

type imageErrorClient struct{ error }

// NewImageErrorClient returns an ImageClient in which every method returns the given error.
func NewImageErrorClient(e error) ImageClient {
	return imageErrorClient{e}
}

func (e imageErrorClient) ListImages(_ images.ListOptsBuilder) ([]images.Image, error) {
	return nil, e.error
}

func (e imageErrorClient) GetImage(_ string) (*images.Image, error) {
	return nil, e.error
}

func (e imageErrorClient) CreateImage(_ context.Context, _ images.CreateOptsBuilder) (*images.Image, error) {
	return nil, e.error
}

func (e imageErrorClient) DeleteImage(_ context.Context, _ string) error {
	return e.error
}

func (e imageErrorClient) UploadData(_ context.Context, _ string, _ io.Reader) error {
	return e.error
}

func (e imageErrorClient) GetImportInfo(_ context.Context) (*imageimport.ImportInfo, error) {
	return nil, e.error
}

func (e imageErrorClient) CreateImport(_ context.Context, _ string, _ imageimport.CreateOptsBuilder) error {
	return e.error
}
