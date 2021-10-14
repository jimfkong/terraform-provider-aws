package storagegateway_test

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/storagegateway"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfstoragegateway "github.com/hashicorp/terraform-provider-aws/internal/service/storagegateway"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_storagegateway_gateway", &resource.Sweeper{
		Name: "aws_storagegateway_gateway",
		F:    sweepGateways,
	})
}

func sweepGateways(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).StorageGatewayConn

	err = conn.ListGatewaysPages(&storagegateway.ListGatewaysInput{}, func(page *storagegateway.ListGatewaysOutput, lastPage bool) bool {
		if len(page.Gateways) == 0 {
			log.Print("[DEBUG] No Storage Gateway Gateways to sweep")
			return true
		}

		for _, gateway := range page.Gateways {
			name := aws.StringValue(gateway.GatewayName)

			log.Printf("[INFO] Deleting Storage Gateway Gateway: %s", name)
			input := &storagegateway.DeleteGatewayInput{
				GatewayARN: gateway.GatewayARN,
			}

			_, err := conn.DeleteGateway(input)
			if err != nil {
				if tfawserr.ErrMessageContains(err, storagegateway.ErrorCodeGatewayNotFound, "") {
					continue
				}
				log.Printf("[ERROR] Failed to delete Storage Gateway Gateway (%s): %s", name, err)
			}
		}

		return !lastPage
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping Storage Gateway Gateway sweep for %s: %s", region, err)
			return nil
		}
		return fmt.Errorf("Error retrieving Storage Gateway Gateways: %w", err)
	}
	return nil
}

