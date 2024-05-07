package checks

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/newrelic/newrelic-lambda-extension/config"
	"github.com/newrelic/newrelic-lambda-extension/lambda/extension/api"
	"github.com/newrelic/newrelic-lambda-extension/util"
)

type handlerConfigs struct {
	handlerName string
	conf        *config.Configuration
}

var handlerPath = util.GetEnvOrDefault("LAMBDA_TASK_ROOT", "/var/task")

func handlerCheck(ctx context.Context, conf *config.Configuration, reg *api.RegistrationResponse, r runtimeConfig) error {
	if r.language != "" {
		h := handlerConfigs{
			handlerName: reg.Handler,
			conf:        conf,
		}

		if !r.check(h) {
			return fmt.Errorf("missing handler file %s (NEW_RELIC_LAMBDA_HANDLER=%s)", h.handlerName, conf.NRHandler)
		}
	}

	return nil
}

func (r runtimeConfig) check(h handlerConfigs) bool {
	if !h.conf.TestingOverride {
		esm := strings.ToLower(os.Getenv("NEW_RELIC_USE_ESM"))
		if esm == "true" {
			return true
		}
		if isDockerEnvironment() {
			return true
		}
	}
	functionHandler := r.getTrueHandler(h)
	p := removePathMethodName(functionHandler)
	if r.language == Node {
		pJS := pathFormatter(p, "js")
		cJS := pathFormatter(p, "cjs")
		pMJS := pathFormatter(p, "mjs")

		if util.PathExists(pJS) || util.PathExists(pMJS) || util.PathExists(cJS) {
			return true
		}
	} else {
		p = pathFormatter(p, r.fileType)
	}
	return util.PathExists(p)
}

func (r runtimeConfig) getTrueHandler(h handlerConfigs) string {
	if h.handlerName != r.wrapperName {
		util.Logln("Warning: handler not set to New Relic layer wrapper", r.wrapperName)
		return h.handlerName
	}

	return h.conf.NRHandler
}

func removePathMethodName(p string) string {
	s := strings.Split(p, ".")
	return strings.Join(s[:len(s)-1], "/")
}

func pathFormatter(functionHandler string, fileType string) string {
	p := fmt.Sprintf("%s/%s.%s", handlerPath, functionHandler, fileType)
	return p
}

func isDockerEnvironment() bool {
	_, dockerHost := os.LookupEnv("DOCKER_HOST")
	_, dockerCgroup := os.LookupEnv("container")
	return dockerHost || dockerCgroup
}
