package azurerm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppServiceSourceControlToken(t *testing.T) {
	resourceName := "azurerm_app_service_source_control_token.test"
	token := strings.ToLower(acctest.RandString(41))
	tokenSecret := strings.ToLower(acctest.RandString(41))

	config := testAccAzureRMAppServiceSourceControlToken(token, tokenSecret)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSourceControlTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "GitHub"),
					resource.TestCheckResourceAttr(resourceName, "token", token),
					resource.TestCheckResourceAttr(resourceName, "token_secret", tokenSecret),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMAppServiceSourceControlToken(token, tokenSecret string) string {
	return fmt.Sprintf(`
resource "azurerm_app_service_source_control_token" "test" {
  type         = "GitHub"
  token        = "%s"
  token_secret = "%s"
}
`, token, tokenSecret)
}

func testCheckAzureRMAppServiceSourceControlTokenDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Web.BaseClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_source_control_token" {
			continue
		}

		scmType := rs.Primary.Attributes["type"]

		resp, err := conn.GetSourceControl(ctx, scmType)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}