func TestAccAWSStorageGatewayGateway_GatewayType_Cached(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_GatewayType_Cached(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "storagegateway", regexp.MustCompile(`gateway/sgw-.+`)),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "gateway_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gateway_timezone", "GMT"),
					resource.TestCheckResourceAttr(resourceName, "gateway_type", "CACHED"),
					resource.TestCheckResourceAttr(resourceName, "medium_changer_type", ""),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "smb_guest_password", ""),
					resource.TestCheckResourceAttr(resourceName, "smb_security_strategy", ""),
					resource.TestCheckResourceAttr(resourceName, "tape_drive_type", ""),
					resource.TestCheckResourceAttrPair(resourceName, "ec2_instance_id", "aws_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "endpoint_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "host_environment", "EC2"),
					resource.TestCheckResourceAttr(resourceName, "gateway_network_interface.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "gateway_network_interface.0.ipv4_address", "aws_instance.test", "private_ip"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_GatewayType_FileFsxSmb(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_GatewayType_FileFSxSMB(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "storagegateway", regexp.MustCompile(`gateway/sgw-.+`)),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "gateway_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gateway_timezone", "GMT"),
					resource.TestCheckResourceAttr(resourceName, "gateway_type", "FILE_FSX_SMB"),
					resource.TestCheckResourceAttr(resourceName, "medium_changer_type", ""),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "smb_guest_password", ""),
					resource.TestCheckResourceAttr(resourceName, "tape_drive_type", ""),
					resource.TestCheckResourceAttrPair(resourceName, "ec2_instance_id", "aws_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "endpoint_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "host_environment", "EC2"),
					resource.TestCheckResourceAttr(resourceName, "gateway_network_interface.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "gateway_network_interface.0.ipv4_address", "aws_instance.test", "private_ip"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_GatewayType_FileS3(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_GatewayType_FileS3(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "storagegateway", regexp.MustCompile(`gateway/sgw-.+`)),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "gateway_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gateway_timezone", "GMT"),
					resource.TestCheckResourceAttr(resourceName, "gateway_type", "FILE_S3"),
					resource.TestCheckResourceAttr(resourceName, "medium_changer_type", ""),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "smb_guest_password", ""),
					resource.TestCheckResourceAttr(resourceName, "tape_drive_type", ""),
					resource.TestCheckResourceAttrPair(resourceName, "ec2_instance_id", "aws_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "endpoint_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "host_environment", "EC2"),
					resource.TestCheckResourceAttr(resourceName, "gateway_network_interface.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "gateway_network_interface.0.ipv4_address", "aws_instance.test", "private_ip"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_GatewayType_Stored(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_GatewayType_Stored(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "storagegateway", regexp.MustCompile(`gateway/sgw-.+`)),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "gateway_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gateway_timezone", "GMT"),
					resource.TestCheckResourceAttr(resourceName, "gateway_type", "STORED"),
					resource.TestCheckResourceAttr(resourceName, "medium_changer_type", ""),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "smb_guest_password", ""),
					resource.TestCheckResourceAttr(resourceName, "tape_drive_type", ""),
					resource.TestCheckResourceAttrPair(resourceName, "ec2_instance_id", "aws_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "endpoint_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "host_environment", "EC2"),
					resource.TestCheckResourceAttr(resourceName, "gateway_network_interface.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "gateway_network_interface.0.ipv4_address", "aws_instance.test", "private_ip"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_GatewayType_Vtl(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_GatewayType_Vtl(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "storagegateway", regexp.MustCompile(`gateway/sgw-.+`)),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "gateway_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gateway_timezone", "GMT"),
					resource.TestCheckResourceAttr(resourceName, "gateway_type", "VTL"),
					resource.TestCheckResourceAttr(resourceName, "medium_changer_type", ""),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "smb_guest_password", ""),
					resource.TestCheckResourceAttr(resourceName, "tape_drive_type", ""),
					resource.TestCheckResourceAttrPair(resourceName, "ec2_instance_id", "aws_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "endpoint_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "host_environment", "EC2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_tags(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayTags1Config(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "storagegateway", regexp.MustCompile(`gateway/sgw-.+`)),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
			{
				Config: testAccGatewayTags2Config(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccGatewayTags1Config(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_GatewayName(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName1 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_GatewayType_FileS3(rName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gateway_name", rName1),
				),
			},
			{
				Config: testAccGatewayConfig_GatewayType_FileS3(rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gateway_name", rName2),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_CloudWatchLogs(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName1 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"
	resourceName2 := "aws_cloudwatch_log_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_Log_Group(rName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttrPair(resourceName, "cloudwatch_log_group_arn", resourceName2, "arn"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_GatewayTimezone(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_GatewayTimezone(rName, "GMT-1:00"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gateway_timezone", "GMT-1:00"),
				),
			},
			{
				Config: testAccGatewayConfig_GatewayTimezone(rName, "GMT-2:00"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gateway_timezone", "GMT-2:00"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_GatewayVpcEndpoint(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"
	vpcEndpointResourceName := "aws_vpc_endpoint.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_GatewayVPCEndpoint(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttrPair(resourceName, "gateway_vpc_endpoint", vpcEndpointResourceName, "dns_entry.0.dns_name"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_SmbActiveDirectorySettings(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"
	domainName := acctest.RandomDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_SMBActiveDirectorySettings(rName, domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.0.domain_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.0.username", "Administrator"),
					resource.TestCheckResourceAttrSet(resourceName, "smb_active_directory_settings.0.active_directory_status"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address", "smb_active_directory_settings"},
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_SmbActiveDirectorySettings_timeout(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"
	domainName := acctest.RandomDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_SMBActiveDirectorySettingsTimeout(rName, domainName, 50),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.0.domain_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.0.timeout_in_seconds", "50"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address", "smb_active_directory_settings"},
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_SmbMicrosoftActiveDirectorySettings(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"
	domainName := acctest.RandomDomainName()
	username := "Admin"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_SMBMicrosoftActiveDirectorySettings(rName, domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.0.domain_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.0.username", username),
					resource.TestCheckResourceAttrSet(resourceName, "smb_active_directory_settings.0.active_directory_status"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address", "smb_active_directory_settings"},
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_SmbMicrosoftActiveDirectorySettings_timeout(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"
	domainName := acctest.RandomDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_SMBMicrosoftActiveDirectorySettingsTimeout(rName, domainName, 50),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.0.domain_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "smb_active_directory_settings.0.timeout_in_seconds", "50"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address", "smb_active_directory_settings"},
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_SmbGuestPassword(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_SMBGuestPassword(rName, "myguestpassword1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "smb_guest_password", "myguestpassword1"),
				),
			},
			{
				Config: testAccGatewayConfig_SMBGuestPassword(rName, "myguestpassword2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "smb_guest_password", "myguestpassword2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address", "smb_guest_password"},
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_SMBSecurityStrategy(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewaySMBSecurityStrategyConfig(rName, "ClientSpecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "smb_security_strategy", "ClientSpecified"),
					resource.TestCheckResourceAttr(resourceName, "smb_file_share_visibility", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
			{
				Config: testAccGatewaySMBSecurityStrategyConfig(rName, "MandatorySigning"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "smb_security_strategy", "MandatorySigning"),
				),
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_SMBVisibility(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewaySMBVisibilityConfig(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "smb_file_share_visibility", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
			{
				Config: testAccGatewaySMBVisibilityConfig(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "smb_file_share_visibility", "false"),
				),
			},
			{
				Config: testAccGatewaySMBVisibilityConfig(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "smb_file_share_visibility", "true"),
				),
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_disappears(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_GatewayType_Cached(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					acctest.CheckResourceDisappears(acctest.Provider, tfstoragegateway.ResourceGateway(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_bandwidthUpload(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayBandwidthUploadConfig(rName, 102400),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "average_upload_rate_limit_in_bits_per_sec", "102400"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
			{
				Config: testAccGatewayBandwidthUploadConfig(rName, 2*102400),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "average_upload_rate_limit_in_bits_per_sec", "204800"),
				),
			},
			{
				Config: testAccGatewayConfig_GatewayType_Cached(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "average_upload_rate_limit_in_bits_per_sec", "0"),
				),
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_bandwidthDownload(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayBandwidthDownloadConfig(rName, 102400),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "average_download_rate_limit_in_bits_per_sec", "102400"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
			{
				Config: testAccGatewayBandwidthDownloadConfig(rName, 2*102400),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "average_download_rate_limit_in_bits_per_sec", "204800"),
				),
			},
			{
				Config: testAccGatewayConfig_GatewayType_Cached(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "average_upload_rate_limit_in_bits_per_sec", "0"),
				),
			},
		},
	})
}

func TestAccAWSStorageGatewayGateway_bandwidthAll(t *testing.T) {
	var gateway storagegateway.DescribeGatewayInformationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_storagegateway_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, storagegateway.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayBandwidthAllConfig(rName, 102400),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "average_upload_rate_limit_in_bits_per_sec", "102400"),
					resource.TestCheckResourceAttr(resourceName, "average_download_rate_limit_in_bits_per_sec", "102400"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activation_key", "gateway_ip_address"},
			},
			{
				Config: testAccGatewayBandwidthAllConfig(rName, 2*102400),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "average_upload_rate_limit_in_bits_per_sec", "204800"),
					resource.TestCheckResourceAttr(resourceName, "average_download_rate_limit_in_bits_per_sec", "204800"),
				),
			},
			{
				Config: testAccGatewayConfig_GatewayType_Cached(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "average_upload_rate_limit_in_bits_per_sec", "0"),
					resource.TestCheckResourceAttr(resourceName, "average_download_rate_limit_in_bits_per_sec", "0"),
				),
			},
		},
	})
}

func testAccCheckGatewayDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).StorageGatewayConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_storagegateway_gateway" {
			continue
		}

		input := &storagegateway.DescribeGatewayInformationInput{
			GatewayARN: aws.String(rs.Primary.ID),
		}

		_, err := conn.DescribeGatewayInformation(input)

		if err != nil {
			if tfstoragegateway.IsErrGatewayNotFound(err) {
				return nil
			}
			return err
		}
	}

	return nil

}

func testAccCheckGatewayExists(resourceName string, gateway *storagegateway.DescribeGatewayInformationOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).StorageGatewayConn
		input := &storagegateway.DescribeGatewayInformationInput{
			GatewayARN: aws.String(rs.Primary.ID),
		}

		output, err := conn.DescribeGatewayInformation(input)

		if err != nil {
			return err
		}

		if output == nil {
			return fmt.Errorf("Gateway %q does not exist", rs.Primary.ID)
		}

		*gateway = *output

		return nil
	}
}

// testAcc_VPCBase provides a publicly accessible subnet
// and security group, suitable for Storage Gateway EC2 instances of any type
func testAcc_VPCBase(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(),
		fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.test.id
  route_table_id         = aws_vpc.test.main_route_table_id
}

resource "aws_security_group" "test" {
  vpc_id = aws_vpc.test.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName))
}

