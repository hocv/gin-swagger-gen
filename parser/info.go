package parser

var (
	version      = "version"
	title        = "title"
	description  = "description"
	contactName  = "contact.name"
	contactEmail = "contact.email"
	contactURL   = "contact.url"
	host         = "host"
	basePath     = "basepath"
	baseInfo     = []string{
		title,
		version,
		description,
		contactName,
		contactEmail,
		contactURL,
		host,
		basePath,
	}
	licenseInfo = []string{
		"license.name",
		"license.url",
	}
	tagInfo = []string{
		"tag.name",
		"tag.description",
		"tag.description.markdown",
		"tag.docs.url",
		"tag.docs.description",
	}
	securityInfo = []string{
		"securitydefinitions.basic",
		"securitydefinitions.apikey",
		"securitydefinitions.oauth2.application",
		"securitydefinitions.oauth2.implicit",
		"securitydefinitions.oauth2.password",
		"securitydefinitions.oauth2.accesscode",
	}
	defaultValue = map[string]string{
		title:        "Swagger Example API",
		version:      "1.0",
		description:  "This is a sample server Petstore server.",
		contactName:  "API Support",
		contactEmail: "http://www.swagger.io/support",
		contactURL:   "support@swagger.io",
		host:         "petstore.swagger.io",
		basePath:     "/v2",
	}
)
