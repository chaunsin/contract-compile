package compile

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/chaunsin/contract-compile/compile/cmd"

	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	uuid "github.com/satori/go.uuid"
)

type dockerClient struct {
	cfg  *Config
	cli  *client.Client
	auth string
}

func NewDocker(cfg *Config) (ContractCompile, error) {
	var opts = []client.Opt{client.WithAPIVersionNegotiation()}

	switch cfg.Mode {
	case "ca":
		// 使用远程证书方式访问
		opts = append(opts, client.WithHost(cfg.Host), client.WithTLSClientConfig(cfg.CaCertPath, cfg.CertPath, cfg.KeyPath))
	case "ssh":
		// todo:考虑可配置私钥,不使用默认 ~/.ssh目录下得私钥
		// 使用ssh远程方式链接
		// 格式 ssh://root@ip:port 链接得时候需要在本机~/.ssh目录下存放私钥文件,然后服务器端存放证书。
		// -y 为ssh执行附加参数相当于 ssh -y root@ip:port 操作,这样操作避免出现建立验证失败问题。
		helper, err := connhelper.GetConnectionHelperWithSSHOpts("ssh://root@121.229.22.205:22", []string{"-y"})
		if err != nil {
			panic(err)
		}
		httpClient := &http.Client{
			Transport: &http.Transport{
				DialContext: helper.Dialer,
			},
		}
		opts = append(opts,
			client.WithHTTPClient(httpClient),
			client.WithHost(helper.Host),
			client.WithDialContext(helper.Dialer))
	case "host":
		fallthrough
	default:
		// 使用本地宿主机服务配置连接方式
		opts = append(opts, client.FromEnv)
	}

	// 根据配置生成访问docker-hub、registry、harbor等镜像仓库凭证
	auth, err := registry.EncodeAuthConfig(cfg.AuthConfig)
	if err != nil {
		return nil, fmt.Errorf("EncodeAuthConfig: %w", err)
	}

	// 建立连接
	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, fmt.Errorf("NewClientWithOpts: %w", err)
	}
	pong, err := cli.Ping(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("Ping: %w", err)
	}
	log.Println("pong:", pong)

	c := &dockerClient{
		cli:  cli,
		cfg:  cfg,
		auth: auth,
	}
	return c, nil
}

func (c *dockerClient) Close(ctx context.Context) error {
	if c.cli != nil {
		return c.cli.Close()
	}
	return nil
}

func (c *dockerClient) Ping(ctx context.Context) error {
	_, err := c.cli.Ping(ctx)
	return err
}

func (c *dockerClient) Execute(ctx context.Context, args cmd.Args) error {
	if err := args.Valid(); err != nil {
		return err
	}
	target, ok := cmd.Get(args.Organization)
	fmt.Printf("org:%s value:%+v\n", args.Organization, target)
	if !ok {
		return fmt.Errorf("%s organization is not found", args.Organization)
	}
	fn, ok := target[args.Images]
	if !ok {
		return fmt.Errorf("%s images is not found", args.Images)
	}
	_, cmdStr, err := fn.Cmd(args)
	if err != nil {
		return fmt.Errorf("Cmd: %w", err)
	}

	var (
		name    = uuid.NewV4().String()
		contain = &container.Config{
			Image:      args.Images.String(),
			Entrypoint: cmdStr,
		}
		hostConfig = &container.HostConfig{
			Privileged: true,
			Mounts: []mount.Mount{
				{
					Type:        mount.TypeBind,
					Source:      args.HostDir, // 只能是绝对路径
					Target:      args.TargetDir,
					ReadOnly:    false,
					Consistency: "",
				},
			},
		}
	)

	// 判断镜像是否存在不存在则拉取镜像
	_, _, err = c.cli.ImageInspectWithRaw(ctx, args.Images.String())
	if err != nil {
		if client.IsErrNotFound(err) {
			// 拉取镜像
			reader, err := c.cli.ImagePull(ctx, args.Images.String(), types.ImagePullOptions{})
			if err != nil {
				return fmt.Errorf("ImagePull(%s): %w", args.Images.String(), err)
			}
			defer reader.Close()
			// io.Copy(os.Stdout, reader)
			info, _ := io.ReadAll(reader)
			log.Printf("CompileContract ImagePull:%s", string(info))
		} else {
			return fmt.Errorf("ImageInspectWithRaw(%s): %w", args.Images.String(), err)
		}
	}

	// // 拉取镜像
	// reader, err := c.cli.ImagePull(ctx, args.Images.String(), types.ImagePullOptions{})
	// if err != nil {
	//	return fmt.Errorf("ImagePull(%s): %w", args.Images.String(), err)
	// }
	// defer reader.Close()
	// io.Copy(os.Stdout, reader)

	// 创建容器
	resp, err := c.cli.ContainerCreate(ctx, contain, hostConfig, nil, nil, name)
	if err != nil {
		panic(err)
	}
	log.Printf("WARNING:%v\n", resp.Warnings)

	// 启动容器
	if err := c.cli.ContainerStart(ctx, name, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	defer func() {
		if err := c.cli.ContainerRemove(ctx, name, types.ContainerRemoveOptions{}); err != nil {
			log.Printf("ContainerRemove(%s): %s", name, err)
		}
	}()

	// 等待容器完成并退出,获取退出码
	wait, werr := c.cli.ContainerWait(ctx, name, container.WaitConditionNotRunning)
	select {
	case w := <-wait:
		fmt.Printf("ContainerWait: %+v\n", w)
	case e := <-werr:
		return fmt.Errorf("ContainerWait: %w", e)
	}

	// 查看容器日志
	out, err := c.cli.ContainerLogs(ctx, name, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return fmt.Errorf("ContainerLogs: %w", err)
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return nil
}
