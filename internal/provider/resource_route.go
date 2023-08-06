package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/matthisholleville/terraform-pritunl-provider/internal/pritunl"
)

var serverLocks = map[string]*sync.Mutex{}
var serverLocksMutex = &sync.Mutex{}

func resourceRoute() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the server",
			},
			"network": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Network address with subnet to route",
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					return validation.IsCIDR(i, s)
				},
			},
			"comment": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Comment for route",
			},
			"nat": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Description: "NAT vpn traffic destined to this network",
				Computed:    true,
			},
		},
		CreateContext: resourceCreateRoute,
		ReadContext:   resourceReadRoute,
		UpdateContext: resourceUpdateRoute,
		DeleteContext: resourceDeleteRoute,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePrepareImportRoute,
		},
	}
}

func resourceCreateRoute(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(pritunl.Client)

	serverId := d.Get("server_id").(string)

	serverLock := getOrCreateServerLock(serverId)

	serverLock.Lock()
	defer serverLock.Unlock()

	routes, err := apiClient.GetRoutesByServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, route := range routes {
		if route.Network == d.Get("network") {
			return diag.FromErr(fmt.Errorf("Route already exist with same network %s on %s server. Route ID %s", d.Get("network"), serverId, route.GetID()))
		}
	}

	routeData := map[string]interface{}{
		"network": d.Get("network"),
		"comment": d.Get("comment"),
		"nat":     d.Get("nat"),
	}

	err = apiClient.StopServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}

	route, err := apiClient.AddRouteToServer(serverId, pritunl.ConvertMapToRoute(routeData))
	if err != nil {
		return diag.FromErr(err)
	}

	err = apiClient.StartServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(route.GetID())

	return resourceReadRoute(ctx, d, meta)
}

func resourcePrepareImportRoute(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	re := regexp.MustCompile(`server/([a-z0-9]+)/route/([a-z0-9]+)`)
	matches := re.FindAllStringSubmatch(d.Id(), -1)

	if matches == nil {
		return nil, errors.New("Incorrect format. Must be server/(?P<server_id>[^/]+)/route/(?P<route_id>[^/]+)")
	}

	var serverId string
	var routeId string

	for _, match := range matches {
		serverId = match[1]
		routeId = match[2]
	}

	d.Set("server_id", serverId)
	d.SetId(routeId)

	return []*schema.ResourceData{d}, nil
}

func resourceReadRoute(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(pritunl.Client)

	serverId := d.Get("server_id").(string)
	routes, err := apiClient.GetRoutesByServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, route := range routes {
		if route.GetID() == d.Id() {
			d.Set("server_id", serverId)
			d.Set("network", route.Network)
			d.Set("comment", route.Comment)
			d.Set("nat", route.Nat)
		}

		return nil
	}
	return diag.FromErr(errors.New("Unable to find route."))
}

func resourceUpdateRoute(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(pritunl.Client)

	serverId := d.Get("server_id").(string)

	serverLock := getOrCreateServerLock(serverId)

	serverLock.Lock()
	defer serverLock.Unlock()

	routes, err := apiClient.GetRoutesByServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, route := range routes {
		if route.GetID() == d.Id() {
			if v, ok := d.GetOk("nat"); ok {
				route.Nat = v.(bool)
			}

			if v, ok := d.GetOk("comment"); ok {
				route.Comment = v.(string)
			}

			err = apiClient.StopServer(serverId)
			if err != nil {
				return diag.FromErr(err)
			}

			route, err := apiClient.UpdateRouteOnServer(serverId, route)
			if err != nil {
				return diag.FromErr(err)
			}

			err = apiClient.StartServer(serverId)
			if err != nil {
				return diag.FromErr(err)
			}

			d.SetId(route.GetID())

			return resourceReadRoute(ctx, d, meta)
		}
	}
	return diag.FromErr(errors.New("Unable to find route."))
}

func resourceDeleteRoute(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(pritunl.Client)

	serverId := d.Get("server_id").(string)

	serverLock := getOrCreateServerLock(serverId)

	serverLock.Lock()
	defer serverLock.Unlock()

	err := apiClient.StopServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = apiClient.DeleteRouteFromServer(serverId, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = apiClient.StartServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil

}

func getOrCreateServerLock(serverID string) *sync.Mutex {
	serverLocksMutex.Lock()
	defer serverLocksMutex.Unlock()

	lock, exists := serverLocks[serverID]
	if !exists {
		lock = &sync.Mutex{}
		serverLocks[serverID] = lock
	}

	return lock
}
