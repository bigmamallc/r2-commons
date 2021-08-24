package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/bigmamallc/r2-commons/api/tcp"
	"github.com/go-redis/redis/v8"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func WithRedis(f func(r *redis.Client)) error {
	port := tcp.MustFindFreePort()
	portStr := strconv.Itoa(port)

	// https://stackoverflow.com/a/55297762
	out, err := exec.Command("docker", "run", "--rm", "-d", "-p", portStr+":6379", "redis", "sh", "-c",
		"rm -f /data/dump.rdb && redis-server").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to launch the redis container (%w): %s", err, string(out))
	}
	ctrID := strings.TrimSpace(string(out))
	defer func() {
		_ = exec.Command("docker", "kill", ctrID).Run()
	}()

	r := redis.NewClient(&redis.Options{
		Addr: "localhost:" + portStr,
	})

	deadline := time.Now().Add(time.Duration(5) * time.Second)
	for {
		if _, err := r.Ping(context.Background()).Result(); err == nil {
			break
		}

		if time.Now().After(deadline) {
			return errors.New("timed out waiting for redis")
		}
	}

	f(r)

	return nil
}

func MustWithRedis(f func(r *redis.Client)) {
	err := WithRedis(f)
	if err != nil {
		panic(err)
	}
}
