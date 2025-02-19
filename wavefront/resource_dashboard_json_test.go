package wavefront

import (
	"fmt"
	"testing"

	"github.com/WavefrontHQ/go-wavefront-management-api/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccWavefrontDashboardJson_Basic(t *testing.T) {
	var record wavefront.Dashboard
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardJSONDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboardJSONBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardJSONExists("wavefront_dashboard_json.test_dashboard_json", &record),
					testAccCheckWavefrontDashboardJSONAttributes(&record),

					resource.TestCheckResourceAttr(
						"wavefront_dashboard_json.test_dashboard_json", "id", "tftestimport"),
				),
			},
		},
	})
}

func TestAccWavefrontDashboardJson_Updated(t *testing.T) {
	var record wavefront.Dashboard
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardJSONDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboardJSONBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardJSONExists("wavefront_dashboard_json.test_dashboard_json", &record),
					testAccCheckWavefrontDashboardJSONAttributes(&record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard_json.test_dashboard_json", "id", "tftestimport"),
				),
			},
			{
				Config: testAccCheckWavefrontDashboardJSONNewValue(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardJSONExists("wavefront_dashboard_json.test_dashboard_json", &record),
					testAccCheckWavefrontDashboardJSONAttributesUpdated(&record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard_json.test_dashboard_json", "id", "tftestimport"),
				),
			},
		},
	})
}

func TestAccWavefrontDashboardJson_Multiple(t *testing.T) {
	var record wavefront.Dashboard

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontDashboardJSONDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontDashboardJSONMultiple(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontDashboardJSONExists("wavefront_dashboard_json.test_dashboard_1", &record),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard_json.test_dashboard_1", "id", "test_dashboard_1"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard_json.test_dashboard_2", "id", "test_dashboard_2"),
					resource.TestCheckResourceAttr(
						"wavefront_dashboard_json.test_dashboard_3", "id", "test_dashboard_3"),
				),
			},
		},
	})
}

func testAccCheckWavefrontDashboardJSONDestroy(s *terraform.State) error {

	dashboards := testAccProvider.Meta().(*wavefrontClient).client.Dashboards()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "wavefront_dashboard" {
			continue
		}

		results, err := dashboards.Find(
			[]*wavefront.SearchCondition{
				{
					Key:            "id",
					Value:          rs.Primary.ID,
					MatchingMethod: "EXACT",
				},
			})
		if err != nil {
			return fmt.Errorf("error finding Wavefront Dashboard. %s", err)
		}
		if len(results) > 0 {
			return fmt.Errorf("dashboard still exists")
		}
	}

	return nil
}

func testAccCheckWavefrontDashboardJSONAttributes(dashboard *wavefront.Dashboard) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if dashboard.Name != "Terraform Test Dashboard Json" {
			return fmt.Errorf("bad value: %s", dashboard.Name)
		}

		return nil
	}
}

func testAccCheckWavefrontDashboardJSONAttributesUpdated(dashboard *wavefront.Dashboard) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if dashboard.Name != "Terraform Test Dashboard Json Updated" {
			return fmt.Errorf("bad value: %s", dashboard.Name)
		}

		return nil
	}
}

func testAccCheckWavefrontDashboardJSONExists(n string, dashboard *wavefront.Dashboard) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		dash := wavefront.Dashboard{
			ID: rs.Primary.ID,
		}

		dashboards := testAccProvider.Meta().(*wavefrontClient).client.Dashboards()
		err := dashboards.Get(&dash)
		if err != nil {
			return fmt.Errorf("did not find Dashboard with id %s, %s", rs.Primary.ID, err)
		}
		*dashboard = dash
		return nil
	}
}

