package provider

import (
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Runtime blocks are live ubus-derived state, exposed only on data sources (not
// resources): they are read-only observed state, never desired config.

type ifaceAddrModel struct {
	Address types.String `tfsdk:"address"`
	Mask    types.Int64  `tfsdk:"mask"`
}

type ifaceRouteModel struct {
	Target  types.String `tfsdk:"target"`
	Mask    types.Int64  `tfsdk:"mask"`
	Nexthop types.String `tfsdk:"nexthop"`
	Source  types.String `tfsdk:"source"`
}

type networkInterfaceRuntimeModel struct {
	Up          types.Bool        `tfsdk:"up"`
	Pending     types.Bool        `tfsdk:"pending"`
	Available   types.Bool        `tfsdk:"available"`
	L3Device    types.String      `tfsdk:"l3_device"`
	Uptime      types.Int64       `tfsdk:"uptime"`
	IPv4Address []ifaceAddrModel  `tfsdk:"ipv4_address"`
	IPv6Address []ifaceAddrModel  `tfsdk:"ipv6_address"`
	IPv6Prefix  []ifaceAddrModel  `tfsdk:"ipv6_prefix"`
	Route       []ifaceRouteModel `tfsdk:"route"`
}

type wirelessInterfaceRuntimeModel struct {
	Ifname         types.String `tfsdk:"ifname"`
	BSSID          types.String `tfsdk:"bssid"`
	Channel        types.Int64  `tfsdk:"channel"`
	Frequency      types.Int64  `tfsdk:"frequency"`
	Signal         types.Int64  `tfsdk:"signal"`
	Noise          types.Int64  `tfsdk:"noise"`
	TxpowerActual  types.Int64  `tfsdk:"txpower_actual"`
	AssoclistCount types.Int64  `tfsdk:"assoclist_count"`
}

func addrList(raw any) []ifaceAddrModel {
	arr, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := make([]ifaceAddrModel, 0, len(arr))
	for _, e := range arr {
		m, ok := e.(map[string]any)
		if !ok {
			continue
		}
		out = append(out, ifaceAddrModel{Address: strVal(m, "address"), Mask: int64Val(m, "mask")})
	}
	return out
}

func parseNetworkInterfaceRuntime(obj map[string]any) *networkInterfaceRuntimeModel {
	rt, ok := obj["runtime"].(map[string]any)
	if !ok {
		rt = map[string]any{}
	}
	rm := &networkInterfaceRuntimeModel{
		Up:          boolVal(rt, "up"),
		Pending:     boolVal(rt, "pending"),
		Available:   boolVal(rt, "available"),
		L3Device:    strVal(rt, "l3_device"),
		Uptime:      int64Val(rt, "uptime"),
		IPv4Address: addrList(rt["ipv4-address"]),
		IPv6Address: addrList(rt["ipv6-address"]),
		IPv6Prefix:  addrList(rt["ipv6-prefix"]),
	}
	if arr, ok := rt["route"].([]any); ok {
		for _, e := range arr {
			m, ok := e.(map[string]any)
			if !ok {
				continue
			}
			rm.Route = append(rm.Route, ifaceRouteModel{
				Target:  strVal(m, "target"),
				Mask:    int64Val(m, "mask"),
				Nexthop: strVal(m, "nexthop"),
				Source:  strVal(m, "source"),
			})
		}
	}
	return rm
}

func parseWirelessInterfaceRuntime(obj map[string]any) *wirelessInterfaceRuntimeModel {
	rt, ok := obj["runtime"].(map[string]any)
	if !ok {
		rt = map[string]any{}
	}
	return &wirelessInterfaceRuntimeModel{
		Ifname:         strVal(rt, "ifname"),
		BSSID:          strVal(rt, "bssid"),
		Channel:        int64Val(rt, "channel"),
		Frequency:      int64Val(rt, "frequency"),
		Signal:         int64Val(rt, "signal"),
		Noise:          int64Val(rt, "noise"),
		TxpowerActual:  int64Val(rt, "txpower_actual"),
		AssoclistCount: int64Val(rt, "assoclist_count"),
	}
}

func networkInterfaceRuntimeAttribute() dsschema.SingleNestedAttribute {
	addr := dsschema.NestedAttributeObject{Attributes: map[string]dsschema.Attribute{
		"address": dsschema.StringAttribute{Computed: true, Description: "Address."},
		"mask":    dsschema.Int64Attribute{Computed: true, Description: "Prefix length."},
	}}
	return dsschema.SingleNestedAttribute{
		Computed:    true,
		Description: "Live ubus-derived runtime state (read-only; reflects actual operation, not config).",
		Attributes: map[string]dsschema.Attribute{
			"up":           dsschema.BoolAttribute{Computed: true, Description: "Whether the interface is up."},
			"pending":      dsschema.BoolAttribute{Computed: true, Description: "Whether the interface is mid-setup."},
			"available":    dsschema.BoolAttribute{Computed: true, Description: "Whether the interface is available."},
			"l3_device":    dsschema.StringAttribute{Computed: true, Description: "Actual L3 kernel device."},
			"uptime":       dsschema.Int64Attribute{Computed: true, Description: "Seconds since the interface came up."},
			"ipv4_address": dsschema.ListNestedAttribute{Computed: true, Description: "Assigned IPv4 addresses.", NestedObject: addr},
			"ipv6_address": dsschema.ListNestedAttribute{Computed: true, Description: "Assigned IPv6 addresses.", NestedObject: addr},
			"ipv6_prefix":  dsschema.ListNestedAttribute{Computed: true, Description: "Delegated IPv6 prefixes.", NestedObject: addr},
			"route": dsschema.ListNestedAttribute{Computed: true, Description: "Active routes.", NestedObject: dsschema.NestedAttributeObject{Attributes: map[string]dsschema.Attribute{
				"target":  dsschema.StringAttribute{Computed: true},
				"mask":    dsschema.Int64Attribute{Computed: true},
				"nexthop": dsschema.StringAttribute{Computed: true},
				"source":  dsschema.StringAttribute{Computed: true},
			}}},
		},
	}
}

// networkInterfaceDSModel mirrors networkInterfaceModel plus the data-source-only
// runtime block. Keep field list in sync with networkInterfaceModel.
type networkInterfaceDSModel struct {
	ID            types.String                  `tfsdk:"id"`
	Managed       types.Bool                    `tfsdk:"managed"`
	ETag          types.String                  `tfsdk:"etag"`
	Device        types.String                  `tfsdk:"device"`
	Proto         types.String                  `tfsdk:"proto"`
	IPAddr        types.String                  `tfsdk:"ipaddr"`
	IPAddrs       types.List                    `tfsdk:"ipaddrs"`
	Netmask       types.String                  `tfsdk:"netmask"`
	Gateway       types.String                  `tfsdk:"gateway"`
	DNS           types.List                    `tfsdk:"dns"`
	IP6Assign     types.String                  `tfsdk:"ip6assign"`
	MTU           types.String                  `tfsdk:"mtu"`
	Auto          types.Bool                    `tfsdk:"auto"`
	PrivateKey    types.String                  `tfsdk:"private_key"`
	HasPrivateKey types.Bool                    `tfsdk:"has_private_key"`
	ListenPort    types.String                  `tfsdk:"listen_port"`
	Addresses     types.List                    `tfsdk:"addresses"`
	Nohostroute   types.Bool                    `tfsdk:"nohostroute"`
	IP4Table      types.String                  `tfsdk:"ip4table"`
	IP6Table      types.String                  `tfsdk:"ip6table"`
	PeerDNS       types.Bool                    `tfsdk:"peerdns"`
	DefaultRoute  types.Bool                    `tfsdk:"defaultroute"`
	Metric        types.String                  `tfsdk:"metric"`
	Hostname      types.String                  `tfsdk:"hostname"`
	ClientID      types.String                  `tfsdk:"clientid"`
	ReqPrefix     types.String                  `tfsdk:"reqprefix"`
	ReqAddress    types.String                  `tfsdk:"reqaddress"`
	IP6Hint       types.String                  `tfsdk:"ip6hint"`
	IP6IfaceID    types.String                  `tfsdk:"ip6ifaceid"`
	Delegate      types.Bool                    `tfsdk:"delegate"`
	Runtime       *networkInterfaceRuntimeModel `tfsdk:"runtime"`
}

func networkInterfaceDS(b networkInterfaceModel, obj map[string]any) networkInterfaceDSModel {
	return networkInterfaceDSModel{
		ID: b.ID, Managed: b.Managed, ETag: b.ETag, Device: b.Device, Proto: b.Proto,
		IPAddr: b.IPAddr, IPAddrs: b.IPAddrs, Netmask: b.Netmask, Gateway: b.Gateway, DNS: b.DNS,
		IP6Assign: b.IP6Assign, MTU: b.MTU, Auto: b.Auto, PrivateKey: b.PrivateKey, HasPrivateKey: b.HasPrivateKey,
		ListenPort: b.ListenPort, Addresses: b.Addresses, Nohostroute: b.Nohostroute, IP4Table: b.IP4Table, IP6Table: b.IP6Table,
		PeerDNS: b.PeerDNS, DefaultRoute: b.DefaultRoute, Metric: b.Metric, Hostname: b.Hostname, ClientID: b.ClientID,
		ReqPrefix: b.ReqPrefix, ReqAddress: b.ReqAddress, IP6Hint: b.IP6Hint, IP6IfaceID: b.IP6IfaceID, Delegate: b.Delegate,
		Runtime: parseNetworkInterfaceRuntime(obj),
	}
}

// wirelessInterfaceDSModel mirrors wirelessInterfaceModel plus runtime.
type wirelessInterfaceDSModel struct {
	ID         types.String                   `tfsdk:"id"`
	Managed    types.Bool                     `tfsdk:"managed"`
	ETag       types.String                   `tfsdk:"etag"`
	Device     types.String                   `tfsdk:"device"`
	Network    types.String                   `tfsdk:"network"`
	Mode       types.String                   `tfsdk:"mode"`
	SSID       types.String                   `tfsdk:"ssid"`
	Encryption types.String                   `tfsdk:"encryption"`
	Disabled   types.Bool                     `tfsdk:"disabled"`
	Hidden     types.Bool                     `tfsdk:"hidden"`
	Isolate    types.Bool                     `tfsdk:"isolate"`
	Key        types.String                   `tfsdk:"key"`
	HasKey     types.Bool                     `tfsdk:"has_key"`
	Runtime    *wirelessInterfaceRuntimeModel `tfsdk:"runtime"`
}

func wirelessInterfaceDS(b wirelessInterfaceModel, obj map[string]any) wirelessInterfaceDSModel {
	return wirelessInterfaceDSModel{
		ID: b.ID, Managed: b.Managed, ETag: b.ETag, Device: b.Device, Network: b.Network,
		Mode: b.Mode, SSID: b.SSID, Encryption: b.Encryption, Disabled: b.Disabled, Hidden: b.Hidden,
		Isolate: b.Isolate, Key: b.Key, HasKey: b.HasKey, Runtime: parseWirelessInterfaceRuntime(obj),
	}
}

func wirelessInterfaceRuntimeAttribute() dsschema.SingleNestedAttribute {
	return dsschema.SingleNestedAttribute{
		Computed:    true,
		Description: "Live iwinfo-derived runtime state (read-only).",
		Attributes: map[string]dsschema.Attribute{
			"ifname":          dsschema.StringAttribute{Computed: true, Description: "Kernel wireless device name."},
			"bssid":           dsschema.StringAttribute{Computed: true, Description: "BSSID."},
			"channel":         dsschema.Int64Attribute{Computed: true, Description: "Operating channel."},
			"frequency":       dsschema.Int64Attribute{Computed: true, Description: "Operating frequency (MHz)."},
			"signal":          dsschema.Int64Attribute{Computed: true, Description: "Signal level (dBm)."},
			"noise":           dsschema.Int64Attribute{Computed: true, Description: "Noise floor (dBm)."},
			"txpower_actual":  dsschema.Int64Attribute{Computed: true, Description: "Actual transmit power (dBm)."},
			"assoclist_count": dsschema.Int64Attribute{Computed: true, Description: "Number of associated clients."},
		},
	}
}
