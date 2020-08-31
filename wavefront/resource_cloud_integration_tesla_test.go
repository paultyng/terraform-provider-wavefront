package wavefront

import (
	"testing"

	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccWavefrontCloudIntegrationTesla_Basic(t *testing.T) {
	var record wavefront.CloudIntegration
	resourcePrefix := "wavefront_cloud_integration_tesla.tesla"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontCloudIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontCloudIntegrationTeslaBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontCloudIntegrationExists(resourcePrefix, &record),
					testAccCheckWavefrontCloudIntegrationAttributes(&record, wfTesla),
					// Check the attributes...
					testAccCheckWavefrontCloudIntegrationResourceAttributes(resourcePrefix, wfTesla),
					resource.TestCheckResourceAttr(resourcePrefix, "email", "email@example.com"),
					resource.TestCheckResourceAttr(resourcePrefix, "password", "password"),
				),
			},
		},
	})
}

func TestAccWavefrontCloudIntegrationTesla_BasicChanged(t *testing.T) {
	var record wavefront.CloudIntegration
	resourcePrefix := "wavefront_cloud_integration_tesla.tesla"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontCloudIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontCloudIntegrationTeslaBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontCloudIntegrationExists(resourcePrefix, &record),
					testAccCheckWavefrontCloudIntegrationAttributes(&record, wfTesla),
					// Check the attributes...
					testAccCheckWavefrontCloudIntegrationResourceAttributes(resourcePrefix, wfTesla),
					resource.TestCheckResourceAttr(resourcePrefix, "email", "email@example.com"),
					resource.TestCheckResourceAttr(resourcePrefix, "password", "password"),
				),
			},
			{
				Config: testAccCheckWavefrontCloudIntegrationTeslaBasicChanged(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontCloudIntegrationExists(resourcePrefix, &record),
					testAccCheckWavefrontCloudIntegrationAttributes(&record, wfTesla),
					// Check the attributes...
					testAccCheckWavefrontCloudIntegrationResourceAttributes(resourcePrefix, wfTesla),
					resource.TestCheckResourceAttr(resourcePrefix, "email", "email@example.com"),
					resource.TestCheckResourceAttr(resourcePrefix, "password", "password2"),
				),
			},
		},
	})
}

func testAccCheckWavefrontCloudIntegrationTeslaBasic() string {
	return `
resource "wavefront_cloud_integration_tesla" "tesla" {
  name              = "Test Integration"
  force_save        = true
  additional_tags = {
    "tag1" = "value1"
    "tag2" = "value2"
  }
  email    = "email@example.com"
  password = "password"
}
`
}

func testAccCheckWavefrontCloudIntegrationTeslaBasicChanged() string {
	return `
resource "wavefront_cloud_integration_tesla" "tesla" {
  name              = "Test Integration"
  force_save        = true
  additional_tags = {
    "tag1" = "value1"
    "tag2" = "value2"
  }
  email    = "email@example.com"
  password = "password2"
}
`
}
