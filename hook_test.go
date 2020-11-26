package logrustash

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

type simpleFmter struct{}

func (f simpleFmter) Format(e *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("msg: %#v", e.Message)), nil
}

func TestFire(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	h := Hook{
		writer:    buffer,
		formatter: simpleFmter{},
	}

	entry := &logrus.Entry{
		Message: "my message",
		Data:    logrus.Fields{},
	}

	require.NoError(t, h.Fire(entry))

	assert.Equal(t, "msg: \"my message\"", buffer.String())
}

type FailFmt struct{}

var errorFailFmt = errors.New("fail format")

func (f FailFmt) Format(e *logrus.Entry) ([]byte, error) {
	return nil, errorFailFmt
}

func TestFireFormatError(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	h := Hook{
		writer:    buffer,
		formatter: FailFmt{},
	}

	assert.Error(t, errorFailFmt, h.Fire(&logrus.Entry{Data: logrus.Fields{}}))
}

type FailWrite struct{}

var errorFailWrite = errors.New("fail write")

func (w FailWrite) Write(d []byte) (int, error) {
	return 0, errorFailWrite
}

func TestFireWriteError(t *testing.T) {
	h := Hook{
		writer:    FailWrite{},
		formatter: &logrus.JSONFormatter{},
	}

	assert.Equal(t, errorFailWrite, h.Fire(&logrus.Entry{Data: logrus.Fields{}}))
}

func TestDefaultFormatterWithFields(t *testing.T) {
	format := DefaultFormatter(logrus.Fields{"ID": 123})

	entry := &logrus.Entry{
		Message: "msg1",
		Data:    logrus.Fields{"f1": "bla"},
	}

	res, err := format.Format(entry)
	require.NoError(t, err)

	expected := []string{
		"f1\":\"bla\"",
		"ID\":123",
		"message\":\"msg1\"",
	}

	for _, exp := range expected {
		assert.True(t, strings.Contains(string(res), exp))
	}
}

func TestDefaultFormatterWithEmptyFields(t *testing.T) {
	now := time.Now()
	formatter := DefaultFormatter(logrus.Fields{})

	entry := &logrus.Entry{
		Message: "message bla bla",
		Level:   logrus.DebugLevel,
		Time:    now,
		Data: logrus.Fields{
			"Key1": "Value1",
		},
	}

	res, err := formatter.Format(entry)
	require.NoError(t, err)

	expected := []string{
		"\"message\":\"message bla bla\"",
		"\"level\":\"debug\"",
		"\"Key1\":\"Value1\"",
		"\"@version\":\"1\"",
		"\"type\":\"log\"",
		fmt.Sprintf("\"@timestamp\":\"%s\"", now.Format(time.RFC3339Nano)),
	}

	for _, exp := range expected {
		assert.True(t, strings.Contains(string(res), exp))
	}
}

func TestLogstashFieldsNotOverridden(t *testing.T) {
	_ = DefaultFormatter(logrus.Fields{"user1": "11"})

	assert.NotContains(t, logstashFields, "user1")
}
