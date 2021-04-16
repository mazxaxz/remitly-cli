package deploy

import (
	log "github.com/sirupsen/logrus"
)

type profileContext struct {
	name string
	http httpSpec
}

type httpSpec struct {
	url      string
	username string
}

func profileContextFrom(source map[string]interface{}, profile string) (profileContext, error) {
	var contexts []interface{}
	if val, exists := source["contexts"]; exists {
		var ok bool
		contexts, ok = val.([]interface{})
		if !ok {
			return profileContext{}, ErrInvalidContextsFileSyntax
		}
	} else {
		return profileContext{}, ErrInvalidContextsFileSyntax
	}

	// wanted to use viper.Get("contexts.0.name") etc, but accessing
	// nested values through arrays does not work as it should
	for _, ctx := range contexts {
		if ctxMap, ok := ctx.(map[interface{}]interface{}); ok {
			var pc profileContext
			if val, exists := ctxMap["name"]; exists {
				if pc.name, ok = val.(string); ok {
					if pc.name != profile {
						continue
					}
				} else {
					log.WithField("context", ctx).Warn("context name has to be a string, entry skipped")
					continue
				}
			} else {
				log.WithField("context", ctx).Warn("context name does not exist, entry skipped")
				continue
			}

			pc.http = httpSpec{}
			if val, exists := ctxMap["http"]; exists {
				if httpMap, ok := val.(map[interface{}]interface{}); ok {

					if v, exists := httpMap["url"]; exists {
						if pc.http.url, ok = v.(string); !ok {
							log.WithField("context", ctx).Warn("contexts[].http.url has to be a string, entry skipped")
							continue
						}
					} else {
						log.WithField("context", ctx).Warn("contexts[].http.url was not specified, entry skipped")
						continue
					}

					if v, exists := httpMap["username"]; exists {
						if pc.http.username, ok = v.(string); !ok {
							log.WithField("context", ctx).Warn("contexts[].http.username has to be a string, entry skipped")
							continue
						}
					} else {
						log.WithField("context", ctx).Warn("contexts[].http.username was not specified, entry skipped")
						continue
					}
				} else {
					log.WithField("context", ctx).Warn("contexts[].http has invalid syntax, entry skipped")
					continue
				}
			} else {
				log.WithField("context", ctx).Warn("contexts[].http was not specified, entry skipped")
				continue
			}

			return pc, nil
		}
	}
	return profileContext{}, ErrProfileNotFound
}