func testAcc_FileGatewayBase(rName string) string {
	return acctest.ConfigCompose(
		testAcc_VPCBase(rName),
		// Reference: https://docs.aws.amazon.com/storagegateway/latest/userguide/Requirements.html
		acctest.AvailableEC2InstanceTypeForAvailabilityZone("aws_subnet.test.availability_zone", "m5.xlarge", "m4.xlarge"),
		fmt.Sprintf(`
# Reference: https://docs.aws.amazon.com/storagegateway/latest/userguide/ec2-gateway-file.html
data "aws_ssm_parameter" "aws_service_storagegateway_ami_FILE_S3_latest" {
  name = "/aws/service/storagegateway/ami/FILE_S3/latest"
}

resource "aws_instance" "test" {
  depends_on = [aws_route.test]

  ami                         = data.aws_ssm_parameter.aws_service_storagegateway_ami_FILE_S3_latest.value
  associate_public_ip_address = true
  instance_type               = data.aws_ec2_instance_type_offering.available.instance_type
  vpc_security_group_ids      = [aws_security_group.test.id]
  subnet_id                   = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName))
}

func testAcc_TapeAndVolumeGatewayBase(rName string) string {
	return acctest.ConfigCompose(
		testAcc_VPCBase(rName),
		// Reference: https://docs.aws.amazon.com/storagegateway/latest/userguide/Requirements.html
		acctest.AvailableEC2InstanceTypeForAvailabilityZone("aws_subnet.test.availability_zone", "m5.xlarge", "m4.xlarge"),
		fmt.Sprintf(`
# Reference: https://docs.aws.amazon.com/storagegateway/latest/userguide/ec2-gateway-common.html
# NOTE: CACHED, STORED, and VTL Gateway Types share the same AMI
data "aws_ssm_parameter" "aws_service_storagegateway_ami_CACHED_latest" {
  name = "/aws/service/storagegateway/ami/CACHED/latest"
}

resource "aws_instance" "test" {
  depends_on = [aws_route.test]

  ami                         = data.aws_ssm_parameter.aws_service_storagegateway_ami_CACHED_latest.value
  associate_public_ip_address = true
  instance_type               = data.aws_ec2_instance_type_offering.available.instance_type
  vpc_security_group_ids      = [aws_security_group.test.id]
  subnet_id                   = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }
}
`, rName))
}

func testAccGatewayConfig_GatewayType_Cached(rName string) string {
	return testAcc_TapeAndVolumeGatewayBase(rName) + fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address = aws_instance.test.public_ip
  gateway_name       = %q
  gateway_timezone   = "GMT"
  gateway_type       = "CACHED"
}
`, rName)
}

