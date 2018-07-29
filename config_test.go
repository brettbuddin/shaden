package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArgs_Available(t *testing.T) {
	var tests = []struct {
		args  []string
		check func(*testing.T, Config)
	}{
		{
			args: []string{"-seed", "1"},
			check: func(t *testing.T, cfg Config) {
				assert.Equal(t, int64(1), cfg.Seed)
			},
		},
		{
			args: []string{"-frame", "256"},
			check: func(t *testing.T, cfg Config) {
				assert.Equal(t, 256, cfg.FrameSize)
			},
		},
		{
			args: []string{"-addr", ":3000"},
			check: func(t *testing.T, cfg Config) {
				assert.Equal(t, ":3000", cfg.HTTPAddr)
			},
		},
		{
			args: []string{"-addr", ":3000"},
			check: func(t *testing.T, cfg Config) {
				assert.Equal(t, ":3000", cfg.HTTPAddr)
			},
		},
		{
			args: []string{"-repl"},
			check: func(t *testing.T, cfg Config) {
				assert.True(t, cfg.REPL)
			},
		},
		{
			args: []string{"-samplerate", "22.05"},
			check: func(t *testing.T, cfg Config) {
				assert.Equal(t, 22050.0, cfg.SampleRate)
			},
		},
		{
			args: []string{"-disable-single-sample"},
			check: func(t *testing.T, cfg Config) {
				assert.True(t, cfg.SingleSampleDisabled)
			},
		},
		{
			args: []string{"-fade-in", "200"},
			check: func(t *testing.T, cfg Config) {
				assert.Equal(t, 200, cfg.FadeIn)
			},
		},
		{
			args: []string{"-device-list"},
			check: func(t *testing.T, cfg Config) {
				assert.True(t, cfg.DeviceList)
			},
		},
		{
			args: []string{"-device-in", "2"},
			check: func(t *testing.T, cfg Config) {
				assert.Equal(t, 2, cfg.DeviceIn)
			},
		},
		{
			args: []string{"-device-out", "3"},
			check: func(t *testing.T, cfg Config) {
				assert.Equal(t, 3, cfg.DeviceOut)
			},
		},
		{
			args: []string{"-device-latency", "high"},
			check: func(t *testing.T, cfg Config) {
				assert.Equal(t, "high", cfg.DeviceLatency)
			},
		},
		{
			args: []string{"-device-frame", "1024"},
			check: func(t *testing.T, cfg Config) {
				assert.Equal(t, 1024, cfg.DeviceFrameSize)
			},
		},
		{
			args: []string{"/path/to/patch.lisp"},
			check: func(t *testing.T, cfg Config) {
				assert.Equal(t, "/path/to/patch.lisp", cfg.ScriptPath)
			},
		},
		{
			args: []string{""},
			check: func(t *testing.T, cfg Config) {
				assert.Equal(t, "", cfg.ScriptPath)
			},
		},
		{
			args: []string{"-gain", "-6"},
			check: func(t *testing.T, cfg Config) {
				assert.Equal(t, -6.0, cfg.Gain)
			},
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, "_"), func(t *testing.T) {
			cfg, err := parseArgs(tt.args)
			assert.NoError(t, err)
			tt.check(t, cfg)
		})
	}
}

func TestParseArgs_Erroneous(t *testing.T) {
	var tests = []struct {
		name string
		args []string
	}{
		{
			name: "empty http addr",
			args: []string{"-addr", ""},
		},
		{
			name: "frame size greater than device frame size",
			args: []string{"-frame", "1024", "-device-frame", "256"},
		},
		{
			name: "frame size not a multiple of device frame size",
			args: []string{"-frame", "100", "-device-frame", "1024"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseArgs(tt.args)
			assert.Error(t, err)
		})
	}
}
