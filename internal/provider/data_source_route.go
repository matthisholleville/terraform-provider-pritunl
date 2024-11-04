package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/matthisholleville/terraform-pritunl-provider/internal/pritunl"
)

func dataSourceRoute() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information about a specific route in a Pritunl server.",
		ReadContext: dataSourceRouteRead,
		Schema: map[string]*schema.Schema{
			"server_id": {
				Description: "ID of the server in which the route is configured.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"network": {
				Description: "Network of the route.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"comment": {
				Description: "Comment associated with the route.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"nat": {
				Description: "Indicates if NAT is enabled for this route.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	serverID := d.Get("server_id").(string)
	network := d.Get("network").(string)

	route, err := SearchRouteByNetwork(meta, serverID, network)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(route.GetID())
	d.Set("network", route.Network)
	d.Set("server_id", serverID)
	d.Set("comment", route.Comment)
	d.Set("nat", route.Nat)

	return nil
}

func SearchRouteByNetwork(meta interface{}, serverID, network string) (pritunl.Route, error) {
	apiClient := meta.(pritunl.Client)

	routes, err := apiClient.GetRoutesByServer(serverID)
	if err != nil {
		return pritunl.Route{}, err
	}

	for _, route := range routes {
		if route.Network == network {
			return route, nil
		}
	}

	return pritunl.Route{}, errors.New("could not find a route with specified server_id and network")
}
