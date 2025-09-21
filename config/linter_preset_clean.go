package config

func LinterPresetClean() []LinterRule {
	return []LinterRule{
		{
			Path: "/internal/services/**",
			Deny: []string{
				"/internal/repositories/**",
				"/internal/clients/**",
				"/internal/services/**",
				"/internal/usecases/**",
				"/internal/controllers/**",
				"/internal/applications/**",
			},
			Type:        "strict",
			OnlyInner:   true,
			Description: "services can't depend on other service, use case, repository, client, controller, application",
		},
		{
			Path: "/internal/services/**",
			Allow: []string{
				"/internal/entities/**",
				"/internal/dto/**",
				"/internal/entities",
				"/internal/dto",
			},
			Type:        "lax",
			OnlyInner:   true,
			Description: "services can depend on entities or dto only",
		},
		{
			Path: "/internal/usecases/**",
			Deny: []string{
				"/internal/repositories/**",
				"/internal/clients/**",
				"/internal/services/**",
				"/internal/usecases/**",
				"/internal/controllers/**",
				"/internal/applications/**",
			},
			Type:        "strict",
			OnlyInner:   true,
			Description: "use cases can't depend on other use case, service, repository, client, controller, application",
		},
		{
			Path: "/internal/usecases/**",
			Allow: []string{
				"/internal/entities/**",
				"/internal/dto/**",
				"/internal/entities",
				"/internal/dto",
			},
			Type:        "lax",
			OnlyInner:   true,
			Description: "use cases can depend on entities or dto only",
		},
		{
			Path: "/internal/repositories/**",
			Deny: []string{
				"/internal/repositories/**",
				"/internal/clients/**",
				"/internal/services/**",
				"/internal/usecases/**",
				"/internal/controllers/**",
				"/internal/applications/**",
			},
			Type:        "strict",
			OnlyInner:   true,
			Description: "repositories can't depend on other repository, client, use case, service, controller, application",
		},
		{
			Path: "/internal/repositories/**",
			Allow: []string{
				"/internal/entities/**",
				"/internal/dto/**",
				"/internal/entities",
				"/internal/dto",
			},
			Type:        "lax",
			OnlyInner:   true,
			Description: "repositories can depend on entities or dto only",
		},
		{
			Path: "/internal/clients/**",
			Deny: []string{
				"/internal/repositories/**",
				"/internal/clients/**",
				"/internal/services/**",
				"/internal/usecases/**",
				"/internal/controllers/**",
				"/internal/applications/**",
			},
			Type:        "strict",
			OnlyInner:   true,
			Description: "clients can't depend on other client, repository, use case, service, controller, application",
		},
		{
			Path: "/internal/clients/**",
			Allow: []string{
				"/internal/entities/**",
				"/internal/dto/**",
				"/internal/entities",
				"/internal/dto",
			},
			Type:        "lax",
			OnlyInner:   true,
			Description: "clients can depend on entities or dto only",
		},
		{
			Path: "/internal/controllers/**",
			Deny: []string{
				"/internal/repositories/**",
				"/internal/clients/**",
				"/internal/services/**",
				"/internal/usecases/**",
				"/internal/controllers/**",
				"/internal/applications/**",
			},
			Type:        "strict",
			OnlyInner:   true,
			Description: "controllers can't depend on other controller, repository, client, use case, service, application",
		},
		{
			Path: "/internal/controllers/**",
			Allow: []string{
				"/internal/entities/**",
				"/internal/dto/**",
				"/internal/entities",
				"/internal/dto",
			},
			Type:        "lax",
			OnlyInner:   true,
			Description: "controllers can depend on entities or dto only",
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
