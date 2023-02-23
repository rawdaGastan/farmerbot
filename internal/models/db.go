// Package models for farmerbot
package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

// RedisManager represents interface for redis DB actions
type RedisManager interface {
	GetFarm() (Farm, error)
	GetPower() (Power, error)
	GetNode(nodeID uint32) (Node, error)
	GetNodes() ([]Node, error)
	UpdatesNodes(node Node) error
	SetNodes(nodes []Node) error
	SetFarm(farm Farm) error
	SetPower(power Power) error
	SaveConfig(config Config) error
	FilterOnNodes() ([]Node, error)
}

// Config is the configuration for farmerbot
type Config struct {
	Farm  Farm   `json:"farm"`
	Nodes []Node `json:"nodes"`
	Power Power  `json:"power"`
}

// RedisDB for saving config for farmerbot
type RedisDB struct {
	redis *redis.Client
}

// NewRedisDB generates new redis db
func NewRedisDB(address string) RedisDB {
	return RedisDB{
		redis: redis.NewClient(&redis.Options{
			Addr: address,
		}),
	}
}

// GetFarm gets farm from the database
func (db *RedisDB) GetFarm() (Farm, error) {
	var dest Farm
	nodes, err := db.redis.Get("farm").Bytes()
	if err != nil {
		return Farm{}, err
	}

	if err := json.Unmarshal(nodes, &dest); err != nil {
		return Farm{}, err
	}

	return dest, nil
}

// GetPower gets power from the database
func (db *RedisDB) GetPower() (Power, error) {
	var dest Power
	nodes, err := db.redis.Get("power").Bytes()
	if err != nil {
		return Power{}, err
	}

	if err := json.Unmarshal(nodes, &dest); err != nil {
		return Power{}, err
	}

	return dest, nil
}

// GetNode gets a node from the database
func (db *RedisDB) GetNode(nodeID uint32) (Node, error) {
	var dest []Node
	nodes, err := db.redis.Get("nodes").Bytes()
	if err != nil {
		return Node{}, err
	}

	if err := json.Unmarshal(nodes, &dest); err != nil {
		return Node{}, err
	}

	for _, n := range dest {
		if n.ID == nodeID {
			return n, nil
		}
	}

	return Node{}, fmt.Errorf("node %d not found", nodeID)
}

// GetNodes gets nodes from the database
func (db *RedisDB) GetNodes() ([]Node, error) {
	var dest []Node
	nodes, err := db.redis.Get("nodes").Bytes()
	if err != nil {
		return []Node{}, err
	}

	if err := json.Unmarshal(nodes, &dest); err != nil {
		return []Node{}, err
	}

	return dest, nil
}

// UpdatesNodes adds or updates a node in the database
func (db *RedisDB) UpdatesNodes(node Node) error {
	nodes, err := db.GetNodes()
	if err != nil {
		return err
	}

	found := false
	for i, n := range nodes {
		if n.ID == node.ID {
			nodes[i] = node
			found = true
		}
	}

	if !found {
		nodes = append(nodes, node)
	}

	n, err := json.Marshal(nodes)
	if err != nil {
		return err
	}

	return db.redis.Set("nodes", n, 0).Err()
}

// SetNodes sets the nodes in the database
func (db *RedisDB) SetNodes(nodes []Node) error {
	n, err := json.Marshal(nodes)
	if err != nil {
		return err
	}

	return db.redis.Set("nodes", n, 0).Err()
}

// SetFarm sets the farm in the database
func (db *RedisDB) SetFarm(farm Farm) error {
	f, err := json.Marshal(farm)
	if err != nil {
		return err
	}
	return db.redis.Set("farm", f, 0).Err()
}

// SetPower sets the power in the database
func (db *RedisDB) SetPower(power Power) error {
	p, err := json.Marshal(power)
	if err != nil {
		return err
	}
	return db.redis.Set("power", p, 0).Err()
}

// SaveConfig saves the configuration in the database
func (db *RedisDB) SaveConfig(config Config) error {
	if err := db.SetFarm(config.Farm); err != nil {
		return err
	}

	if err := db.SetPower(config.Power); err != nil {
		return err
	}

	return db.SetNodes(config.Nodes)
}

// FilterOnNodes filters db ON nodes
func (db *RedisDB) FilterOnNodes() ([]Node, error) {
	nodes, err := db.GetNodes()
	if err != nil {
		return []Node{}, errors.New("failed to get nodes from db")
	}

	out := make([]Node, 0)
	for _, node := range nodes {
		if node.PowerState.ON {
			out = append(out, node)
		}
	}
	return out, nil
}
