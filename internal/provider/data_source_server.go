package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/matthisholleville/terraform-pritunl-provider/internal/pritunl"
)

func dataSourceServer() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information about the Pritunl server.",
		ReadContext: dataSourceServerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_name": {
				Description: "Server name",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceServerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	server_name := d.Get("server_name")
	apiClient := meta.(pritunl.Client)

	servers, err := apiClient.GetServers()
	if err != nil {
		return diag.FromErr(err)
	}

	server, err := SearchServer(servers, server_name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(server.ID)
	d.Set("server_name", server.Name)

	return nil
}

func SearchServer(servers []pritunl.Server, serverName interface{}) (pritunl.Server, error) {
	for _, server := range servers {
		if server.Name == serverName {
			return server, nil
		}
	}
	return pritunl.Server{}, errors.New("Server not found")
}
