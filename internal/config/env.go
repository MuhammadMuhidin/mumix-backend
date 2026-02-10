package config

import "os"

type AppConfig struct {
	Port            string
	DatabaseURL     string
	GinMode         string
	SupabaseURL     string
	SupabaseRoleKey string
}

func Load() AppConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = "debug"
	}

	return AppConfig{
		Port:            port,
		GinMode:         mode,
		DatabaseURL:     os.Getenv("SUPABASE_DB_URL"),
		SupabaseURL:     os.Getenv("SUPABASE_URL"),
		SupabaseRoleKey: os.Getenv("SUPABASE_ROLE_KEY"),
	}
}