func testAccCheckWavefrontDashboardJSONBasic() string {
	return `
data "wavefront_default_user_group" "everyone" { 
}

resource "wavefront_dashboard_json" "test_dashboard_json" {
	dashboard_json = <<EOF
{
  "name": "Terraform Test Dashboard Json",
  "description": "a",
  "eventFilterType": "BYCHART",
  "eventQuery": "",
  "defaultTimeWindow": "",
  "url": "tftestimport",
  "displayDescription": false,
  "displaySectionTableOfContents": true,
  "displayQueryParameters": false,
  "acl": {
    "canModify": [
       "${data.wavefront_default_user_group.everyone.group_id}"
     ]
  },
  "sections": [
    {
      "name": "section 1",
      "rows": [
        {
          "charts": [
            {
              "name": "chart 1",
              "sources": [
                {
                  "name": "source 1",
                  "query": "ts()",
                  "scatterPlotSource": "Y",
                  "querybuilderEnabled": false,
                  "sourceDescription": ""
                }
              ],
              "units": "someunit",
              "base": 0,
              "noDefaultEvents": false,
              "interpolatePoints": false,
              "includeObsoleteMetrics": false,
              "description": "This is chart 1, showing something",
              "chartSettings": {
                "type": "markdown-widget",
                "max": 100,
                "expectedDataSpacing": 120,
                "windowing": "full",
                "windowSize": 10,
                "autoColumnTags": false,
                "columnTags": "deprecated",
                "tagMode": "all",
                "numTags": 2,
                "customTags": [
                  "tag1",
                  "tag2"
                ],
                "groupBySource": true,
                "y1Max": 100,
                "y1Units": "units",
                "y0ScaleSIBy1024": true,
                "y1ScaleSIBy1024": true,
                "y0UnitAutoscaling": true,
                "y1UnitAutoscaling": true,
                "fixedLegendEnabled": true,
                "fixedLegendUseRawStats": true,
                "fixedLegendPosition": "RIGHT",
                "fixedLegendDisplayStats": [
                  "stat1",
                  "stat2"
                ],
                "fixedLegendFilterSort": "TOP",
                "fixedLegendFilterLimit": 1,
                "fixedLegendFilterField": "CURRENT",
                "plainMarkdownContent": "markdown content"
              },
              "chartAttributes": {
                "dashboardLayout": {
                  "x": 0,
                  "y": 0,
                  "w": 4,
                  "h": 7
                }
              },
              "summarization": "MEAN"
            }
          ],
          "heightFactor": 50
        }
      ]
    }
  ],
  "parameterDetails": {
    "param": {
      "hideFromView": false,
      "description": null,
      "allowAll": null,
      "tagKey": null,
      "queryValue": null,
      "dynamicFieldType": null,
      "reverseDynSort": null,
      "parameterType": "SIMPLE",
      "label": "test",
      "defaultValue": "Label",
      "valuesToReadableStrings": {
        "Label": "test"
      },
      "selectedLabel": "Label",
      "value": "test"
    }
  },
  "tags" :{
    "customerTags":  ["terraform"]
  }
}
EOF
}
`
}