func testAccGatewayConfig_GatewayType_FileFSxSMB(rName string) string {
	return testAcc_FileGatewayBase(rName) + fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address = aws_instance.test.public_ip
  gateway_name       = %q
  gateway_timezone   = "GMT"
  gateway_type       = "FILE_FSX_SMB"
}
`, rName)
}

func testAccGatewayConfig_GatewayType_FileS3(rName string) string {
	return testAcc_FileGatewayBase(rName) + fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address = aws_instance.test.public_ip
  gateway_name       = %q
  gateway_timezone   = "GMT"
  gateway_type       = "FILE_S3"
}
`, rName)
}

func testAccGatewayConfig_Log_Group(rName string) string {
	return testAcc_FileGatewayBase(rName) + fmt.Sprintf(`
resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address       = aws_instance.test.public_ip
  gateway_name             = %[1]q
  gateway_timezone         = "GMT"
  gateway_type             = "FILE_S3"
  cloudwatch_log_group_arn = aws_cloudwatch_log_group.test.arn
}
`, rName)
}

func testAccGatewayConfig_GatewayType_Stored(rName string) string {
	return testAcc_TapeAndVolumeGatewayBase(rName) + fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address = aws_instance.test.public_ip
  gateway_name       = %q
  gateway_timezone   = "GMT"
  gateway_type       = "STORED"
}
`, rName)
}

func testAccGatewayConfig_GatewayType_Vtl(rName string) string {
	return testAcc_TapeAndVolumeGatewayBase(rName) + fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address = aws_instance.test.public_ip
  gateway_name       = %q
  gateway_timezone   = "GMT"
  gateway_type       = "VTL"
}
`, rName)
}

