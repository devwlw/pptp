package docker

import (
	sc "context"
	"encoding/json"
	"log"
	"mail/config"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/client"
)

type SDK struct {
	version string
	client  *client.Client
	ctx     sc.Context
	cfg     *config.DeployConfig
}

func NewSDK(version string, cfg *config.DeployConfig) *SDK {
	c, err := client.NewClientWithOpts(client.WithVersion(version))
	if err != nil {
		panic(err)
	}
	ctx := sc.Background()
	return &SDK{
		version: version,
		client:  c,
		ctx:     ctx,
		cfg:     cfg,
	}
}

func (s *SDK) Containers() ([]types.Container, error) {
	return s.client.ContainerList(s.ctx, types.ContainerListOptions{All: true})
}

func (s *SDK) CreateContainer(name string) error {
	cfg := s.cfg
	env := []string{
		"SERVER=" + cfg.PPTP.Endpoint,
		"TUNNEL=vps",
		"USERNAME=" + cfg.PPTP.UserName,
		"PASSWORD=" + cfg.PPTP.Password,
		"HOSTIP=" + cfg.HostIp,
	}
	runConfig := &container.Config{
		Env:   env,
		Image: "pptp",
	}
	hostConfig := &container.HostConfig{
		NetworkMode: "macvlan",
		Privileged:  true,
	}
	body, err := s.client.ContainerCreate(s.ctx, runConfig, hostConfig, nil, name)
	if err != nil {
		return err
	}
	dd, _ := json.Marshal(body)
	log.Println(string(dd))
	return s.client.ContainerStart(s.ctx, body.ID, types.ContainerStartOptions{})
	// docker run -d --privileged
	/*
		args := []string{
			"run",
			"-d",
			"--privileged",
			"-e SERVER=" + cfg.PPTP.Endpoint,
			"-e TUNNEL=vps",
			"-e USERNAME=" + cfg.PPTP.UserName,
			"-e PASSWORD=" + cfg.PPTP.Password,
			"-e HOSTIP=" + cfg.HostIp,
			"--network=macvlan",
			"--name=" + name,
			"pptp",
		}
		log.Println("cmd: docker " + strings.Join(args, " "))
		cmd := exec.Command("docker", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()*/
}

func (s *SDK) DeleteContainer(id string) error {
	return s.client.ContainerRemove(s.ctx, id, types.ContainerRemoveOptions{
		Force: true,
	})
}

func (s *SDK) StartContainer(id string) error {
	return s.client.ContainerStart(s.ctx, id, types.ContainerStartOptions{})
}

func (s *SDK) StopContainer(id string) error {
	return s.client.ContainerStop(s.ctx, id, nil)
}

func (s *SDK) SendMail(id, mailType, receiver, title, body, username, password, mode string) error {
	res, err := s.client.ContainerExecCreate(s.ctx, id, types.ExecConfig{
		AttachStdin:  false,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          []string{"/deliver", "send", "-type=" + mailType, "-receiver=" + receiver, "-title=" + title, "-body=" + body, "-username=" + username, "-password=" + password, "-mode=" + mode},
	})
	if err != nil {
		return err
	}
	log.Println("execId:", res.ID)
	return s.client.ContainerExecStart(s.ctx, res.ID, types.ExecStartCheck{
		Detach: false,
		Tty:    false,
	})
}

/*func (s *SDK) Do(method, action string, body io.Reader) (*sjson.Json, error) {
	urlStr := fmt.Sprintf("%s/%s", s.EndPoint, action)
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	sj, err := sjson.NewFromReader(res.Body)
	if res.StatusCode != 200 {
		return nil, errors.New(sj.Get("message").MustString())
	}
	return sj, nil
}*/
