package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestRedisConfig_PrintingMasksPassword(t *testing.T) {
	cfg := RedisConfig{RedisHost: "redis:6379", RedisPassword: "s3cr3t", RedisDb: 2}

	outputs := map[string]string{
		"%v":  fmt.Sprintf("%v", cfg),
		"%+v": fmt.Sprintf("%+v", cfg),
		"%#v": fmt.Sprintf("%#v", cfg),
		"%s":  fmt.Sprintf("%s", cfg),
		"ptr": fmt.Sprintf("%+v", &cfg),
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}

	outputs["json"] = string(data)

	for name, out := range outputs {
		if strings.Contains(out, "s3cr3t") || !strings.Contains(out, "***") || !strings.Contains(out, "redis:6379") {
			t.Errorf("%s = %s — want the password masked, other fields kept", name, out)
		}
	}

	if !strings.HasPrefix(outputs["%#v"], "config.RedisConfig{") {
		t.Errorf("%%#v = %s", outputs["%#v"])
	}

	if cfg.RedisPassword != "s3cr3t" {
		t.Fatal("printing must not change the config")
	}

	var back map[string]any
	if err := json.Unmarshal(data, &back); err != nil || back["redis_db"] != float64(2) {
		t.Fatalf("json keys: %v %v", back, err)
	}

	if out := fmt.Sprintf("%+v", RedisConfig{}); strings.Contains(out, "***") {
		t.Errorf("an empty password must stay empty: %s", out)
	}
}
