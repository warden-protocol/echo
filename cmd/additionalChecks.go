package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/rs/zerolog"
)

type NodeStatus struct {
	Result struct {
		SyncInfo struct {
			LatestBlockHeight string `json:"latest_block_height"`
		} `json:"sync_info"`
	} `json:"result"`
}

const (
	httpTimeout = 5
)

func performAdditionalChecks(cfg Config) map[string]bool {
	additionalChecks := make(map[string]bool)
	additionalChecks["node_behind"] = nodeBehind(cfg.Peers, cfg.BehindThreshold)
	return additionalChecks
}

func nodeBehind(nodes []string, threshold int64) bool {
	var remoteHeight int64

	log := log.New(
		log.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(log.InfoLevel).With().Timestamp().Logger()

	// get local height, if failed return false because behind cannot be confirmed
	localHeigth, err := getHeight("http://localhost:26657")
	if err != nil {
		log.Error().Msgf("Failed to get height: %v", err)
		return false
	}

	for _, node := range nodes {
		remoteHeight, err = getHeight(node)
		if err != nil {
			log.Error().Msgf("Failed to get height: %v", err)
			continue
		}
		if localHeigth < remoteHeight-threshold {
			log.Error().Msgf("Local node (%d) is behind %s (%d)", localHeigth, node, remoteHeight)
			return true
		}
	}

	log.Debug().Msgf("Current height: %d", localHeigth)
	return false
}

func getHeight(host string) (int64, error) {
	nodeStatus := NodeStatus{}
	cli := &http.Client{
		Timeout: httpTimeout * time.Second,
	}
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, host+"/status", nil)
	if err != nil {
		return 0, err
	}
	resp, err := cli.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if err = json.Unmarshal(b, &nodeStatus); err != nil {
		return 0, err
	}

	height, err := strconv.ParseInt(nodeStatus.Result.SyncInfo.LatestBlockHeight, 10, 64)
	if err != nil {
		return 0, err
	}

	return height, nil
}
