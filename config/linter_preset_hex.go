package config

func LinterPresetHex() []LinterRule {
	return []LinterRule{
		{
			Path: "/internal/usecases/**",
			Deny: []string{
				"/internal/adapters/**",
				"/internal/usecases/**",
				"/internal/controllers/**",
				"/internal/applications/**",
			},
			Type:        "strict",
			OnlyInner:   true,
			Description: "use cases can't depend on other use case, adapter, controller, application",
		},
		{
			Path: "/internal/usecases/**",
			Allow: []string{
				"/internal/entities/**",
				"/internal/entities",
			},
			Type:        "lax",
			OnlyInner:   true,
			Description: "use cases can depend on entities only",
		},
		{
			Path: "/internal/adapters/**",
			Deny: []string{
				"/internal/adapters/**",
				"/internal/usecases/**",
				"/internal/controllers/**",
				"/internal/applications/**",
			},
			Type:        "strict",
			OnlyInner:   true,
			Description: "adapters can't depend on other adapter, use case, controller, application",
		},
		{
			Path: "/internal/adapters/**",
			Allow: []string{
				"/internal/entities/**",
				"/internal/entities",
			},
			Type:        "lax",
			OnlyInner:   true,
			Description: "adapters can depend on entities only",
		},
		{
			Path: "/internal/controllers/**",
			Deny: []string{
				"/internal/adapters/**",
				"/internal/usecases/**",
				"/internal/controllers/**",
				"/internal/applications/**",
			},
			Type:        "strict",
			OnlyInner:   true,
			Description: "controllers can't depend on other controller, adapter, use case, application",
		},
		{
			Path: "/internal/controllers/**",
			Allow: []string{
				"/internal/entities/**",
				"/internal/entities",
			},
			Type:        "lax",
			OnlyInner:   true,
			Description: "controllers can depend on entities only",
		},
		{
			Path: "/internal/applications/**",
			Deny: []string{
				"/internal/applications/**",
			},
			Type:        "strict",
			OnlyInner:   true,
			Description: "applications can't depend on other application",
		},
		{
			Path: "/cmd/**",
			Allow: []string{
				"/internal/applications/**",
			},
			Deny: []string{
				"/cmd/**",
			},
			Type:        "lax",
			OnlyInner:   true,
			Description: "cmd can depend on application only",
		},
	}
}
