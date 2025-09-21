package config

func LinterPresetHexLite() []LinterRule {
	return []LinterRule{
		{
			Path: "/internal/usecases/**",
			Deny: []string{
				"/internal/repositories/**",
				"/internal/clients/**",
				"/internal/usecases/**",
				"/internal/controllers/**",
				"/internal/applications/**",
			},
			Type:        "strict",
			OnlyInner:   true,
			Description: "use cases can't depend on other use case, repository, client, controller, application",
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
			Path: "/internal/repositories/**",
			Deny: []string{
				"/internal/repositories/**",
				"/internal/clients/**",
				"/internal/usecases/**",
				"/internal/controllers/**",
				"/internal/applications/**",
			},
			Type:        "strict",
			OnlyInner:   true,
			Description: "repositories can't depend on other repository, client, use case, controller, application",
		},
		{
			Path: "/internal/repositories/**",
			Allow: []string{
				"/internal/entities/**",
				"/internal/entities",
			},
			Type:        "lax",
			OnlyInner:   true,
			Description: "repositories can depend on entities only",
		},
		{
			Path: "/internal/clients/**",
			Deny: []string{
				"/internal/repositories/**",
				"/internal/clients/**",
				"/internal/usecases/**",
				"/internal/controllers/**",
				"/internal/applications/**",
			},
			Type:        "strict",
			OnlyInner:   true,
			Description: "clients can't depend on other client, repository, use case, controller, application",
		},
		{
			Path: "/internal/clients/**",
			Allow: []string{
				"/internal/entities/**",
				"/internal/entities",
			},
			Type:        "lax",
			OnlyInner:   true,
			Description: "clients can depend on entities only",
		},
		{
			Path: "/internal/controllers/**",
			Deny: []string{
				"/internal/repositories/**",
				"/internal/clients/**",
				"/internal/usecases/**",
				"/internal/controllers/**",
				"/internal/applications/**",
			},
			Type:        "strict",
			OnlyInner:   true,
			Description: "controllers can't depend on other controller, repository, client, use case, application",
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
