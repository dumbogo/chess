// +build integration

package messagebroker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testConfig = Config{
		URL: "localhost:4222",
	}
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	// nats.DefaultURL
	mb, err := New(testConfig)
	defer mb.Close()
	assert.Nil(err)

	_, err2 := New(Config{URL: "badcon:4222"})
	assert.NotNil(err2)
}

func TestPublish(t *testing.T) {
	assert := assert.New(t)
	mb, err := New(testConfig)
	defer mb.Close()
	assert.Nil(err)

	err = mb.Publish("sometopic", Message{Payload: []byte("someinfo")})
	assert.Nil(err)
}

func TestSubscribe(t *testing.T) {
	assert := assert.New(t)
	mb, err := New(testConfig)
	defer mb.Close()
	assert.Nil(err)

	assert.Nil(err)
	subCh, err := mb.Subscribe("sometopic")
	assert.Nil(err)

	message := Message{Payload: []byte("someinfo")}
	err = mb.Publish("sometopic", message)

	select {
	case m := <-subCh:
		assert.EqualValues(message, m)
	case <-time.After(time.Millisecond * 500):
		return
	}
}