func testAccGatewayConfig_GatewayTimezone(rName, gatewayTimezone string) string {
	return testAcc_FileGatewayBase(rName) + fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address = aws_instance.test.public_ip
  gateway_name       = %q
  gateway_timezone   = %q
  gateway_type       = "FILE_S3"
}
`, rName, gatewayTimezone)
}

func testAccGatewayConfig_GatewayVPCEndpoint(rName string) string {
	return testAcc_TapeAndVolumeGatewayBase(rName) + fmt.Sprintf(`
data "aws_vpc_endpoint_service" "storagegateway" {
  service = "storagegateway"
}

resource "aws_vpc_endpoint" "test" {
  security_group_ids = [aws_security_group.test.id]
  service_name       = data.aws_vpc_endpoint_service.storagegateway.service_name
  subnet_ids         = [aws_subnet.test.id]
  vpc_endpoint_type  = data.aws_vpc_endpoint_service.storagegateway.service_type
  vpc_id             = aws_vpc.test.id
}

resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address   = aws_instance.test.public_ip
  gateway_name         = %[1]q
  gateway_timezone     = "GMT"
  gateway_type         = "CACHED"
  gateway_vpc_endpoint = aws_vpc_endpoint.test.dns_entry[0].dns_name
}
`, rName)
}

func testAccGatewayConfig_DirectoryServiceSimpleDirectory(rName, domainName string) string {
	return fmt.Sprintf(`
resource "aws_directory_service_directory" "test" {
  name     = %[2]q
  password = "SuperSecretPassw0rd"
  size     = "Small"

  vpc_settings {
    subnet_ids = aws_subnet.test[*].id
    vpc_id     = aws_vpc.test.id
  }

  tags = {
    Name = %[1]q
  }
}

`, rName, domainName)
}

func testAccGatewayConfig_DirectoryServiceMicrosoftAD(rName, domainName string) string {
	return fmt.Sprintf(`
resource "aws_directory_service_directory" "test" {
  edition  = "Standard"
  name     = %[2]q
  password = "SuperSecretPassw0rd"
  type     = "MicrosoftAD"

  vpc_settings {
    subnet_ids = aws_subnet.test[*].id
    vpc_id     = aws_vpc.test.id
  }

  tags = {
    Name = %[1]q
  }
}