func testAccCheckWavefrontDashboardJSONNewValue() string {
	return `
data "wavefront_default_user_group" "everyone" { 
}

resource "wavefront_dashboard_json" "test_dashboard_json" {
	dashboard_json = <<EOF
{
  "name": "Terraform Test Dashboard Json Updated",
  "description": "a",
  "eventFilterType": "BYCHART",
  "eventQuery": "",
  "defaultTimeWindow": "",
  "url": "tftestimport",
  "displayDescription": false,
  "displaySectionTableOfContents": true,
  "displayQueryParameters": false,
  "acl": {
    "canModify": [
       "${data.wavefront_default_user_group.everyone.group_id}"
     ]
  },
  "sections": [
    {
      "name": "section 1",
      "rows": [
        {
          "charts": [
            {
              "name": "chart 1",
              "sources": [
                {
                  "name": "source 1",
                  "query": "ts()",
                  "scatterPlotSource": "Y",
                  "querybuilderEnabled": false,
                  "sourceDescription": ""
                }
              ],
              "units": "someunit",
              "base": 0,
              "noDefaultEvents": false,
              "interpolatePoints": false,
              "includeObsoleteMetrics": false,
              "description": "This is chart 1, showing something",
              "chartSettings": {
                "type": "markdown-widget",
                "max": 100,
                "expectedDataSpacing": 120,
                "windowing": "full",
                "windowSize": 10,
                "autoColumnTags": false,
                "columnTags": "deprecated",
                "tagMode": "all",
                "numTags": 2,
                "customTags": [
                  "tag1",
                  "tag2"
                ],
                "groupBySource": true,
                "y1Max": 100,
                "y1Units": "units",
                "y0ScaleSIBy1024": true,
                "y1ScaleSIBy1024": true,
                "y0UnitAutoscaling": true,
                "y1UnitAutoscaling": true,
                "fixedLegendEnabled": true,
                "fixedLegendUseRawStats": true,
                "fixedLegendPosition": "RIGHT",
                "fixedLegendDisplayStats": [
                  "stat1",
                  "stat2"
                ],
                "fixedLegendFilterSort": "TOP",
                "fixedLegendFilterLimit": 1,
                "fixedLegendFilterField": "CURRENT",
                "plainMarkdownContent": "markdown content"
              },
              "chartAttributes": {
                "dashboardLayout": {
                  "x": 0,
                  "y": 0,
                  "w": 4,
                  "h": 7
                }
              },
              "summarization": "MEAN"
            }
          ],
          "heightFactor": 50
        }
      ]
    }
  ],
  "parameterDetails": {
    "param": {
      "hideFromView": false,
      "parameterType": "SIMPLE",
      "label": "test",
      "defaultValue": "Label",
      "valuesToReadableStrings": {
        "Label": "test"
      },
      "selectedLabel": "Label",
      "value": "test"
    }
  },
  "tags" :{
    "customerTags":  ["terraform"]
  }
}
EOF
}
`
}

func testAccCheckWavefrontDashboardJSONMultiple() string {
	return `
data "wavefront_default_user_group" "everyone" { 
}

resource "wavefront_dashboard_json" "test_dashboard_1" {
  dashboard_json = <<EOF
{
  "name": "test_dashboard_1",
  "eventFilterType": "BYCHART",
  "url": "test_dashboard_1",
  "displayDescription": false,
  "displaySectionTableOfContents": true,
  "displayQueryParameters": false,
  "acl": {
    "canModify": [
       "${data.wavefront_default_user_group.everyone.group_id}"
     ]
  },
  "sections": [
    {
      "name": "New Section",
      "rows": []
    }
  ],
  "parameterDetails": {}
}
EOF
}
resource "wavefront_dashboard_json" "test_dashboard_2" {
  dashboard_json = <<EOF
{
  "name": "test_dashboard_2",
  "eventFilterType": "BYCHART",
  "url": "test_dashboard_2",
  "displayDescription": false,
  "displaySectionTableOfContents": true,
  "displayQueryParameters": false,
  "acl": {
    "canModify": [
       "${data.wavefront_default_user_group.everyone.group_id}"
     ]
  },
  "sections": [
    {
      "name": "New Section",
      "rows": []
    }
  ],
  "parameterDetails": {}
}
EOF
}
resource "wavefront_dashboard_json" "test_dashboard_3" {
  dashboard_json = <<EOF
{
  "name": "test_dashboard_3",
  "eventFilterType": "BYCHART",
  "url": "test_dashboard_3",
  "displayDescription": false,
  "displaySectionTableOfContents": true,
  "displayQueryParameters": false,
  "acl": {
    "canModify": [
       "${data.wavefront_default_user_group.everyone.group_id}"
     ]
  },
  "sections": [
    {
      "name": "New Section",
      "rows": []
    }
  ],
  "parameterDetails": {}
}
EOF
}
`
}
