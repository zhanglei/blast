// Copyright (c) 2019 Minoru Osuka
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mosuka/blast/config"
	"github.com/mosuka/blast/dispatcher"
	"github.com/mosuka/blast/logutils"
	"github.com/urfave/cli"
)

func distributorNodeStart(c *cli.Context) error {
	managerAddr := c.String("cluster-grpc-address")

	grpcAddr := c.String("grpc-address")
	httpAddr := c.String("http-address")

	logLevel := c.GlobalString("log-level")
	logFilename := c.GlobalString("log-file")
	logMaxSize := c.GlobalInt("log-max-size")
	logMaxBackups := c.GlobalInt("log-max-backups")
	logMaxAge := c.GlobalInt("log-max-age")
	logCompress := c.GlobalBool("log-compress")

	grpcLogLevel := c.GlobalString("grpc-log-level")
	grpcLogFilename := c.GlobalString("grpc-log-file")
	grpcLogMaxSize := c.GlobalInt("grpc-log-max-size")
	grpcLogMaxBackups := c.GlobalInt("grpc-log-max-backups")
	grpcLogMaxAge := c.GlobalInt("grpc-log-max-age")
	grpcLogCompress := c.GlobalBool("grpc-log-compress")

	httpLogFilename := c.GlobalString("http-log-file")
	httpLogMaxSize := c.GlobalInt("http-log-max-size")
	httpLogMaxBackups := c.GlobalInt("http-log-max-backups")
	httpLogMaxAge := c.GlobalInt("http-log-max-age")
	httpLogCompress := c.GlobalBool("http-log-compress")

	// create logger
	logger := logutils.NewLogger(
		logLevel,
		logFilename,
		logMaxSize,
		logMaxBackups,
		logMaxAge,
		logCompress,
	)

	// create logger
	grpcLogger := logutils.NewGRPCLogger(
		grpcLogLevel,
		grpcLogFilename,
		grpcLogMaxSize,
		grpcLogMaxBackups,
		grpcLogMaxAge,
		grpcLogCompress,
	)

	// create HTTP access logger
	httpAccessLogger := logutils.NewApacheCombinedLogger(
		httpLogFilename,
		httpLogMaxSize,
		httpLogMaxBackups,
		httpLogMaxAge,
		httpLogCompress,
	)

	// create cluster config
	clusterConfig := config.DefaultClusterConfig()
	if managerAddr != "" {
		clusterConfig.ManagerAddr = managerAddr
	}

	// create node config
	nodeConfig := &config.NodeConfig{
		GRPCAddr: grpcAddr,
		HTTPAddr: httpAddr,
	}

	svr, err := dispatcher.NewServer(clusterConfig, nodeConfig, logger, grpcLogger, httpAccessLogger)
	if err != nil {
		return err
	}

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Kill, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go svr.Start()

	<-quitCh

	svr.Stop()

	return nil
}