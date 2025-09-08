package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"supershell/internal/store"
)

type command string

const (
	cmdAdd     command = "add"
	cmdUpdate  command = "update"
	cmdDelete  command = "delete"
	cmdList    command = "list"
	cmdGet     command = "get"
	cmdConnect command = "connect"
)

type Connection struct {
	Nickname   string `json:"nickname"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	User       string `json:"user"`
	AuthMethod string `json:"auth_method"` // key|password
	KeyPath    string `json:"key_path,omitempty"`
	Password   string `json:"password,omitempty"`
}

func Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: supershell <add|update|delete|list|get|connect> [flags]")
	}

	s, err := store.New()
	if err != nil {
		return err
	}

	switch command(args[0]) {
	case cmdAdd:
		return runAdd(s, args[1:])
	case cmdUpdate:
		return runUpdate(s, args[1:])
	case cmdDelete:
		return runDelete(s, args[1:])
	case cmdList:
		return runList(s)
	case cmdGet:
		return runGet(s, args[1:])
	case cmdConnect:
		return runConnect(s, args[1:])
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func parseCommonFlags(fs *flag.FlagSet, c *Connection) {
	fs.StringVar(&c.Nickname, "name", "", "nickname for the connection (required)")
	fs.StringVar(&c.Host, "host", "", "host or IP (required)")
	fs.IntVar(&c.Port, "port", 22, "ssh port")
	fs.StringVar(&c.User, "user", os.Getenv("USER"), "ssh username")
	fs.StringVar(&c.AuthMethod, "auth", "key", "auth method: key or password")
	fs.StringVar(&c.KeyPath, "key", os.ExpandEnv("$HOME/.ssh/id_rsa"), "path to private key")
	fs.StringVar(&c.Password, "password", "", "ssh password (not recommended)")
}

func runAdd(s *store.Store, args []string) error {
	var c Connection
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(new(strings.Builder))
	parseCommonFlags(fs, &c)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if c.Nickname == "" || c.Host == "" {
		return errors.New("add requires --name and --host")
	}
	if c.AuthMethod != "key" && c.AuthMethod != "password" {
		return errors.New("--auth must be key or password")
	}
	if c.AuthMethod == "key" && c.KeyPath == "" {
		return errors.New("--key required when --auth=key")
	}
	if c.AuthMethod == "password" && c.Password == "" {
		return errors.New("--password required when --auth=password")
	}
	return s.Add(store.Record{
		Nickname:   c.Nickname,
		Host:       c.Host,
		Port:       c.Port,
		User:       c.User,
		AuthMethod: c.AuthMethod,
		KeyPath:    c.KeyPath,
		Password:   c.Password,
	})
}

func runUpdate(s *store.Store, args []string) error {
	var c Connection
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	fs.SetOutput(new(strings.Builder))
	parseCommonFlags(fs, &c)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if c.Nickname == "" {
		return errors.New("update requires --name")
	}
	return s.Update(c.Nickname, func(r *store.Record) {
		if c.Host != "" {
			r.Host = c.Host
		}
		if c.Port != 0 {
			r.Port = c.Port
		}
		if c.User != "" {
			r.User = c.User
		}
		if c.AuthMethod != "" {
			r.AuthMethod = c.AuthMethod
		}
		if c.KeyPath != "" {
			r.KeyPath = c.KeyPath
		}
		if c.Password != "" {
			r.Password = c.Password
		}
	})
}

func runDelete(s *store.Store, args []string) error {
	fs := flag.NewFlagSet("delete", flag.ContinueOnError)
	fs.SetOutput(new(strings.Builder))
	name := fs.String("name", "", "nickname to delete (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *name == "" {
		return errors.New("delete requires --name")
	}
	return s.Delete(*name)
}

func runList(s *store.Store) error {
	recs := s.List()
	if len(recs) == 0 {
		fmt.Println("No connections saved.")
		return nil
	}
	for _, r := range recs {
		fmt.Printf("%s\t%s@%s:%d\t%s\n", r.Nickname, r.User, r.Host, r.Port, r.AuthMethod)
	}
	return nil
}

func runGet(s *store.Store, args []string) error {
	fs := flag.NewFlagSet("get", flag.ContinueOnError)
	fs.SetOutput(new(strings.Builder))
	name := fs.String("name", "", "nickname to get (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *name == "" {
		return errors.New("get requires --name")
	}
	r, err := s.Get(*name)
	if err != nil {
		return err
	}
	fmt.Printf("%s\t%s@%s:%d\t%s\n", r.Nickname, r.User, r.Host, r.Port, r.AuthMethod)
	if r.KeyPath != "" {
		fmt.Printf("key:\t%s\n", r.KeyPath)
	}
	if r.Password != "" {
		fmt.Printf("password:\t%s\n", strings.Repeat("*", len(r.Password)))
	}
	return nil
}

func runConnect(s *store.Store, args []string) error {
	fs := flag.NewFlagSet("connect", flag.ContinueOnError)
	fs.SetOutput(new(strings.Builder))
	name := fs.String("name", "", "nickname to connect to (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *name == "" {
		return errors.New("connect requires --name")
	}
	r, err := s.Get(*name)
	if err != nil {
		return err
	}
	sshArgs := []string{"-p", fmt.Sprintf("%d", r.Port)}
	if r.AuthMethod == "key" && r.KeyPath != "" {
		sshArgs = append(sshArgs, "-i", r.KeyPath)
	}
	sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", r.User, r.Host))
	cmd := exec.Command("ssh", sshArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
