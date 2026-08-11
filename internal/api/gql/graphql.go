// Package gql implements a GraphQL API for httptoolkit-go that mirrors the
// TypeScript graphql-api.ts from httptoolkit-server.
package gql

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// Provider is the interface the GraphQL layer uses to talk to the rest of the
// server.  The concrete *api.Server satisfies this interface.
type Provider interface {
	Version() string
	CertificatePath() string
	CertificateContent() string
	CertificateFingerprint() string
	NetworkInterfaces() map[string]any
	SystemProxy() map[string]any
	DNSServers(proxyPort int) []string
	RuleParameterKeys() []string
	GetInterceptors(proxyPort int) []map[string]any
	GetInterceptor(id string) (map[string]any, error)
	IsInterceptorActive(id string, proxyPort int) bool
	InterceptorMetadata(id string, metaType string) (any, error)
	ActivateInterceptor(id string, proxyPort int, options map[string]any) (any, error)
	DeactivateInterceptor(id string, proxyPort int, options map[string]any) (bool, error)
	TriggerUpdate()
	TriggerShutdown()
}

// Handler returns an http.Handler that serves the GraphQL endpoint.
func Handler(p Provider) http.Handler {
	schema := buildSchema(p)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST requests are supported", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Query         string         `json:"query"`
			Variables     map[string]any `json:"variables"`
			OperationName string         `json:"operationName"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  req.Query,
			VariableValues: req.Variables,
			OperationName:  req.OperationName,
			Context:        context.Background(),
		})

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	})
}

// ──────────────────────────────────────────────
// Schema definition
// ──────────────────────────────────────────────

var jsonScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:         "Json",
	Description:  "An arbitrary JSON value",
	Serialize:    func(value any) any { return value },
	ParseValue:   func(value any) any { return value },
	ParseLiteral: func(valueAST ast.Value) any { return nil },
})

var voidScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:         "Void",
	Description:  "Represents no meaningful return value",
	Serialize:    func(value any) any { return nil },
	ParseValue:   func(value any) any { return nil },
	ParseLiteral: func(valueAST ast.Value) any { return nil },
})

var metadataTypeEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "MetadataType",
	Values: graphql.EnumValueConfigMap{
		"SUMMARY":  {Value: "summary"},
		"DETAILED": {Value: "detailed"},
	},
})

var proxyType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Proxy",
	Fields: graphql.Fields{
		"proxyUrl": {Type: graphql.NewNonNull(graphql.String)},
		"noProxy":  {Type: graphql.NewList(graphql.NewNonNull(graphql.String))},
	},
})

var interceptionConfigType = graphql.NewObject(graphql.ObjectConfig{
	Name: "InterceptionConfig",
	Fields: graphql.Fields{
		"certificatePath":        {Type: graphql.NewNonNull(graphql.String)},
		"certificateContent":     {Type: graphql.NewNonNull(graphql.String)},
		"certificateFingerprint": {Type: graphql.NewNonNull(graphql.String)},
	},
})

func buildInterceptorType(p Provider) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Interceptor",
		Fields: graphql.Fields{
			"id":          {Type: graphql.NewNonNull(graphql.ID)},
			"version":     {Type: graphql.NewNonNull(graphql.String)},
			"isActivable": {Type: graphql.NewNonNull(graphql.Boolean)},
			"isActive": {
				Type: graphql.NewNonNull(graphql.Boolean),
				Args: graphql.FieldConfigArgument{
					"proxyPort": {Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (any, error) {
					interceptor, _ := params.Source.(map[string]any)
					id, _ := interceptor["id"].(string)
					proxyPort, _ := params.Args["proxyPort"].(int)
					return p.IsInterceptorActive(id, proxyPort), nil
				},
			},
			"metadata": {
				Type: jsonScalar,
				Args: graphql.FieldConfigArgument{
					"type": {Type: metadataTypeEnum},
				},
				Resolve: func(params graphql.ResolveParams) (any, error) {
					interceptor, _ := params.Source.(map[string]any)
					id, _ := interceptor["id"].(string)
					metaType := "summary"
					if t, ok := params.Args["type"].(string); ok && t != "" {
						metaType = t
					}
					return p.InterceptorMetadata(id, metaType)
				},
			},
		},
	})
}

func buildSchema(p Provider) graphql.Schema {
	interceptorType := buildInterceptorType(p)

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"version": {
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(_ graphql.ResolveParams) (any, error) {
					return p.Version(), nil
				},
			},
			"config": {
				Type: graphql.NewNonNull(interceptionConfigType),
				Resolve: func(_ graphql.ResolveParams) (any, error) {
					return map[string]any{
						"certificatePath":        p.CertificatePath(),
						"certificateContent":     p.CertificateContent(),
						"certificateFingerprint": p.CertificateFingerprint(),
					}, nil
				},
			},
			"interceptors": {
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(interceptorType))),
				Resolve: func(params graphql.ResolveParams) (any, error) {
					proxyPort := 0
					if pp, ok := params.Args["proxyPort"].(int); ok {
						proxyPort = pp
					}
					return p.GetInterceptors(proxyPort), nil
				},
			},
			"interceptor": {
				Type: graphql.NewNonNull(interceptorType),
				Args: graphql.FieldConfigArgument{
					"id": {Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(params graphql.ResolveParams) (any, error) {
					id, _ := params.Args["id"].(string)
					return p.GetInterceptor(id)
				},
			},
			"networkInterfaces": {
				Type: jsonScalar,
				Resolve: func(_ graphql.ResolveParams) (any, error) {
					return p.NetworkInterfaces(), nil
				},
			},
			"systemProxy": {
				Type: proxyType,
				Resolve: func(_ graphql.ResolveParams) (any, error) {
					sp := p.SystemProxy()
					if sp == nil {
						return nil, nil
					}
					return sp, nil
				},
			},
			"dnsServers": {
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.String))),
				Args: graphql.FieldConfigArgument{
					"proxyPort": {Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (any, error) {
					proxyPort, _ := params.Args["proxyPort"].(int)
					servers := p.DNSServers(proxyPort)
					// Convert []string -> []any for graphql-go
					out := make([]any, len(servers))
					for i, s := range servers {
						out[i] = s
					}
					return out, nil
				},
			},
			"ruleParameterKeys": {
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.String))),
				Resolve: func(_ graphql.ResolveParams) (any, error) {
					keys := p.RuleParameterKeys()
					out := make([]any, len(keys))
					for i, k := range keys {
						out[i] = k
					}
					return out, nil
				},
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"activateInterceptor": {
				Type: jsonScalar,
				Args: graphql.FieldConfigArgument{
					"id":        {Type: graphql.NewNonNull(graphql.ID)},
					"proxyPort": {Type: graphql.NewNonNull(graphql.Int)},
					"options":   {Type: jsonScalar},
				},
				Resolve: func(params graphql.ResolveParams) (any, error) {
					id, _ := params.Args["id"].(string)
					proxyPort, _ := params.Args["proxyPort"].(int)
					var options map[string]any
					if o, ok := params.Args["options"]; ok && o != nil {
						options, _ = o.(map[string]any)
					}
					return p.ActivateInterceptor(id, proxyPort, options)
				},
			},
			"deactivateInterceptor": {
				Type: graphql.NewNonNull(graphql.Boolean),
				Args: graphql.FieldConfigArgument{
					"id":        {Type: graphql.NewNonNull(graphql.ID)},
					"proxyPort": {Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (any, error) {
					id, _ := params.Args["id"].(string)
					proxyPort, _ := params.Args["proxyPort"].(int)
					return p.DeactivateInterceptor(id, proxyPort, nil)
				},
			},
			"triggerUpdate": {
				Type: voidScalar,
				Resolve: func(_ graphql.ResolveParams) (any, error) {
					p.TriggerUpdate()
					return nil, nil
				},
			},
			"shutdown": {
				Type: voidScalar,
				Resolve: func(_ graphql.ResolveParams) (any, error) {
					p.TriggerShutdown()
					return nil, nil
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	return schema
}

// networkInterfacesToMap converts net.Interfaces to the JSON map format the UI expects.
func networkInterfacesToMap() map[string]any {
	ifaces, _ := net.Interfaces()
	out := make(map[string]any, len(ifaces))
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		entries := make([]any, 0, len(addrs))
		for _, a := range addrs {
			var ip net.IP
			switch v := a.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil {
				continue
			}
			family := "IPv6"
			if v4 := ip.To4(); v4 != nil {
				ip = v4
				family = "IPv4"
			}
			entries = append(entries, map[string]any{
				"address":  ip.String(),
				"family":   family,
				"internal": ip.IsLoopback() || ip.IsLinkLocalUnicast(),
			})
		}
		if len(entries) > 0 {
			out[iface.Name] = entries
		}
	}
	return out
}
