package main

type ImageConfig struct {
	Entrypoint string
	Cmd        string
	Env        []string
	WorkingDir string
	User       string
}

func main() {
	imageConfig := ImageConfig{
		Entrypoint: "/whoami",
		Cmd:        "",
		Env: []string{
			"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
			"WHOAMI_NAME=/blanc/init.json",
		},
		WorkingDir: "/",
		User:       "",
	}
}
