package client

// noisyKeys are HATEOAS/metadata fields from Ziti controller API responses
// that consume tokens without providing value to an AI agent.
var noisyKeys = map[string]bool{
	"_links":                    true,
	"authPolicy":                true,
	"type":                      true,
	"authenticators":            true,
	"envInfo":                   true,
	"sdkInfo":                   true,
	"serviceHostingCosts":       true,
	"serviceHostingPrecedences": true,
	"enrollment":                true,
	"appData":                   true,
	"interfaces":                true,
	"tags":                      true,
	"filterableFields":          true,
}

// StripNoise recursively removes noisy keys from controller API responses
// to reduce token consumption for AI agents.
func StripNoise(value any) any {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case map[string]any:
		cleaned := make(map[string]any)
		for k, val := range v {
			if noisyKeys[k] {
				continue
			}
			stripped := StripNoise(val)
			if stripped != nil {
				cleaned[k] = stripped
			}
		}
		if len(cleaned) == 0 {
			return nil
		}
		return cleaned

	case []any:
		var cleaned []any
		for _, item := range v {
			stripped := StripNoise(item)
			if stripped != nil {
				cleaned = append(cleaned, stripped)
			}
		}
		if len(cleaned) == 0 {
			return nil
		}
		return cleaned

	default:
		return value
	}
}