`, rName, domainName)
}

func testAccGatewaySMBActiveDirectorySettingsBaseConfig(rName string) string {
	return acctest.ConfigCompose(
		// Reference: https://docs.aws.amazon.com/storagegateway/latest/userguide/Requirements.html
		acctest.AvailableEC2InstanceTypeForAvailabilityZone("aws_subnet.test[0].availability_zone", "m5.xlarge", "m4.xlarge"),
		acctest.ConfigAvailableAZsNoOptIn(),
		fmt.Sprintf(`
# Directory Service Directories must be deployed across multiple EC2 Availability Zones
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  count = 2

  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = "10.0.${count.index}.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.test.id
  route_table_id         = aws_vpc.test.main_route_table_id
}

resource "aws_security_group" "test" {
  vpc_id = aws_vpc.test.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc_dhcp_options" "test" {
  domain_name         = aws_directory_service_directory.test.name
  domain_name_servers = aws_directory_service_directory.test.dns_ip_addresses

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc_dhcp_options_association" "test" {
  dhcp_options_id = aws_vpc_dhcp_options.test.id
  vpc_id          = aws_vpc.test.id
}

# Reference: https://docs.aws.amazon.com/storagegateway/latest/userguide/ec2-gateway-file.html
data "aws_ssm_parameter" "aws_service_storagegateway_ami_FILE_S3_latest" {
  name = "/aws/service/storagegateway/ami/FILE_S3/latest"
}

resource "aws_instance" "test" {
  depends_on = [aws_route.test, aws_vpc_dhcp_options_association.test]

  ami                         = data.aws_ssm_parameter.aws_service_storagegateway_ami_FILE_S3_latest.value
  associate_public_ip_address = true
  instance_type               = data.aws_ec2_instance_type_offering.available.instance_type
  vpc_security_group_ids      = [aws_security_group.test.id]
  subnet_id                   = aws_subnet.test[0].id

  tags = {
    Name = %[1]q
  }
}
`, rName))
}

func testAccGatewayConfig_SMBActiveDirectorySettings(rName, domainName string) string {
	return acctest.ConfigCompose(
		testAccGatewaySMBActiveDirectorySettingsBaseConfig(rName),
		testAccGatewayConfig_DirectoryServiceSimpleDirectory(rName, domainName),
		fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address = aws_instance.test.public_ip
  gateway_name       = %[1]q
  gateway_timezone   = "GMT"
  gateway_type       = "FILE_S3"

  smb_active_directory_settings {
    domain_name = aws_directory_service_directory.test.name
    password    = aws_directory_service_directory.test.password
    username    = "Administrator"
  }
}
`, rName))
}

func testAccGatewayConfig_SMBActiveDirectorySettingsTimeout(rName, domainName string, timeout int) string {
	return acctest.ConfigCompose(
		testAccGatewaySMBActiveDirectorySettingsBaseConfig(rName),
		testAccGatewayConfig_DirectoryServiceSimpleDirectory(rName, domainName),
		fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address = aws_instance.test.public_ip
  gateway_name       = %[1]q
  gateway_timezone   = "GMT"
  gateway_type       = "FILE_S3"

  smb_active_directory_settings {
    domain_name        = aws_directory_service_directory.test.name
    password           = aws_directory_service_directory.test.password
    username           = "Administrator"
    timeout_in_seconds = %[2]d
  }
}
`, rName, timeout))
}

func testAccGatewayConfig_SMBMicrosoftActiveDirectorySettings(rName, domainName string) string {
	return acctest.ConfigCompose(
		testAccGatewaySMBActiveDirectorySettingsBaseConfig(rName),
		testAccGatewayConfig_DirectoryServiceMicrosoftAD(rName, domainName),
		fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address = aws_instance.test.public_ip
  gateway_name       = %[1]q
  gateway_timezone   = "GMT"
  gateway_type       = "FILE_S3"

  smb_active_directory_settings {
    domain_name = aws_directory_service_directory.test.name
    password    = aws_directory_service_directory.test.password
    username    = "Admin"
  }
}
`, rName))
}

