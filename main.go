package main

import (
	"context"
	"fmt"
	"encoding/json"
	//"fmt"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-11-01/network"
	//"github.com/Azure/azure-sdk-for-go/sdk/to"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	//"fmt"
)
//func AzureCredentais() {

//}
func List_All_ResourceGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
	//clientId := os.Getenv("AZURE_CLIENT_ID")
	//clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	//tenantId :=  os.Getenv("AZURE_TENANT_ID")
		RGlist := make([]string, 0)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Printf("Error while creating an new authentication, %v ", err)

	}
	rgClient := resources.NewGroupsClient(subscriptionId)
	rgClient.Authorizer = authorizer
	for list, err2 := rgClient.ListComplete(context.Background(), "", nil); list.NotDone(); err = list.Next() {
		if err != nil {
			log.Printf("error traverising RG list, %v ", err2)

		}

		rgName := *list.Value().Name
		RGlist = append(RGlist, rgName)

	}
	json.NewEncoder(w).Encode(RGlist)
}
//}
//
	func createGroup() (group resources.Group, err error) {
		subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
		authorizer, err := auth.NewAuthorizerFromEnvironment()
		if err != nil {
			log.Printf("Error while creating an new authentication, %v ", err)

		}
		groupsClient := resources.NewGroupsClient(subscriptionId)
		groupsClient.Authorizer = authorizer

		return groupsClient.CreateOrUpdate(
			context.Background(),
			"click2cloud-HarshChatte",
			resources.Group{
				Location: to.StringPtr("South India")})
	}
func create_Resource_Group( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	group, err := createGroup()
	if err != nil {
		log.Fatalf("failed to create group: %v", err)
	}
	json.NewEncoder(w).Encode(*group.Name)
	json.NewEncoder(w).Encode("Resource Group Created Successfully")

}

func CreateNetworkSecurityGroup(ctx context.Context, nsgName string)(nsg network.SecurityGroup, err error) {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Printf("Error while creating an new authentication, %v ", err)

	}

	nsgClient := network.NewSecurityGroupsClient(subscriptionId)

	//nsgClient := getNsgClient()
	nsgClient.Authorizer = authorizer
	future, err := nsgClient.CreateOrUpdate(

		context.Background(),
		"click2cloud-HarshChatte",
		"MyNSG",
		network.SecurityGroup{
			Location: to.StringPtr("South India"),
			SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
				SecurityRules: &[]network.SecurityRule{
					{
						Name: to.StringPtr("allow_ssh"),
						SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
							Protocol:                 network.SecurityRuleProtocolTCP,
							SourceAddressPrefix:      to.StringPtr("0.0.0.0/0"),
							SourcePortRange:          to.StringPtr("1-65535"),
							DestinationAddressPrefix: to.StringPtr("0.0.0.0/0"),
							DestinationPortRange:     to.StringPtr("22"),
							Access:                   network.SecurityRuleAccessAllow,
							Direction:                network.SecurityRuleDirectionInbound,
							Priority:                 to.Int32Ptr(100),
						},
					},
					{
						Name: to.StringPtr("allow_https"),
						SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
							Protocol:                 network.SecurityRuleProtocolTCP,
							SourceAddressPrefix:      to.StringPtr("0.0.0.0/0"),
							SourcePortRange:          to.StringPtr("1-65535"),
							DestinationAddressPrefix: to.StringPtr("0.0.0.0/0"),
							DestinationPortRange:     to.StringPtr("443"),
							Access:                   network.SecurityRuleAccessAllow,
							Direction:                network.SecurityRuleDirectionInbound,
							Priority:                 to.Int32Ptr(200),
						},
					},
				},
			},
		},
	)

	if err != nil {
		return nsg,fmt.Errorf("cannot create nsg: %v", err)
	}

	err = future.WaitForCompletionRef(context.Background(), nsgClient.Client)
	if err != nil {
		return nsg,fmt.Errorf("cannot get nsg create or update future response: %v", err)
	}

	return future.Result(nsgClient)
}
func Create_Network_Security_Group( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	 nsg,err:= CreateNetworkSecurityGroup(context.Background(),"MyNSG")
	if err != nil {
		log.Fatalf("failed to create NetworkSecurityGroup: %v", err)
	}
	json.NewEncoder(w).Encode(*nsg.Name)
	json.NewEncoder(w).Encode("NetworkSecurityGroup Created Successfully")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/resources", List_All_ResourceGroups).Methods("GET")
	router.HandleFunc("/resources/New", create_Resource_Group).Methods("POST")
	router.HandleFunc("/resources/New1",Create_Network_Security_Group).Methods("POST")
	http.ListenAndServe(":8080", router)
}
//func Network_Security_Group( w http.ResponseWriter, r *http.Request) {
//	nsgClient := network.NewSecurityGroupsClient(subscriptionId)
//	nsgClient.Authorizer = authorizer
//