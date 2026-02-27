package util

import (
	"fmt"
	"sort"
	"strings"
)

func FlattenMap(input map[string]any) map[string]string {
	out := make(map[string]string)
	flatten("", input, out)
	return out
}

func flatten(prefix string, value any, out map[string]string) {
	switch typed := value.(type) {
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for k := range typed {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			next := key
			if prefix != "" {
				next = prefix + "." + key
			}
			flatten(next, typed[key], out)
		}
	case []any:
		parts := make([]string, 0, len(typed))
		for _, item := range typed {
			parts = append(parts, fmt.Sprint(item))
		}
		out[prefix] = strings.Join(parts, ",")
	default:
		out[prefix] = fmt.Sprint(typed)
	}
}

func IsSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	sensitiveWords := []string{"password", "secret", "token", "apikey", "private-key", "access-key", "credential"}
	for _, word := range sensitiveWords {
		if strings.Contains(lower, word) {
			return true
		}
	}
	return false
}

func MaskValue(value string) string {
	if value == "" {
		return value
	}
	if len(value) <= 4 {
		return "****"
	}
	return value[:2] + "****" + value[len(value)-2:]
}

func IsCriticalConfigKey(key string) bool {
	criticalKeys := []string{
		"spring.datasource.url",
		"spring.datasource.username",
		"spring.datasource.host",
		"spring.datasource.port",
		"spring.datasource.database",
		"spring.cloud.nacos.server-addr",
		"spring.cloud.nacos.namespace",
		"spring.cloud.nacos.group",
		"spring.cloud.nacos.discovery.server-addr",
		"spring.cloud.nacos.discovery.namespace",
		"spring.cloud.nacos.discovery.group",
		"spring.cloud.nacos.config.server-addr",
		"spring.cloud.nacos.config.namespace",
		"spring.cloud.nacos.config.group",
	}
	lower := strings.ToLower(key)
	for _, item := range criticalKeys {
		if lower == item {
			return true
		}
	}
	if IsSensitiveKey(lower) {
		return true
	}
	return false
}