func testAccGatewayConfig_SMBMicrosoftActiveDirectorySettingsTimeout(rName, domainName string, timeout int) string {
	return acctest.ConfigCompose(
		testAccGatewaySMBActiveDirectorySettingsBaseConfig(rName),
		testAccGatewayConfig_DirectoryServiceMicrosoftAD(rName, domainName),
		fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address = aws_instance.test.public_ip
  gateway_name       = %[1]q
  gateway_timezone   = "GMT"
  gateway_type       = "FILE_S3"

  smb_active_directory_settings {
    domain_name        = aws_directory_service_directory.test.name
    password           = aws_directory_service_directory.test.password
    username           = "Admin"
    timeout_in_seconds = %[2]d
  }
}
`, rName, timeout))
}

func testAccGatewayConfig_SMBGuestPassword(rName, smbGuestPassword string) string {
	return testAcc_FileGatewayBase(rName) + fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address = aws_instance.test.public_ip
  gateway_name       = %q
  gateway_timezone   = "GMT"
  gateway_type       = "FILE_S3"
  smb_guest_password = %q
}
`, rName, smbGuestPassword)
}

func testAccGatewaySMBSecurityStrategyConfig(rName, strategy string) string {
	return testAcc_FileGatewayBase(rName) + fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address    = aws_instance.test.public_ip
  gateway_name          = %q
  gateway_timezone      = "GMT"
  gateway_type          = "FILE_S3"
  smb_security_strategy = %q
}
`, rName, strategy)
}

func testAccGatewaySMBVisibilityConfig(rName string, visible bool) string {
	return testAcc_FileGatewayBase(rName) + fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address        = aws_instance.test.public_ip
  gateway_name              = %[1]q
  gateway_timezone          = "GMT"
  gateway_type              = "FILE_S3"
  smb_file_share_visibility = %[2]t
}
`, rName, visible)
}

func testAccGatewayTags1Config(rName, tagKey1, tagValue1 string) string {
	return testAcc_TapeAndVolumeGatewayBase(rName) + fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address = aws_instance.test.public_ip
  gateway_name       = %q
  gateway_timezone   = "GMT"
  gateway_type       = "CACHED"

  tags = {
    %q = %q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccGatewayTags2Config(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return testAcc_TapeAndVolumeGatewayBase(rName) + fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address = aws_instance.test.public_ip
  gateway_name       = %q
  gateway_timezone   = "GMT"
  gateway_type       = "CACHED"

  tags = {
    %q = %q
    %q = %q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccGatewayBandwidthUploadConfig(rName string, rate int) string {
	return testAcc_TapeAndVolumeGatewayBase(rName) + fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address                        = aws_instance.test.public_ip
  gateway_name                              = %[1]q
  gateway_timezone                          = "GMT"
  gateway_type                              = "CACHED"
  average_upload_rate_limit_in_bits_per_sec = %[2]d
}
`, rName, rate)
}

func testAccGatewayBandwidthDownloadConfig(rName string, rate int) string {
	return testAcc_TapeAndVolumeGatewayBase(rName) + fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address                          = aws_instance.test.public_ip
  gateway_name                                = %[1]q
  gateway_timezone                            = "GMT"
  gateway_type                                = "CACHED"
  average_download_rate_limit_in_bits_per_sec = %[2]d
}
`, rName, rate)
}

func testAccGatewayBandwidthAllConfig(rName string, rate int) string {
	return testAcc_TapeAndVolumeGatewayBase(rName) + fmt.Sprintf(`
resource "aws_storagegateway_gateway" "test" {
  gateway_ip_address                          = aws_instance.test.public_ip
  gateway_name                                = %[1]q
  gateway_timezone                            = "GMT"
  gateway_type                                = "CACHED"
  average_upload_rate_limit_in_bits_per_sec   = %[2]d
  average_download_rate_limit_in_bits_per_sec = %[2]d
}
`, rName, rate)
}
