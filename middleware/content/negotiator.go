package content

import (
	"net/http"
	"strconv"
	"strings"
)

type AcceptRange struct {
	Type       string
	SubType    string
	Weight     float64
	Parameters map[string]string
	raw        string
}

func (r *AcceptRange) Raw() string {
	return r.raw
}

func AcceptMediaTypes(request *http.Request) []AcceptRange {
	result := []AcceptRange{}

	for _, value := range request.Header["Accept"] {
		result = append(result, ParseAcceptRanges(value)...)
	}

	return result
}

func ParseAcceptRanges(accepts string) []AcceptRange {
	result := []AcceptRange{}
	remaining := accepts

	for {
		var accept string

		accept, remaining = extractFieldAndSkipToken(remaining, ',')
		result = append(result, ParseAcceptRange(accept))

		if len(remaining) == 0 {
			break
		}
	}

	return result
}

func ParseAcceptRange(accept string) AcceptRange {
	typeAndSub, rawParams := extractFieldAndSkipToken(accept, ';')

	mainType, subType := extractFieldAndSkipToken(typeAndSub, '/')
	parameters := extractParams(rawParams)
	weight := extractWeight(parameters)

	return AcceptRange{
		Type:       mainType,
		SubType:    subType,
		Parameters: parameters,
		Weight:     weight,
		raw:        accept,
	}
}

func extractWeight(parameters map[string]string) float64 {
	if weight, ok := parameters["q"]; ok {
		if result, err := strconv.ParseFloat(weight, 64); err == nil {
			return result
		}
	}
	return 1
}

func extractParams(raw string) map[string]string {
	result := map[string]string{}
	rest := raw

	for {
		var param string
		param, rest = extractFieldAndSkipToken(rest, ';')

		if len(param) > 0 {
			key, value := extractFieldAndSkipToken(param, '=')
			result[key] = value
		}

		if len(rest) == 0 {
			break
		}
	}

	return result
}

func extractFieldAndSkipToken(value string, separator rune) (string, string) {
	f, r := extractField(value, separator)

	if len(r) > 0 {
		r = r[1:]
	}

	return f, r
}

func extractField(value string, separator rune) (field, rest string) {
	field = value

	for i, word := range value {
		if word == separator {
			field = strings.TrimSpace(value[:i])
			rest = strings.TrimSpace(value[i:])
			break
		}
	}

	return
}

func negotiateContentType(accepts []AcceptRange, offers []AcceptRange, defaultOffer AcceptRange) string {
	best := defaultOffer.Raw()
	weight := defaultOffer.Weight
	params := 0

	for _, offer := range offers {
		for _, accept := range accepts {
			booster := float64(0)

			if accept.Type != "*" {
				booster++

				if accept.SubType != "*" {
					booster++
				}
			}

			if weight > (accept.Weight + booster) {
				continue
			} else if accept.Type == "*" && accept.SubType == "*" && ((accept.Weight + booster) > weight) {
				best = offer.Raw()
				weight = accept.Weight + booster
			} else if accept.SubType == "*" && offer.Type == accept.Type && ((accept.Weight + booster) > weight) {
				best = offer.Raw()
				weight = accept.Weight + booster
			} else if accept.Type == offer.Type && accept.SubType == offer.SubType {
				paramCount := compareParams(accept.Parameters, offer.Parameters)

				if paramCount >= params {
					best = offer.Raw()
					weight = accept.Weight + booster
					params = paramCount
				}
			}
		}
	}

	return best
}

func compareParams(a map[string]string, b map[string]string) (count int) {
	for key, value := range a {
		if v, ok := b[key]; ok && value == v {
			count++
		}
	}

	return count
}

func NegotiateContentType(request *http.Request, offers []string, defaultOffer string) string {
	accepts := AcceptMediaTypes(request)
	ranges := []AcceptRange{}

	for _, offer := range offers {
		ranges = append(ranges, ParseAcceptRange(offer))
	}

	return negotiateContentType(accepts, ranges, ParseAcceptRange(defaultOffer))
}
