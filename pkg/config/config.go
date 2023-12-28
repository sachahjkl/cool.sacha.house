package config

import "github.com/spf13/viper"

type Config struct {
	Port                 string `mapstructure:"PORT"`
	DBUrl                string `mapstructure:"DB_URL"`
	EncryptionKey        string `mapstructure:"ENCRYPTION_KEY"`
	Version              string `mapstructure:"VERSION"`
	CommitSha            string `mapstructure:"COMMIT_SHA"`
	GithubBearerToken    string `mapstructure:"GITHUB_BEARER_TOKEN"`
	GitlabBearerToken    string `mapstructure:"GITLAB_BEARER_TOKEN"`
	ProxycurlBearerToken string `mapstructure:"PROXYCURL_BEARER_TOKEN"`
	AdminIps             string `mapstructure:"ADMIN_IPS"`
	GitlabApiEndpoint    string `mapstructure:"GITLAB_API_ENDPOINT"`
	GithubApiEndpoint    string `mapstructure:"GITHUB_API_ENDPOINT"`
	ProxycurlApiEndpoint string `mapstructure:"PROXYCURL_API_ENDPOINT"`
	HygraphApiEndpoint   string `mapstructure:"HYGRAPH_API_ENDPOINT"`
	LinkedinGistId       string `mapstructure:"LINKEDIN_GIST_ID"`
	GitRepoId            string `mapstructure:"GIT_REPO_ID"`
}

func LoadConfig() (c *Config, err error) {

	viper.AddConfigPath("./configs/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// need to manually bind otherwise values are empty (why ?)
	viper.BindEnv("VERSION")
	viper.BindEnv("COMMIT_SHA")

	err = viper.ReadInConfig()

	// if the default files could not be set, we define some default values
	if err != nil {
		viper.SetDefault("PORT", "3000")
		viper.SetDefault("DB_URL", "./last_resort.db")
		viper.SetDefault("ENCRYPTION_KEY", "ABCDEFGHJKLMNOPQRSTUVWXYZ")
		viper.SetDefault("VERSION", "v0.0.0+empty")
		viper.SetDefault("COMMIT_SHA", "abcdefg")
	}

	err = viper.Unmarshal(&c)

	return
}
