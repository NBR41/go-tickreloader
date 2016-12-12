package gotickreloader

import (
	"sync"
	"time"
)

// Client struct for reloader client
type Client struct {
	fGetter      getterFunc
	getterParams []interface{}
	reloadTick   time.Duration
	value        interface{}
	lastError    error
	loaded       bool
	exitChan     chan bool
	sync.Mutex
}

type getterFunc func(...interface{}) (interface{}, error)

// NewClient return a new Client instance
func NewClient(reloadTick time.Duration, fGetter getterFunc, getterParams ...interface{}) *Client {
	return &Client{
		fGetter:      fGetter,
		getterParams: getterParams,
		reloadTick:   reloadTick,
		exitChan:     make(chan bool),
	}
}

// StartTickReload start the reload
func (c *Client) StartTickReload() {
	go c.reload()
}

// StopTickReload stops the reload
func (c *Client) StopTickReload() {
	c.exitChan <- true
}

// Get returns the value
func (c *Client) Get() (interface{}, error) {
	c.Lock()
	defer c.Unlock()

	if !c.loaded {
		c.value, c.lastError = c.fGetter(c.getterParams...)
		c.loaded = true
	}

	return c.value, c.lastError
}

// reload
func (c *Client) reload() {
	tick := time.Tick(c.reloadTick)
	for {
		select {
		case <-c.exitChan:
			return
		case <-tick:
			c.Lock()
			c.value, c.lastError = c.fGetter(c.getterParams...)
			c.Unlock()
		}
	}
}
